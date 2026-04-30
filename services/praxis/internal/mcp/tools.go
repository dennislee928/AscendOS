package mcp

import "context"

// Tool is the minimal interface expected by Praxis MCP tool handlers.
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input map[string]any) (map[string]any, error)
}

// Registry stores tool implementations for later MCP server wiring.
type Registry interface {
	Register(tool Tool)
	Tools() []Tool
}

// InMemoryRegistry is a lightweight placeholder implementation.
type InMemoryRegistry struct {
	tools []Tool
}

func (r *InMemoryRegistry) Register(tool Tool) {
	if tool == nil {
		return
	}
	r.tools = append(r.tools, tool)
}

func (r *InMemoryRegistry) Tools() []Tool {
	return append([]Tool(nil), r.tools...)
}
