package upgrade

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AssetName_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		goos    string
		goarch  string
		want    string
		wantErr bool
	}{
		{"linux amd64", osLinux, archAmd64, "verv_linux_amd64", false},
		{"darwin arm64", osDarwin, archArm64, "verv_darwin_arm64", false},
		{"windows amd64", "windows", archAmd64, "", true},
		{"linux 386", osLinux, "386", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := assetName(tc.goos, tc.goarch)
			if tc.wantErr {
				require.ErrorIs(t, err, errUnsupportedTarget)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_AssetUrl_BuildsReleaseDownloadPath(t *testing.T) {
	t.Parallel()

	got := assetUrl("v0.0.52", "verv_linux_amd64")

	require.Equal(t, "https://github.com/vervstack/verv/releases/download/v0.0.52/verv_linux_amd64", got)
}
