package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// YahooFinanceClient is a client for Yahoo Finance API
type YahooFinanceClient struct {
	Client *http.Client
	BaseURL string
}

// NewYahooFinanceClient creates a new Yahoo Finance API client
func NewYahooFinanceClient() *YahooFinanceClient {
	return &YahooFinanceClient{
}
		Client: &http.Client{Timeout: 30 * time.Second},
		BaseURL: "https://query1.finance.yahoo.com",
	}
}

// HandleGetQuote fetches current quote data for a stock symbol
func HandleGetQuote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required")
}

	client := NewYahooFinanceClient()
	apiURL := fmt.Sprintf("%s/v7/finance/quote?symbols=%s", client.BaseURL, url.QueryEscape(strings.ToUpper(symbol)))

	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, fetchErr := client.Client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(decodeErr.Error())
}

	quoteResult, found := result["quoteResponse"].(map[string]interface{})
	if !found {
		return err("invalid response format")
}

	resultList, found := quoteResult["result"].([]interface{})
	if !ok || len(resultList) == 0 {
		return err("no quote data found for symbol: " + symbol)
}

	quote, found := resultList[0].(map[string]interface{})
	if !found {
		return err("invalid quote data format")
}

	output := fmt.Sprintf("Symbol: %s\nPrice: %.2f\nChange: %.2f (%.2f%%)\nMarket Cap: %v\nVolume: %v\n52 Week High: %.2f\n52 Week Low: %.2f",
		getStringValue(quote, "symbol"),
		getFloatValue(quote, "regularMarketPrice"),
		getFloatValue(quote, "regularMarketChange"),
		getFloatValue(quote, "regularMarketChangePercent"),
		formatLargeNumber(getFloatValue(quote, "marketCap")),
		int64(getFloatValue(quote, "regularMarketVolume")),
		getFloatValue(quote, "fiftyTwoWeekHigh"),
		getFloatValue(quote, "fiftyTwoWeekLow"),
	)

	return ok(output)

// HandleGetHistoricalData fetches historical price data for a stock symbol
func HandleGetHistoricalData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required")
}

	period, _ :=getString(args, "period")
	if period == "" {
		period = "1mo"
	}

	interval, _ :=getString(args, "interval")
	if interval == "" {
		interval = "1d"
	}

	client := NewYahooFinanceClient()
	apiURL := fmt.Sprintf("%s/v8/finance/chart/%s?range=%s&interval=%s",
		client.BaseURL,
		url.QueryEscape(strings.ToUpper(symbol)),
		url.QueryEscape(period),
		url.QueryEscape(interval),
	)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, fetchErr := client.Client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(decodeErr.Error())
}

	chart, found := result["chart"].(map[string]interface{})
	if !found {
		return err("invalid response format")
}

	resultArr, found := chart["result"].([]interface{})
	if !ok || len(resultArr) == 0 {
		return err("no historical data found for symbol: " + symbol)
}

	data, found := resultArr[0].(map[string]interface{})
	if !found {
		return err("invalid data format")
}

	meta, _ := data["meta"].(map[string]interface{})
	timestamps, _ := data["timestamp"].([]interface{})
	quotes, _ := data["indicators"].(map[string]interface{})["quote"].([]interface{})

	if meta == nil || timestamps == nil || len(quotes) == 0 {
		return err("incomplete data returned")
}

	quoteData, found := quotes[0].(map[string]interface{})
	if !found {
		return err("invalid quote data format")
}

	closes, _ := quoteData["close"].([]interface{})

	output := fmt.Sprintf("Symbol: %s | Period: %s | Interval: %s\n", 
		getStringValue(meta, "symbol"), period, interval)
	output += "Date | Close Price\n"
	output += strings.Repeat("-", 30) + "\n"

	count := len(timestamps)
	if count > 10 {
		count = 10
	}

	for i := 0; i < count; i++ {
		ts, tsOk := timestamps[i].(float64)
		closeVal := getSliceFloatValue(closes, i)
		if tsOk {
			t := time.Unix(int64(ts), 0)
			output += fmt.Sprintf("%s | %.2f\n", t.Format("2006-01-02"), closeVal)

	}

	if len(timestamps) > 10 {
		output += fmt.Sprintf("... and %d more data points", len(timestamps)-10)

	return ok(output)
}

}
}

// HandleSearchSymbols searches for stock symbols by company name or keyword
func HandleSearchSymbols(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query is required")
}

	client := NewYahooFinanceClient()
	apiURL := fmt.Sprintf("%s/v1/finance/search?q=%s&quotesCount=10&newsCount=0",
		client.BaseURL,
		url.QueryEscape(query),
	)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, fetchErr := client.Client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(decodeErr.Error())
}

	quotes, found := result["quotes"].([]interface{})
	if !found {
		return err("no search results found")
}

	if len(quotes) == 0 {
		return ok("No results found for: " + query)
}

	output := fmt.Sprintf("Search results for: %s\n\n", query)
	output += "Symbol | Name | Exchange | Type\n"
	output += strings.Repeat("-", 60) + "\n"

	for i, q := range quotes {
		if i >= 10 {
			break
		}
		quote, found := q.(map[string]interface{})
		if !found {
			continue
		}
		sym := getStringValue(quote, "symbol")
		name := getStringValue(quote, "shortname")
		if name == "" {
			name = getStringValue(quote, "longname")

		exch := getStringValue(quote, "exchange")
		typ := getStringValue(quote, "quoteType")
		output += fmt.Sprintf("%s | %s | %s | %s\n", sym, name, exch, typ)

	return ok(output)
}

}
}

