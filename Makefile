
.PHONY: all generate test lint install clean

all: generate lint test install

generate:
	go generate ./...

test:
	go test ./...

lint:
	golint ./...

install:
	go install

clean:
	go clean
	rm -f service-catalog-mcp-server
