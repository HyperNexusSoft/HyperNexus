package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"regexp"
	"sort"
)

// ToolResponse, ok, e, getString, getInt, getBool は parity.go で定義されていると仮定します。

func HandleX(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	apiKey, _ :=getString(args, "api_key")
	symbol, _ :=getString(args, "symbol")

	// URLの作成
	url := fmt.Sprintf("https://api.yfinance.com/quote/%s?apikey=%s", symbol, apiKey)

	// HTTPクライアントの設定
	client := http.DefaultClient

	// リクエストの送信
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	// リクエストの実行
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// JSONのパース
	var data map[string]interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return nil, parseErr
	}

	// データの整形
	price, found := data["regularMarketPrice"].(float64)
	if !found {
		return nil, fmt.Errorf("failed to get price")
}

	// レスポンスの作成
	return ok(fmt.Sprintf("Price for %s: %.2f", symbol, price))
}

func HandleY(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	apiKey, _ :=getString(args, "api_key")
	symbol, _ :=getString(args, "symbol")

	// URLの作成
	url := fmt.Sprintf("https://api.yfinance.com/quote/%s/history?apikey=%s", symbol, apiKey)

	// HTTPクライアントの設定
	client := http.DefaultClient

	// リクエストの送信
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	// リクエストの実行
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// JSONのパース
	var data map[string]interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return nil, parseErr
	}

	// データの整形
	var prices []float64
	for _, item := range data["prices"].([]interface{}) {
		price, found := item.(map[string]interface{})["regularMarketPrice"].(float64)
		if found {
			prices = append(prices, price)

	}

	// レスポンスの作成
	return ok(fmt.Sprintf("Historical prices for %s: %v", symbol, prices))
}

}

func HandleZ(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	apiKey, _ :=getString(args, "api_key")
	symbol, _ :=getString(args, "symbol")

	// URLの作成
	url := fmt.Sprintf("https://api.yfinance.com/quote/%s/profile?apikey=%s", symbol, apiKey)

	// HTTPクライアントの設定
	client := http.DefaultClient

	// リクエストの送信
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	// リクエストの実行
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// JSONのパース
	var data map[string]interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return nil, parseErr
	}

	// データの整形
	companyName, found := data["longName"].(string)
	if !found {
		return nil, fmt.Errorf("failed to get company name")
}

	// レスポンスの作成
	return ok(fmt.Sprintf("Company name for %s: %s", symbol, companyName))
}

func HandleA(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	apiKey, _ :=getString(args, "api_key")
	symbol, _ :=getString(args, "symbol")

	// URLの作成
	url := fmt.Sprintf("https://api.yfinance.com/quote/%s/options?apikey=%s", symbol, apiKey)

	// HTTPクライアントの設定
	client := http.DefaultClient

	// リクエストの送信
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	// リクエストの実行
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// JSONのパース
	var data map[string]interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return nil, parseErr
	}

	// データの整形
	var options []string
	for _, item := range data["options"].([]interface{}) {
		options = append(options, item.(map[string]interface{})["symbol"].(string))

	// レスポンスの作成
	return ok(fmt.Sprintf("Options for %s: %v", symbol, options))
}

}

func HandleB(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// 引数の取得
	apiKey, _ :=getString(args, "api_key")
	symbol, _ :=getString(args, "symbol")

	// URLの作成
	url := fmt.Sprintf("https://api.yfinance.com/quote/%s/technicals?apikey=%s", symbol, apiKey)

	// HTTPクライアントの設定
	client := http.DefaultClient

	// リクエストの送信
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}

	// リクエストの実行
	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return nil, fetchErr
	}
	defer resp.Body.Close()

	// レスポンスの読み取り
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	// JSONのパース
	var data map[string]interface{}
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return nil, parseErr
	}

	// データの整形
	var technicals map[string]interface{}
	for key, value := range data {
		if key == "symbol" {
			continue
		}
		technicals[key] = value
	}

	// レスポンスの作成
	return ok(fmt.Sprintf("Technical indicators for %s: %v", symbol, technicals))
}