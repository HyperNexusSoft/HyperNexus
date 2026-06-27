package tools

import (
	"context"
)

// Handler functions
func HandleTestRadar(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("text")
}

func HandleVisualTest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("text")
}