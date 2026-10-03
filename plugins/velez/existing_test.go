package velez

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const inspectOutput = `[{
  "Config": {"Image": "vervstack/velez:v1.4.2"},
  "HostConfig": {"PortBindings": {"53890/tcp": [{"HostIp": "", "HostPort": "8080"}]}},
  "Mounts": [
    {"Source": "/var/run/docker.sock", "Destination": "/var/run/docker.sock"},
    {"Source": "/home/user/velez", "Destination": "/tmp/velez"}
  ]
}]`

func Test_ParseInspect_ReusesPortKeyPathAndVersion(t *testing.T) {
	got, err := parseInspect(inspectOutput)
	require.NoError(t, err)

	require.Equal(t, Existing{Version: "v1.4.2", Port: 8080, KeyPath: "/home/user/velez"}, got)
}

func Test_ParseInspect_Failures(t *testing.T) {
	cases := []struct {
		name   string
		output string
		want   error
	}{
		{"empty list", `[]`, errEmptyInspect},
		{
			"no port binding",
			`[{"Config": {"Image": "vervstack/velez:v1"}, "HostConfig": {"PortBindings": {}}, "Mounts": []}]`,
			errNoPortBinding,
		},
		{
			"no keys mount",
			`[{"Config": {"Image": "vervstack/velez:v1"},
			  "HostConfig": {"PortBindings": {"53890/tcp": [{"HostPort": "8080"}]}}, "Mounts": []}]`,
			errNoKeyMount,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseInspect(tc.output)
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func Test_ParseInspect_InvalidJson(t *testing.T) {
	_, err := parseInspect("not json")
	require.Error(t, err)
}

func Test_ToExisting_ImageWithoutTagHasEmptyVersion(t *testing.T) {
	var c inspectedContainer

	c.Config.Image = "vervstack/velez"
	c.HostConfig.PortBindings = map[string][]struct {
		HostPort string `json:"HostPort"`
	}{"53890/tcp": {{HostPort: "9000"}}}
	c.Mounts = []struct {
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
	}{{Source: "/k", Destination: "/tmp/velez"}}

	got, err := toExisting(c)
	require.NoError(t, err)

	require.Equal(t, Existing{Version: "", Port: 9000, KeyPath: "/k"}, got)
}
