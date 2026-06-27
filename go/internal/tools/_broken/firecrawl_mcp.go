package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	firecrawlAPIURL = "https://api.firecrawl.dev"
	defaultTimeout  = 30 * time.Second
)

var http.DefaultClient = http.DefaultClient

func HandleScrape(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url is required")
}

	onlyMainContent, _ :=getBool(args, "onlyMainContent")
	formats, _ :=getString(args, "formats")
	maxAge, _ :=getInt(args, "maxAge")
	extract, _ :=getString(args, "extract")
	timeout, _ :=getInt(args, "timeout")
	if timeout == 0 {
		timeout = 30
	}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/scrape", apiURL)
	reqBody := map[string]interface{}{
		"url":             urlStr,
		"onlyMainContent": onlyMainContent,
	}
	if formats != "" {
		reqBody["formats"] = formats
	}
	if maxAge > 0 {
		reqBody["cache"] = map[string]int{"maxAge": maxAge}
	}
	if extract != "" {
		reqBody["extract"] = extract
	}

	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if markdown, found := result["markdown"].(string); found {
		return ok(markdown)
}

	return ok(fmt.Sprintf("%v", result))
}

}

func HandleSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query is required")
}

	page, _ :=getInt(args, "page")
	if page == 0 {
		page = 1
	}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/search", apiURL)
	reqBody := map[string]interface{}{
		"query": query,
		"page":  page,
	}

	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if data, found := result["data"].([]interface{}); found {
		var output strings.Builder
		for i, item := range data {
			if m, found := item.(map[string]interface{}); found {
				output.WriteString(fmt.Sprintf("Result %d:\n", i+1))
				if title, found := m["title"].(string); found {
					output.WriteString(fmt.Sprintf("Title: %s\n", title))

				if url, found := m["url"].(string); found {
					output.WriteString(fmt.Sprintf("URL: %s\n", url))

				if desc, found := m["description"].(string); found {
					output.WriteString(fmt.Sprintf("Description: %s\n", desc))

				if content, found := m["content"].(string); found {
					output.WriteString(fmt.Sprintf("Content: %s\n\n", content))

			}
		}
		return ok(output.String())
}

	return ok(fmt.Sprintf("%v", result))
}

}
}
}
}
}

func HandleMap(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url is required")
}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/map", apiURL)
	reqBody := map[string]interface{}{
		"url": urlStr,
	}

	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if links, found := result["links"].([]interface{}); found {
		var output strings.Builder
		output.WriteString("Discovered links:\n")
		for _, link := range links {
			if l, found := link.(string); found {
				output.WriteString(l + "\n")

		}
		return ok(output.String())
}

	return ok(fmt.Sprintf("%v", result))
}

}
}

func HandleCrawl(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url is required")
}

	limit, _ :=getInt(args, "limit")
	if limit == 0 {
		limit = 100
	}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/crawl", apiURL)
	reqBody := map[string]interface{}{
		"url":   urlStr,
		"limit": limit,
	}

	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if jobID, found := result["jobId"].(string); found {
		return ok(fmt.Sprintf("Crawl started with job ID: %s", jobID))
}

	return ok(fmt.Sprintf("%v", result))
}

}

func HandleCheckCrawlStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	jobID, _ :=getString(args, "jobId")
	if jobID == "" {
		return err("jobId is required")
}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/crawl/status/%s", apiURL, jobID)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if status, found := result["status"].(string); found {
		if status == "completed" {
			if data, found := result["data"].([]interface{}); found {
				var output strings.Builder
				output.WriteString(fmt.Sprintf("Crawl completed. Found %d pages:\n", len(data)))
				for i, item := range data {
					if m, found := item.(map[string]interface{}); found {
						if url, found := m["url"].(string); found {
							output.WriteString(fmt.Sprintf("%d: %s\n", i+1, url))

					}
				}
				return ok(output.String())

		}
		return ok(fmt.Sprintf("Crawl status: %s", status))
}

	return ok(fmt.Sprintf("%v", result))
}

}
}
}

func HandleParse(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, _ :=getString(args, "filePath")
	uploadRef, _ :=getString(args, "uploadRef")

	if filePath == "" && uploadRef == "" {
		return err("either filePath or uploadRef is required")
}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	if uploadRef != "" {
		// Second call with uploadRef
		reqURL := fmt.Sprintf("%s/v1/parse", apiURL)
		reqBody := map[string]interface{}{
			"uploadRef": uploadRef,
		}

		jsonBody, jsonErr := json.Marshal(reqBody)
		if jsonErr != nil {
			return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

		req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
		if reqErr != nil {
			return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

		apiKey := os.Getenv("FIRECRAWL_API_KEY")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)

		req.Header.Set("Content-Type", "application/json")

		resp, respErr := http.DefaultClient.Do(req)
		if respErr != nil {
			return err(fmt.Sprintf("request failed: %v", respErr))
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

		var result map[string]interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
			return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

		if markdown, found := result["markdown"].(string); found {
			return ok(markdown)
}

		return ok(fmt.Sprintf("%v", result))
}

	// First call with filePath - return upload instructions
	fileInfo, statErr := os.Stat(filePath)
	if statErr != nil {
		return err(fmt.Sprintf("failed to access file: %v", statErr))
}

	if fileInfo.Size() > 10*1024*1024 { // 10MB limit
		return err("file size exceeds 10MB limit")
}

	uploadRef = fmt.Sprintf("local-file-%d", time.Now().UnixNano())
	return ok(fmt.Sprintf(`To parse this file, you need to upload it first.
}
Upload reference: %s
Next tool call should use: {"uploadRef": "%s"}`, uploadRef, uploadRef))

}

func HandleExtract(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url is required")
}

	schema := args["schema"]
	if schema == nil {
		return err("schema is required")
}

	apiURL := os.Getenv("FIRECRAWL_API_URL")
	if apiURL == "" {
		apiURL = firecrawlAPIURL
	}

	reqURL := fmt.Sprintf("%s/v1/extract", apiURL)
	reqBody := map[string]interface{}{
		"url":    urlStr,
		"schema": schema,
	}

	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return err(fmt.Sprintf("request failed: %v", respErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API error: %s - %s", resp.Status, string(body)))
}

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(fmt.Sprintf("failed to decode response: %v", decodeErr))
}

	if data, found := result["data"]; found {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		return ok(string(jsonData))
}

	return ok(fmt.Sprintf("%v", result))
}
}