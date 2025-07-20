package simple

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func SetupSimpleTools(s *server.MCPServer) {

	// Add a simple greeting tool
	greetTool := mcp.NewTool(
		"greet",
		mcp.WithDescription("Greet someone with a personalized message"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
		mcp.WithString("language",
			mcp.Description("Language for greeting (en, es, fr)"),
			mcp.DefaultString("en"),
		),
	)
	s.AddTool(greetTool, greetHandler)

	// Add a file reader tool
	fileReaderTool := mcp.NewTool(
		"read_file",
		mcp.WithDescription("Read contents of a text file"),
		mcp.WithString("filepath",
			mcp.Required(),
			mcp.Description("Path to the file to read"),
		),
	)
	s.AddTool(fileReaderTool, fileReaderHandler)

	// Add a system info tool
	sysInfoTool := mcp.NewTool(
		"system_info",
		mcp.WithDescription("Get basic system information"),
	)
	s.AddTool(sysInfoTool, systemInfoHandler)

	// Add a calculator tool with multiple parameters
	calcTool := mcp.NewTool(
		"calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("Operation to perform: add, subtract, multiply, divide"),
		),
	)
	s.AddTool(calcTool, calculatorHandler)
}

// greetHandler handles greeting requests
func greetHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'name': %v", err)), nil
	}

	language := request.GetString("language", "en")

	var greeting string
	switch language {
	case "es":
		greeting = fmt.Sprintf("¡Hola, %s! ¿Cómo estás?", name)
	case "fr":
		greeting = fmt.Sprintf("Bonjour, %s! Comment allez-vous?", name)
	default:
		greeting = fmt.Sprintf("Hello, %s! How are you doing today?", name)
	}

	return mcp.NewToolResultText(greeting), nil
}

// fileReaderHandler reads file contents
func fileReaderHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filepath, err := request.RequireString("filepath")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'filepath': %v", err)), nil
	}

	// Basic security check - prevent directory traversal
	if filepath == "" || filepath[0] == '/' {
		return mcp.NewToolResultError("Invalid file path: absolute paths not allowed"), nil
	}

	// For demo purposes, we'll simulate reading a file
	// In a real implementation, you'd use os.ReadFile(filepath)
	content := fmt.Sprintf("Simulated content of file: %s\n\nThis is where the actual file contents would appear.\nTimestamp: %s",
		filepath, time.Now().Format(time.RFC3339))

	return mcp.NewToolResultText(content), nil
}

// systemInfoHandler provides system information
func systemInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := fmt.Sprintf(`System Information:
- Server: File Operations MCP Server v1.0.0
- Protocol: Model Context Protocol
- Transport: stdio
- Time: %s
- Status: Running
- Go Version: Go 1.21+
- Features: Tools, Resources, Prompts`, time.Now().Format(time.RFC3339))

	return mcp.NewToolResultText(info), nil
}

// calculatorHandler performs arithmetic operations
func calculatorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a, err := request.RequireFloat("a")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameter 'a': %v", err)), nil
	}

	b, err := request.RequireFloat("b")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid parameter 'b': %v", err)), nil
	}

	operation, err := request.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'operation': %v", err)), nil
	}

	var result float64
	var resultText string

	switch operation {
	case "add":
		result = a + b
		resultText = fmt.Sprintf("%.2f + %.2f = %.2f", a, b, result)
	case "subtract":
		result = a - b
		resultText = fmt.Sprintf("%.2f - %.2f = %.2f", a, b, result)
	case "multiply":
		result = a * b
		resultText = fmt.Sprintf("%.2f × %.2f = %.2f", a, b, result)
	case "divide":
		if b == 0 {
			return mcp.NewToolResultError("Cannot divide by zero"), nil
		}
		result = a / b
		resultText = fmt.Sprintf("%.2f ÷ %.2f = %.2f", a, b, result)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("Unknown operation: %s. Supported operations: add, subtract, multiply, divide", operation)), nil
	}

	return mcp.NewToolResultText(resultText), nil
}
