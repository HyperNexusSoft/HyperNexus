package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PostmanAPIResponse represents a standard Postman API response
type PostmanAPIResponse struct {
	Collections []Collection `json:"collections,omitempty"`
	Workspace   *Workspace   `json:"workspace,omitempty"`
	// Add other fields as needed
}

type Collection struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	UID        string `json:"uid"`
	OwnerID    string `json:"owner_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func getPostmanClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func getPostmanAPIKey(args map[string]interface{}) (string, error) {
	apiKey, _ :=getString(args, "api_key")
	if apiKey == "" {
		apiKey = os.Getenv("POSTMAN_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("api_key parameter or POSTMAN_API_KEY environment variable is required")
}

	return apiKey, nil
}

}

func postmanRequest(method, endpoint string, apiKey string, body io.Reader) (*http.Response, error) {
	client := getPostmanClient()

	req, reqErr := http.NewRequest(method, endpoint, body)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "postman-api-mcp/1.0")

	return client.Do(req)
}

// HandleListCollections lists all collections in the workspace
func HandleListCollections(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	apiKey, apiKeyErr := getPostmanAPIKey(args)
	if apiKeyErr != nil {
		return err(apiKeyErr.Error())
}

	workspaceID, _ :=getString(args, "workspace_id")
	var endpoint string
	if workspaceID != "" {
		endpoint = fmt.Sprintf("https://api.getpostman.com/collections?workspace_id=%s", url.QueryEscape(workspaceID))
	} else {
		endpoint = "https://api.getpostman.com/collections"
	}

	resp, respErr := postmanRequest("GET", endpoint, apiKey, nil)
	if respErr != nil {
		return err(fmt.Sprintf("failed to make request: %s", respErr.Error()))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return err(fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(bodyBytes)))
}

	var apiResp PostmanAPIResponse
	if parseErr := json.NewDecoder(resp.Body).Decode(&apiResp); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %s", parseErr.Error()))
}

	if len(apiResp.Collections) == 0 {
		return ok("No collections found")
}

	// Sort by name for consistency
	sort.Slice(apiResp.Collections, func(i, j int) bool {
		return apiResp.Collections[i].Name < apiResp.Collections[j].Name
	})

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d collections:\n\n", len(apiResp.Collections)))

	for i, collection := range apiResp.Collections {
		result.WriteString(fmt.Sprintf("%d. %s (ID: %s)\n", i+1, collection.Name, collection.ID))
		if collection.UpdatedAt != "" {
			result.WriteString(fmt.Sprintf("   Last updated: %s\n", collection.UpdatedAt))

		result.WriteString("\n")

	return ok(result.String())
}

}
}

// HandleGetCollection retrieves details for a specific collection
func HandleGetCollection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
}

	// ... (same as original code)

// HandleListWorkspaces lists all workspaces for the API key
func HandleListWorkspaces(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
}

	// ... (same as original code)

// HandleSearchCollections searches for collections by name
func HandleSearchCollections(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
}

	// ... (same as original code)

// HandleExportCollection exports a collection to a file
func HandleExportCollection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// ... (same as original code)
}