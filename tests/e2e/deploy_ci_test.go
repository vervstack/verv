package e2e

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	deploySecret       = "s3cret"
	deployAuthHeader   = "Grpc-Metadata-Authorization"
	versionPath        = "/api/version"
	registryListPath   = "/api/container-registry/list"
	registryCredsInfix = "/credentials"
	deployCreatePath   = "/api/service/deploy/create"
)

type standInBehaviour struct {
	versionStatus int
	arch          string
	registries    []string
}

type velezStandIn struct {
	server      *httptest.Server
	behaviour   standInBehaviour
	mu          sync.Mutex
	paths       []string
	authHeaders []string
}

type deployResult struct {
	output   string
	exitCode int
}

func newForbiddenBehaviour() standInBehaviour {
	return standInBehaviour{versionStatus: http.StatusForbidden}
}

func newNoRegistryBehaviour() standInBehaviour {
	return standInBehaviour{versionStatus: http.StatusOK, arch: "amd64"}
}

func newNoArchBehaviour() standInBehaviour {
	return standInBehaviour{versionStatus: http.StatusOK, registries: []string{"main"}}
}

func newVelezStandIn(t *testing.T, behaviour standInBehaviour) *velezStandIn {
	t.Helper()

	standIn := &velezStandIn{behaviour: behaviour}
	standIn.server = httptest.NewServer(http.HandlerFunc(standIn.handle))

	t.Cleanup(standIn.server.Close)

	return standIn
}

func (s *velezStandIn) handle(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.paths = append(s.paths, r.URL.Path)
	s.authHeaders = append(s.authHeaders, r.Header.Get(deployAuthHeader))
	s.mu.Unlock()

	switch {
	case r.Method == http.MethodGet && r.URL.Path == versionPath:
		s.handleVersion(w)
	case r.Method == http.MethodPost && r.URL.Path == registryListPath:
		s.handleRegistryList(w)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *velezStandIn) handleVersion(w http.ResponseWriter) {
	if s.behaviour.versionStatus != http.StatusOK {
		w.WriteHeader(s.behaviour.versionStatus)

		return
	}

	body := map[string]string{"version": "v0.0.0-e2e", "os": "linux", "arch": s.behaviour.arch}

	writeStandInJson(w, body)
}

func (s *velezStandIn) handleRegistryList(w http.ResponseWriter) {
	instances := make([]map[string]string, 0, len(s.behaviour.registries))
	for _, name := range s.behaviour.registries {
		instances = append(instances, map[string]string{"name": name})
	}

	body := map[string]any{"instances": instances}

	writeStandInJson(w, body)
}

func writeStandInJson(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(body)
}

func (s *velezStandIn) recordedPaths() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string(nil), s.paths...)
}

func newDeployProject(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	vervDir := filepath.Join(dir, ".verv")

	err := os.MkdirAll(vervDir, 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(vervDir, "vervonomicon.yaml"), []byte("name: github.com/example/demo\n"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM scratch\n"), 0o644)
	require.NoError(t, err)

	return dir
}

func runDeploy(t *testing.T, projectDir, nodeUrl string) deployResult {
	t.Helper()

	bin := buildCLI(t)

	cmd := exec.CommandContext(t.Context(), bin, "deploy", "--environment", "prod")

	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(),
		"CI=true",
		"VERV_NODE_URL="+nodeUrl,
		"VERV_NODE_SECRET="+deploySecret,
	)

	out, err := cmd.CombinedOutput()

	result := deployResult{output: string(out)}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.exitCode = exitErr.ExitCode()

		return result
	}

	require.NoError(t, err)

	return result
}

func requireDeployFailed(t *testing.T, result deployResult) {
	t.Helper()

	require.NotZero(t, result.exitCode, "output:\n%s", result.output)
}

func requireOutputMentions(t *testing.T, result deployResult, fragment string) {
	t.Helper()

	require.Contains(t, strings.ToLower(result.output), strings.ToLower(fragment))
}

func requireSecretNotPrinted(t *testing.T, result deployResult) {
	t.Helper()

	require.NotContains(t, result.output, deploySecret)
}

func requireNoRequests(t *testing.T, standIn *velezStandIn) {
	t.Helper()

	require.Empty(t, standIn.recordedPaths())
}

func requireOnlyProbeRequests(t *testing.T, standIn *velezStandIn) {
	t.Helper()

	for _, path := range standIn.recordedPaths() {
		require.NotContains(t, path, registryCredsInfix)
		require.NotEqual(t, deployCreatePath, path)
	}
}

func requireNoRegistryRequests(t *testing.T, standIn *velezStandIn) {
	t.Helper()

	require.NotContains(t, standIn.recordedPaths(), registryListPath)

	requireOnlyProbeRequests(t, standIn)
}

func Test_DeployCi_FailsBeforeBuild(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("builds the real CLI; skipped in -short mode")
	}

	t.Run("node url missing", func(t *testing.T) {
		t.Parallel()

		standIn := newVelezStandIn(t, newNoRegistryBehaviour())
		projectDir := newDeployProject(t)

		result := runDeploy(t, projectDir, "")

		requireDeployFailed(t, result)
		requireOutputMentions(t, result, "VERV_NODE_URL")
		requireSecretNotPrinted(t, result)
		requireNoRequests(t, standIn)
	})

	t.Run("node rejects credentials", func(t *testing.T) {
		t.Parallel()

		standIn := newVelezStandIn(t, newForbiddenBehaviour())
		projectDir := newDeployProject(t)

		result := runDeploy(t, projectDir, standIn.server.URL)

		requireDeployFailed(t, result)
		requireSecretNotPrinted(t, result)
		requireNoRegistryRequests(t, standIn)
	})

	t.Run("node has no registry", func(t *testing.T) {
		t.Parallel()

		standIn := newVelezStandIn(t, newNoRegistryBehaviour())
		projectDir := newDeployProject(t)

		result := runDeploy(t, projectDir, standIn.server.URL)

		requireDeployFailed(t, result)
		requireOutputMentions(t, result, "registry")
		requireOnlyProbeRequests(t, standIn)
	})

	t.Run("node reports no arch", func(t *testing.T) {
		t.Parallel()

		standIn := newVelezStandIn(t, newNoArchBehaviour())
		projectDir := newDeployProject(t)

		result := runDeploy(t, projectDir, standIn.server.URL)

		requireDeployFailed(t, result)
		requireOutputMentions(t, result, "upgrade Velez")
		requireOnlyProbeRequests(t, standIn)
	})
}
