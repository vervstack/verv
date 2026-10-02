package upgrade

import (
	"go.redsock.ru/rerrors"
)

const (
	releaseDownloadUrl = "https://github.com/vervstack/verv/releases/download/"

	osLinux   = "linux"
	osDarwin  = "darwin"
	archAmd64 = "amd64"
	archArm64 = "arm64"
)

// supportedOs and supportedArch mirror the release workflow's build matrix.
var (
	supportedOs   = map[string]bool{osLinux: true, osDarwin: true}
	supportedArch = map[string]bool{archAmd64: true, archArm64: true}
)

// assetName returns the release asset published for the given target. The target comes from
// runtime.GOOS / runtime.GOARCH — constants baked in at compile time by the release matrix's
// GOOS/GOARCH — so the binary always asks for exactly the build it already is.
func assetName(goos, goarch string) (string, error) {
	if !supportedOs[goos] || !supportedArch[goarch] {
		return "", rerrors.Wrap(errUnsupportedTarget, goos+"/"+goarch)
	}

	return "verv_" + goos + "_" + goarch, nil
}

func assetUrl(tag, asset string) string {
	return releaseDownloadUrl + tag + "/" + asset
}
