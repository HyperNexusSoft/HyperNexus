package tools

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ArXivEntry represents an entry in arXiv Atom feed
type ArXivEntry struct {
	ID       string `xml:"id"`
	Title    string `xml:"title"`
	Summary  string `xml:"summary"`
	Author   []struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Category  []struct {
		Term string `xml:"term,attr"`
	} `xml:"category"`
	Link []struct {
		Href  string `xml:"href,attr"`
		Type  string `xml:"type,attr"`
		Rel   string `xml:"rel,attr"`
		Title string `xml:"title,attr"`
	} `xml:"link"`
}

type ArXivFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []ArXivEntry `xml:"entry"`
}

// CrossRefWork represents a paper from CrossRef API
type CrossRefWork struct {
	Status       string `json:"status"`
	MessageType  string `json:"message-type"`
	Message      struct {
		DOI       string `json:"DOI"`
		Title     []string `json:"title"`
		Author    []struct {
			Given  string `json:"given"`
			Family string `json:"family"`
		} `json:"author"`
		Published struct {
			DateParts [][]int `json:"date-parts"`
		} `json:"published"`
		Abstract string `json:"abstract"`
		URL      string `json:"URL"`
	} `json:"message"`
}

// HandleSearchPapers searches for academic papers using arXiv API
func HandleSearchPapers(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}

	startIndex, _ :=getInt(args, "start_index")
	if startIndex < 0 {
		startIndex = 0
	}

	// Build arXiv API URL
	baseURL := "http://export.arxiv.org/api/query"
	params := url.Values{}
	params.Add("search_query", "all:"+query)
	params.Add("start", strconv.Itoa(startIndex))
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "relevance")
	params.Add("sortOrder", "descending")

	apiURL := baseURL + "?" + params.Encode()

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
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

	var feed ArXivFeed
	parseErr := xml.Unmarshal(body, &feed)
	if parseErr != nil {
		return err(parseErr.Error())
}

	var results []map[string]interface{}
	for _, entry := range feed.Entries {
		var authors []string
		for _, author := range entry.Author {
			authors = append(authors, author.Name)

		var categories []string
		for _, cat := range entry.Category {
			categories = append(categories, cat.Term)

		pdfURL := ""
		for _, link := range entry.Link {
			if link.Title == "pdf" || (link.Type == "application/pdf" && link.Rel == "alternate") {
				pdfURL = link.Href
				break
			}
		}

		result := map[string]interface{}{
			"id":          extractArXivID(entry.ID),
			"title":       strings.TrimSpace(entry.Title),
			"summary":     strings.TrimSpace(entry.Summary),
			"authors":     authors,
			"published":   entry.Published,
			"categories":  categories,
			"pdf_url":     pdfURL,
			"source":      "arxiv",
		}
		results = append(results, result)

	return ok(map[string]interface{}{
}
		"query":      query,
		"total_found": len(results),
		"papers":     results,
	})

}
}
}

// HandleGetPaper retrieves details of a specific paper by arXiv ID
func HandleGetPaper(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	paperID, _ :=getString(args, "id")
	if paperID == "" {
		return err("id parameter is required")
}

	// Clean and format arXiv ID
	paperID = strings.TrimSpace(paperID)
	paperID = strings.ReplaceAll(paperID, "arxiv:", "")
	paperID = strings.ReplaceAll(paperID, "arXiv:", "")

	baseURL := "http://export.arxiv.org/api/query"
	params := url.Values{}
	params.Add("id_list", paperID)

	apiURL := baseURL + "?" + params.Encode()

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
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

	var feed ArXivFeed
	parseErr := xml.Unmarshal(body, &feed)
	if parseErr != nil {
		return err(parseErr.Error())
}

	if len(feed.Entries) == 0 {
		return err("paper not found")
}

	entry := feed.Entries[0]

	var authors []string
	for _, author := range entry.Author {
		authors = append(authors, author.Name)

	var categories []string
	for _, cat := range entry.Category {
		categories = append(categories, cat.Term)

	pdfURL := ""
	absURL := ""
	for _, link := range entry.Link {
		if link.Title == "pdf" {
			pdfURL = link.Href
		}
		if link.Rel == "alternate" && link.Type == "text/html" {
			absURL = link.Href
		}
	}

	return ok(map[string]interface{}{
}
		"id":         extractArXivID(entry.ID),
		"title":      strings.TrimSpace(entry.Title),
		"summary":    strings.TrimSpace(entry.Summary),
		"authors":    authors,
		"published":  entry.Published,
		"updated":    entry.Updated,
		"categories": categories,
		"pdf_url":    pdfURL,
		"abstract_url": absURL,
		"source":     "arxiv",
	})

}
}

