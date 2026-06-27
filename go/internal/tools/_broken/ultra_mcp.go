package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var http.DefaultClient = http.DefaultClient

// HandleUltraFetch performs a high-performance HTTP GET request with timeout support.
func HandleUltraFetch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	target, _ :=getString(args, "url")
	if target == "" {
		return err("url is required")
}

	method, _ :=getString(args, "method")
	if method == "" {
		method = "GET"
	}

	req, reqErr := http.NewRequestWithContext(ctx, method, target, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	// Add custom headers if provided
	if headersRaw, exists := args["headers"]; exists {
		if headersMap, found := headersRaw.(map[string]interface{}); found {
			for k, v := range headersMap {
				if valStr, found := v.(string); found {
					req.Header.Add(k, valStr)

			}
		}
	}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	result := map[string]interface{}{
		"status":  resp.StatusCode,
		"headers": resp.Header,
		"body":    string(body),
	}

	jsonBytes, jsonErr := json.MarshalIndent(result, "", "  ")
	if jsonErr != nil {
		return err(jsonErr.Error())
}

	return ok(string(jsonBytes))
}

}

// HandleUltraExec executes a shell command and returns the combined output.
func HandleUltraExec(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	command, _ :=getString(args, "command")
	if command == "" {
		return err("command is required")
}

	// Split command string into arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return err("invalid command format")
}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return err(fmt.Sprintf("execution failed: %s\nOutput: %s", execErr.Error(), string(output)))
}

	return ok(string(output))
}

// HandleUltraSearch searches for a regex pattern within files in a directory.
func HandleUltraSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pattern, _ :=getString(args, "pattern")
	dir, _ :=getString(args, "dir")

	if pattern == "" {
		return err("pattern is required")
}

	if dir == "" {
		dir = "."
	}

	re, reErr := regexp.Compile(pattern)
	if reErr != nil {
		return err(reErr.Error())
}

	var results []string

	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		// Simple check for text files based on extension
		if !strings.HasSuffix(strings.ToLower(path), ".txt") &&
			!strings.HasSuffix(strings.ToLower(path), ".go") &&
			!strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil // Skip unreadable files
		}

		matches := re.FindAllString(string(content), -1)
		if len(matches) > 0 {
			results = append(results, fmt.Sprintf("File: %s\nMatches: %d", path, len(matches)))

		return nil
	})

	if walkErr != nil {
		return err(walkErr.Error())
}

	if len(results) == 0 {
		return ok("No matches found.")
}

	return ok(strings.Join(results, "\n\n"))
}

}

// HandleUltraRead reads the content of a specific file.
func HandleUltraRead(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("path is required")
}

	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(string(content))
}

// HandleUltraWrite writes content to a specific file.
func HandleUltraWrite(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	content, _ :=getString(args, "content")

	if path == "" {
		return err("path is required")
}

	writeErr := os.WriteFile(path, []byte(content), 0644)
	if writeErr != nil {
		return err(writeErr.Error())
}

	return ok(fmt.Sprintf("Successfully wrote to %s", path))
}