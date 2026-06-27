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

// HandleFetchURL fetches the content of a given URL and returns it as text.
func HandleFetchURL(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing 'url' parameter")
}

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch URL: status code %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(string(body))
}

// HandleParseJSON parses a JSON string and returns a specific key's value.
func HandleParseJSON(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	jsonStr, _ :=getString(args, "json")
	key, _ :=getString(args, "key")
	if jsonStr == "" || key == "" {
		return err("missing 'json' or 'key' parameter")
}

	var data map[string]interface{}
	parseErr := json.Unmarshal([]byte(jsonStr), &data)
	if parseErr != nil {
		return err(parseErr.Error())
}

	value, exists := data[key]
	if !exists {
		return err(fmt.Sprintf("key '%s' not found in JSON", key))
}

	return ok(fmt.Sprintf("%v", value))
}

// HandleExtractDomain extracts the domain from a given URL.
func HandleExtractDomain(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing 'url' parameter")
}

	parsedURL, parseErr := url.Parse(urlStr)
	if parseErr != nil {
		return err(parseErr.Error())
}

	domain := parsedURL.Hostname()
	if domain == "" {
		return err("failed to extract domain from URL")
}

	return ok(domain)
}

// HandleSearchText searches for a keyword in a given text and returns the count of occurrences.
func HandleSearchText(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	text, _ :=getString(args, "text")
	keyword, _ :=getString(args, "keyword")
	if text == "" || keyword == "" {
		return err("missing 'text' or 'keyword' parameter")
}

	count := strings.Count(text, keyword)
	return ok(fmt.Sprintf("keyword '%s' found %d times", keyword, count))
}

// HandleValidateURL validates if a given string is a properly formatted URL.
func HandleValidateURL(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing 'url' parameter")
}

	_, parseErr := url.ParseRequestURI(urlStr)
	if parseErr != nil {
		return ok("false")
}

	return ok("true")
}