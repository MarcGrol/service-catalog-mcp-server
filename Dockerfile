# Build the application from source
FROM acr-main.is.adyen.com/containers/golang-base AS build-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./
COPY data/service-catalog.sqlite data
COPY data/slos.sqlite data

# Run the tests in the container
RUN go test -v ./...

# Build the executable
RUN go build .


# Deploy the application binary into a lean image
FROM acr-main.is.adyen.com/containers/golang-base AS runtime-stage

WORKDIR /

COPY --from=build-stage /app/service-catalog-mcp-server /service-catalog-mcp-server

EXPOSE 8080

#USER nonroot:nonroot

ENTRYPOINT ["/service-catalog-mcp-server"]