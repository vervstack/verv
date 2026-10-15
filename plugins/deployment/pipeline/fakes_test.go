package pipeline_test

import (
	"context"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
)

type callLog struct {
	calls []string
}

type fakeNode struct {
	log *callLog

	platform    nodeclient.Platform
	versionErr  error
	names       []string
	namesErr    error
	credentials nodeclient.RegistryCredentials
	deployErr   error

	credentialsName string
	deployRequests  []nodeclient.DeployRequest
}

func (f *fakeNode) Version(context.Context) (nodeclient.Platform, error) {
	f.log.calls = append(f.log.calls, "node.Version")
	return f.platform, f.versionErr
}

func (f *fakeNode) ListRegistryNames(context.Context) ([]string, error) {
	f.log.calls = append(f.log.calls, "node.ListRegistryNames")
	return f.names, f.namesErr
}

func (f *fakeNode) RegistryCredentials(_ context.Context, name string) (nodeclient.RegistryCredentials, error) {
	f.log.calls = append(f.log.calls, "node.RegistryCredentials")
	f.credentialsName = name
	return f.credentials, nil
}

func (f *fakeNode) CreateDeploy(_ context.Context, req nodeclient.DeployRequest) error {
	f.log.calls = append(f.log.calls, "node.CreateDeploy")
	f.deployRequests = append(f.deployRequests, req)
	return f.deployErr
}

type fakeProject struct {
	log *callLog

	serviceName string
	tag         string
}

func (f *fakeProject) ServiceName(string) (string, error) {
	f.log.calls = append(f.log.calls, "project.ServiceName")
	return f.serviceName, nil
}

func (f *fakeProject) Tag(context.Context, string) (string, error) {
	f.log.calls = append(f.log.calls, "project.Tag")
	return f.tag, nil
}

type fakeBuilder struct {
	log *callLog

	err      error
	requests []imagebuild.Request
}

func (f *fakeBuilder) Build(_ context.Context, req imagebuild.Request) error {
	f.log.calls = append(f.log.calls, "builder.Build")
	f.requests = append(f.requests, req)
	return f.err
}

type nopIO struct {
	lines []string
}

func (n *nopIO) Println(in ...string) {
	n.lines = append(n.lines, in...)
}

func (n *nopIO) Print(string) {}

func (n *nopIO) PrintlnColored(_ colors.Color, in ...string) {
	n.lines = append(n.lines, in...)
}

func (n *nopIO) PrintColored(colors.Color, string) {}

func (n *nopIO) Error(string) {}

func (n *nopIO) GetInput() (string, error) { return "", nil }

func (n *nopIO) GetInputOneOf([]string) string { return "" }