// HandleGetPaperByDOI retrieves paper details using DOI from CrossRef API
func HandleGetPaperByDOI(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	doi, _ :=getString(args, "doi")
	if doi == "" {
		return err("doi parameter is required")
}

	doi = strings.TrimSpace(doi)
	doi = strings.TrimPrefix(doi, "https://doi.org/")
	doi = strings.TrimPrefix(doi, "http://doi.org/")
	doi = strings.TrimPrefix(doi, "doi:")

	apiURL := "https://api.crossref.org/works/" + url.PathEscape(doi)

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "PaperSearchServer/1.0 (mailto:support@papersearch.example)")

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
		return err(fmt.Sprintf("CrossRef API returned status %d", resp.StatusCode))
}

	var work CrossRefWork
	parseErr := json.Unmarshal(body, &work)
	if parseErr != nil {
		return err(parseErr.Error())
}

	var authors []string
	for _, author := range work.Message.Author {
		name := strings.TrimSpace(author.Given + " " + author.Family)
		authors = append(authors, name)

	title := ""
	if len(work.Message.Title) > 0 {
		title = work.Message.Title[0]
	}

	year := ""
	month := ""
	day := ""
	if len(work.Message.Published.DateParts) > 0 && len(work.Message.Published.DateParts[0]) >= 1 {
		year = strconv.Itoa(work.Message.Published.DateParts[0][0])

	if len(work.Message.Published.DateParts) > 0 && len(work.Message.Published.DateParts[0]) >= 2 {
		month = strconv.Itoa(work.Message.Published.DateParts[0][1])

	if len(work.Message.Published.DateParts) > 0 && len(work.Message.Published.DateParts[0]) >= 3 {
		day = strconv.Itoa(work.Message.Published.DateParts[0][2])

	publishedDate := year
	if month != "" {
		publishedDate = fmt.Sprintf("%s-%s", year, month)

	if day != "" {
		publishedDate = fmt.Sprintf("%s-%s-%s", year, month, day)

	// Clean HTML from abstract
	abstract := cleanHTMLAbstract(work.Message.Abstract)

	return ok(map[string]interface{}{
}
		"doi":         doi,
		"title":       title,
		"authors":     authors,
		"published":   publishedDate,
		"abstract":    abstract,
		"url":         work.Message.URL,
		"source":      "crossref",
	})

}
}
}
}
}
}

// HandleSearchByAuthor searches for papers by author name
func HandleSearchByAuthor(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	author, _ :=getString(args, "author")
	if author == "" {
		return err("author parameter is required")
}

	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}

	startIndex, _ :=getInt(args, "start_index")
	if startIndex < 0 {
		startIndex = 0
	}

	// Build arXiv API URL for author search
	baseURL := "http://export.arxiv.org/api/query"
	params := url.Values{}
	params.Add("search_query", "au:"+author)
	params.Add("start", strconv.Itoa(startIndex))
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")

	apiURL := baseURL + "?" + params.Encode()

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
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

	var feed ArXivFeed
	parseErr := xml.Unmarshal(body, &feed)
	if parseErr != nil {
		return err(parseErr.Error())
}

	var results []map[string]interface{}
	for _, entry := range feed.Entries {
		var authors []string
		for _, a := range entry.Author {
			authors = append(authors, a.Name)

		result := map[string]interface{}{
			"id":        extractArXivID(entry.ID),
			"title":     strings.TrimSpace(entry.Title),
			"summary":   strings.TrimSpace(entry.Summary),
			"authors":   authors,
			"published": entry.Published,
			"source":    "arxiv",
		}
		results = append(results, result)

	return ok(map[string]interface{}{
}
		"author":      author,
		"total_found": len(results),
		"papers":      results,
	})

}
}

// HandleSearchByCategory searches for papers in a specific category
func HandleSearchByCategory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ :=getString(args, "category")
	if category == "" {
		return err("category parameter is required")
}

	maxResults, _ :=getInt(args, "max_results")
	if maxResults == 0 {
		maxResults = 10
	}

	startIndex, _ :=getInt(args, "start_index")
	if startIndex < 0 {
		startIndex = 0
	}

	// Build arXiv API URL for category search
	baseURL := "http://export.arxiv.org/api/query"
	params := url.Values{}
	params.Add("search_query", "cat:"+category)
	params.Add("start", strconv.Itoa(startIndex))
	params.Add("max_results", strconv.Itoa(maxResults))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")

	apiURL := baseURL + "?" + params.Encode()

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
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

	var feed ArXivFeed
	parseErr := xml.Unmarshal(body, &feed)
	if parseErr != nil {
		return err(parseErr.Error())
}

	var results []map[string]interface{}
	for _, entry := range feed.Entries {
		var authors []string
		for _, a := range entry.Author {
			authors = append(authors, a.Name)

		result := map[string]interface{}{
			"id":        extractArXivID(entry.ID),
			"title":     strings.TrimSpace(entry.Title),
			"summary":   strings.TrimSpace(entry.Summary),
			"authors":   authors,
			"published": entry.Published,
			"source":    "arxiv",
		}
		results = append(results, result)

	return ok(map[string]interface{}{
}
		"category":    category,
		"total_found": len(results),
		"papers":      results,
	})

}
}

// extractArXivID extracts the numeric ID from an arXiv URL
func extractArXivID(id string) string {
	re := regexp.MustCompile(`(\d+\.\d+)`)
	matches := re.FindStringSubmatch(id)
	if len(matches) > 1 {
		return matches[1]
	}
	return id
}

// cleanHTMLAbstract removes HTML tags from abstract text
func cleanHTMLAbstract(abstract string) string {
	if abstract == "" {
		return ""
	}
	// Remove common HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(abstract, "")
	// Decode common HTML entities
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&#39;", "'")
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")
	return strings.TrimSpace(cleaned)
}