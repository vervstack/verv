package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const (
	sysboxRuntimeName = "sysbox-runc"

	dockerBin = "docker"
	archAmd64 = "amd64"
	archArm64 = "arm64"

	dockerRuntimesFormat        = `{{range $k, $v := .Runtimes}}{{$k}} {{end}}`
	dockerRootDirFormat         = "{{.DockerRootDir}}"
	dockerSecurityOptionsFormat = "{{json .SecurityOptions}}"

	snapDockerRootPrefix = "/var/snap/"
	rootlessOption       = "rootless"

	minKernelMajor = 5
	minKernelMinor = 12

	declinedByUserReason = "declined by the user"

	sysboxDownloadTimeout = 10 * time.Minute
	dockerRestartTimeout  = 2 * time.Minute
	runtimeWaitTimeout    = 30 * time.Second
	runtimeWaitInterval   = 2 * time.Second

	debFileMode = 0o644
)

var (
	supportedArchs = []string{archAmd64, archArm64}

	kernelReleasePattern = regexp.MustCompile(`^(\d+)\.(\d+)`)
)

// skipError marks a step that cannot run on this machine; Setup warns and moves on instead of failing.
type skipError struct {
	reason string
}

func (e *skipError) Error() string {
	return e.reason
}

func newSkipError(reason string) error {
	skip := &skipError{reason: reason}

	return rerrors.Wrap(skip)
}

// IsSysboxSetUp reports whether Docker has the sysbox-runc runtime registered. Read-only, no root needed.
func IsSysboxSetUp(ctx context.Context) bool {
	return isSysboxSetUp(ctx, Options{})
}

// CountRunningContainers returns the number of running containers; 0 when Docker is missing or unreachable.
func CountRunningContainers(_ context.Context) int {
	req := cmd.Request{Tool: dockerBin, Args: []string{"ps", "-q"}}

	out, err := cmd.Execute(req)
	if err != nil {
		return 0
	}

	return countNonEmptyLines(out)
}

func isSysboxSetUp(_ context.Context, _ Options) bool {
	out, err := dockerInfo(dockerRuntimesFormat)
	if err != nil {
		return false
	}

	return hasSysboxRuntime(out)
}

func installSysbox(ctx context.Context, opts Options) (string, error) {
	arch, err := checkSysboxPreflight(ctx, opts)
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	releases, err := fetchSysboxReleases(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "error fetching sysbox releases")
	}

	release, err := selectSysboxRelease(releases, opts.SysboxVersion)
	if err != nil {
		return "", rerrors.Wrap(err, "error selecting sysbox release")
	}

	debUrl, err := findSysboxDebUrl(release, arch)
	if err != nil {
		return "", rerrors.Wrap(err, "error finding sysbox package")
	}

	debPath, err := downloadSysboxDeb(ctx, debUrl)
	if err != nil {
		return "", rerrors.Wrap(err, "error downloading sysbox package")
	}

	defer func() {
		_ = os.Remove(debPath)
	}()

	err = aptGet("update")
	if err != nil {
		return "", rerrors.Wrap(err, "error updating package index")
	}

	err = aptGet("install", "-y", debPath)
	if err != nil {
		return "", rerrors.Wrap(err, "error installing sysbox package")
	}

	err = ensureSysboxRuntime(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "error registering sysbox runtime")
	}

	return fmt.Sprintf("Sysbox %s installed (runtime %s registered)", release.Version, sysboxRuntimeName), nil
}

// checkSysboxPreflight returns the dpkg architecture, or a skipError when sysbox cannot be installed here.
func checkSysboxPreflight(ctx context.Context, opts Options) (arch string, err error) {
	if opts.IsSysboxSkipped {
		return "", newSkipError(declinedByUserReason)
	}

	arch, err = checkSysboxHost()
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	err = checkSysboxDocker()
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	count := CountRunningContainers(ctx)
	if count > 0 && !opts.IsDockerRestartAllowed {
		reason := fmt.Sprintf("%d containers are running and the Docker restart was not confirmed", count)

		return "", newSkipError(reason)
	}

	return arch, nil
}

