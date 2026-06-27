package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// enscanBaseURL is the default API endpoint for ENScan_GO
var enscanBaseURL = "http://localhost:31000"

// doEnscanGet performs a GET request to the ENScan API and returns the raw body as string
func doEnscanGet(ctx context.Context, path string, params url.Values) (string, error) {
	u := enscanBaseURL + path
	if len(params) > 0 {
		u = u + "?" + params.Encode()

	req, reqErr := http.NewRequestWithContext(ctx, "GET", u, nil)
	if reqErr != nil {
		return "", fmt.Errorf("creating request: %w", reqErr)
}

	client := http.DefaultClient
	resp, doErr := client.Do(req)
	if doErr != nil {
		return "", fmt.Errorf("executing request: %w", doErr)
}

	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", fmt.Errorf("reading response: %w", readErr)
}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
}

	return string(body), nil
}

}

// HandleEnscanSearch searches for enterprise information using the ENScan_GO API
func HandleEnscanSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("name parameter is required")
}

	params := url.Values{}
	params.Set("name", name)

	// Optional type parameter (data source: aqc, tyc, kc, rb, all, miit, coolapk, qimai)
	if t := getString(args, "type"); t != "" {
		params.Set("type", t)

	// Optional field parameter (icp, weibo, wechat, app, job, wx_app, copyright, supplier)
	if f := getString(args, "field"); f != "" {
		params.Set("field", f)

	// Optional invest parameter (investment ratio filter)
	if inv := getInt(args, "invest"); inv != 0 {
		params.Set("invest", strconv.Itoa(inv))

	// Optional depth parameter (recursive search depth)
	if d := getInt(args, "depth"); d != 0 {
		params.Set("depth", strconv.Itoa(d))

	// Optional holds parameter
	if h := getBool(args, "holds"); h {
		params.Set("holds", "true")

	// Optional supplier parameter
	if s := getBool(args, "supplier"); s {
		params.Set("supplier", "true")

	// Optional branch parameter
	if b := getBool(args, "branch"); b {
		params.Set("branch", "true")

	result, apiErr := doEnscanGet(ctx, "/api/info", params)
	if apiErr != nil {
		return err(fmt.Sprintf("ENScan search failed: %v", apiErr))
}

	// Try to pretty-print if JSON
	var prettyJSON map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(result), &prettyJSON); jsonErr == nil {
		formatted, fmtErr := json.MarshalIndent(prettyJSON, "", "  ")
		if fmtErr == nil {
			return ok(string(formatted))

	}

	return ok(result)
}

}
}
}
}
}
}
}
}

// HandleEnscanAdvanceFilter uses the pro/advance_filter endpoint to search companies
func HandleEnscanAdvanceFilter(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("name parameter is required")
}

	dataType, _ :=getString(args, "type")
	if dataType == "" {
		return err("type parameter is required (e.g., aqc, tyc, kc, rb)")
}

	params := url.Values{}
	params.Set("name", name)
	params.Set("type", dataType)

	path := "/api/pro/advance_filter"
	result, apiErr := doEnscanGet(ctx, path, params)
	if apiErr != nil {
		return err(fmt.Sprintf("ENScan advance filter failed: %v", apiErr))
}

	var prettyJSON map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(result), &prettyJSON); jsonErr == nil {
		formatted, fmtErr := json.MarshalIndent(prettyJSON, "", "  ")
		if fmtErr == nil {
			return ok(string(formatted))

	}

	return ok(result)
}

}

