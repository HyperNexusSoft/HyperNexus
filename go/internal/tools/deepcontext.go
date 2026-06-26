package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	validContextName = regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,64}$`)
)

// HandleCreateContext creates a new context with the given name and metadata
func HandleCreateContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("context name is required")
}

	if !validContextName.MatchString(name) {
		return err("invalid context name format")
}

	metadata := make(map[string]interface{})
	if meta, found := args["metadata"].(map[string]interface{}); found {
		metadata = meta
	}

	// In a real implementation, this would persist the context
	// For MCP tool purposes, we'll just simulate success
	return ok(fmt.Sprintf("context '%s' created successfully", name))
}

// HandleGetContext retrieves context information by name
func HandleGetContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("context name is required")
}

	if !validContextName.MatchString(name) {
		return err("invalid context name format")
}

	// Simulate retrieving context data
	result := map[string]interface{}{
		"name":     name,
		"created":  time.Now().Format(time.RFC3339),
		"metadata": map[string]string{"status": "active"},
	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal context data: %v", jsonErr))
}

	return ok(string(jsonData))
}

// HandleListContexts lists available contexts with optional filtering
func HandleListContexts(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Simulate listing contexts
	contexts := []map[string]interface{}{
		{
			"name":    "default",
			"created": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			"name":    "production",
			"created": time.Now().Add(-72 * time.Hour).Format(time.RFC3339),
		},
	}

	filter, _ :=getString(args, "filter")
	if filter != "" {
		filtered := []map[string]interface{}{}
		for _, ctx := range contexts {
			if strings.Contains(strings.ToLower(ctx["name"].(string)), strings.ToLower(filter)) {
				filtered = append(filtered, ctx)

		}
		contexts = filtered
	}

	jsonData, jsonErr := json.Marshal(contexts)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal contexts list: %v", jsonErr))
}

	return ok(string(jsonData))
}

}

// HandleAnalyzeContext performs analysis on a specific context
func HandleAnalyzeContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("context name is required")
}

	depth, _ :=getInt(args, "depth")
	if depth == 0 {
		depth = 3 // default depth
	}

	// Simulate analysis results
	analysis := map[string]interface{}{
		"context": name,
		"depth":   depth,
		"findings": []map[string]interface{}{
			{
				"type":    "dependency",
				"item":    "database",
				"status":  "healthy",
				"details": "Connection pool optimized",
			},
			{
				"type":    "configuration",
				"item":    "timeout",
				"status":  "warning",
				"details": "Timeout value could be increased",
			},
		},
	}

	jsonData, jsonErr := json.Marshal(analysis)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal analysis results: %v", jsonErr))
}

	return ok(string(jsonData))
}

// HandleDeleteContext removes a context by name
func HandleDeleteContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("context name is required")
}

	if !validContextName.MatchString(name) {
		return err("invalid context name format")
}

	// Simulate deletion
	return ok(fmt.Sprintf("context '%s' deleted successfully", name))
}