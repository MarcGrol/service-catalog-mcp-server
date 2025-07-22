# Service Catalog MCP Server

## Introduction

This project implements a Service Catalog for the MCP (Multi-Cloud Platform) server. It provides a centralized repository for managing and discovering services, modules, interfaces, and their relationships within a complex microservices architecture. The goal is to improve visibility, facilitate understanding of dependencies, and streamline development and operations.

## Features

- **Service Discovery**: Easily find and understand available services.
- **Module Management**: Organize and track software modules.
- **Interface Cataloging**: Document and manage API interfaces.
- **Dependency Mapping**: Visualize relationships between services, modules, and interfaces.
- **Complexity Analysis**: Identify and analyze the complexity of interfaces and modules.
- **Team-based Views**: Filter services and modules by owning teams.
- **Search Functionality**: Efficiently search the catalog for specific entities.

## Installation

To get started with the Service Catalog MCP Server, follow these steps:

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-org/service-catalog-mcp-server.git
    cd service-catalog-mcp-server
    ```

2.  **Build the project**:
    ```bash
    go mod tidy
    go build -o mcp-server .
    ```

## Usage

Once built, you can run the server:

```bash
./mcp-server
```

The server will expose various endpoints for querying the service catalog. Refer to the `internal/servicecatalog/handlers` directory for available API endpoints and their functionalities.

## Project Structure

- `main.go`: Entry point of the application.
- `internal/app`: Application initialization and setup.
- `internal/config`: Configuration management.
- `internal/mystore`: Data storage and persistence layer.
- `internal/servicecatalog`: Core service catalog logic, including handlers, repository, and search.
- `internal/transport`: Handles communication protocols (e.g., HTTP, SSE).

## Contributing

Contributions are welcome! Please see the `CONTRIBUTING.md` for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.