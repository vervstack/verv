lint:
	golangci-lint run --fix

test:
	GOPROXY=https://proxy.golang.org,direct go test ./...

gen:
	go generate ./