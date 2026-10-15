package imagebuild_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.redsock.ru/rerrors"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
)

var errFakeBuild = rerrors.New("fake build failure")

type silentPrinter struct {
	errors []string
}

func (p *silentPrinter) Println(...string)                      {}
func (p *silentPrinter) Print(string)                           {}
func (p *silentPrinter) PrintlnColored(colors.Color, ...string) {}
func (p *silentPrinter) PrintColored(colors.Color, string)      {}
func (p *silentPrinter) Error(in string)                        { p.errors = append(p.errors, in) }
func (p *silentPrinter) GetInput() (string, error)              { return "", nil }
func (p *silentPrinter) GetInputOneOf(options []string) string  { return "" }

type fakeExecutor struct {
	calls            []imagebuild.Command
	configDir        string
	configJson       []byte
	output           string
	err              error
	isConfigDirThere bool
}

func (f *fakeExecutor) Run(_ context.Context, command imagebuild.Command) (string, error) {
	f.calls = append(f.calls, command)

	for _, env := range command.Env {
		value, ok := strings.CutPrefix(env, "DOCKER_CONFIG=")
		if ok {
			f.configDir = value
		}
	}

	raw, err := os.ReadFile(filepath.Join(f.configDir, "config.json"))
	if err == nil {
		f.configJson = raw
		f.isConfigDirThere = true
	}

	return f.output, f.err
}

func newRealDockerConfig(t *testing.T) string {
	t.Helper()

	realDir := t.TempDir()
	t.Setenv("DOCKER_CONFIG", realDir)

	pluginsDir := filepath.Join(realDir, "cli-plugins")
	err := os.MkdirAll(pluginsDir, 0o700)
	require.NoError(t, err)

	config := `{"auths":{"old.example.com":{"auth":"b2xk"}},"credsStore":"desktop",` +
		`"credHelpers":{"x":"y"},"currentContext":"colima"}`
	err = os.WriteFile(filepath.Join(realDir, "config.json"), []byte(config), 0o644)
	require.NoError(t, err)

	return realDir
}

func newBuildRequest() imagebuild.Request {
	return imagebuild.Request{
		ProjectDir: "/project",
		Platform:   "linux/arm64",
		ImageRef:   "registry.example.com/artel:abc",
		Credentials: imagebuild.Credentials{
			Username:    "node",
			Password:    "secret",
			RegistryUrl: "https://registry.example.com/",
		},
	}
}

func Test_Build_RunsBuildxOnceWithTemporaryDockerConfig(t *testing.T) {
	realDir := newRealDockerConfig(t)
	executor := &fakeExecutor{}
	printer := &silentPrinter{}

	err := imagebuild.Build(context.Background(), printer, executor, newBuildRequest())

	require.NoError(t, err)
	require.Len(t, executor.calls, 1)

	command := executor.calls[0]
	require.Equal(t, "docker", command.Tool)
	require.Equal(t, []string{
		"buildx", "build",
		"--platform", "linux/arm64",
		"-f", "Dockerfile",
		"-t", "registry.example.com/artel:abc",
		"--push",
		".",
	}, command.Args)
	require.Equal(t, "/project", command.WorkDir)
	require.Equal(t, []string{"DOCKER_CONFIG=" + executor.configDir}, command.Env)
	require.NotEqual(t, realDir, executor.configDir)

	require.True(t, executor.isConfigDirThere)
	_, err = os.Stat(executor.configDir)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func Test_Build_TemporaryDockerConfigContent(t *testing.T) {
	realDir := newRealDockerConfig(t)
	executor := &fakeExecutor{}

	err := imagebuild.Build(context.Background(), &silentPrinter{}, executor, newBuildRequest())
	require.NoError(t, err)

	var fields map[string]json.RawMessage
	err = json.Unmarshal(executor.configJson, &fields)
	require.NoError(t, err)

	require.NotContains(t, fields, "credsStore")
	require.NotContains(t, fields, "credHelpers")
	require.JSONEq(t, `"colima"`, string(fields["currentContext"]))

	var auths map[string]map[string]string
	err = json.Unmarshal(fields["auths"], &auths)
	require.NoError(t, err)

	wantAuth := base64.StdEncoding.EncodeToString([]byte("node:secret"))
	require.Equal(t, map[string]map[string]string{"registry.example.com": {"auth": wantAuth}}, auths)

	original, err := os.ReadFile(filepath.Join(realDir, "config.json"))
	require.NoError(t, err)
	require.Contains(t, string(original), "credsStore")
}

func Test_Build_TemporaryDockerConfigLayout(t *testing.T) {
	realDir := newRealDockerConfig(t)
	inspector := &layoutInspector{}

	err := imagebuild.Build(context.Background(), &silentPrinter{}, inspector, newBuildRequest())
	require.NoError(t, err)

	require.Equal(t, os.FileMode(0o600), inspector.configMode)
	require.Equal(t, os.FileMode(0o700), inspector.dirMode)
	require.Equal(t, filepath.Join(realDir, "cli-plugins"), inspector.pluginsLinkTarget)
}

type layoutInspector struct {
	configMode        os.FileMode
	dirMode           os.FileMode
	pluginsLinkTarget string
}

func (l *layoutInspector) Run(_ context.Context, command imagebuild.Command) (string, error) {
	dir := strings.TrimPrefix(command.Env[0], "DOCKER_CONFIG=")

	configInfo, err := os.Stat(filepath.Join(dir, "config.json"))
	if err != nil {
		return "", rerrors.Wrap(err)
	}
	l.configMode = configInfo.Mode().Perm()

	dirInfo, err := os.Stat(dir)
	if err != nil {
		return "", rerrors.Wrap(err)
	}
	l.dirMode = dirInfo.Mode().Perm()

	l.pluginsLinkTarget, err = os.Readlink(filepath.Join(dir, "cli-plugins"))
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	return "", nil
}

func Test_Build_FailureReportsOnlyOutputTail(t *testing.T) {
	newRealDockerConfig(t)

	var lines []string
	for i := 1; i <= 30; i++ {
		lines = append(lines, fmt.Sprintf("line-%02d", i))
	}
	executor := &fakeExecutor{output: strings.Join(lines, "\n") + "\n", err: errFakeBuild}
	printer := &silentPrinter{}

	err := imagebuild.Build(context.Background(), printer, executor, newBuildRequest())

	require.Error(t, err)
	require.NotContains(t, err.Error(), lines[len(lines)-1])
	require.Len(t, printer.errors, 1)
	require.Equal(t, 20, strings.Count(strings.TrimRight(printer.errors[0], "\n"), "\n")+1)
	require.Contains(t, printer.errors[0], lines[len(lines)-1])
	require.NotContains(t, printer.errors[0], lines[0]+"\n")
}
