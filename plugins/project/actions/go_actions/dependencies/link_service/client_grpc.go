package link_service

import (
	stderrs "errors"
	"path"
	"strings"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"go.vervstack.ru/verv/internal/cmd"
	rscliconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies/link_service/grpc_discovery"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/clients_generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/config_generators"
)

type GrpcClient struct {
	Modules []string

	Cfg *rscliconfig.VervConfig
	Io  io.IO
}

func (g GrpcClient) GetFolderName() string {
	return "grpc"
}

func (g GrpcClient) AppendToProject(proj project.IProject) error {
	if len(g.Modules) == 0 {
		return nil
	}

	succeeded, failed := g.getPackages(g.Modules)
	if len(failed) != 0 {
		g.Io.Error("some packages weren't found: " + strings.Join(failed, ","))
	}

	var errs error
	for _, item := range succeeded {
		err := g.applyLink(proj, item)
		if err != nil {
			errs = stderrs.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}

	grpcClientConnFilePath := path.Join(g.Cfg.Env.PathsToClients[0], patterns.GRPCServer, patterns.ConnFileName)
	if proj.GetFolder().GetByPath(grpcClientConnFilePath) == nil {
		content, err := clients_generators.GenerateGRPCConn()
		if err != nil {
			return rerrors.Wrap(err, "error generating grpc conn file")
		}

		proj.GetFolder().Add(&folder.Folder{
			Name:    grpcClientConnFilePath,
			Content: content,
		})
	}

	return nil
}

func (g GrpcClient) getPackages(args []string) (succeeded, failed []string) {
	for _, packageName := range args {
		if g.getPackage(packageName) {
			succeeded = append(succeeded, packageName)
		} else {
			failed = append(failed, packageName)
		}
	}

	return succeeded, failed
}

func (g GrpcClient) getPackage(packageName string) (ok bool) {
	if !strings.Contains(packageName, "@") {
		packageName += "@latest"
	}

	_, err := cmd.Execute(cmd.Request{
		Tool: "go",
		Args: []string{"get", packageName},
	})

	return err == nil
}

func (g GrpcClient) applyLink(proj project.IProject, packageName string) error {
	discovery := grpc_discovery.GrpcDiscovery{Cfg: g.Cfg}

	pkg, err := discovery.DiscoverPackage(packageName)
	if err != nil {
		return rerrors.Wrap(err, "error discovering package")
	}

	if pkg == nil {
		return nil
	}

	if idx := strings.Index(packageName, "@"); idx > -1 {
		packageName = packageName[:idx]
	}
	grpcPkgPath := path.Join(g.Cfg.Env.PathsToClients[0], patterns.GRPCServer)

	grpcClientsFolder := proj.GetFolder().GetByPath(grpcPkgPath)
	if grpcClientsFolder == nil {
		grpcClientsFolder = &folder.Folder{
			Name: grpcPkgPath,
		}

		proj.GetFolder().Add(grpcClientsFolder)
	}

	resourceName := resources.GrpcResourceName + "_" + generators.NormalizeResourceName(path.Base(packageName))

	_, err = proj.GetConfig().GRPC(resourceName)
	if err != nil {
		if !rerrors.Is(err, matreshka.ErrNotFound) {
			return rerrors.Wrap(err, "error getting grpc resource from config")
		}

		grpcResource := &resources.GRPC{
			Name:             resources.Name(resourceName),
			Module:           packageName,
			ConnectionString: "0.0.0.0:50051",
		}
		proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, grpcResource)
	}

	grpcClientFile, err := config_generators.GenerateGRPCClient(*pkg)
	if err != nil {
		return rerrors.Wrap(err, "error generating grpc client")
	}

	grpcClientsFolder.Add(
		&folder.Folder{
			Name:    path.Base(packageName) + ".go",
			Content: grpcClientFile,
		})

	return nil
}
