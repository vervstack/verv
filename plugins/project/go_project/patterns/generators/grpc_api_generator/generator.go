package grpc_api_generator

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project"
)

type serviceProtoApiArgs struct {
	PackageName    string
	GoPackageName  string
	NpmPackageName string
}

func GenerateServiceApiProto(project project.IProject) (*folder.Folder, error) {
	args := serviceProtoApiArgs{
		PackageName: project.GetShortName(),
	}

	protoFile := &rw.RW{}

	err := basicApiProtoTemplate.Execute(protoFile, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating service api proto")
	}

	return &folder.Folder{
		Name:    project.GetShortName() + ".proto",
		Content: protoFile.Bytes(),
	}, nil
}
