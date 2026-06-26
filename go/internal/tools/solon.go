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

func ok(content string) ToolResponse {
	return ToolResponse{Content: content, Error: nil}
}

func err(e error) ToolResponse {
	return ToolResponse{Content: "", Error: e}
}

func getString(args map[string]interface{}, key string) string {
	value, found := args[key]
	if !found {
		return ""
	}
	str, found := value.(string)
	if !found {
		return ""
	}
	return str
}

func HandleIssueTemplate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	templateType, _ :=getString(args, "type")
	if templateType == "" {
		return err("template type is required")
}

	client := http.DefaultClient

	var templateURL string
	switch templateType {
	case "bug_report":
		templateURL = "https://raw.githubusercontent.com/opensolon/solon/main/.github/ISSUE_TEMPLATE/bug_report.md"
	case "feature_request":
		templateURL = "https://raw.githubusercontent.com/opensolon/solon/main/.github/ISSUE_TEMPLATE/feature_request.md"
	case "problem_support":
		templateURL = "https://raw.githubusercontent.com/opensolon/solon/main/.github/ISSUE_TEMPLATE/problem_support.md"
	default:
		return err("unsupported template type")
}

	resp, fetchErr := client.Get(templateURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch template: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch template: status %d", resp.StatusCode))
}

	var content string
	if e := json.NewDecoder(resp.Body).Decode(&content); e != nil {
		return err(fmt.Sprintf("failed to decode template: %v", e))
}

	return ok(content)
}

func HandleContributingGuide(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient

	contributingURL := "https://raw.githubusercontent.com/opensolon/solon/main/CONTRIBUTING.md"
	resp, fetchErr := client.Get(contributingURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch contributing guide: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch contributing guide: status %d", resp.StatusCode))
}

	var content string
	if e := json.NewDecoder(resp.Body).Decode(&content); e != nil {
		return err(fmt.Sprintf("failed to decode contributing guide: %v", e))
}

	return ok(content)
}

func HandleRepositoryInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	repoName, _ :=getString(args, "repo")
	if repoName == "" {
		return err("repository name is required")
}

	client := http.DefaultClient

	apiURL := fmt.Sprintf("https://api.github.com/repos/opensolon/%s", repoName)
	resp, fetchErr := client.Get(apiURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch repository info: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch repository info: status %d", resp.StatusCode))
}

	var repoInfo map[string]interface{}
	if e := json.NewDecoder(resp.Body).Decode(&repoInfo); e != nil {
		return err(fmt.Sprintf("failed to decode repository info: %v", e))
}

	description, found := repoInfo["description"].(string)
	if !found {
		description = ""
	}
	htmlURL, found := repoInfo["html_url"].(string)
	if !found {
		htmlURL = ""
	}
	stargazersCount, found := repoInfo["stargazers_count"].(int)
	if !found {
		stargazersCount = 0
	}

	info := fmt.Sprintf("Repository: %s\nDescription: %s\nStars: %d\nURL: %s",
		repoName, description, stargazersCount, htmlURL)

	return ok(info)
}

func HandleIssueSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("search query is required")
}

	client := http.DefaultClient

	searchURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s+repo:opensolon/solon", url.QueryEscape(query))
	resp, fetchErr := client.Get(searchURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to search issues: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to search issues: status %d", resp.StatusCode))
}

	var searchResults map[string]interface{}
	if e := json.NewDecoder(resp.Body).Decode(&searchResults); e != nil {
		return err(fmt.Sprintf("failed to decode search results: %v", e))
}

	issues, found := searchResults["items"].([]interface{})
	if !found {
		return err("failed to decode search results")
}

	if len(issues) == 0 {
		return ok("No issues found matching your query")
}

	var result strings.Builder
	result.WriteString("Search results for '" + query + "':\n\n")

	for i, issue := range issues {
		issueMap, found := issue.(map[string]interface{})
		if !found {
			continue
		}
		title, found := issueMap["title"].(string)
		if !found {
			continue
		}
		number, found := issueMap["number"].(int)
		if !found {
			continue
		}
		htmlURL, found := issueMap["html_url"].(string)
		if !found {
			continue
		}

		result.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, title, htmlURL))

	return ok(result.String())
}

}

func HandleLicenseInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient

	licenseURL := "https://raw.githubusercontent.com/opensolon/solon/main/LICENSE"
	resp, fetchErr := client.Get(licenseURL)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch license: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch license: status %d", resp.StatusCode))
}

	var content string
	if e := json.NewDecoder(resp.Body).Decode(&content); e != nil {
		return err(fmt.Sprintf("failed to decode license: %v", e))
}

	return ok(content)
}