dev-build:
	go build -o $$(go env GOPATH)/bin/verv-dev .

lint:
	golangci-lint run --fix

test:
	GOPROXY=https://proxy.golang.org,direct go test ./...

gen:
	go generate ./