// HandleEnscanGetBaseInfo retrieves basic company information by PID using the pro/get_base_info endpoint
func HandleEnscanGetBaseInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pid, _ :=getString(args, "pid")
	if pid == "" {
		return err("pid parameter is required")
}

	dataType, _ :=getString(args, "type")
	if dataType == "" {
		return err("type parameter is required (e.g., aqc, tyc, kc, rb)")
}

	params := url.Values{}
	params.Set("pid", pid)
	params.Set("type", dataType)

	path := "/api/pro/get_base_info"
	result, apiErr := doEnscanGet(ctx, path, params)
	if apiErr != nil {
		return err(fmt.Sprintf("ENScan get base info failed: %v", apiErr))
}

	var prettyJSON map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(result), &prettyJSON); jsonErr == nil {
		formatted, fmtErr := json.MarshalIndent(prettyJSON, "", "  ")
		if fmtErr == nil {
			return ok(string(formatted))

	}

	return ok(result)
}

}

// HandleEnscanGetPage retrieves paginated data for a specific field using the pro/get_page endpoint
func HandleEnscanGetPage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pid, _ :=getString(args, "pid")
	if pid == "" {
		return err("pid parameter is required")
}

	dataType, _ :=getString(args, "type")
	if dataType == "" {
		return err("type parameter is required (e.g., aqc, tyc, kc, rb)")
}

	field, _ :=getString(args, "field")
	if field == "" {
		return err("field parameter is required (e.g., icp, weibo, wechat, app, job, wx_app, copyright, supplier)")
}

	params := url.Values{}
	params.Set("pid", pid)
	params.Set("type", dataType)
	params.Set("field", field)

	if p := getInt(args, "page"); p > 0 {
		params.Set("page", strconv.Itoa(p))

	path := "/api/pro/get_page"
	result, apiErr := doEnscanGet(ctx, path, params)
	if apiErr != nil {
		return err(fmt.Sprintf("ENScan get page failed: %v", apiErr))
}

	var prettyJSON map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(result), &prettyJSON); jsonErr == nil {
		formatted, fmtErr := json.MarshalIndent(prettyJSON, "", "  ")
		if fmtErr == nil {
			return ok(string(formatted))

	}

	return ok(result)
}

}
}

// HandleEnscanGetENMap retrieves the field mapping (ENMap) for a given data source
func HandleEnscanGetENMap(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	dataType, _ :=getString(args, "type")
	if dataType == "" {
		return err("type parameter is required (e.g., aqc, tyc, kc, rb)")
}

	params := url.Values{}
	params.Set("type", dataType)

	path := "/api/pro/get_ensd"
	result, apiErr := doEnscanGet(ctx, path, params)
	if apiErr != nil {
		return err(fmt.Sprintf("ENScan get ENMap failed: %v", apiErr))
}

	var prettyJSON map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(result), &prettyJSON); jsonErr == nil {
		formatted, fmtErr := json.MarshalIndent(prettyJSON, "", "  ")
		if fmtErr == nil {
			return ok(string(formatted))

	}

	return ok(result)
}

}

// HandleEnscanCheckHealth checks if the ENScan_GO API server is reachable
func HandleEnscanCheckHealth(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", enscanBaseURL+"/api/info?name=test", nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create health check request: %v", reqErr))
}

	resp, doErr := client.Do(req)
	if doErr != nil {
		return err(fmt.Sprintf("ENScan_GO API server at %s is not reachable: %v", enscanBaseURL, doErr))
}

	defer resp.Body.Close()

	statusMsg := fmt.Sprintf("ENScan_GO API server at %s is reachable (status %d)", enscanBaseURL, resp.StatusCode)
	return ok(statusMsg)
}

// setEnscanBaseURL allows overriding the API base URL (e.g., from environment)
func init() {
	if envURL := strings.TrimSpace(osGetEnv("ENSCAN_API_URL")); envURL != "" {
		enscanBaseURL = strings.TrimRight(envURL, "/")

}

}

// osGetEnv wraps os.Getenv for testability
func osGetEnv(key string) string {
	return strings.TrimSpace(osLookupEnv(key))
}

// osLookupEnv wraps os.Getenv
func osLookupEnv(key string) string {
	val, _ := osGetEnvRaw(key)
	return val
}

// osGetEnvRaw wraps os.LookupEnv
func osGetEnvRaw(key string) (string, bool) {
	return os.LookupEnv(key)
}