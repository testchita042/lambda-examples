// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import (
	"net/http"
	"strings"
)

// Server represents an MCP protocol server
type Server struct {
	Name        string
	Version     string
	Description string
	Tools       []ToolDescription
}

// NewServer creates a new MCP server with the given parameters
func NewServer(name, version, description string) *Server {
	return &Server{
		Name:        name,
		Version:     version,
		Description: description,
		Tools:       []ToolDescription{},
	}
}

// RegisterTool adds a tool to the server's available tools
func (s *Server) RegisterTool(tool ToolDescription) {
	s.Tools = append(s.Tools, tool)
}

// SetCORSHeaders sets standard CORS headers to allow MCP Inspector to connect
func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// SetSSEHeaders sets standard Server-Sent Events headers
func SetSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

// HandleInitialize creates the initialize response data
func (s *Server) HandleInitialize() map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":        s.Name,
			"version":     s.Version,
			"description": s.Description,
		},
	}
}

// HandleTools creates the tools list response data
func (s *Server) HandleTools() map[string]interface{} {
	return map[string]interface{}{
		"tools": s.Tools,
	}
}

// GetMethodFromPath extracts the method name from the request path
func GetMethodFromPath(path string) string {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) > 0 {
		return pathParts[len(pathParts)-1]
	}
	return ""
}
