package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "net/url"
    "strings"
)

func ok(text string) ToolResponse {
    return ToolResponse{Text: text, OK: true}
}

func err(e error) ToolResponse {
    return ToolResponse{Text: e.Error(), OK: false}
}

func getString(args map[string]interface{}, key string) (string, bool) {
    val, found := args[key]
    if !found {
        return "", false
    }
    strVal, found := val.(string)
    if !found {
        return "", false
    }
    return strVal, true
}

func getInt(args map[string]interface{}, key string) (int, bool) {
    val, found := args[key]
    if !found {
        return 0, false
    }
    intVal, found := val.(int)
    if !found {
        return 0, false
    }
    return intVal, true
}

func getBool(args map[string]interface{}, key string) (bool, bool) {
    val, found := args[key]
    if !found {
        return false, false
    }
    boolVal, found := val.(bool)
    if !found {
        return false, false
    }
    return boolVal, true
}

// HandlePing handles the ping tool
func HandlePing(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    return ok("Pong")
}

// HandleEcho handles the echo tool
func HandleEcho(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    message, _ :=getString(args, "message")
    if !found {
        return err("Missing required argument: message")
    }
    return ok(message)
}

// HandleTest handles the test tool
func HandleTest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    urlValue, _ :=getString(args, "url")
    if !found {
        return err("Missing required argument: url")
    }
    u, parseErr := url.Parse(urlValue)
    if parseErr != nil {
        return err(parseErr.Error())
    }
    client := http.DefaultClient
    req, reqErr := http.NewRequest("GET", u.String(), nil)
    if reqErr != nil {
        return err(reqErr.Error())
    }
    resp, fetchErr := client.Do(req)
    if fetchErr != nil {
        return err(fetchErr.Error())
    }
    defer resp.Body.Close()
    return ok(fmt.Sprintf("%s %s", resp.Status, resp.Status))
}

func HandleVersion(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    return ok("1.0.0")
}

func HandleHelp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    return ok("Available tools: ping, echo, test, version")
}