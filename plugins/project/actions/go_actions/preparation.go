package go_actions

import (
	stderrs "errors"
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies/link_service"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/dockerfile_generator"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/main_generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/middleware_generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/server_generators/impl_gen"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/transport_generators"
)

type PrepareProjectStructure struct {
}

func (a PrepareProjectStructure) Do(p project.IProject) error {
	rootF := p.GetFolder()

	mainFileContent, err := main_generators.GenerateMain()
	if err != nil {
		return rerrors.Wrap(err, "error generating main.go")
	}

	cmd := &folder.Folder{Name: patterns.CmdFolder}
	cmd.Add(&folder.Folder{
		Name:    path.Join(patterns.ServiceFolder, patterns.MainFileName),
		Content: mainFileContent,
	})
	rootF.Add(cmd)

	configFolder := &folder.Folder{Name: patterns.ConfigsFolder}
	rootF.Add(configFolder)
	rootF.Add(&folder.Folder{Name: patterns.InternalFolder})

	rootF.Add(
		patterns.Readme.Copy(),
		patterns.GitIgnore.Copy(),
	)

	if rootF.GetByPath(patterns.Linter.Name) == nil {
		rootF.Add(patterns.Linter.Copy())
	}

	return nil
}
func (a PrepareProjectStructure) NameInAction() string {
	return "Preparing project structure"
}

type PrepareClients struct {
	C  *vervconfig.VervConfig
	IO io.IO
}

func (a PrepareClients) Do(p project.IProject) error {
	if a.C == nil {
		a.C = vervconfig.GetConfig()
	}

	if a.IO == nil {
		a.IO = io.StdIO{}
	}

	var (
		simpleClients []string
		grpcClients   []string
	)

	cfg := p.GetConfig()

	for _, r := range cfg.DataSources {
		grpcC, ok := r.(*resources.GRPC)
		if ok {
			grpcClients = append(grpcClients, grpcC.Module)
		} else {
			simpleClients = append(simpleClients, r.GetName())
		}
	}

	var errs []error

	deps := dependencies.GetDependencies(a.C, simpleClients)
	if len(deps) != 0 {
		for _, item := range deps {
			err := item.AppendToProject(p)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	grpcClient := link_service.GrpcClient{
		Modules: grpcClients,
		Cfg:     a.C,
		Io:      a.IO,
	}

	err := grpcClient.AppendToProject(p)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return stderrs.Join(errs...)
	}

	return nil
}

func (a PrepareClients) NameInAction() string {
	return "Generating clients"
}

type PrepareServer struct{}

func (a PrepareServer) Do(p project.IProject) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	rootF := p.GetFolder()

	if rootF.GetByPath(patterns.Moti.Name) == nil {
		rootF.Add(patterns.Moti.Copy())
	}

	internalFolder := rootF.GetByPath(patterns.InternalFolder)
	if internalFolder == nil {
		internalFolder = &folder.Folder{Name: patterns.InternalFolder}
		rootF.Add(internalFolder)
	}

	transportFolder := internalFolder.GetByPath(patterns.TransportFolder)
	if transportFolder == nil {
		transportFolder = &folder.Folder{Name: patterns.TransportFolder}
		internalFolder.Add(transportFolder)
	}

	middlewareFolder := internalFolder.GetByPath(patterns.MiddlewareFolder)
	if middlewareFolder == nil {
		middlewareFolder = &folder.Folder{Name: patterns.MiddlewareFolder}
		internalFolder.Add(middlewareFolder)
	}

	err := generateTransportFiles(transportFolder, p.GetName())
	if err != nil {
		return err
	}

	generateMiddlewareFiles(middlewareFolder)

	implFolders, err := impl_gen.GenerateImpl(vervconfig.GetConfig(), p)
	if err != nil {
		return rerrors.Wrap(err, "error during stub generation")
	}

	addMissingImplFolders(transportFolder, implFolders)

	return nil
}

