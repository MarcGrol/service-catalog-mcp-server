
.PHONY: all generate test lint install clean

all: tidy generate fmt lint test install

tidy:
	go mod tidy

generate:
	go generate ./...

fmt:
	go fmt ./...

lint:
	golint ./...

test:
	go test ./...

install:
	go install

clean:
	go clean
	rm -f service-catalog-mcp-server
