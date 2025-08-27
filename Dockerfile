# Build the application from source
FROM acr-main.is.adyen.com/containers/golang-base AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download ...

COPY *.go ./
# Build the executable
RUN go build

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...


# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/service-catalog-mcp-server /service-catalog-mcp-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/service-catalog-mcp-server"]