package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ok(v string) ToolResponse {
	return ToolResponse{Ok: v}
}

func err(v string) ToolResponse {
	return ToolResponse{Error: v}
}

func getString(args map[string]interface{}, key string) string {
	if value, found := args[key]; found {
		if str, found := value.(string); found {
			return str
		}
	}
	return ""
}

func HandleNavigate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	parsedURL, parseErr := url.Parse(urlStr)
	if parseErr != nil {
		return err(fmt.Sprintf("invalid url: %s", parseErr.Error()))
}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %s", reqErr.Error()))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch url: %s", fetchErr.Error()))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %s", readErr.Error()))
}

	contentType := resp.Header.Get("Content-Type")
	result := fmt.Sprintf("Status: %d\nContent-Type: %s\nBody Length: %d bytes", resp.StatusCode, contentType, len(body))

	if strings.Contains(contentType, "application/json") {
		var prettyJSON map[string]interface{}
		if jsonErr := json.Unmarshal(body, &prettyJSON); jsonErr == nil {
			prettyBytes, marshalErr := json.MarshalIndent(prettyJSON, "", "  ")
			if marshalErr == nil {
				result += fmt.Sprintf("\n\nJSON Content:\n%s", string(prettyBytes))

		}
	}

	return ok(result)
}

}

func HandleSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("missing required argument: query")
}

	baseURL := "https://www.example.com/search"
	params := url.Values{}
	params.Add("q", query)
	searchURL := baseURL + "?" + params.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create search request: %s", reqErr.Error()))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to perform search: %s", fetchErr.Error()))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read search results: %s", readErr.Error()))
}

	return ok(fmt.Sprintf("Search performed for: %s\nStatus: %d\nResults Length: %d bytes", query, resp.StatusCode, len(body)))
}

func HandleExtract(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create extraction request: %s", reqErr.Error()))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch page for extraction: %s", fetchErr.Error()))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read page content: %s", readErr.Error()))
}

	textContent := string(body)

	if len(textContent) > 2000 {
		textContent = textContent[:2000] + "... (truncated)"
	}

	return ok(fmt.Sprintf("Extracted content from %s:\n\n%s", urlStr, textContent))
}

func HandleStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("missing required argument: url")
}

	req, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create status request: %s", reqErr.Error()))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to check status: %s", fetchErr.Error()))
}

	defer resp.Body.Close()

	statusMsg := fmt.Sprintf("URL: %s\nStatus Code: %d\nStatus: %s", urlStr, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusMsg += "\nResult: OK"
	} else {
		statusMsg += "\nResult: Error"
	}

	return ok(statusMsg)
}