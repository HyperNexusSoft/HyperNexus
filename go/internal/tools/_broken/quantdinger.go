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

func HandleGetMarketData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required")
}

	interval, _ :=getString(args, "interval")
	if interval == "" {
		interval = "1d"
	}

	startDate, _ :=getString(args, "start_date")
	endDate, _ :=getString(args, "end_date")

	client := http.DefaultClient
	apiUrl := fmt.Sprintf("https://api.example.com/marketdata/%s?interval=%s", symbol, interval)

	if startDate != "" {
		apiUrl += fmt.Sprintf("&start_date=%s", url.QueryEscape(startDate))

	if endDate != "" {
		apiUrl += fmt.Sprintf("&end_date=%s", url.QueryEscape(endDate))

	resp, reqErr := client.Get(apiUrl)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to fetch market data: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var result map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Market data for %s (%s): %v", symbol, interval, result))
}

}
}

func HandleRunBacktest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	strategy, _ :=getString(args, "strategy")
	if strategy == "" {
		return err("strategy is required")
}

	symbol, _ :=getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required")
}

	startDate, _ :=getString(args, "start_date")
	if startDate == "" {
		return err("start_date is required")
}

	endDate, _ :=getString(args, "end_date")
	if endDate == "" {
		return err("end_date is required")
}

	initialCapital, _ :=getInt(args, "initial_capital")
	if initialCapital <= 0 {
		initialCapital = 10000
	}

	client := http.DefaultClient
	apiUrl := fmt.Sprintf("https://api.example.com/backtest?strategy=%s&symbol=%s&start_date=%s&end_date=%s&initial_capital=%d",
		url.QueryEscape(strategy),
		url.QueryEscape(symbol),
		url.QueryEscape(startDate),
		url.QueryEscape(endDate),
		initialCapital)

	resp, reqErr := client.Get(apiUrl)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to run backtest: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var result map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Backtest results: %v", result))
}

func HandleGetStrategyList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient
	resp, reqErr := client.Get("https://api.example.com/strategies")
	if reqErr != nil {
		return err(fmt.Sprintf("failed to fetch strategies: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var strategies []string
	parseErr := json.NewDecoder(resp.Body).Decode(&strategies)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Available strategies: %s", strings.Join(strategies, ", ")))
}

func HandleGetAccountBalance(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient
	resp, reqErr := client.Get("https://api.example.com/account/balance")
	if reqErr != nil {
		return err(fmt.Sprintf("failed to fetch account balance: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var balance map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&balance)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Account balance: %v", balance))
}

func HandleGetOrderHistory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	limit, _ :=getInt(args, "limit")
	if limit <= 0 {
		limit = 10
	}

	client := http.DefaultClient
	apiUrl := fmt.Sprintf("https://api.example.com/orders?symbol=%s&limit=%d",
		url.QueryEscape(symbol),
		limit)

	resp, reqErr := client.Get(apiUrl)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to fetch order history: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var orders []map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&orders)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Order history for %s: %v", symbol, orders))
}

func HandleGetNews(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ :=getString(args, "symbol")
	source, _ :=getString(args, "source")
	if source == "" {
		source = "all"
	}

	client := http.DefaultClient
	apiUrl := fmt.Sprintf("https://api.example.com/news?symbol=%s&source=%s",
		url.QueryEscape(symbol),
		url.QueryEscape(source))

	resp, reqErr := client.Get(apiUrl)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to fetch news: %v", reqErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned status %d", resp.StatusCode))
}

	var news []map[string]interface{}
	parseErr := json.NewDecoder(resp.Body).Decode(&news)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	return ok(fmt.Sprintf("Latest news for %s: %v", symbol, news))
}