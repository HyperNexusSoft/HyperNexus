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

// HandleUnlaStatus checks the status of Unla services
func HandleUnlaStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	endpoint, _ :=getString(args, "endpoint")
	if endpoint == "" {
		endpoint = "http://localhost:5234/health"
	}

	client := http.Client{Timeout: 30 * time.Second}
	req, reqErr := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	status := fmt.Sprintf("Status: %d\nBody: %s", resp.StatusCode, string(body))
	return ok(status)
}

// HandleUnlaConfigReload triggers a configuration reload via SIGHUP or HTTP
func HandleUnlaConfigReload(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	method, _ :=getString(args, "method")
	if method == "" {
		method = "signal"
	}

	pidFile, _ :=getString(args, "pid_file")
	if pidFile == "" {
		pidFile = "/tmp/mcp-gateway.pid"
	}

	switch method {
	case "signal":
		pidBytes, readErr := os.ReadFile(pidFile)
		if readErr != nil {
			return err("Failed to read PID file: " + readErr.Error())
}

		pidStr := strings.TrimSpace(string(pidBytes))
		pid, parseErr := strconv.Atoi(pidStr)
		if parseErr != nil {
			return err("Failed to parse PID: " + parseErr.Error())
}

		process, procErr := os.FindProcess(pid)
		if procErr != nil {
			return err("Failed to find process: " + procErr.Error())
}

		sendErr := process.Signal(os.Signal(1)) // SIGHUP = 1
		if sendErr != nil {
			return err("Failed to send SIGHUP: " + sendErr.Error())
}

		return ok(fmt.Sprintf("SIGHUP sent to process %d", pid))
}

	case "http":
		endpoint, _ :=getString(args, "endpoint")
		if endpoint == "" {
			endpoint = "http://localhost:5234/reload"
		}

		client := http.Client{Timeout: 30 * time.Second}
		req, reqErr := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
		if reqErr != nil {
			return err(reqErr.Error())
}

		resp, fetchErr := client.Do(req)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return err(readErr.Error())
}

		return ok(fmt.Sprintf("HTTP reload response (%d): %s", resp.StatusCode, string(body)))
}

	default:
		return err("Unknown reload method: " + method)

}

// HandleUnlaValidateConfig validates an MCP gateway configuration file
func HandleUnlaValidateConfig(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	configPath, _ :=getString(args, "config_path")
	if configPath == "" {
		return err("config_path is required")
}

	absPath, absErr := filepath.Abs(configPath)
	if absErr != nil {
		return err(absErr.Error())
}

	info, statErr := os.Stat(absPath)
	if statErr != nil {
		return err("Failed to stat config file: " + statErr.Error())
}

	if info.IsDir() {
		return err("config_path must be a file, not a directory")
}

	data, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return err("Failed to read config file: " + readErr.Error())
}

	content := string(data)
	if strings.TrimSpace(content) == "" {
		return err("Config file is empty")
}

	// Basic YAML structure validation
	lines := strings.Split(content, "\n")
	hasServers := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "servers:") || trimmed == "servers" {
			hasServers = true
			break
		}
	}

	if !hasServers {
		return err("Config missing required 'servers' section")
}

	result := fmt.Sprintf("Config file validated successfully:\n- Path: %s\n- Size: %d bytes\n- Has servers section: yes", absPath, len(data))
	return ok(result)

// HandleUnlaBuildInfo retrieves build information by running the binary
func HandleUnlaBuildInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	binaryPath, _ :=getString(args, "binary_path")
	if binaryPath == "" {
		// Try common paths
		candidates := []string{
			"./mcp-gateway",
			"./cmd/mcp-gateway/mcp-gateway",
			"/usr/local/bin/mcp-gateway",
		}
		for _, candidate := range candidates {
			if _, statErr := os.Stat(candidate); statErr == nil {
				binaryPath = candidate
				break
			}
		}
	}

	if binaryPath == "" {
		return err("binary_path is required and could not be auto-detected")
}

	cmd := exec.CommandContext(ctx, binaryPath, "version")
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return err("Failed to run binary: " + runErr.Error() + "\nOutput: " + string(output))
}

	return ok(strings.TrimSpace(string(output)))
}

// HandleUnlaProxyTest tests an MCP proxy configuration by making a request
func HandleUnlaProxyTest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	gatewayURL, _ :=getString(args, "gateway_url")
	if gatewayURL == "" {
		return err("gateway_url is required")
}

	parsedURL, parseErr := url.Parse(gatewayURL)
	if parseErr != nil {
		return err("Invalid gateway_url: " + parseErr.Error())
}

	endpoint, _ :=getString(args, "endpoint")
	if endpoint == "" {
		endpoint = "/mcp/user/mcp"
	}

	fullURL := parsedURL.String() + endpoint

	client := http.Client{Timeout: 30 * time.Second}

	// Prepare a simple MCP initialize request
	mcpRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "unla-test-client",
				"version": "1.0.0",
			},
		},
	}

	reqBody, marshalErr := json.Marshal(mcpRequest)
	if marshalErr != nil {
		return err("Failed to marshal request: " + marshalErr.Error())
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", fullURL, strings.NewReader(string(reqBody)))
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("Content-Type", "application/json")

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	result := fmt.Sprintf("Proxy test result:\n- URL: %s\n- Status: %d\n- Response: %s", fullURL, resp.StatusCode, string(body))
	return ok(result)
}

// HandleUnlaServerList lists configured MCP servers from a config file
func HandleUnlaServerList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	configPath, _ :=getString(args, "config_path")
	if configPath == "" {
		configPath = "configs/mcp-gateway.yaml"
	}

	absPath, absErr := filepath.Abs(configPath)
	if absErr != nil {
		return err(absErr.Error())
}

	data, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return err("Failed to read config: " + readErr.Error())
}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Extract server names from YAML
	var servers []string
	inServers := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "servers:" {
			inServers = true
			continue
		}
		if inServers {
			if strings.HasPrefix(trimmed, "-") {
				// Extract name from list item
				name := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
				if name != "" {
					servers = append(servers, name)

			} else if trimmed != "" && !strings.Contains(trimmed, ":") {
				break
			}
		}
	}

	if len(servers) == 0 {
		return ok("No servers found in configuration")
}

	result := "Configured MCP Servers:\n"
	for i, server := range servers {
		result += fmt.Sprintf("%d. %s\n", i+1, server)

	return ok(strings.TrimSpace(result))
}
}
}