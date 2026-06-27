package tools

import (
	"context"
	"fmt"
)

// We assume the following from parity.go (we don't redeclare):
// 
// 
// func ok(text string) (ToolResponse, error) { ... }
// func err(text string) (ToolResponse, error) { ... }
// func getString(args map[string]interface{}, key string) string { ... }
// func getInt(args map[string]interface{}, key string) int { ... }
// func getBool(args map[string]interface{}, key string) bool { ... }

func HandleAddMemory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	content, _ :=getString(args, "content")
	if content == "" {
		return err("content is required")
}

	// In a real implementation, we would call an API to add memory.
	// For now, we simulate success.
	return ok(fmt.Sprintf("Added memory: %s", content))
}

func HandleSearchMemory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query is required")
}

	limit, _ :=getInt(args, "limit")
	if limit <= 0 {
		limit = 10 // default
	}
	// Simulate search: in real implementation, we would call an API.
	// We return a mock result.
	return ok(fmt.Sprintf("Found %d memories for query: %s", limit, query))
}