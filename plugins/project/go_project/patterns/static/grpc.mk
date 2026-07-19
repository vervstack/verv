gen-server-grpc: .prepare-grpc-folders .deps-grpc .gen-server-grpc

.prepare-grpc-folders:
	mkdir -p pkg

.deps-grpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	moti install

.gen-server-grpc:
	moti generate
