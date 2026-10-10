package setup

import (
	"context"
	"os"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/server"
)

const (
	UserFlag      = "user"
	SshKeyUrlFlag = "ssh-key-url"

	SysboxVersionFlag = "sysbox-version"
	SkipSysboxFlag    = "skip-sysbox"

	userNameEnv      = "USER_NAME"
	userPwdEnv       = "USER_PWD"
	sshKeyUrlEnv     = "SSH_KEY_URL"
	sysboxVersionEnv = "SYSBOX_VERSION"
	skipSysboxEnv    = "SKIP_SYSBOX"
	assumeYesEnv     = "SETUP_ASSUME_YES"

	defaultUserName  = "deployer"
	defaultSshKeyUrl = "https://github.com/alexskilled.keys"
)

type serverSetup struct {
	io        io.IO
	deployCmd *cobra.Command
}

// NewCommand builds the `verv setup-server` command. deployCmd is the
// deploy-velez command offered once the server is prepared.
func NewCommand(basicProc processor.Processor, deployCmd *cobra.Command) *cobra.Command {
	proc := serverSetup{
		io:        basicProc.IO,
		deployCmd: deployCmd,
	}

	c := &cobra.Command{
		Use:   "setup-server",
		Short: "Prepares this server: deploy user, ssh, docker, sysbox",
		Long: "Creates the deploy user, enables public key ssh authentication, installs Docker, " +
			"sysbox and the user's ssh keys. Ubuntu only, must run as root",

		RunE: proc.run,

		Annotations: map[string]string{
			"verv:emoji":              "🖥️",
			"verv:group":              menu.GroupVelez,
			"verv:showsSetupProgress": "true",
		},

		Hidden: runtime.GOOS != "linux",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().String(UserFlag, defaultUserName, "name of the user to create (env USER_NAME)")
	c.Flags().String(SshKeyUrlFlag, defaultSshKeyUrl, "url of the public ssh keys to authorize (env SSH_KEY_URL)")
	c.Flags().String(SysboxVersionFlag, "", "sysbox version to install, empty for the latest (env SYSBOX_VERSION)")
	c.Flags().Bool(SkipSysboxFlag, false, "do not install sysbox (env SKIP_SYSBOX)")

	return c
}

// SetupProgress reports how many setup-server steps are complete for the default (or env-configured)
// deploy user; (0, 0) off Linux, where setup-server does not apply.
func SetupProgress(ctx context.Context) (done, total int) {
	if runtime.GOOS != "linux" {
		return 0, 0
	}

	opts := server.Options{
		UserName:  resolveValue(false, defaultUserName, os.Getenv(userNameEnv)),
		SshKeyUrl: resolveValue(false, defaultSshKeyUrl, os.Getenv(sshKeyUrlEnv)),
	}

	return server.Progress(ctx, opts)
}

func (p *serverSetup) run(cmd *cobra.Command, _ []string) error {
	opts, aborted, err := p.resolveOptions(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error resolving options")
	}

	if aborted {
		return nil
	}

	err = server.Setup(p.io, opts)
	if err != nil {
		return rerrors.Wrap(err, "error setting up server")
	}

	isConfirmed, aborted, err := confirmDeploy()
	if err != nil {
		return rerrors.Wrap(err, "error confirming velez deploy")
	}

	if aborted || !isConfirmed {
		return nil
	}

	p.deployCmd.SetContext(cmd.Context())

	err = p.deployCmd.RunE(p.deployCmd, nil)
	if err != nil {
		return rerrors.Wrap(err, "error deploying velez")
	}

	return nil
}

func (p *serverSetup) resolveOptions(cmd *cobra.Command) (opts server.Options, aborted bool, err error) {
	userFlag, err := cmd.Flags().GetString(UserFlag)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error reading user flag")
	}

	keyUrlFlag, err := cmd.Flags().GetString(SshKeyUrlFlag)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error reading ssh-key-url flag")
	}

	versionFlag, err := cmd.Flags().GetString(SysboxVersionFlag)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error reading sysbox-version flag")
	}

	skipFlag, err := cmd.Flags().GetBool(SkipSysboxFlag)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error reading skip-sysbox flag")
	}

	opts.UserName = resolveValue(cmd.Flags().Changed(UserFlag), userFlag, os.Getenv(userNameEnv))
	opts.SshKeyUrl = resolveValue(cmd.Flags().Changed(SshKeyUrlFlag), keyUrlFlag, os.Getenv(sshKeyUrlEnv))

	isVersionChanged := cmd.Flags().Changed(SysboxVersionFlag)

	opts.SysboxVersion = resolveValue(isVersionChanged, versionFlag, os.Getenv(sysboxVersionEnv))
	opts.IsSysboxSkipped = resolveBool(cmd.Flags().Changed(SkipSysboxFlag), skipFlag, os.Getenv(skipSysboxEnv))
	opts.Password = os.Getenv(userPwdEnv)

	if opts.Password == "" {
		isExistingUser := server.UserExists(opts.UserName)

		opts.Password, aborted, err = promptPassword(p.io, opts.UserName, isExistingUser)
		if err != nil {
			return server.Options{}, false, rerrors.Wrap(err, "error prompting for password")
		}

		if aborted {
			return opts, true, nil
		}
	}

	opts, aborted, err = p.resolveSysbox(opts)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error resolving sysbox options")
	}

	return opts, aborted, nil
}

func resolveValue(isFlagChanged bool, flagValue, envValue string) string {
	if isFlagChanged {
		return flagValue
	}

	if envValue != "" {
		return envValue
	}

	return flagValue
}

// resolveBool mirrors resolveValue for a bool flag; an env value strconv.ParseBool rejects counts as unset.
func resolveBool(isFlagChanged, flagValue bool, envValue string) bool {
	if isFlagChanged {
		return flagValue
	}

	envBool, err := strconv.ParseBool(envValue)
	if err != nil {
		return flagValue
	}

	return envBool
}
