
.PHONY: all generate test lint install clean dockerbuild dockerrun dockerview dockertest

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

dockerbuild:
	docker  --log-level debug build \
	    --no-cache \
	    -t service-catalog-mcp-server:local \
	    -f Dockerfile \
	    .

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
