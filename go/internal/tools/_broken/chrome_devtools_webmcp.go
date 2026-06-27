package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Note: ToolResponse, ok, e, getString, getInt, getBool, TextContent are defined in parity.go

// Helper to make HTTP requests with timeout
func makeRequest(ctx context.Context, method, urlStr string, body io.Reader) ([]byte, error) {
	client := http.DefaultClient
	req, e := http.NewRequestWithContext(ctx, method, urlStr, body)
	if e != nil {
		return nil, e
	}
	req.Header.Set("User-Agent", "chrome-devtools-webmcp/1.0")
	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// HandleNavigate navigates the browser to a specific URL
func HandleNavigate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	// Validate URL format
	if _, parseErr := url.Parse(urlStr); parseErr != nil {
		return err(fmt.Sprintf("invalid URL: %v", parseErr))
}

	// Simulate navigation by making a request to verify accessibility
	// In a real implementation, this would interact with a Chrome DevTools Protocol endpoint
	body, fetchErr := makeRequest(ctx, "GET", urlStr, nil)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to navigate/verify URL: %v", fetchErr))
}

	return ok(fmt.Sprintf("Successfully navigated to %s. Response length: %d bytes", urlStr, len(body)))
}

// HandleGetPageInfo retrieves basic information about the current page
func HandleGetPageInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// In a real scenario, this would query the CDP for page metadata
	// Here we simulate fetching metadata from the provided URL if given
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	body, fetchErr := makeRequest(ctx, "GET", urlStr, nil)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch page info: %v", fetchErr))
}

	// Simple heuristic to extract title (very basic HTML parsing)
	title := "Unknown Title"
	re := regexp.MustCompile(`(?i)<title>([^<]+)</title>`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		title = strings.TrimSpace(matches[1])

	return ok(fmt.Sprintf("Page Title: %s\nContent Length: %d bytes", title, len(body)))
}

}

// HandleExecuteScript executes a simple JavaScript snippet (simulated)
func HandleExecuteScript(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	script, _ :=getString(args, "script")
	if script == "" {
		return err("missing required argument: script")
}

	// Basic validation: ensure script doesn't contain obviously malicious patterns
	// In a real CDP implementation, this would be sent to the runtime.evaluate method
	if strings.Contains(script, "alert(") || strings.Contains(script, "confirm(") {
		return err("blocked: script contains potentially unsafe functions")
}

	// Simulate execution result
	result := fmt.Sprintf("Executed: %s", script)
	if strings.Contains(script, "document.title") {
		result = "Result: <current-page-title>"
	}

	return ok(result)
}

// HandleGetCookies retrieves cookies for a specific domain
func HandleGetCookies(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	domain, _ :=getString(args, "domain")
	if domain == "" {
		return err("missing required argument: domain")
}

	// Simulate fetching cookies (in real CDP, use Network.getCookies)
	// We'll just return a mock structure for demonstration
	cookies := []map[string]string{
		{"name": "session_id", "value": "abc123", "domain": domain},
		{"name": "user_pref", "value": "dark_mode", "domain": domain},
	}

	jsonData, marshalErr := json.Marshal(cookies)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to serialize cookies: %v", marshalErr))
}

	return ok(fmt.Sprintf("Found %d cookies for %s:\n%s", len(cookies), domain, string(jsonData)))
}

// HandleSearchPage searches for text within the current page
func HandleSearchPage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("missing required argument: query")
}

	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	body, fetchErr := makeRequest(ctx, "GET", urlStr, nil)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch page for search: %v", fetchErr))
}

	content := string(body)
	count := strings.Count(strings.ToLower(content), strings.ToLower(query))

	return ok(fmt.Sprintf("Found %d occurrences of '%s' in %s", count, query, urlStr))
}

// HandleGetResources lists resources (scripts, stylesheets) on the page
func HandleGetResources(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	body, fetchErr := makeRequest(ctx, "GET", urlStr, nil)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch page resources: %v", fetchErr))
}

	content := string(body)
	var resources []string

	// Extract script sources
	scriptRe := regexp.MustCompile(`<script[^>]+src=["']([^"']+)["']`)")
	for _, match := range scriptRe.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			resources = append(resources, match[1])

	}

	// Extract stylesheet links
	linkRe := regexp.MustCompile(`<link[^>]+href=["']([^"']+)["'][^>]*rel=["']stylesheet["']`)")
	for _, match := range linkRe.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			resources = append(resources, match[1])

	}

	// Sort for consistent output
	sort.Strings(resources)

	return ok(fmt.Sprintf("Found %d resources:\n- %s", len(resources), strings.Join(resources, "\n- ")))
}
}
}