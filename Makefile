
.PHONY: all install generate test lint install clean dockerbuild dockerrun dockerview dockertest

all: generate fmt lint test build tidy

install:
	go install go.uber.org/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/lint/golint@latest

generate:
	go generate ./...

fmt: generate
	find . -name "*.go" -exec goimports -l -w -local github.com/MarcGrol/service-catalog-mcp-server {} \;

lint: fmt
	golint ./...

test: lint
	go test ./...

build: test
	go install ./...

tidy:
	go mod tidy

dockerbuild:
	docker build \
		--log-level debug  \
	    --no-cache \
	    -t acr-main.is.adyen.com/is/service-catalog-mcp-server:0.1 \
	    -f Dockerfile \
	    .
dockerpush:
	docker tag service-catalog-mcp-server:0.1 acr-main.is.adyen.com/is/service-catalog-mcp-server:0.1
	docker login acr-main.is.adyen.com/is
	docker push acr-main.is.adyen.com/is/service-catalog-mcp-server:0.1

dockerrun:
	docker run \
	-p 8000:8000 \
	--rm  docker.io/library/service-catalog-mcp-server:local

dockerview:
	docker container ls

dockertest:
	curl -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"method":"tools/call","params":{"name":"suggest_candidates","arguments":{"keyword":"partner"}},"jsonrpc":"2.0","id":9}' http://localhost:8000/tools/call/

clean:
	go clean
	rm -f service-catalog-mcp-server
