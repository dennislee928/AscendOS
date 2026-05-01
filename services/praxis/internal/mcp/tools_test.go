package mcp

import (
	"context"
	"testing"
)

type stubTool struct {
	name string
}

func (t stubTool) Name() string { return t.name }

func (t stubTool) Description() string { return "stub" }

func (t stubTool) Execute(context.Context, map[string]any) (map[string]any, error) {
	return map[string]any{"ok": true}, nil
}

func TestInMemoryRegistryIgnoresNilAndCopiesTools(t *testing.T) {
	var registry InMemoryRegistry
	registry.Register(nil)
	registry.Register(stubTool{name: "alpha"})

	tools := registry.Tools()
	if len(tools) != 1 {
		t.Fatalf("len(Tools()) = %d, want %d", len(tools), 1)
	}
	if got := tools[0].Name(); got != "alpha" {
		t.Fatalf("Tools()[0].Name() = %q, want %q", got, "alpha")
	}

	tools[0] = stubTool{name: "mutated"}
	if got := registry.Tools()[0].Name(); got != "alpha" {
		t.Fatalf("registry.Tools()[0].Name() = %q, want %q", got, "alpha")
	}
}
