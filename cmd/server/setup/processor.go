package setup

import (
	"os"
	"runtime"

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

	userNameEnv  = "USER_NAME"
	userPwdEnv   = "USER_PWD"
	sshKeyUrlEnv = "SSH_KEY_URL"

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
		Short: "Prepares this server: deploy user, ssh, docker",
		Long: "Creates the deploy user, enables public key ssh authentication, installs Docker " +
			"and the user's ssh keys. Ubuntu only, must run as root",

		RunE: proc.run,

		Annotations: map[string]string{"verv:emoji": "🖥️", "verv:group": menu.GroupVelez},

		Hidden: runtime.GOOS != "linux",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().String(UserFlag, defaultUserName, "name of the user to create (env USER_NAME)")
	c.Flags().String(SshKeyUrlFlag, defaultSshKeyUrl, "url of the public ssh keys to authorize (env SSH_KEY_URL)")

	return c
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

	opts.UserName = resolveValue(cmd.Flags().Changed(UserFlag), userFlag, os.Getenv(userNameEnv))
	opts.SshKeyUrl = resolveValue(cmd.Flags().Changed(SshKeyUrlFlag), keyUrlFlag, os.Getenv(sshKeyUrlEnv))
	opts.Password = os.Getenv(userPwdEnv)

	if opts.Password != "" {
		return opts, false, nil
	}

	opts.Password, aborted, err = promptPassword(p.io, opts.UserName)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error prompting for password")
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
