package version

import (
	"runtime/debug"
	"strings"
)

const (
	devVersion = "dev"
)

// version is set at release build time: -ldflags "-X go.vervstack.ru/verv/version.version=<tag>".
var version string

// GetVersion returns the release tag baked in at build time. A `go install ...@vX` build carries
// no ldflags, so it falls back to the module version Go embeds; any other build reports "dev".
func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		info = &debug.BuildInfo{}
	}

	return resolveVersion(version, info.Main.Version)
}

// resolveVersion prefers the baked-in tag, then a clean tagged module version. Pseudo-versions
// ("v0.0.50-0.2026...") and dirty builds ("+dirty") are local builds, not releases.
func resolveVersion(baked, moduleVersion string) string {
	if baked != "" {
		return baked
	}

	if moduleVersion == "" || moduleVersion == "(devel)" || strings.ContainsAny(moduleVersion, "-+") {
		return devVersion
	}

	return moduleVersion
}
