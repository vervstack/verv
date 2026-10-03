package deploy

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/velez"
)

var errNoExistingContainer = rerrors.New("no velez container found on this machine — use deploy-velez")

type velezUpgrade struct {
	velezDeploy
}

// NewUpgradeCommand builds the `verv upgrade-velez` command. It replaces the
// running Velez container with a newer image tag, reusing the host port and
// keys path of the existing container so the user is only asked for a version.
// It starts hidden: the interactive picker un-hides it, and hides deploy-velez,
// only when a Velez container already exists.
func NewUpgradeCommand(basicProc processor.Processor) *cobra.Command {
	proc := velezUpgrade{velezDeploy{io: basicProc.IO}}

	c := &cobra.Command{
		Use:   "upgrade-velez",
		Short: "Upgrades the Velez node running on this machine",
		Long:  "Pulls a chosen vervstack/velez image tag and restarts the existing container with its current port and keys path",

		RunE: proc.run,

		Annotations: map[string]string{"verv:emoji": "⬆️", "verv:group": menu.GroupVelez},

		Hidden: true,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().Bool(DebugFlag, false, "allow running outside Linux, for local dev/testing")

	return c
}

func (p *velezUpgrade) run(cmd *cobra.Command, _ []string) error {
	debug, err := cmd.Flags().GetBool(DebugFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading debug flag")
	}

	if !platformAllowed(runtime.GOOS, debug) {
		return rerrors.Wrap(errUnsupportedPlatform)
	}

	existing, found, err := velez.FindExisting()
	if err != nil {
		return rerrors.Wrap(err, "error reading existing velez container")
	}

	if !found {
		return rerrors.Wrap(errNoExistingContainer)
	}

	tags, err := p.fetchTags(cmd.Context())
	if err != nil {
		return rerrors.Wrap(err, "error fetching velez tags")
	}

	title := fmt.Sprintf("Which Velez version would you like to upgrade to? (current: %s)", existing.Version)

	version, aborted, err := selectVersion(p.io, title, tags)
	if err != nil {
		return rerrors.Wrap(err, "error selecting velez version")
	}

	if aborted {
		return nil
	}

	err = velez.Deploy(p.io, version, existing.Port, existing.KeyPath)
	if err != nil {
		return rerrors.Wrap(err, "error upgrading velez")
	}

	return nil
}

// PreferUpgrade swaps deploy-velez for upgrade-velez in the interactive picker
// when a Velez container already runs on this machine. Any failure to read the
// container leaves deploy-velez in place.
func PreferUpgrade(deployCmd, upgradeCmd *cobra.Command) {
	if runtime.GOOS != "linux" {
		return
	}

	_, found, err := velez.FindExisting()
	if err != nil || !found {
		return
	}

	deployCmd.Hidden = true
	upgradeCmd.Hidden = false
}
