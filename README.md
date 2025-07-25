# Service Catalog MCP Server

## Introduction

The mcp-server is an AI-integrated service catalog plugin that exposes deep structural insights into our large codebase. 
It allows AI agents like Claude Desktop, Gemini CLI, or other LLM-powered tools to explore the system architecture in a programmatic, structured way.

All information is harvested from the source code of our production platform — making it a powerful assistant for refactoring, impact analysis, ownership mapping, and integration discovery.

## Features

- **Service Discovery**: Easily find and understand available services.
- **Module Management**: Organize and track software modules.
- **Interface Cataloging**: Document and manage API interfaces.
- **Dependency Mapping**: Visualize relationships between services, modules, and interfaces.
- **Complexity Analysis**: Identify and analyze the complexity of interfaces and modules.
- **Team-based Views**: Filter services and modules by owning teams.
- **Search Functionality**: Efficiently search the catalog for specific entities.
- **SLO Discovery**: Efficiently search the SLOs for all applications.

## Installation

To get started with the Service Catalog MCP Server, follow these steps:

1.  **Clone the repository**:
    ```bash
    cd
    mkdir -p src
    cd src
    git clone https://github.com/your-org/service-catalog-mcp-server.git
    cd service-catalog-mcp-server

    # Make sure the sqlite databases (distributed separately) are in place:
    ./data/service-catalog.sqlite
    ./data/slos.sqlite
    ```

2.  **Build the project**:
    ```bash
    # tests, builds and installs the "service-catalog-mcp-server"-executable in ~/go/bin/
    
    make
    ```

## Usage

Once built, you can run the server:

```bash
# show help
~/go/bin/service-catalog-mcp-server -h

# start with default settings: stdio
~/go/bin/service-catalog-mcp-server

```


### Quick Verification

By default, the server runs in `stdio` mode. You can test it by pasting JSON-RPC 2.0 requests directly into the terminal. For example, to search for candidates related to "partner", paste the following JSON and press Enter:

```json
{"method":"tools/call","params":{"name":"suggest_candidates","arguments":{"keyword":"partner"}},"jsonrpc":"2.0","id":9}
```
The server will print the JSON-RPC response to standard output. For more examples, see `examples.md`.

## Integration with Claude-desktop

To integrate the `service-catalog-mcp-server` with Claude-desktop using `stdio` transport, follow these steps:

1.  **Build the project**:
    Ensure you have built the `service-catalog-mcp-server` executable as described in the "Installation" section.

2.  **Configure `claude_desktop_config.json`**:
    Locate your `claude_desktop_config.json` file (its location varies by operating system, but it's typically in your user's configuration directory for Claude-desktop). Add the following entry:

    ```json
    {
      "mcpServers": {
        "service-catalog": {
          "command": "/path/to/your/mcp-server",
          "args": [],
          "env": {}
        }
      }
    }
    ```
    **Important**: Replace `/path/to/your/service-catalog-mcp-server` with the actual absolute path to your `service-catalog-mcp-server` executable.

3.  **Restart Claude-desktop**:
    After saving the `claude_desktop_config.json` file, restart Claude-desktop for the changes to take effect. The `service-catalog` MCP server should now be available for use.

## Project Structure

- `main.go`: Entry point of the application.
- `internal/core/`: Framework to configure and start mcp-services using differnt transports
- `internal/plugin/servicecatalog`: Core service catalog logic, including handlers, repository, and search.

## Contributing

Contributions are welcome! Please see the `CONTRIBUTING.md` for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