// generateTransportFiles scaffolds internal/transport's fixed files once. None of them carry a
// "Code generated" header - like custom.go (see app_struct_generators.GenerateAppFiles), they're
// a starting point meant to be hand-edited after the first scaffold, so each is only generated
// when it doesn't already exist: overwriting on every tidy silently discarded prior edits (see
// the bug report that motivated this guard).
func generateTransportFiles(transportFolder *folder.Folder, fullProjPath string) error {
	if transportFolder.GetByPath(patterns.ServerManagerFileName) == nil {
		serverManagerContent, err := transport_generators.GenerateServerManager()
		if err != nil {
			return rerrors.Wrap(err, "error generating server manager")
		}

		transportFolder.Add(&folder.Folder{Name: patterns.ServerManagerFileName, Content: serverManagerContent})
	}

	if transportFolder.GetByPath(patterns.GrpcServerFileName) == nil {
		grpcServerContent := transport_generators.GenerateGrpcServer()

		transportFolder.Add(&folder.Folder{Name: patterns.GrpcServerFileName, Content: grpcServerContent})
	}

	if transportFolder.GetByPath(patterns.HttpServerFileName) == nil {
		httpServerContent, err := transport_generators.GenerateHttpServer()
		if err != nil {
			return rerrors.Wrap(err, "error generating http server")
		}

		transportFolder.Add(&folder.Folder{Name: patterns.HttpServerFileName, Content: httpServerContent})
	}

	if transportFolder.GetByPath(patterns.GatewayMuxFileName) == nil {
		gatewayMuxContent, err := transport_generators.GenerateGatewayMux(fullProjPath)
		if err != nil {
			return rerrors.Wrap(err, "error generating gateway mux")
		}

		transportFolder.Add(&folder.Folder{Name: patterns.GatewayMuxFileName, Content: gatewayMuxContent})
	}

	return nil
}

// generateMiddlewareFiles scaffolds internal/middleware's fixed files once - same "generate only
// if missing" rationale as generateTransportFiles above.
func generateMiddlewareFiles(middlewareFolder *folder.Folder) {
	if middlewareFolder.GetByPath(patterns.CookieNamesFileName) == nil {
		cookieNamesContent := middleware_generators.GenerateCookieNames()
		middlewareFolder.Add(&folder.Folder{Name: patterns.CookieNamesFileName, Content: cookieNamesContent})
	}

	if middlewareFolder.GetByPath(patterns.CookieAnnotatorFileName) == nil {
		cookieAnnotatorContent := middleware_generators.GenerateCookieAnnotator()
		middlewareFolder.Add(&folder.Folder{Name: patterns.CookieAnnotatorFileName, Content: cookieAnnotatorContent})
	}

	if middlewareFolder.GetByPath(patterns.CookieResponseFileName) == nil {
		cookieResponseContent := middleware_generators.GenerateCookieResponse()
		middlewareFolder.Add(&folder.Folder{Name: patterns.CookieResponseFileName, Content: cookieResponseContent})
	}

	if middlewareFolder.GetByPath(patterns.CSRFInterceptorFileName) == nil {
		csrfInterceptorContent := middleware_generators.GenerateCSRFInterceptor()
		middlewareFolder.Add(&folder.Folder{Name: patterns.CSRFInterceptorFileName, Content: csrfInterceptorContent})
	}

	if middlewareFolder.GetByPath(patterns.RequestSchemeAnnotatorFileName) == nil {
		requestSchemeAnnotatorContent := middleware_generators.GenerateRequestSchemeAnnotator()
		middlewareFolder.Add(
			&folder.Folder{Name: patterns.RequestSchemeAnnotatorFileName, Content: requestSchemeAnnotatorContent},
		)
	}
}

func addMissingImplFolders(transportFolder *folder.Folder, implFolders []*folder.Folder) {
	for _, implF := range implFolders {
		exists := false

		for _, tF := range transportFolder.Inner {
			if tF.Name == implF.Name {
				exists = true

				break
			}
		}

		if !exists {
			transportFolder.Add(implF)
		}
	}
}

func (a PrepareServer) NameInAction() string {
	return "Preparing server files"
}

type PrepareDockerfile struct{}

func (a PrepareDockerfile) Do(p project.IProject) (err error) {
	dockerFile := p.GetFolder().GetByPath(patterns.DockerfileFile)
	if dockerFile != nil {
		return nil
	}

	dockerFile = &folder.Folder{
		Name: patterns.DockerfileFile,
	}

	dockerFile.Content, err = dockerfile_generator.GenerateDockerfile(p)
	if err != nil {
		return rerrors.Wrap(err)
	}

	p.GetFolder().Add(dockerFile)

	return nil
}
func (a PrepareDockerfile) NameInAction() string {
	return "Preparing Dockerfile"
}
