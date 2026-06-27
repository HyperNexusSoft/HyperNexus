package tools". No explanation.")

We should also check for any other issues: The code uses getString, getInt, getBool, ok, e, TextContent, etc. Those are pre-defined in parity.go, so fine.

Now, let's produce the fixed code:

- Remove "fmt" import (unused)
- Remove "os/exec" import (unused)
- Remove any stray characters (the backslash and comma are already not present in the provided code, but just in case, ensure the import block is clean)
- Fix the typo "erno" in HandleXHSBatchDownload: delete that word.

Also, check for "io" imports: used in io.ReadAll. "net/http" used. "net/url" used. "os" used in os.UserHomeDir, os.Stat. "path/filepath" used. "strconv" used. "strings" used. "time" used.

So only "fmt" and "os/exec" are unused.

Now, produce the final code.

But note: The original code has "fmt" imported but not used. In the provided code, it's there. The error message didn't mention fmt, but it's a compilation error. We'll remove it.

Also, ensure no other typos.

Let's write the fixed code.package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const xhsBaseURL = "http://127.0.0.1:5556"

var xhsClient = http.DefaultClient

type xhsDetailRequest struct {
	URL      string `json:"url"`
	Download bool   `json:"download,omitempty"`
	Index    []int  `json:"index,omitempty"`
}

type xhsDetailResponse struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Author      string   `json:"author,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Images      []string `json:"images,omitempty"`
	Videos      []string `json:"videos,omitempty"`
	Downloaded  bool     `json:"downloaded,omitempty"`
	Error       string   `json:"error,omitempty"`
}

func xhsAPICall(ctx context.Context, endpoint string, payload interface{}) (map[string]interface{}, error) {
	var body io.Reader
	if payload != nil {
		jsonData, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, marshalErr
		}
		body = bytes.NewBuffer(jsonData)

	req, reqErr := http.NewRequestWithContext(ctx, "POST", xhsBaseURL+endpoint, body)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := xhsClient.Do(req)
	if apiErr != nil {
		return nil, apiErr
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(respBody, &result); parseErr != nil {
		return nil, parseErr
	}
	return result, nil
}

}

// HandleXHSDownload downloads XHS content from a given URL
func HandleXHSDownload(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url parameter is required")
}

	download, _ :=getBool(args, "download")

	var index []int
	if idxRaw, found := args["index"]; found {
		switch v := idxRaw.(type) {
		case []interface{}:
			for _, item := range v {
				switch n := item.(type) {
				case float64:
					index = append(index, int(n))
				case int:
					index = append(index, n)

			}
		}
	}

	reqPayload := xhsDetailRequest{
		URL:      urlStr,
		Download: download,
		Index:    index,
	}

	result, apiErr := xhsAPICall(ctx, "/xhs/detail", reqPayload)
	if apiErr != nil {
		return err("XHS API call failed: " + apiErr.Error())
}

	if errMsg, found := result["error"].(string); ok && errMsg != "" {
		return err("XHS API error: " + errMsg)
}

	resultJSON, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal response: " + marshalErr.Error())
}

	return ok(string(resultJSON))
}

}

// HandleXHSExtract extracts XWWW information without downloading
func HandleXHSExtract(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url parameter is required")
}

	reqPayload := xhsDetailRequest{
		URL:      urlStr,
		Download: false,
	}

	result, apiErr := xhsAPICall(ctx, "/xhs/detail", reqPayload)
	if apiErr != nil {
		return err("XHS API call failed: " + apiErr.Error())
}

	if errMsg, found := result["error"].(string); ok && errMsg != "" {
		return err("XHS API error: " + errMsg)
}

	resultJSON, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal response: " + marshalErr.Error())
}

	return ok(string(resultJSON))
}

// HandleXHSBatchDownload downloads multiple XHS URLs
func HandleXHSBatchDownload(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlsRaw, found := args["urls"]
	if !found {
		return err("urls parameter is required")
}

	var urls []string
	switch v := urlsRaw.(type) {
	case []interface{}:
		for _, item := range v {
			if s, found := item.(string); ok && s != "" {
				urls = append(urls, s)

		}
	case string:
		for _, s := range strings.Split(v, " ") {
			s = strings.TrimSpace(s)
			if s != "" {
				urls = append(urls, s)

		}
	}

	if len(urls) == 0 {
		return err("no valid URLs provided")
}

	download, _ :=getBool(args, "download")

	results := make([]map[string]interface{}, 0, len(urls))
	var failed []string

	for _, urlStr := range urls {
		reqPayload := xhsDetailRequest{
			URL:      urlStr,
			Download: download,
		}

		result, apiErr := xhsAPICall(ctx, "/xhs/detail", reqPayload)
		if apiErr != nil {
			failed = append(failed, urlStr+": "+apiErr.Error())
			continue
		}

		results = append(results, result)

	output := map[string]interface{}{
		"results": results,
		"count":   len(results),
	}

	if len(failed) > 0 {
		output["failed"] = failed
	}

	resultJSON, marshalErr := json.MarshalIndent(output, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal response: " + marshalErr.Error())
}

	return ok(string(resultJSON))
}

}
}
}

// HandleXHSGetDownloadPath returns the default download path for XHS-Downloader
func HandleXHSGetDownloadPath(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return err("failed to get user home directory: " + homeErr.Error())
}

	// Check common paths
	paths := []string{
		filepath.Join(homeDir, "XHS-Downloader", "Download"),
		filepath.Join(homeDir, "Downloads", "XHS-Downloader"),
		filepath.Join(homeDir, ".local", "share", "XHS-Downloader", "Download"),
	}

	// Try to find existing path
	for _, p := range paths {
		if _, statErr := os.Stat(p); statErr == nil {
			return ok(p)

	}

	// Return default if none exists
	return ok(paths[0])
}

}

// HandleXHSSearch searches XHS content (uses local API if available)
func HandleXHSSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	limit, _ :=getInt(args, "limit")
	if limit <= 0 {
		limit = 10
	}

	// Build search URL for XHS web
	searchURL := "https://www.xiaohongshu.com/search_result?keyword=" + url.QueryEscape(query)

	// Return search URL and instructions
	result := map[string]interface{}{
		"search_url": searchURL,
		"query":      query,
		"limit":      limit,
		"note":       "XHS-Downloader API does not expose search directly. Use the search URL to find content, then extract individual URLs with xhs_extract or xhs_download.",
	}

	resultJSON, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal response: " + marshalErr.Error())
}

	return ok(string(resultJSON))
}

// HandleXHSCheckServer checks if XHS-Downloader server is running
func HandleXHSCheckServer(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	req, reqErr := http.NewRequestWithContext(ctx, "GET", xhsBaseURL+"/docs", nil)
	if reqErr != nil {
		return err("failed to create request: " + reqErr.Error())
}

	resp, apiErr := xhsClient.Do(req)
	if apiErr != nil {
		// Try to start server if not running
		return ok("XHS-Downloader server is not running at " + xhsBaseURL + ". Please start it with: python main.py api")
}

	defer resp.Body.Close()

	status := "running"
	if resp.StatusCode != 200 {
		status = "unavailable (status: " + strconv.Itoa(resp.StatusCode) + ")"
	}

	result := map[string]interface{}{
		"server_url": xhsBaseURL,
		"status":     status,
		"docs_url":   xhsBaseURL + "/docs",
	}

	resultJSON, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal response: " + marshalErr.Error())
}

	return ok(string(resultJSON))
}