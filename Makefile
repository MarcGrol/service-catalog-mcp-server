
.PHONY: all generate test lint install clean

all: generate fmt lint test install

generate:
	go generate ./...

fmt:
	go fmt ./...

test:
	go test ./...

lint:
	golint ./...

install:
	go install

clean:
	go clean
	rm -f service-catalog-mcp-server
