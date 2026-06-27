package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HandleSearch searches a codebase using semble for fast, token-efficient code search.
// It supports local paths and git URLs, with configurable content type and result count.
func HandleSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	path, _ :=getString(args, "path")
	content, _ :=getString(args, "content")
	topKStr, _ :=getString(args, "top_k")

	if query == "" {
		return err("query is required")
}

	if path == "" {
		path = "."
	}

	// Build the command arguments
	cmdArgs := []string{"search", query, path}
	if content != "" {
		cmdArgs = append(cmdArgs, "--content", content)

	if topKStr != "" {
		cmdArgs = append(cmdArgs, "--top-k", topKStr)

	// Try uvx first, then semble directly
	output, runErr := runSemble(ctx, cmdArgs)
	if runErr != nil {
		return err(fmt.Sprintf("semble search failed: %s", runErr.Error()))
}

	return ok(strings.TrimSpace(string(output)))
}

}
}

// HandleFindRelated finds code similar to a known file and line location.
// Useful for discovering related implementations, interfaces, or usages.
func HandleFindRelated(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	file, _ :=getString(args, "file")
	lineStr, _ :=getString(args, "line")
	path, _ :=getString(args, "path")

	if file == "" {
		return err("file is required")
}

	if lineStr == "" {
		return err("line is required")
}

	if path == "" {
		path = "."
	}

	line, convErr := strconv.Atoi(lineStr)
	if convErr != nil {
		return err(fmt.Sprintf("invalid line number: %s", convErr.Error()))
}

	cmdArgs := []string{"find-related", file, strconv.Itoa(line), path}

	output, runErr := runSemble(ctx, cmdArgs)
	if runErr != nil {
		return err(fmt.Sprintf("semble find-related failed: %s", runErr.Error()))
}

	return ok(strings.TrimSpace(string(output)))
}

// HandleSavings shows token savings statistics across all semble searches.
// Reports total tokens saved, call counts, and efficiency ratios.
func HandleSavings(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cmdArgs := []string{"savings"}

	output, runErr := runSemble(ctx, cmdArgs)
	if runErr != nil {
		return err(fmt.Sprintf("semble savings failed: %s", runErr.Error()))
}

	return ok(strings.TrimSpace(string(output)))
}

// HandleIndexInfo returns information about the semble index for a given path,
// including file count, chunk count, and index metadata.
func HandleIndexInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	content, _ :=getString(args, "content")

	if path == "" {
		path = "."
	}

	// Check if semble is available and get version info
	versionArgs := []string{"--version"}
	versionOutput, verErr := runSemble(ctx, versionArgs)
	versionStr := "unknown"
	if verErr == nil {
		versionStr = strings.TrimSpace(string(versionOutput))

	// Check if the path exists
	absPath, pathErr := filepath.Abs(path)
	if pathErr != nil {
		return err(fmt.Sprintf("invalid path: %s", pathErr.Error()))
}

	info := map[string]interface{}{
		"path":         absPath,
		"version":      versionStr,
		"content_type": content,
	}

	// Check if path exists
	fileInfo, statErr := os.Stat(absPath)
	if statErr != nil {
		info["exists"] = false
		info["error"] = statErr.Error()
	} else {
		info["exists"] = true
		info["is_directory"] = fileInfo.IsDir()

	// Check for .sembleignore
	ignorePath := filepath.Join(absPath, ".sembleignore")
	if _, ignoreErr := os.Stat(ignorePath); ignoreErr == nil {
		info["has_sembleignore"] = true
	} else {
		info["has_sembleignore"] = false
	}

	// Check for .gitignore
	gitignorePath := filepath.Join(absPath, ".gitignore")
	if _, gitignoreErr := os.Stat(gitignorePath); gitignoreErr == nil {
		info["has_gitignore"] = true
	} else {
		info["has_gitignore"] = false
	}

	// Check if it's a git repo
	gitDir := filepath.Join(absPath, ".git")
	if _, gitErr := os.Stat(gitDir); gitErr == nil {
		info["is_git_repo"] = true
	} else {
		info["is_git_repo"] = false
	}

	jsonBytes, jsonErr := json.MarshalIndent(info, "", "  ")
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal info: %s", jsonErr.Error()))
}

	return ok(string(jsonBytes))
}

}
}

// runSemble tries to run semble via uvx or directly, returning the output.
func runSemble(ctx context.Context, args []string) ([]byte, error) {
	// Try running semble directly first
	if _, lookErr := exec.LookPath("semble"); lookErr == nil {
		return runCommand(ctx, "semble", args)
}

	// Fall back to uvx
	if _, lookErr := exec.LookPath("uvx"); lookErr == nil {
		uvArgs := append([]string{"--from", "semble[mcp]", "semble"}, args...)
		return runCommand(ctx, "uvx", uvArgs)
}

	// Try uv tool run
	if _, lookErr := exec.LookPath("uv"); lookErr == nil {
		uvArgs := append([]string{"tool", "run", "--from", "semble[mcp]", "semble"}, args...)
		return runCommand(ctx, "uv", uvArgs)
}

	return nil, fmt.Errorf("semble not found: install with 'uv tool install semble' or 'pip install semble'")
}

