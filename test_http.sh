#!/bin/bash

# Test script for STATELESS Streamable HTTP MCP Server
# This should work without session management!

BASE_URL="http://localhost:8080/mcp"

echo "=== Testing STATELESS MCP Server ==="
echo "Server should be running with: ./mcp-server -http"
echo ""

# Test 1: Direct tool list (should work without initialization!)
echo "=== Test 1: Direct Tools List (No Init Required) ==="
curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list"
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}'

echo -e "\n"

# Test 2: Direct tool call (should work without initialization!)
echo "=== Test 2: Direct Tool Call - Greet ==="
curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "greet",
      "arguments": {
        "name": "Stateless User",
        "language": "es"
      }
    }
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "greet", "arguments": {"name": "Stateless User", "language": "es"}}}'

echo -e "\n"

# Test 3: Calculator tool
echo "=== Test 3: Calculator Tool ==="
curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "calculate",
      "arguments": {
        "a": 15,
        "b": 7,
        "operation": "multiply"
      }
    }
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "calculate", "arguments": {"a": 15, "b": 7, "operation": "multiply"}}}'

echo -e "\n"

# Test 4: Server info
echo "=== Test 4: Server Info ==="
curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "server_info",
      "arguments": {}
    }
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "server_info", "arguments": {}}}'

echo -e "\n"

# Test 5: Optional - Test initialization (should still work but not required)
echo "=== Test 5: Optional Initialization (Not Required in Stateless Mode) ==="
curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": 5,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {
        "roots": {"listChanged": false},
        "sampling": {}
      },
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc": "2.0", "id": 5, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": false}, "sampling": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}'

echo -e "\n=== Tests Complete ==="
echo "If you see JSON responses above (not errors), the stateless server is working!"

# One-liner test commands for easy copy-paste:
echo ""
echo "=== Quick One-Liner Test Commands ==="
echo ""
echo "# List tools:"
echo "curl -X POST http://localhost:3000/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'"
echo ""
echo "# Call greet tool:"
echo "curl -X POST http://localhost:3000/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"greet\",\"arguments\":{\"name\":\"Test User\",\"language\":\"fr\"}}}'"
echo ""
echo "# Call calculator:"
echo "curl -X POST http://localhost:3000/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":3,\"method\":\"tools/call\",\"params\":{\"name\":\"calculate\",\"arguments\":{\"a\":20,\"b\":4,\"operation\":\"divide\"}}}'"
