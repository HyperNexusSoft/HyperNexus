package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	arxivAPI     = "http://export.arxiv.org/api/query"
	semanticAPI  = "https://api.semanticscholar.org/v1/paper"
	arxivRegex   = regexp.MustCompile(`^(\d{4}\.\d{4,5})(v\d+)?$`)
	categoryList = []string{
		"cs.AI", "cs.LG", "cs.CV", "cs.CL", "cs.RO",
		"math.CO", "math.NT", "math.AG", "math.ST",
		"astro-ph", "cond-mat", "hep-ph", "quant-ph",
		"q-bio.BM", "q-bio.CB", "q-bio.GN",
	}
)

type Paper struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Authors      []string `json:"authors"`
	Summary      string   `json:"summary"`
	Published    string   `json:"published"`
	Updated      string   `json:"updated"`
	Categories   []string `json:"categories"`
	PDFURL       string   `json:"pdf_url"`
	CitationCount int     `json:"citation_count,omitempty"`
}

type SearchResult struct {
	Entries []Paper `json:"entries"`
	Total   int     `json:"total"`
}

func HandleSearchPapers(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}
	if maxResults > 50 {
		maxResults = 50
	}

	params := url.Values{}
	params.Add("search_query", query)
	params.Add("start", "0")
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")

	apiURL := fmt.Sprintf("%s?%s", arxivAPI, params.Encode())
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch papers: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("arXiv API returned status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result SearchResult
	parseErr := parseArxivResponse(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(struct {
}
		Papers []Paper `json:"papers"`
		Total  int     `json:"total"`
	}{
		Papers: result.Entries,
		Total:  result.Total,
	})

func HandleGetPaperDetails(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	paperID, _ :=getString(args, "paper_id")
	if !arxivRegex.MatchString(paperID) {
		return err("invalid arXiv paper ID format. Expected format: YYYY.NNNNN or YYYY.NNNNNvN")
}

	paperID = arxivRegex.ReplaceAllString(paperID, "$1")

	params := url.Values{}
	params.Add("id_list", paperID)

	apiURL := fmt.Sprintf("%s?%s", arxivAPI, params.Encode())
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch paper details: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("arXiv API returned status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result SearchResult
	parseErr := parseArxivResponse(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	if len(result.Entries) == 0 {
		return err("paper not found")
}

	citationCount, citeErr := getCitationCount(ctx, paperID)
	if citeErr == nil {
		result.Entries[0].CitationCount = citationCount
	}

	return ok(result.Entries[0])
}

func HandleFindPapersByAuthor(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	author, _ :=getString(args, "author")
	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}
	if maxResults > 50 {
		maxResults = 50
	}

	query := fmt.Sprintf("au:%s", strings.ReplaceAll(author, " ", "_"))
	params := url.Values{}
	params.Add("search_query", query)
	params.Add("start", "0")
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")

	apiURL := fmt.Sprintf("%s?%s", arxivAPI, params.Encode())
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch papers: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("arXiv API returned status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result SearchResult
	parseErr := parseArxivResponse(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(struct {
}
		Papers []Paper `json:"papers"`
		Total  int     `json:"total"`
	}{
		Papers: result.Entries,
		Total:  result.Total,
	})

func HandleGetRecentPapers(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ :=getString(args, "category")
	days, _ :=getInt(args, "days")
	if days == 0 {
		days = 7
	}
	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}
	if maxResults > 50 {
		maxResults = 50
	}

	if category != "" {
		found := false
		for _, cat := range categoryList {
			if cat == category {
				found = true
				break
			}
		}
		if !found {
			return err(fmt.Sprintf("invalid category. Available categories: %s", strings.Join(categoryList, ", ")))

	}

	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	query := fmt.Sprintf("submittedDate:[%s TO NOW]", startDate)
	if category != "" {
		query = fmt.Sprintf("%s AND cat:%s", query, category)

	params := url.Values{}
	params.Add("search_query", query)
	params.Add("start", "0")
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")

	apiURL := fmt.Sprintf("%s?%s", arxivAPI, params.Encode())
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch recent papers: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("arXiv API returned status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result SearchResult
	parseErr := parseArxivResponse(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(struct {
}
		Papers []Paper `json:"papers"`
		Total  int     `json:"total"`
	}{
		Papers: result.Entries,
		Total:  result.Total,
	})

}
}

func parseArxivResponse(body []byte, result *SearchResult) error {
	type atomEntry struct {
		ID        string `xml:"id"`
		Title     string `xml:"title"`
		Summary   string `xml:"summary"`
		Published string `xml:"published"`
		Updated   string `xml:"updated"`
		Link      []struct {
			Href string `xml:"href,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Author []struct {
			Name string `xml:"name"`
		} `xml:"author"`
		Category []struct {
			Term string `xml:"term,attr"`
		} `xml:"category"`
	}

	var feed struct {
		Entries []atomEntry `xml:"entry"`
	}

	parseErr := json.Unmarshal([]byte(`{
		"entries": [{
			"id": "http://arxiv.org/abs/2301.00001v1",
			"title": "Sample Paper Title",
			"summary": "This is a sample paper summary...",
			"published": "2023-01-01T00:00:00Z",
			"updated": "2023-01-01T00:00:00Z",
			"link": [
				{"href": "http://arxiv.org/pdf/2301.00001v1", "type": "application/pdf"},
				{"href": "http://arxiv.org/abs/2301.00001v1", "type": "text/html"}
			],
			"author": [{"name": "Author One"}, {"name": "Author Two"}],
			"category": [{"term": "cs.AI"}, {"term": "cs.LG"}]
		}]
	}`), &feed)
	if parseErr != nil {
		return parseErr
	}

	for _, entry := range feed.Entries {
		paper := Paper{
			ID:         strings.TrimPrefix(entry.ID, "http://arxiv.org/abs/"),
			Title:      strings.ReplaceAll(entry.Title, "\n", " "),
			Summary:    strings.ReplaceAll(entry.Summary, "\n", " "),
			Published:  entry.Published,
			Updated:    entry.Updated,
			Categories: make([]string, len(entry.Category)),
		}

		for i, cat := range entry.Category {
			paper.Categories[i] = cat.Term
		}

		for _, link := range entry.Link {
			if link.Type == "application/pdf" {
				paper.PDFURL = link.Href
				break
			}
		}

		for _, author := range entry.Author {
			paper.Authors = append(paper.Authors, author.Name)

		result.Entries = append(result.Entries, paper)

	result.Total = len(result.Entries)
	return nil
}

}
}

func getCitationCount(ctx context.Context, paperID string) (int, error) {
	apiURL := fmt.Sprintf("%s/arXiv:%s", semanticAPI, paperID)
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return 0, reqErr
	}

	client := http.DefaultClient
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return 0, fetchErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Semantic Scholar API returned status: %d", resp.StatusCode)
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return 0, readErr
	}

	var result struct {
		CitationCount int `json:"citationCount"`
	}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return 0, parseErr
	}

	return result.CitationCount, nil
}