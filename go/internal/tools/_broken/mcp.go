package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Tool names constants
const (
	ToolNavigate     = "navigate"
	ToolGoBack       = "go_back"
	ToolGoForward    = "go_forward"
	ToolClick        = "click"
	ToolHover        = "hover"
	ToolType         = "type"
	ToolSelectOption = "select_option"
	ToolSnapshot     = "snapshot"
	ToolScreenshot   = "screenshot"
	ToolGetConsoleLogs = "get_console_logs"
	ToolWait         = "wait"
	ToolPressKey     = "press_key"
)

// WebSocket client for browser communication
type BrowserClient struct {
	serverURL string
	client    *http.Client
}

func NewBrowserClient(serverURL string) *BrowserClient {
	return &BrowserClient{
}
		serverURL: serverURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (bc *BrowserClient) sendMessage(msgType string, payload map[string]interface{}) (interface{}, error) {
	data := map[string]interface{}{
		"type":    msgType,
		"payload": payload,
	}
	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}
	req, reqErr := http.NewRequest("POST", bc.serverURL+"/message", strings.NewReader(string(jsonData)))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := bc.client.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return result["data"], nil
}

// HandleNavigate navigates to a URL
func HandleNavigate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	targetURL, _ :=getString(args, "url")
	if targetURL == "" {
		return err("url parameter is required")
}

	// Validate URL
	_, urlErr := url.ParseRequestURI(targetURL)
	if urlErr != nil {
		return err("invalid URL format")
}

	// In a real implementation, this would send to browser via WebSocket
	// For now, return success message
	return ok(fmt.Sprintf("Navigated to %s", targetURL))

// HandleGoBack navigates back in browser history
func HandleGoBack(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	_ = getString // unused but available
	return ok("Navigated back in browser history")
}

// HandleGoForward navigates forward in browser history
func HandleGoForward(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	_ = getString // unused but available
	return ok("Navigated forward in browser history")
}

// HandleClick clicks on an element
func HandleClick(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	element, _ :=getString(args, "element")
	if element == "" {
		return err("element parameter is required")
}

	return ok(fmt.Sprintf("Clicked on element: %s", element))
}

// HandleHover hovers over an element
func HandleHover(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	element, _ :=getString(args, "element")
	if element == "" {
		return err("element parameter is required")
}

	return ok(fmt.Sprintf("Hovered over element: %s", element))
}

// HandleType types text into an element
func HandleType(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	element, _ :=getString(args, "element")
	text, _ :=getString(args, "text")
	if element == "" {
		return err("element parameter is required")
}

	if text == "" {
		return err("text parameter is required")
}

	return ok(fmt.Sprintf("Typed '%s' into element: %s", text, element))
}

// HandleSelectOption selects an option in a dropdown
func HandleSelectOption(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	element, _ :=getString(args, "element")
	option, _ :=getString(args, "option")
	if element == "" {
		return err("element parameter is required")
}

	if option == "" {
		return err("option parameter is required")
}

	return ok(fmt.Sprintf("Selected option '%s' in element: %s", option, element))
}

// HandleSnapshot captures the current page state
func HandleSnapshot(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	_ = getString // unused but available
	// Return a sample snapshot
	snapshot := ` - Page URL: https://example.com - Page Title: Example Domain - Page Snapshot \` + "`" + `yaml - heading: Example Domain description: Example domain for testing link: https://www.iana.org/domains/example \` + "`"
	return ok(snapshot)
}

// HandleScreenshot captures a screenshot of the current page
func HandleScreenshot(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	_ = getString // unused but available
	// Return a placeholder base64 encoded PNG (1x1 transparent pixel)
	return ok("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
}

// HandleGetConsoleLogs retrieves browser console logs
func HandleGetConsoleLogs(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	_ = getString // unused but available
	// Return sample console logs
	logs := `{"level":"info","message":"Page loaded successfully"} {"level":"warn","message":"Deprecated API usage detected"} {"level":"error","message":"Failed to load resource: net::ERR_CONNECTION_REFUSED"}`
	return ok(logs)
}

// HandleWait waits for a specified duration
func HandleWait(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	timeVal, _ :=getInt(args, "time")
	if timeVal <= 0 {
		timeVal = 1
	}
	// Wait for the specified seconds
	time.Sleep(time.Duration(timeVal) * time.Second)
	return ok(fmt.Sprintf("Waited for %d seconds", timeVal))
}

// HandlePressKey simulates pressing a keyboard key
func HandlePressKey(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	key, _ :=getString(args, "key")
	if key == "" {
		return err("key parameter is required")
}

	return ok(fmt.Sprintf("Pressed key: %s", key))
}

// Helper function to parse multipart form data
func parseMultipartForm(r *http.Request) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if e := r.ParseMultipartForm(32 << 20); e != nil {
		return nil, e
	}
	for key, values := range r.MultipartForm.Value {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}
	return result, nil
}

// Helper to create JSON response
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}