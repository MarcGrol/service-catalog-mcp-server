package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/mystore"
)

func main() {
	ctx := context.Background()

	var (
		useSSE        = flag.Bool("sse", false, "Use SSE transport instead of stdio")
		useStreamable = flag.Bool("http", false, "Use Streamable HTTP transport (easier for testing)")
		port          = flag.String("port", "8080", "Port for SSE server")
		baseURL       = flag.String("baseurl", "http://localhost", "Base URL for SSE server")
	)
	flag.Parse()

	// Create a new MCP server
	s := server.NewMCPServer(
		"Marc's project MCP Server", // Server name
		"1.0.0",                     // Version
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging())
	//setupSimpleTools(s)

	projectStore, projectStoreCleanup, err := mystore.New[ProjectConfig](ctx)
	if err != nil {
		log.Fatalf("Error creating basket store: %s", err)
	}
	defer projectStoreCleanup()

	projectService := New(s, projectStore)
	projectService.initialize(ctx)

	if *useStreamable {
		// Option 2: Streamable HTTP Transport
		streamableServer := server.NewStreamableHTTPServer(s,
			// Streamable HTTP Transport (stateless - no session management needed)
			server.WithStateLess(true), // This is the key! Disables session management
		)
		log.Printf("Starting MCP server with Streamable HTTP transport on :%s", *port)
		log.Printf("HTTP endpoint: http://localhost:%s/mcp (direct JSON-RPC calls)", *port)
		log.Println("Test with: curl -X POST http://localhost:" + *port + "/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'")

		if err := streamableServer.Start(":" + *port); err != nil {
			log.Fatalf("Streamable HTTP server error: %v", err)
		}
	} else if *useSSE {
		// SSE Transport (HTTP-based)
		fullBaseURL := fmt.Sprintf("%s:%s", *baseURL, *port)

		// Create SSE server
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(fullBaseURL),
		)

		log.Printf("Starting MCP server with SSE transport on %s", fullBaseURL)
		log.Printf("SSE endpoint: %s/sse", fullBaseURL)
		log.Printf("Message endpoint: %s/message", fullBaseURL)

		// Start SSE server - this blocks
		if err := sseServer.Start(":" + *port); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}
	} else {
		// Default stdio transport
		log.Println("Starting MCP server with stdio transport...")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Stdio server error: %v", err)
		}
	}
}
