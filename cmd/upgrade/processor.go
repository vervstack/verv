package upgrade

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	vervio "go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/version"
)

const (
	downloadTimeout = 2 * time.Minute
	binaryMode      = 0o755
	tempPattern     = ".verv-upgrade-*"
)

type upgrader struct {
	io        vervio.IO
	latestTag string
}

// NewCommand builds `verv upgrade`. latestTag is the newer release tag, or "" when already up to
// date — in that case the command is hidden from the picker (Hidden only affects listing, so
// `verv upgrade` still works by name and reports there is nothing to do).
func NewCommand(basicProc processor.Processor, latestTag string) *cobra.Command {
	proc := upgrader{
		io:        basicProc.IO,
		latestTag: latestTag,
	}

	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades verv to the latest release",
		Long:  "Downloads the latest release binary for this OS and architecture and replaces the running verv binary",

		RunE: proc.run,

		Annotations: map[string]string{"verv:emoji": "⬆️", "verv:group": menu.GroupSettings},

		Hidden: latestTag == "",

		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

func (u *upgrader) run(cmd *cobra.Command, _ []string) error {
	if u.latestTag == "" {
		u.io.Println("Already on the latest version (" + version.GetVersion() + ")")

		return nil
	}

	asset, err := assetName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return rerrors.Wrap(err)
	}

	dest, err := runningBinaryPath()
	if err != nil {
		return rerrors.Wrap(err, "error locating running binary")
	}

	u.io.Println(fmt.Sprintf(
		"Upgrading %s -> %s for %s/%s, replacing %s",
		version.GetVersion(), u.latestTag, runtime.GOOS, runtime.GOARCH, dest,
	))

	spinner := vervio.NewSpinner(u.io)
	spinner.Start("Downloading " + asset)

	err = install(cmd.Context(), assetUrl(u.latestTag, asset), dest)
	if err != nil {
		spinner.Stop(false, "Downloading "+asset+" — failed")

		return rerrors.Wrap(err, "error installing release")
	}

	spinner.Stop(true, "Upgraded to "+u.latestTag)

	return nil
}

func runningBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", rerrors.Wrap(err, "error getting executable path")
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving executable symlinks")
	}

	return resolved, nil
}

// install downloads url into a temp file beside dest and renames it over dest. The rename is
// atomic and safe for a running binary on unix; the temp file is removed on any failure.
func install(ctx context.Context, url, dest string) error {
	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return rerrors.Wrap(err, "error building download request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rerrors.Wrap(err, "error downloading release asset")
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return rerrors.Wrap(errDownloadStatus, fmt.Sprintf("status: %d", resp.StatusCode))
	}

	tmp, err := os.CreateTemp(filepath.Dir(dest), tempPattern)
	if err != nil {
		return rerrors.Wrap(err, "error creating temp file next to the binary (is the directory writable?)")
	}

	tmpPath := tmp.Name()

	err = writeExecutable(tmp, resp.Body)
	if err != nil {
		_ = os.Remove(tmpPath)

		return rerrors.Wrap(err, "error writing release asset")
	}

	err = os.Rename(tmpPath, dest)
	if err != nil {
		_ = os.Remove(tmpPath)

		return rerrors.Wrap(err, "error replacing binary")
	}

	return nil
}

func writeExecutable(f *os.File, body io.Reader) error {
	_, err := io.Copy(f, body)
	if err != nil {
		_ = f.Close()

		return rerrors.Wrap(err, "error copying body")
	}

	err = f.Close()
	if err != nil {
		return rerrors.Wrap(err, "error closing temp file")
	}

	err = os.Chmod(f.Name(), binaryMode)
	if err != nil {
		return rerrors.Wrap(err, "error making temp file executable")
	}

	return nil
}