// runCommand executes a command with the given arguments and returns its output.
func runCommand(ctx context.Context, name string, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, cmdErr := cmd.Output()
	if cmdErr != nil {
		if exitErr, found := cmdErr.(*exec.ExitError); found {
			return nil, fmt.Errorf("%s: %s", cmdErr.Error(), strings.TrimSpace(string(exitErr.Stderr)))
}

		return nil, cmdErr
	}
	return output, nil
}

// HandleSearchRemote searches a remote git repository by cloning it on demand.
// The repository URL is provided and semble handles cloning and caching.
func HandleSearchRemote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	repoURL, _ :=getString(args, "repo_url")
	content, _ :=getString(args, "content")
	topKStr, _ :=getString(args, "top_k")

	if query == "" {
		return err("query is required")
}

	if repoURL == "" {
		return err("repo_url is required")
}

	// Validate URL
	if !strings.HasPrefix(repoURL, "https://") && !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "git@") {
		return err("repo_url must be a valid git URL (https://, http://, or git@)")
}

	cmdArgs := []string{"search", query, repoURL}
	if content != "" {
		cmdArgs = append(cmdArgs, "--content", content)

	if topKStr != "" {
		cmdArgs = append(cmdArgs, "--top-k", topKStr)

	output, runErr := runSemble(ctx, cmdArgs)
	if runErr != nil {
		return err(fmt.Sprintf("semble remote search failed: %s", runErr.Error()))
}

	return ok(strings.TrimSpace(string(output)))
}

}
}

// HandleInstall sets up semble integration with coding agents.
// Detects installed agents and configures MCP server, instructions, or sub-agent.
func HandleInstall(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	mode, _ :=getString(args, "mode")

	cmdArgs := []string{"install"}
	if mode != "" {
		cmdArgs = append(cmdArgs, "--mode", mode)

	// install is interactive, so we provide guidance instead
	// Check if semble is available
	if _, lookErr := exec.LookPath("semble"); lookErr != nil {
		// Check uvx
		if _, uvxErr := exec.LookPath("uvx"); uvxErr != nil {
			return err("semble is not installed. Install with: uv tool install semble")

	}

	// Return installation guidance since install is interactive
	guidance := `Semble Install Guide:

To install semble integrations, run interactively in your terminal:
  semble install

This detects installed coding agents (Claude Code, Codex, OpenCode, etc.)
and lets you choose which integrations to enable:

1. MCP server - lets the agent call semble directly as a tool
2. Instructions - adds CLI usage guidance to AGENTS.md / CLAUDE.md
3. Sub-agent - installs a dedicated semble-search sub-agent

To undo setup:
  semble uninstall

Manual MCP config example:
{
  "mcpServers": {
    "semble": {
      "command": "uvx",
      "args": ["--from", "semble[mcp]", "semble", "mcp"]
    }
  }
}

For detailed instructions, see: https://github.com/MinishLab/semble/blob/main/docs/installation.md`

	return ok(guidance)
}

}
}

// HandleCacheInfo returns information about the semble cache location and contents.
func HandleCacheInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cacheDir := os.Getenv("SEMBLE_CACHE_LOCATION")
	if cacheDir == "" {
		// Determine OS-specific default
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			home = "~"
		}
		// Check platform
		goos := os.Getenv("GOOS")
		if goos == "" {
			goos = "linux"
		}
		switch goos {
		case "darwin":
			cacheDir = filepath.Join(home, "Library", "Caches", "semble")
		case "windows":
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData == "" {
				localAppData = filepath.Join(home, "AppData", "Local")

			cacheDir = filepath.Join(localAppData, "semble", "Cache")
		default:
			cacheDir = filepath.Join(home, ".cache", "semble")

	}

	info := map[string]interface{}{
		"cache_location": cacheDir,
		"env_override":   os.Getenv("SEMBLE_CACHE_LOCATION") != "",
	}

	// Check if cache directory exists
	if stat, statErr := os.Stat(cacheDir); statErr == nil {
		info["exists"] = true
		info["is_directory"] = stat.IsDir()

		// Try to list contents
		entries, readErr := os.ReadDir(cacheDir)
		if readErr == nil {
			names := make([]string, 0, len(entries))
			for _, entry := range entries {
				names = append(names, entry.Name())

			info["contents"] = names
			info["item_count"] = len(names)

	} else {
		info["exists"] = false
	}

	jsonBytes, jsonErr := json.MarshalIndent(info, "", "  ")
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal cache info: %s", jsonErr.Error()))
}

	return ok(string(jsonBytes))
}

// Unused import guard - these are used by helper functions
var _ = io.EOF
var _ = http.Client{Timeout: 30 * time.Second}
var _ = url.Values{}
}
}
}
}