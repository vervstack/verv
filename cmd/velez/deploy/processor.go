package deploy

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/velez"
)

const (
	PortFlag    = "port"
	KeyPathFlag = "key-path"
	DebugFlag   = "debug"

	defaultPort    = 53890
	defaultKeyPath = "~/velez"
)

var errUnsupportedPlatform = rerrors.New("deploy-velez requires linux (use --debug to override)")

type velezDeploy struct {
	io io.IO
}

// NewCommand builds the `verv deploy-velez` command. It is hidden from the
// interactive command picker on non-Linux platforms (Hidden only affects
// listing — `verv deploy-velez --debug` still works directly by name), and
// RunE re-checks the platform at runtime since Hidden is fixed at
// construction time and can't see the --debug flag yet.
func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := velezDeploy{
		io: basicProc.IO,
	}

	c := &cobra.Command{
		Use:   "deploy-velez",
		Short: "Deploys a Velez node on this machine",
		Long:  "Pulls a chosen vervstack/velez image tag and (re)starts it as a Docker container on this machine",

		RunE: proc.run,

		Annotations: map[string]string{"verv:emoji": "🛰️"},

		Hidden: runtime.GOOS != "linux",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().IntP(PortFlag, "p", defaultPort, "host port to bind the Velez node to")
	c.Flags().String(KeyPathFlag, defaultKeyPath, "path on this machine to store Velez's keys")
	c.Flags().Bool(DebugFlag, false, "allow running outside Linux, for local dev/testing")

	return c
}

func (p *velezDeploy) run(cmd *cobra.Command, _ []string) error {
	debug, err := cmd.Flags().GetBool(DebugFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading debug flag")
	}

	if !platformAllowed(runtime.GOOS, debug) {
		return rerrors.Wrap(errUnsupportedPlatform)
	}

	tags, err := p.fetchTags(cmd.Context())
	if err != nil {
		return rerrors.Wrap(err, "error fetching velez tags")
	}

	version, aborted, err := selectVersion(tags)
	if err != nil {
		return rerrors.Wrap(err, "error selecting velez version")
	}

	if aborted {
		return nil
	}

	port, aborted, err := p.resolvePort(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error resolving port")
	}

	if aborted {
		return nil
	}

	keyPath, aborted, err := p.resolveKeyPath(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error resolving key path")
	}

	if aborted {
		return nil
	}

	err = velez.Deploy(p.io, version, port, keyPath)
	if err != nil {
		return rerrors.Wrap(err, "error deploying velez")
	}

	return nil
}

func (p *velezDeploy) fetchTags(ctx context.Context) ([]string, error) {
	spinner := io.NewSpinner(p.io)
	spinner.Start("Fetching latest Velez versions from Docker Hub")

	tags, err := velez.FetchLatestTags(ctx)
	if err != nil {
		spinner.Stop(false, "Fetching latest Velez versions — failed")

		return nil, rerrors.Wrap(err)
	}

	spinner.Stop(true, fmt.Sprintf("Found %d Velez version(s)", len(tags)))

	return tags, nil
}

func (p *velezDeploy) resolvePort(cmd *cobra.Command) (port int, aborted bool, err error) {
	port, err = cmd.Flags().GetInt(PortFlag)
	if err != nil {
		return 0, false, rerrors.Wrap(err, "error reading port flag")
	}

	if cmd.Flags().Changed(PortFlag) {
		return port, false, nil
	}

	portStr, aborted, err := promptWithDefault("Host port for the Velez node", strconv.Itoa(port))
	if err != nil {
		return 0, false, rerrors.Wrap(err, "error prompting for port")
	}

	if aborted {
		return 0, true, nil
	}

	port, err = strconv.Atoi(strings.TrimSpace(portStr))
	if err != nil {
		return 0, false, rerrors.Wrap(err, "error parsing port")
	}

	return port, false, nil
}

func (p *velezDeploy) resolveKeyPath(cmd *cobra.Command) (keyPath string, aborted bool, err error) {
	keyPath, err = cmd.Flags().GetString(KeyPathFlag)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error reading key-path flag")
	}

	if !cmd.Flags().Changed(KeyPathFlag) {
		keyPath, aborted, err = promptWithDefault("Path to store Velez's keys", keyPath)
		if err != nil {
			return "", false, rerrors.Wrap(err, "error prompting for key path")
		}

		if aborted {
			return "", true, nil
		}
	}

	keyPath, err = expandKeyPath(keyPath)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error expanding key path")
	}

	return keyPath, false, nil
}
