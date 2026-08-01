include verv.mk

# generates folders and installs dependencies
warmup:
	make .prepare-grpc-folders
	make .deps-grpc

# generates code on warm project
codegen:
	moti generate

lint:
	golangci-lint run ./...