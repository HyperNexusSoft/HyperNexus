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
	"strings"
	"time"
)

// shared HTTP client
var http.DefaultClient = http.DefaultClient

// HandlePing returns a simple acknowledgement.
func HandlePing(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("pong")
}

// HandleRepoInfo fetches basic repository information from GitHub.
func HandleRepoInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	owner, _ :=getString(args, "owner")
	repo, _ :=getString(args, "repo")
	if owner == "" || repo == "" {
		return err("owner and repo must be provided")
}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", url.PathEscape(owner), url.PathEscape(repo))

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("GitHub API returned status %d", resp.StatusCode))
}

	var data map[string]interface{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&data)
	if decodeErr != nil {
		return err(decodeErr.Error())
}

	pretty, marshalErr := json.MarshalIndent(data, "", "  ")
	if marshalErr != nil {
		return err(marshalErr.Error())
}

	return ok(string(pretty))
}

// HandleSearchIssues searches issues in a repository using GitHub's search API.
func HandleSearchIssues(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	owner, _ :=getString(args, "owner")
	repo, _ :=getString(args, "repo")
	query, _ :=getString(args, "query")
	if owner == "" || repo == "" || query == "" {
		return err("owner, repo and query must be provided")
}

	// Build query: repo:owner/repo query
	searchQuery := fmt.Sprintf("repo:%s/%s %s", owner, repo, query)
	values := url.Values{}
	values.Set("q", searchQuery)
	apiURL := "https://api.github.com/search/issues?" + values.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("GitHub search API returned status %d", resp.StatusCode))
}

	var result struct {
		Items []struct {
			Number  int    `json:"number"`
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
		} `json:"items"`
	}
	decodeErr := json.NewDecoder(resp.Body).Decode(&result)
	if decodeErr != nil {
		return err(decodeErr.Error())
}

	if len(result.Items) == 0 {
		return ok("no matching issues found")
}

	var sb strings.Builder
	for _, item := range result.Items {
		fmt.Fprintf(&sb, "#%d: %s (%s)\n", item.Number, item.Title, item.HTMLURL)

	return ok(sb.String())
}

}

// HandleListReleases lists releases for a repository.
func HandleListReleases(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	owner, _ :=getString(args, "owner")
	repo, _ :=getString(args, "repo")
	if owner == "" || repo == "" {
		return err("owner and repo must be provided")
}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", url.PathEscape(owner), url.PathEscape(repo))

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("GitHub releases API returned status %d", resp.StatusCode))
}

	var releases []struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}
	decodeErr := json.NewDecoder(resp.Body).Decode(&releases)
	if decodeErr != nil {
		return err(decodeErr.Error())
}

	if len(releases) == 0 {
		return ok("no releases found")
}

	var sb strings.Builder
	for _, rel := range releases {
		fmt.Fprintf(&sb, "%s - %s (%s)\n", rel.TagName, rel.Name, rel.HTMLURL)

	return ok(sb.String())
}

}

// HandleFileContent fetches raw file content from a repository.
func HandleFileContent(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	owner, _ :=getString(args, "owner")
	repo, _ :=getString(args, "repo")
	path, _ :=getString(args, "path")
	ref, _ :=getString(args, "ref") // branch, tag or commit SHA
	if owner == "" || repo == "" || path == "" {
		return err("owner, repo and path must be provided")
}

	if ref == "" {
		ref = "main"
	}
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(ref), url.PathEscape(path))

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch file, status %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(string(body))
}

// HandleLocalFile reads a file from the local filesystem (sandboxed to a base directory).
func HandleLocalFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	baseDir := "/tmp/osaurus" // sandbox root
	relPath, _ :=getString(args, "path")
	if relPath == "" {
		return err("path must be provided")
}

	cleanPath := filepath.Clean(relPath)
	absPath := filepath.Join(baseDir, cleanPath)

	// Prevent path traversal outside the sandbox
	if !strings.HasPrefix(absPath, baseDir) {
		return err("access to the requested path is denied")
}

	file, openErr := os.Open(absPath)
	if openErr != nil {
		return err(openErr.Error())
}

	defer file.Close()
	content, readErr := io.ReadAll(file)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(string(content))
}