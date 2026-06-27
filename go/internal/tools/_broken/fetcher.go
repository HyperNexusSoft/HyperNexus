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
	"time"
)

// HandleFetchUrl fetches raw text content from a specified URL via HTTP GET
func HandleFetchUrl(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	targetURL, _ :=getString(args, "url")
	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
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

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("Request failed with status %d: %s", resp.StatusCode, string(body)))
}

	return ok(string(body))
}

// HandleFetchWithParams fetches content from a URL with custom query parameters
func HandleFetchWithParams(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	baseURL, _ :=getString(args, "url")
	rawParams, paramsExist := args["params"]
	if !paramsExist {
		return err("Missing required 'params' argument: expected map of string key-value pairs")
}

	paramsMap, typeOk := rawParams.(map[string]interface{})
	if !typeOk {
		return err("'params' must be a map of string key-value pairs")
}

	queryVals := url.Values{}
	for key, val := range paramsMap {
		strVal, isStr := val.(string)
		if !isStr {
			strVal = fmt.Sprintf("%v", val)

		queryVals.Set(key, strVal)

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryVals.Encode())
	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
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

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("Request failed with status %d: %s", resp.StatusCode, string(body)))
}

	return ok(string(body))
}

}
}

// HandleExtractLinks extracts all hyperlinks from provided HTML content
func HandleExtractLinks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	htmlContent, _ :=getString(args, "html")
	hrefRegex := regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`)")
	matches := hrefRegex.FindAllStringSubmatch(htmlContent, -1)
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			links = append(links, match[1])

	}
	sort.Strings(links)
	linksJSON, marshalErr := json.Marshal(links)
	if marshalErr != nil {
		return err(marshalErr.Error())
}

	return ok(string(linksJSON))
}
}