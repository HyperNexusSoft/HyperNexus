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

func HandleListPages(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	const apiURL = "http://localhost:9222/json/list"
	client := http.DefaultClient

	resp, fetchErr := client.Get(apiURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch pages: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var pages []map[string]interface{}
	if parseErr := json.Unmarshal(body, &pages); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	var result strings.Builder
	for i, page := range pages {
		if i > 0 {
			result.WriteString("\n")

		title, _ :=getString(page, "title")
		url, _ :=getString(page, "url")
		id, _ :=getString(page, "id")
		result.WriteString(fmt.Sprintf("Page %d: %s (%s) [id: %s]", i+1, title, url, id))

	return ok(result.String())
}

}
}

func HandleBrowsePage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	targetURL, _ :=getString(args, "url")
	if targetURL == "" {
		return err("url parameter is required")
}

	parsedURL, parseErr := url.ParseRequestURI(targetURL)
	if parseErr != nil {
		return err(fmt.Sprintf("invalid URL: %v", parseErr))
}

	if !strings.HasPrefix(parsedURL.Scheme, "http") {
		return err("URL must start with http:// or https://")
}

	const apiURL = "http://localhost:9222/json/new"
	client := http.DefaultClient

	resp, fetchErr := client.Get(fmt.Sprintf("%s?%s", apiURL, url.QueryEscape(targetURL)))
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to create new page: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	pageID, _ :=getString(result, "id")
	return ok(fmt.Sprintf("Successfully navigated to %s (page ID: %s)", targetURL, pageID))
}

func HandleTakeScreenshot(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pageID, _ :=getString(args, "page_id")
	if pageID == "" {
		return err("page_id parameter is required")
}

	const apiURL = "http://localhost:9222/json/captureScreenshot"
	client := http.DefaultClient

	reqBody := map[string]interface{}{
		"id": pageID,
	}
	jsonBody, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", marshalErr))
}

	resp, fetchErr := client.Post(apiURL, "application/json", strings.NewReader(string(jsonBody)))
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to take screenshot: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	data, _ :=getString(result, "data")
	if data == "" {
		return err("no screenshot data returned")
}

	return ok(fmt.Sprintf("Screenshot captured (base64 data length: %d)", len(data)))
}

func HandleGetPerformanceMetrics(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pageID, _ :=getString(args, "page_id")
	if pageID == "" {
		return err("page_id parameter is required")
}

	const apiURL = "http://localhost:9222/json/getPerformanceMetrics"
	client := http.DefaultClient

	reqBody := map[string]interface{}{
		"id": pageID,
	}
	jsonBody, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", marshalErr))
}

	resp, fetchErr := client.Post(apiURL, "application/json", strings.NewReader(string(jsonBody)))
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to get performance metrics: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	metrics := getSlice(result, "metrics")
	if metrics == nil {
		return err("no metrics returned")
}

	var output strings.Builder
	for _, metric := range metrics {
		m := metric.(map[string]interface{})
		name, _ :=getString(m, "name")
		value := getFloat64(m, "value")
		output.WriteString(fmt.Sprintf("%s: %.2f\n", name, value))

	return ok(output.String())
}

}

func HandleEvaluateJavascript(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pageID, _ :=getString(args, "page_id")
	if pageID == "" {
		return err("page_id parameter is required")
}

	expression, _ :=getString(args, "expression")
	if expression == "" {
		return err("expression parameter is required")
}

	const apiURL = "http://localhost:9222/json/evaluate"
	client := http.DefaultClient

	reqBody := map[string]interface{}{
		"id":         pageID,
		"expression": expression,
	}
	jsonBody, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", marshalErr))
}

	resp, fetchErr := client.Post(apiURL, "application/json", strings.NewReader(string(jsonBody)))
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to evaluate JavaScript: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	resultValue, _ :=getString(result, "result")
	if resultValue == "" {
		return err("no result returned from JavaScript evaluation")
}

	return ok(fmt.Sprintf("JavaScript evaluation result: %s", resultValue))
}

func HandleGetNetworkRequests(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pageID, _ :=getString(args, "page_id")
	if pageID == "" {
		return err("page_id parameter is required")
}

	const apiURL = "http://localhost:9222/json/getNetworkRequests"
	client := http.DefaultClient

	reqBody := map[string]interface{}{
		"id": pageID,
	}
	jsonBody, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", marshalErr))
}

	resp, fetchErr := client.Post(apiURL, "application/json", strings.NewReader(string(jsonBody)))
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to get network requests: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON response: %v", parseErr))
}

	requests := getSlice(result, "requests")
	if requests == nil {
		return err("no network requests returned")
}

	var output strings.Builder
	for i, req := range requests {
		r := req.(map[string]interface{})
		if i > 0 {
			output.WriteString("\n")

		method, _ :=getString(r, "method")
		url, _ :=getString(r, "url")
		status, _ :=getInt(r, "status")
		output.WriteString(fmt.Sprintf("Request %d: %s %s [%d]", i+1, method, url, status))

	return ok(output.String())
}

}
}

// Helper functions for type conversion
func getSlice(m map[string]interface{}, key string) []interface{} {
	if val, found := m[key]; found {
		if slice, found := val.([]interface{}); found {
			return slice
		}
	}
	return nil
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if val, found := m[key]; found {
		switch v := val.(type) {
		case float64:
			return v
}
		case float32:
			return float64(v)
}
		case int:
			return float64(v)
}
		case int32:
			return float64(v)
		case int64:
			return float64(v)

	}
	return 0
}

func getInt(m map[string]interface{}, key string) int {
	if val, found := m[key]; found {
		switch v := val.(type) {
		case int:
			return v
}
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		case float32:
			return int(v)
		case json.Number:
			if i, e := v.Int64(); e == nil {
				return int(i)

		}
	}
	return 0
}