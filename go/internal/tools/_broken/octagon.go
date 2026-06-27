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

func HandleOctagonInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	baseURL := "https://api.octagon-service.internal/v1/info"
	client := http.Client{Timeout: 30 * time.Second}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, apiErr := client.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("API request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	version, _ :=getString(result, "version")
	status, _ :=getString(result, "status")
	uptime, _ :=getString(result, "uptime")

	response := fmt.Sprintf("Octagon Service Info:\nVersion: %s\nStatus: %s\nUptime: %s",
		version, status, uptime)

	return ok(response)
}

func HandleOctagonSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	safeQuery := regexp.MustCompile(`[^a-zA-Z0-9\s\-_]`).ReplaceAllString(query, "")

	baseURL := "https://api.octagon-service.internal/v1/search"
	params := url.Values{}
	params.Add("q", safeQuery)
	params.Add("limit", "10")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	client := http.Client{Timeout: 30 * time.Second}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, apiErr := client.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("API request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result struct {
		Results []map[string]interface{} `json:"results"`
	}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	if len(result.Results) == 0 {
		return ok("No results found for your query")
}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("Search results for '%s':\n", query))
	for i, item := range result.Results {
		id, _ :=getString(item, "id")
		name, _ :=getString(item, "name")
		description, _ :=getString(item, "description")

		response.WriteString(fmt.Sprintf("%d. %s (ID: %s)\n   %s\n", i+1, name, id, description))

	return ok(response.String())
}

}

func HandleOctagonMetrics(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	metricType, _ :=getString(args, "type")
	if metricType == "" {
		metricType = "system"
	}

	baseURL := fmt.Sprintf("https://api.octagon-service.internal/v1/metrics/%s", metricType)
	client := http.Client{Timeout: 30 * time.Second}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, apiErr := client.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("API request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("Octagon %s Metrics:\n", metricType))

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)

	sort.Strings(keys)

	for _, k := range keys {
		v := result[k]
		switch val := v.(type) {
		case float64:
			response.WriteString(fmt.Sprintf("%s: %.2f\n", k, val))
		case string:
			response.WriteString(fmt.Sprintf("%s: %s\n", k, val))
		default:
			response.WriteString(fmt.Sprintf("%s: %v\n", k, val))

	}

	return ok(response.String())
}

}
}

func HandleOctagonConfig(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	action, _ :=getString(args, "action")
	key, _ :=getString(args, "key")
	value, _ :=getString(args, "value")

	if action == "" {
		return err("action parameter is required (get/set)")
}

	baseURL := "https://api.octagon-service.internal/v1/config"
	client := http.Client{Timeout: 30 * time.Second}

	var req *http.Request
	var reqErr error

	switch action {
	case "get":
		if key == "" {
			return err("key parameter is required for get action")
}

		fullURL := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(key))
		req, reqErr = http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	case "set":
		if key == "" || value == "" {
			return err("key and value parameters are required for set action")
}

		data := map[string]string{"value": value}
		jsonData, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			return err(fmt.Sprintf("failed to marshal request: %v", jsonErr))
}

		req, reqErr = http.NewRequestWithContext(ctx, "POST", baseURL+"/"+url.PathEscape(key), strings.NewReader(string(jsonData)))
		if req != nil {
			req.Header.Set("Content-Type", "application/json")

	default:
		return err("invalid action, must be 'get' or 'set'")
}

	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, apiErr := client.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("API request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("unexpected status code: %d - %s", resp.StatusCode, string(body)))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	if action == "get" {
		currentValue, _ :=getString(result, "value")
		return ok(fmt.Sprintf("Configuration value for '%s': %s", key, currentValue))
}

	return ok(fmt.Sprintf("Successfully set '%s' to '%s'", key, value))
}
}