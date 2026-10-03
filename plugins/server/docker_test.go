package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseOsRelease_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		content      string
		wantId       string
		wantCodename string
	}{
		{"quoted values", "ID=\"ubuntu\"\nVERSION_CODENAME=\"noble\"\n", "ubuntu", "noble"},
		{"unquoted values", "NAME=Ubuntu\nID=ubuntu\nVERSION_CODENAME=jammy\n", "ubuntu", "jammy"},
		{"missing codename", "ID=debian\n", "debian", ""},
		{"ignores similar keys", "ID_LIKE=debian\nVERSION_ID=\"24.04\"\n", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, codename := parseOsRelease(tc.content)

			require.Equal(t, tc.wantId, id)
			require.Equal(t, tc.wantCodename, codename)
		})
	}
}

func Test_DockerRepoLine_FormatsAptSource(t *testing.T) {
	t.Parallel()

	got := dockerRepoLine("amd64", "noble")

	want := "deb [arch=amd64 signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu noble stable\n"
	require.Equal(t, want, got)
}
