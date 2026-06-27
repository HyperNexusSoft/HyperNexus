package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ToolResponse, ok, e, getString, getInt, getBool, TextContent は parity.go で定義されていると仮定

// HandleSearchXiaohongshu handles search queries for Xiaohongshu.
func HandleSearchXiaohongshu(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return ToolResponse{Error: "missing required parameter: query"}, nil
	}

	apiURL := "https://api.xiaohongshu.com/search"
	params := url.Values{}
	params.Set("keyword", query)

	req, reqErr := http.NewRequest("GET", fmt.Sprintf("%s?%s", apiURL, params.Encode()), nil)
	if reqErr != nil {
		return ToolResponse{}, fmt.Errorf("request creation error: %w", reqErr)
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return ToolResponse{}, fmt.Errorf("request execution error: %w", fetchErr)
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ToolResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

	var result map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		return ToolResponse{}, fmt.Errorf("response parsing error: %w", parseErr)
}

	responseText, found := result["data"].(string)
	if !found {
		return ToolResponse{}, fmt.Errorf("invalid response format")
}

	return ToolResponse{Text: responseText}, nil
}

// HandleGetPostDetails retrieves details of a specific post from Xiaohongshu.
func HandleGetPostDetails(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	postID, _ :=getString(args, "post_id")
	if postID == "" {
		return ToolResponse{Error: "missing required parameter: post_id"}, nil
	}

	apiURL := fmt.Sprintf("https://api.xiaohongshu.com/posts/%s", postID)
	req, reqErr := http.NewRequest("GET", apiURL, nil)
	if reqErr != nil {
		return ToolResponse{}, fmt.Errorf("request creation error: %w", reqErr)
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return ToolResponse{}, fmt.Errorf("request execution error: %w", fetchErr)
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ToolResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

	var result map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		return ToolResponse{}, fmt.Errorf("response parsing error: %w", parseErr)
}

	responseText, found := result["data"].(string)
	if !found {
		return ToolResponse{}, fmt.Errorf("invalid response format")
}

	return ToolResponse{Text: responseText}, nil
}

// HandleGetUserProfile retrieves the profile of a specific user from Xiaohongshu.
func HandleGetUserProfile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	userID, _ :=getString(args, "user_id")
	if userID == "" {
		return ToolResponse{Error: "missing required parameter: user_id"}, nil
	}

	apiURL := fmt.Sprintf("https://api.xiaohongshu.com/users/%s", userID)
	req, reqErr := http.NewRequest("GET", apiURL, nil)
	if reqErr != nil {
		return ToolResponse{}, fmt.Errorf("request creation error: %w", reqErr)
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return ToolResponse{}, fmt.Errorf("request execution error: %w", fetchErr)
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ToolResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

	var result map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		return ToolResponse{}, fmt.Errorf("response parsing error: %w", parseErr)
}

	responseText, found := result["data"].(string)
	if !found {
		return ToolResponse{}, fmt.Errorf("invalid response format")
}

	return ToolResponse{Text: responseText}, nil
}