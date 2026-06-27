package tools

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func HandleHealthCheck(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    return ok("robertpelloni-com: service is healthy")
}

func HandleFetchURL(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    targetURL, _ :=getString(args, "url")
    if targetURL == "" {
        return err("url parameter is required")
}

    client := http.DefaultClient
    resp, fetchErr := client.Get(targetURL)
    if fetchErr != nil {
        return err(fetchErr.Error())
}

    defer resp.Body.Close()

    result := fmt.Sprintf("Status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
    return ok(result)
}

func HandleGetInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    info := "robertpelloni-com MCP Server\n"
    info += "Available tools: health_check, fetch_url, get_info\n"
    info += "Version: 1.0.0"
    return ok(info)
}