package nodeclient_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
)

type recordedRequest struct {
	method string
	path   string
	header string
	body   string
}

func newRecordingServer(t *testing.T, status int, response string) (*httptest.Server, *recordedRequest) {
	t.Helper()

	recorded := &recordedRequest{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		recorded.method = r.Method
		recorded.path = r.URL.EscapedPath()
		recorded.header = r.Header.Get("Grpc-Metadata-Authorization")
		recorded.body = string(raw)

		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server, recorded
}

func Test_Version_SendsAuthAndDecodes(t *testing.T) {
	server, recorded := newRecordingServer(t, http.StatusOK, `{"version":"1.2.3","os":"linux","arch":"arm64"}`)
	client := nodeclient.New(server.URL+"/", "abc")

	platform, err := client.Version(context.Background())
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, recorded.method)
	require.Equal(t, "/api/version", recorded.path)
	require.Equal(t, "abc", recorded.header)
	require.Equal(t, "1.2.3", platform.Version)
	require.Equal(t, "linux", platform.Os)
	require.Equal(t, "arm64", platform.Arch)
}

func Test_Version_OldNodeWithoutPlatform(t *testing.T) {
	server, _ := newRecordingServer(t, http.StatusOK, `{"version":"0.1.0"}`)
	client := nodeclient.New(server.URL, "abc")

	platform, err := client.Version(context.Background())
	require.NoError(t, err)
	require.Equal(t, "0.1.0", platform.Version)
	require.Empty(t, platform.Os)
	require.Empty(t, platform.Arch)
}

func Test_ListRegistryNames_PostsEmptyBodyAndKeepsOrder(t *testing.T) {
	server, recorded := newRecordingServer(t, http.StatusOK, `{"instances":[{"name":"b"},{"name":"a"}]}`)
	client := nodeclient.New(server.URL, "abc")

	names, err := client.ListRegistryNames(context.Background())
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, recorded.method)
	require.Equal(t, "/api/container-registry/list", recorded.path)
	require.JSONEq(t, `{}`, recorded.body)
	require.Equal(t, []string{"b", "a"}, names)
}

func Test_RegistryCredentials_DecodesBothKeyStyles(t *testing.T) {
	cases := []struct {
		name     string
		response string
	}{
		{"camel case", `{"username":"u","password":"p","registryUrl":"reg.example.com"}`},
		{"snake case", `{"username":"u","password":"p","registry_url":"reg.example.com"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, _ := newRecordingServer(t, http.StatusOK, tc.response)
			client := nodeclient.New(server.URL, "abc")

			credentials, err := client.RegistryCredentials(context.Background(), "main")
			require.NoError(t, err)
			require.Equal(t, "u", credentials.Username)
			require.Equal(t, "p", credentials.Password)
			require.Equal(t, "reg.example.com", credentials.RegistryUrl)
		})
	}
}

func Test_RegistryCredentials_EscapesName(t *testing.T) {
	server, recorded := newRecordingServer(t, http.StatusOK, `{}`)
	client := nodeclient.New(server.URL, "abc")

	_, err := client.RegistryCredentials(context.Background(), "team/main")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, recorded.method)
	require.Equal(t, "/api/container-registry/team%2Fmain/credentials", recorded.path)
}

func Test_CreateDeploy_SendsSnakeCaseBody(t *testing.T) {
	server, recorded := newRecordingServer(t, http.StatusOK, `ignored, not json`)
	client := nodeclient.New(server.URL, "abc")

	request := nodeclient.DeployRequest{
		ServiceName: "api",
		Environment: "prod",
		Image:       "reg.example.com/api:1",
	}
	err := client.CreateDeploy(context.Background(), request)
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, recorded.method)
	require.Equal(t, "/api/service/deploy/create", recorded.path)
	require.Equal(t, "abc", recorded.header)
	expectedBody := `{"service_name":"api","environment":"prod","vervonomicon":{"image":"reg.example.com/api:1"}}`
	require.JSONEq(t, expectedBody, recorded.body)
}

func Test_Requests_MapErrorStatuses(t *testing.T) {
	cases := []struct {
		name   string
		status int
		want   error
	}{
		{"unauthorized", http.StatusUnauthorized, nodeclient.ErrUnauthorized},
		{"forbidden", http.StatusForbidden, nodeclient.ErrUnauthorized},
		{"server error", http.StatusInternalServerError, nodeclient.ErrUnexpectedStatus},
		{"not found", http.StatusNotFound, nodeclient.ErrUnexpectedStatus},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, _ := newRecordingServer(t, tc.status, `secret body text`)
			client := nodeclient.New(server.URL, "abc")

			_, err := client.Version(context.Background())
			require.ErrorIs(t, err, tc.want)
			require.NotContains(t, err.Error(), "secret body text")

			request := nodeclient.DeployRequest{}
			err = client.CreateDeploy(context.Background(), request)
			require.ErrorIs(t, err, tc.want)
		})
	}
}
