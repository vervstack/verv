package server

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/rerrors"
)

const (
	version066 = "0.6.6"
	version070 = "0.7.0"

	releasesFixture = `[
  {"tag_name": "v0.6.6", "draft": false, "prerelease": false, "assets": [
    {"name": "sysbox-ce_0.6.6-0.linux_amd64.deb", "browser_download_url": "https://x/0.6.6-amd64.deb"},
    {"name": "sysbox-ce_0.6.6-0.linux_arm64.deb", "browser_download_url": "https://x/0.6.6-arm64.deb"}
  ]},
  {"tag_name": "v0.7.0", "draft": false, "prerelease": false, "assets": [
    {"name": "sysbox-ce_0.7.0.linux_amd64.deb", "browser_download_url": "https://x/0.7.0-amd64.deb"},
    {"name": "sysbox-ce_0.7.0.linux_arm64.deb", "browser_download_url": "https://x/0.7.0-arm64.deb"},
    {"name": "sysbox-ce_0.7.0.tar.gz", "browser_download_url": "https://x/0.7.0.tar.gz"}
  ]},
  {"tag_name": "v0.6.10", "draft": false, "prerelease": false, "assets": [
    {"name": "sysbox-ce_0.6.10.linux_amd64.deb", "browser_download_url": "https://x/0.6.10-amd64.deb"}
  ]},
  {"tag_name": "v0.8.0", "draft": true, "prerelease": false, "assets": []},
  {"tag_name": "v0.9.0", "draft": false, "prerelease": true, "assets": []},
  {"tag_name": "v0.7.0-rc1", "draft": false, "prerelease": false, "assets": []},
  {"tag_name": "latest", "draft": false, "prerelease": false, "assets": []}
]`
)

func Test_IsKernelSupported_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		release string
		want    bool
	}{
		{"5.4 is too old", "5.4.0", false},
		{"5.11 is too old", "5.11.22", false},
		{"5.12 is the minimum", "5.12.0", true},
		{"ubuntu style release", "5.15.0-91-generic\n", true},
		{"6.8 is newer", "6.8.0", true},
		{"garbage", "not-a-kernel", false},
		{"empty release", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isKernelSupported(tc.release))
		})
	}
}

func Test_HasSysboxRuntime_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		output string
		want   bool
	}{
		{"registered", "io.containerd.runc.v2 runc sysbox-runc \n", true},
		{"only runc", "io.containerd.runc.v2 runc \n", false},
		{"similar name", "sysbox-runc-old \n", false},
		{"empty output", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, hasSysboxRuntime(tc.output))
		})
	}
}

func Test_IsRootlessSecurityOptions_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opts string
		want bool
	}{
		{"rootful", `["name=apparmor","name=seccomp,profile=builtin","name=cgroupns"]`, false},
		{"rootless", `["name=seccomp,profile=builtin","name=rootless"]`, true},
		{"empty list", `[]`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isRootlessSecurityOptions(tc.opts))
		})
	}
}

func Test_IsSnapDockerRoot_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		root string
		want bool
	}{
		{"default root", "/var/lib/docker\n", false},
		{"snap root", "/var/snap/docker/common/var-lib-docker\n", true},
		{"empty root", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isSnapDockerRoot(tc.root))
		})
	}
}

func Test_IsArchSupported_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		arch string
		want bool
	}{
		{"amd64", archAmd64 + "\n", true},
		{"arm64", archArm64, true},
		{"armhf", "armhf", false},
		{"i386", "i386", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isArchSupported(tc.arch))
		})
	}
}

func Test_CountNonEmptyLines_Scenarios(t *testing.T) {
	t.Parallel()

	require.Equal(t, 0, countNonEmptyLines(""))
	require.Equal(t, 0, countNonEmptyLines("\n\n"))
	require.Equal(t, 2, countNonEmptyLines("abc123\ndef456\n"))
}

func Test_ParseSysboxReleases_FiltersAndSortsBySemver(t *testing.T) {
	t.Parallel()

	releases, err := parseSysboxReleases([]byte(releasesFixture))
	require.NoError(t, err)

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		versions = append(versions, release.Version)
	}

	require.Equal(t, []string{version070, "0.6.10", version066}, versions)
}

func Test_ParseSysboxReleases_InvalidJson(t *testing.T) {
	t.Parallel()

	_, err := parseSysboxReleases([]byte("{"))
	require.Error(t, err)
}

func Test_FindSysboxDebUrl_Scenarios(t *testing.T) {
	t.Parallel()

	releases, err := parseSysboxReleases([]byte(releasesFixture))
	require.NoError(t, err)

	byVersion := map[string]sysboxRelease{}
	for _, release := range releases {
		byVersion[release.Version] = release
	}

	cases := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr error
	}{
		{"plain name amd64", version070, archAmd64, "https://x/0.7.0-amd64.deb", nil},
		{"plain name arm64", version070, archArm64, "https://x/0.7.0-arm64.deb", nil},
		{"dash zero name", version066, archAmd64, "https://x/0.6.6-amd64.deb", nil},
		{"missing arch", "0.6.10", archArm64, "", errSysboxAssetMissing},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := findSysboxDebUrl(byVersion[tc.version], tc.arch)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_SelectSysboxRelease_Scenarios(t *testing.T) {
	t.Parallel()

	releases, err := parseSysboxReleases([]byte(releasesFixture))
	require.NoError(t, err)

	cases := []struct {
		name    string
		version string
		want    string
		wantErr error
	}{
		{"empty is newest", "", version070, nil},
		{"bare version", version066, version066, nil},
		{"v prefixed version", "v" + version066, version066, nil},
		{"unknown version", "9.9.9", "", errSysboxVersionUnknown},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := selectSysboxRelease(releases, tc.version)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got.Version)
		})
	}
}

func Test_SelectSysboxRelease_NoReleases(t *testing.T) {
	t.Parallel()

	_, err := selectSysboxRelease(nil, "")
	require.ErrorIs(t, err, errSysboxNoReleases)
}

func Test_NormalizeSysboxVersion_Scenarios(t *testing.T) {
	t.Parallel()

	require.Equal(t, normalizeSysboxVersion("0.7.0"), normalizeSysboxVersion("v0.7.0"))
	require.Equal(t, "0.7.0", normalizeSysboxVersion(" v0.7.0 "))
	require.Empty(t, normalizeSysboxVersion(""))
}

func Test_SkipError_FoundThroughWrap(t *testing.T) {
	t.Parallel()

	err := newSkipError("kernel is too old")

	err = rerrors.Wrap(err, "error during step")

	var skip *skipError

	require.ErrorAs(t, err, &skip)
	require.Equal(t, "kernel is too old", skip.reason)
}
