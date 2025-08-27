
.PHONY: all generate test lint install clean

all: tidy generate fmt lint test install

tidy:
	go mod tidy

generate:
	go generate ./...

fmt:
	find . -name "*.go" -exec goimports -l -w -local github.com/MarcGrol/service-catalog-mcp-server {} \;

lint:
	golint ./...

test:
	go test ./...

install:
	go install

docker:
	docker build \
	    -f Dockerfile \
	    -t service-catalog-mcp-server:local .

clean:
	go clean
	rm -f service-catalog-mcp-server