// HandleGetCompanyInfo fetches detailed company information
func HandleGetCompanyInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required")
}

	client := NewYahooFinanceClient()
	apiURL := fmt.Sprintf("%s/v10/finance/quoteSummary/%s?modules=summaryDetail,defaultKeyStatistics,assetProfile",
		client.BaseURL,
		url.QueryEscape(strings.ToUpper(symbol)),
	)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, fetchErr := client.Client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(decodeErr.Error())
}

	quoteSummary, found := result["quoteSummary"].(map[string]interface{})
	if !found {
		return err("invalid response format")
}

	resultArr, found := quoteSummary["result"].([]interface{})
	if !ok || len(resultArr) == 0 {
		return err("no company info found for symbol: " + symbol)
}

	data, found := resultArr[0].(map[string]interface{})
	if !found {
		return err("invalid data format")
}

	profile, _ := data["assetProfile"].(map[string]interface{})
	summary, _ := data["summaryDetail"].(map[string]interface{})
	stats, _ := data["defaultKeyStatistics"].(map[string]interface{})

	output := fmt.Sprintf("Company Information for: %s\n", symbol)
	output += strings.Repeat("=", 40) + "\n\n"

	if profile != nil {
		output += fmt.Sprintf("Industry: %s\n", getStringValue(profile, "industry"))
		output += fmt.Sprintf("Sector: %s\n", getStringValue(profile, "sector"))
		output += fmt.Sprintf("Country: %s\n", getStringValue(profile, "country"))
		output += fmt.Sprintf("Employees: %d\n", int(getFloatValue(profile, "fullTimeEmployees")))
		output += fmt.Sprintf("Website: %s\n", getStringValue(profile, "website"))
		desc := getStringValue(profile, "longBusinessSummary")
		if len(desc) > 300 {
			desc = desc[:300] + "..."
		}
		output += fmt.Sprintf("Summary: %s\n", desc)

	output += "\n--- Financial Metrics ---\n"
	if summary != nil {
		output += fmt.Sprintf("Market Cap: %v\n", formatLargeNumber(getFloatValue(summary, "marketCap")))
		output += fmt.Sprintf("P/E Ratio: %.2f\n", getFloatValue(summary, "trailingPE"))
		output += fmt.Sprintf("Dividend Yield: %.2f%%\n", getFloatValue(summary, "dividendYield")*100)
		output += fmt.Sprintf("52 Week Change: %.2f%%\n", getFloatValue(summary, "fiftyTwoWeekChange")*100)

	if stats != nil {
		output += fmt.Sprintf("EPS (TTM): %.2f\n", getFloatValue(stats, "trailingEps"))
		output += fmt.Sprintf("Beta: %.2f\n", getFloatValue(stats, "beta"))

	return ok(output)
}

}
}
}

// HandleGetMarketSummary fetches current market summary data
func HandleGetMarketSummary(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := NewYahooFinanceClient()
	apiURL := fmt.Sprintf("%s/v6/finance/quote/marketSummary", client.BaseURL)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(reqErr.Error())
}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, fetchErr := client.Client.Do(req)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
		return err(decodeErr.Error())
}

	marketSummary, found := result["marketSummaryResponse"].(map[string]interface{})
	if !found {
		return err("invalid response format")
}

	resultArr, found := marketSummary["result"].([]interface{})
	if !found {
		return err("no market data available")
}

	output := "Market Summary\n"
	output += strings.Repeat("=", 50) + "\n\n"

	for _, item := range resultArr {
		quote, found := item.(map[string]interface{})
		if !found {
			continue
		}

		shortName := getStringValue(quote, "shortName")
		if shortName == "" {
			continue
		}

		price := getFloatValue(quote, "regularMarketPrice")
		change := getFloatValue(quote, "regularMarketChange")
		changePercent := getFloatValue(quote, "regularMarketChangePercent")

		output += fmt.Sprintf("%s\n", shortName)
		output += fmt.Sprintf("  Price: %.2f | Change: %.2f (%.2f%%)\n", price, change, changePercent)
		output += "\n"
	}

	return ok(output)
}

// Helper functions

func getStringValue(m map[string]interface{}, key string) string {
	if v, found := m[key]; found {
		if s, found := v.(string); found {
			return s
		}
	}
	return ""
}

func getFloatValue(m map[string]interface{}, key string) float64 {
	if v, found := m[key]; found {
		switch val := v.(type) {
		case float64:
			return val
}
		case int:
			return float64(val)
}
		case string:
			if f, parseErr := strconv.ParseFloat(val, 64); parseErr == nil {
				return f
			}
		}
	}
	return 0
}

func getSliceFloatValue(slice []interface{}, index int) float64 {
	if slice == nil || index < 0 || index >= len(slice) {
		return 0
	}
	switch v := slice[index].(type) {
	case float64:
		return v
	case int:
		return float64(v)

	return 0
}

func formatLargeNumber(val float64) string {
	if val == 0 {
		return "N/A"
	}
	if val >= 1e12 {
		return fmt.Sprintf("%.2fT", val/1e12)
}

	if val >= 1e9 {
		return fmt.Sprintf("%.2fB", val/1e9)
}

	if val >= 1e6 {
		return fmt.Sprintf("%.2fM", val/1e6)
}

	if val >= 1e3 {
		return fmt.Sprintf("%.2fK", val/1e3)
}

	return fmt.Sprintf("%.2f", val)
}