func checkSysboxHost() (arch string, err error) {
	unameReq := cmd.Request{Tool: "uname", Args: []string{"-r"}}

	release, err := cmd.Execute(unameReq)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading kernel release")
	}

	if !isKernelSupported(release) {
		reason := fmt.Sprintf("kernel %s is older than %d.%d", strings.TrimSpace(release), minKernelMajor, minKernelMinor)

		return "", newSkipError(reason)
	}

	dpkgReq := cmd.Request{Tool: "dpkg", Args: []string{"--print-architecture"}}

	archOut, err := cmd.Execute(dpkgReq)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading architecture")
	}

	arch = strings.TrimSpace(archOut)
	if !isArchSupported(arch) {
		return "", newSkipError("architecture " + arch + " is not supported (amd64 and arm64 only)")
	}

	return arch, nil
}

func checkSysboxDocker() error {
	rootDir, err := dockerInfo(dockerRootDirFormat)
	if err != nil {
		return rerrors.Wrap(err, "error reading docker root dir")
	}

	if isSnapDockerRoot(rootDir) {
		return newSkipError("Docker is installed as a snap, which sysbox does not support")
	}

	securityOptions, err := dockerInfo(dockerSecurityOptionsFormat)
	if err != nil {
		return rerrors.Wrap(err, "error reading docker security options")
	}

	if isRootlessSecurityOptions(securityOptions) {
		return newSkipError("Docker runs rootless, which sysbox does not support")
	}

	return nil
}

func downloadSysboxDeb(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, sysboxDownloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", rerrors.Wrap(err, "error building request")
	}

	req.Header.Set("User-Agent", vervUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error calling server")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", rerrors.Wrap(errUnexpectedStatus, fmt.Sprintf("unexpected status: %d", resp.StatusCode))
	}

	file, err := os.CreateTemp("", "sysbox-*.deb")
	if err != nil {
		return "", rerrors.Wrap(err, "error creating temp file")
	}

	path, err := filepath.Abs(file.Name())
	if err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())

		return "", rerrors.Wrap(err, "error resolving temp file path")
	}

	err = copyToDeb(file, resp.Body)
	if err != nil {
		_ = os.Remove(path)

		return "", rerrors.Wrap(err, "error saving package")
	}

	return path, nil
}

// copyToDeb leaves the file world-readable so apt's unprivileged _apt user can read it.
func copyToDeb(file *os.File, body io.Reader) error {
	_, err := io.Copy(file, body)
	if err != nil {
		_ = file.Close()

		return rerrors.Wrap(err, "error streaming body")
	}

	err = file.Chmod(debFileMode)
	if err != nil {
		_ = file.Close()

		return rerrors.Wrap(err, "error setting file mode")
	}

	err = file.Close()
	if err != nil {
		return rerrors.Wrap(err, "error closing file")
	}

	return nil
}

// ensureSysboxRuntime verifies the runtime is registered, restarting Docker once if the package
// installer left that to us.
func ensureSysboxRuntime(ctx context.Context) error {
	if isSysboxSetUp(ctx, Options{}) {
		return nil
	}

	req := cmd.Request{Tool: "systemctl", Args: []string{"restart", dockerBin}, Timeout: dockerRestartTimeout}

	_, err := cmd.Execute(req)
	if err != nil {
		return rerrors.Wrap(err, "error restarting docker")
	}

	deadline := time.Now().Add(runtimeWaitTimeout)

	for time.Now().Before(deadline) {
		if isSysboxSetUp(ctx, Options{}) {
			return nil
		}

		select {
		case <-ctx.Done():
			return rerrors.Wrap(ctx.Err())
		case <-time.After(runtimeWaitInterval):
		}
	}

	return rerrors.Wrap(errSysboxRuntimeMissing)
}

func dockerInfo(format string) (string, error) {
	req := cmd.Request{Tool: dockerBin, Args: []string{"info", "--format", format}}

	out, err := cmd.Execute(req)
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	return out, nil
}

func hasSysboxRuntime(output string) bool {
	return slices.Contains(strings.Fields(output), sysboxRuntimeName)
}

func isKernelSupported(release string) bool {
	match := kernelReleasePattern.FindStringSubmatch(strings.TrimSpace(release))
	if match == nil {
		return false
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}

	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return false
	}

	return major > minKernelMajor || (major == minKernelMajor && minor >= minKernelMinor)
}

func isArchSupported(arch string) bool {
	return slices.Contains(supportedArchs, strings.TrimSpace(arch))
}

func isSnapDockerRoot(rootDir string) bool {
	return strings.HasPrefix(strings.TrimSpace(rootDir), snapDockerRootPrefix)
}

func isRootlessSecurityOptions(securityOptions string) bool {
	return strings.Contains(securityOptions, rootlessOption)
}

func countNonEmptyLines(out string) int {
	count := 0

	for line := range strings.SplitSeq(out, "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}

	return count
}
