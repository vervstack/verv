package setup

import (
	"context"
	"os"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/server"
)

const (
	fetchingVersionsLabel = "Fetching sysbox versions from GitHub"
)

// resolveSysbox fills the options only the sysbox step needs. All prompts are skipped once sysbox is set up
// or skipped.
func (p *serverSetup) resolveSysbox(opts server.Options) (server.Options, bool, error) {
	ctx := context.Background()

	if opts.IsSysboxSkipped || server.IsSysboxSetUp(ctx) {
		return opts, false, nil
	}

	if opts.SysboxVersion == "" && isStdinTerminal() {
		isConfirmed, aborted, err := confirmSysboxInstall(p.io)
		if err != nil {
			return server.Options{}, false, rerrors.Wrap(err, "error confirming sysbox install")
		}

		if aborted {
			return opts, true, nil
		}

		if !isConfirmed {
			opts.IsSysboxSkipped = true

			return opts, false, nil
		}

		version, aborted, err := p.pickSysboxVersion(ctx)
		if err != nil {
			return server.Options{}, false, rerrors.Wrap(err, "error picking sysbox version")
		}

		if aborted {
			return opts, true, nil
		}

		opts.SysboxVersion = version
	}

	containers := server.CountRunningContainers(ctx)
	if containers == 0 {
		return opts, false, nil
	}

	if os.Getenv(assumeYesEnv) == "1" {
		opts.IsDockerRestartAllowed = true

		return opts, false, nil
	}

	isAllowed, aborted, err := confirmDockerRestart(containers)
	if err != nil {
		return server.Options{}, false, rerrors.Wrap(err, "error confirming docker restart")
	}

	if aborted {
		return opts, true, nil
	}

	opts.IsDockerRestartAllowed = isAllowed

	return opts, false, nil
}

// pickSysboxVersion falls back to the latest version (empty) when the releases cannot be fetched.
func (p *serverSetup) pickSysboxVersion(ctx context.Context) (version string, aborted bool, err error) {
	spinner := io.NewSpinner(p.io)
	spinner.Start(fetchingVersionsLabel)

	versions := fetchSysboxVersionsOrNil(ctx)
	if len(versions) == 0 {
		spinner.StopInfo(fetchingVersionsLabel + " — failed")

		p.io.PrintlnColored(colors.ColorYellow, "⚠ could not fetch sysbox versions, the latest will be installed")

		return "", false, nil
	}

	spinner.Stop(true, "Sysbox versions fetched")

	version, aborted, err = promptSysboxVersion(p.io, versions)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error prompting for sysbox version")
	}

	return version, aborted, nil
}

func fetchSysboxVersionsOrNil(ctx context.Context) []string {
	versions, err := server.FetchSysboxVersions(ctx)
	if err != nil {
		return nil
	}

	return versions
}

func isStdinTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
