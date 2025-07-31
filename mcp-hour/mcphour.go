package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"mcp-hour/adapters"
	"mcp-hour/domain"
	"mcp-hour/lib/mcp"
)

// HourResponse represents the response structure
type HourResponse struct {
	Hour        int    `json:"hour"`
	AmPm        string `json:"amPm"`
	Message     string `json:"message"`
	CurrentTime string `json:"currentTime"`
}

// createServer initializes and configures an MCP server for our hour service
func createServer() *mcp.Server {
	server := mcp.NewServer("HourMCP", "1.0.0", "MCP server that provides current hour information")
	server.RegisterTool(createMCPToolDefinition())
	return server
}

// createMCPToolDefinition creates an MCP-compatible tool definition
func createMCPToolDefinition() mcp.ToolDescription {
	return mcp.ToolDescription{
		Name:        "get_hour",
		Description: "Get the current hour in 12-hour format with AM/PM indicator",
		InputSchema: mcp.Schema{
			Type: "object",
			Properties: map[string]mcp.ParameterProperty{
				"timezone": {
					Type:        "string",
					Description: "Optional timezone (defaults to system timezone)",
				},
			},
			Required: []string{},
		},
		OutputSchema: mcp.Schema{
			Type: "object",
			Properties: map[string]mcp.ParameterProperty{
				"hour": {
					Type:        "integer",
					Description: "Current hour in 12-hour format",
				},
				"amPm": {
					Type:        "string",
					Description: "AM or PM indicator",
				},
				"message": {
					Type:        "string",
					Description: "Message containing the current hour and AM/PM indicator",
				},
				"currentTime": {
					Type:        "string",
					Description: "Current time in ISO format",
				},
			},
			Required: []string{"hour", "amPm", "message", "currentTime"},
		},
	}
}

// getHourInfo returns the current hour information using domain services
func getHourInfo() HourResponse {
	// Create the domain service with a system clock adapter
	clockAdapter := adapters.NewSystemClock()
	hourService := domain.NewHourService(clockAdapter)

	// Get the hour information from domain service
	hour, amPm, currentTime := hourService.GetHourInfo()
	message := fmt.Sprintf("Current hour is %s", currentTime)

	// currentTime is already formatted as ISO8601 from the adapter

	// Return a response with the hour data
	return HourResponse{
		Hour:        hour,
		AmPm:        amPm,
		Message:     message,
		CurrentTime: currentTime,
	}
}

// Handler is the lambda entry point for Chita Cloud
func Handler(r *http.Request, w http.ResponseWriter, req mcp.MCPRequest) (io.ReadCloser, error) {
	fmt.Println("MCP request:", req.JSONRPC, req.ID, req.Method)
	fmt.Println("Request headers:", r.Header)

	// Set CORS headers to allow MCP Inspector to connect
	mcp.SetCORSHeaders(w)

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return nil, nil
	}

	// Set SSE headers
	mcp.SetSSEHeaders(w)

	// Prepare the response based on path
	var responseBody []byte
	var err error

	// Define default request ID for JSON-RPC 2.0
	requestID := req.ID

	// All paths are handled as MCP protocol endpoints
	// Create MCP protocol response based on path
	var responseData interface{}

	// Extract method from path
	pathMethod := mcp.GetMethodFromPath(r.URL.Path)
	fmt.Printf("Path method: %s\n", pathMethod)
	paramsB, err := json.Marshal(req.Params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Params: %s\n", string(paramsB))

	if id, ok := req.Params["requestId"]; ok {
		if idStr, ok := id.(string); ok {
			requestID, _ = strconv.Atoi(idStr)
		} else if idInt, ok := id.(int); ok {
			requestID = idInt
		}
	}

	// Use request method if available, otherwise use path-derived method
	method := req.Method
	if method == "" {
		method = pathMethod
		if method == "" {
			method = "response"
		}
	}

	// Create server with tools
	server := createServer()

	// Handle different MCP protocol paths
	switch method {
	case "initialize":
		// Initialize request - return server capabilities
		responseData = server.HandleInitialize()
		fmt.Println("Sending initialize response")

	case "tools/list":
		// List tools request
		responseData = server.HandleTools()
		fmt.Println("Sending tools list response")

	case "tools/call":
		switch req.Params["name"] {
		case "get_hour":
			hourInfo := getHourInfo()
			responseData = map[string]any{
				"content": []map[string]any{
					{
						"type": "text",
						"text": hourInfo.Message,
					},
				},
				"structuredContent": hourInfo,
			}
			fmt.Println("Sending get_hour response:", hourInfo)
		default:
			responseData = nil
			fmt.Println("Unknown tool name:", req.Params["name"])
		}
	default:
		// Default path - for compatibility with legacy clients
		responseData = getHourInfo()
		fmt.Println("Sending default path response")
	}

	// Format as JSON-RPC 2.0 response for MCP
	responseBody, err = mcp.FormatMCPServerResponse(requestID, method, responseData)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	// Format as SSE
	var buffer strings.Builder

	// Add data
	buffer.WriteString("data: ")
	buffer.Write(responseBody)
	buffer.WriteString("\n\n")

	return io.NopCloser(strings.NewReader(buffer.String())), nil
}
