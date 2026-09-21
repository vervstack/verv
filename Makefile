dev-build:
	go build -o $$(go env GOPATH)/bin/verv-dev .

lint:
	golangci-lint run --fix

test:
	GOPROXY=https://proxy.golang.org,direct go test ./...

gen:
	go generate ./

deploy-proxy:
	scp local/*.ssl.conf germ:~/web/configs/
	ssh germ 'cd web && make reload'