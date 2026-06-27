package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ToolResponse, ok, e, getString, getInt, getBool, TextContent are defined in parity.go

const semgrepAPIBase = "https://semgrep.dev/api/v1"

func HandleScan(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	limit, _ :=getInt(args, "limit")
	if limit == 0 {
		limit = 10
	}

	client := http.DefaultClient

	// Construct URL with query parameters
	baseURL := fmt.Sprintf("%s/rules/search", semgrepAPIBase)
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", strconv.Itoa(limit))

	fullURL := baseURL + "?" + params.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch data: %v", fetchErr))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)))
}

	// Parse JSON response
	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON: %v", parseErr))
}

	// Format output
	output, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format output: %v", marshalErr))
}

	return ok(string(output))
}

func HandleRuleInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ruleID, _ :=getString(args, "rule_id")
	if ruleID == "" {
		return err("rule_id parameter is required")
}

	client := http.DefaultClient

	fullURL := fmt.Sprintf("%s/rules/%s", semgrepAPIBase, url.QueryEscape(ruleID))

	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch rule info: %v", fetchErr))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)))
}

	return ok(string(body))
}

func HandleListLanguages(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient

	fullURL := fmt.Sprintf("%s/languages", semgrepAPIBase)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch languages: %v", fetchErr))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)))
}

	// Pretty print the JSON
	var data interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON: %v", parseErr))
}

	output, marshalErr := json.MarshalIndent(data, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format output: %v", marshalErr))
}

	return ok(string(output))
}

func HandleSearchRules(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	keyword, _ :=getString(args, "keyword")
	category, _ :=getString(args, "category")
	language, _ :=getString(args, "language")

	client := http.DefaultClient

	baseURL := fmt.Sprintf("%s/rules/search", semgrepAPIBase)
	params := url.Values{}

	if keyword != "" {
		params.Set("q", keyword)

	if category != "" {
		params.Set("category", category)

	if language != "" {
		params.Set("language", language)

	fullURL := baseURL + "?" + params.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to search rules: %v", fetchErr))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)))
}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON: %v", parseErr))
}

	output, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format output: %v", marshalErr))
}

	return ok(string(output))
}
}
}
}