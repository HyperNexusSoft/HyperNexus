package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// HandlePing implements the ping tool to verify connectivity.
func HandlePing(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("pong")
}

// HandleListDir implements the list_directory tool to list files in a directory.
func HandleListDir(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("path argument is required")
}

	entries, apiErr := os.ReadDir(path)
	if apiErr != nil {
		return err(apiErr.Error())
}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())

	sort.Strings(names)

	result := fmt.Sprintf("Directory contents of %s:\n%s", path, strings.Join(names, "\n"))
	return ok(result)
}

}

// HandleReadFile implements the read_file tool to read file contents.
func HandleReadFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("path argument is required")
}

	content, apiErr := os.ReadFile(path)
	if apiErr != nil {
		return err(apiErr.Error())
}

	return ok(string(content))
}

// HandleWriteFile implements the write_file tool to write content to a file.
func HandleWriteFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	content, _ :=getString(args, "content")

	if path == "" {
		return err("path argument is required")
}

	apiErr := os.WriteFile(path, []byte(content), 0644)
	if apiErr != nil {
		return err(apiErr.Error())
}

	return ok(fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), path))
}

// HandleRunCommand implements the run_command tool to execute shell commands.
func HandleRunCommand(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	command, _ :=getString(args, "command")
	if command == "" {
		return err("command argument is required")
}

	// Parse command into args
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return err("command cannot be empty")
}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, apiErr := cmd.CombinedOutput()
	if apiErr != nil {
		// Return the output even if there was an error, as it might contain useful info
		return ok(fmt.Sprintf("Command failed with error: %v\nOutput:\n%s", apiErr, string(output)))
}

	return ok(string(output))
}

// HandleFetch implements the fetch tool to make HTTP requests.
func HandleFetch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url argument is required")
}

	// Validate URL
	if _, apiErr := url.Parse(urlStr); apiErr != nil {
		return err("invalid url format")
}

	client := http.DefaultClient
	req, apiErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if apiErr != nil {
		return err(apiErr.Error())
}

	resp, apiErr := client.Do(req)
	if apiErr != nil {
		return err(apiErr.Error())
}

	defer resp.Body.Close()

	body, apiErr := io.ReadAll(resp.Body)
	if apiErr != nil {
		return err(apiErr.Error())
}

	return ok(fmt.Sprintf("Status: %s\nBody:\n%s", resp.Status, string(body)))
}

// HandleSearchFile implements the search_file tool to find files by pattern.
func HandleSearchFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	root, _ :=getString(args, "root")
	pattern, _ :=getString(args, "pattern")

	if root == "" {
		root = "."
	}
	if pattern == "" {
		return err("pattern argument is required")
}

	compiledPattern, apiErr := regexp.Compile(pattern)
	if apiErr != nil {
		return err(fmt.Sprintf("invalid regex pattern: %v", apiErr))
}

	var matches []string
	apiErr = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if compiledPattern.MatchString(info.Name()) {
			matches = append(matches, path)

		return nil
	})

	if apiErr != nil {
		return err(apiErr.Error())
}

	if len(matches) == 0 {
		return ok("No matches found.")
}

	sort.Strings(matches)
	return ok(strings.Join(matches, "\n"))
}

}

// HandleGetEnv implements the get_env tool to retrieve environment variables.
func HandleGetEnv(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	key, _ :=getString(args, "key")
	if key == "" {
		return err("key argument is required")
}

	val, exists := os.LookupEnv(key)
	if !exists {
		return err(fmt.Sprintf("Environment variable %s not found", key))
}

	return ok(val)
}

// HandleSetEnv implements the set_env tool to set environment variables for the current process.
func HandleSetEnv(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	key, _ :=getString(args, "key")
	value, _ :=getString(args, "value")

	if key == "" {
		return err("key argument is required")
}

	apiErr := os.Setenv(key, value)
	if apiErr != nil {
		return err(apiErr.Error())
}

	return ok(fmt.Sprintf("Environment variable %s set to %s", key, value))
}

// HandleSleep implements the sleep tool to pause execution.
func HandleSleep(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	durationStr, _ :=getString(args, "seconds")
	if durationStr == "" {
		return err("seconds argument is required")
}

	seconds, parseErr := strconv.ParseFloat(durationStr, 64)
	if parseErr != nil {
		return err("invalid seconds value")
}

	time.Sleep(time.Duration(seconds * float64(time.Second)))
	return ok(fmt.Sprintf("Slept for %.2f seconds", seconds))
}