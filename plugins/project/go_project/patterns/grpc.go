package patterns

import (
	_ "embed"

	"go.vervstack.ru/verv/internal/io/folder"
)

// Proto contract
var (
	//go:embed static/api/grpc/api.proto
	protoContract []byte
	ProtoContract = &folder.Folder{
		Name:    "api.proto",
		Content: protoContract,
	}
)

// Dependencies and generator
var (
	//go:embed static/moti.yaml
	moti []byte
	Moti = &folder.Folder{
		Name:    "moti.yaml",
		Content: moti,
	}

	//go:embed static/grpc.mk
	GrpcServerGenMK []byte
)
