package tools

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

// HandleEcho sends back the message provided.
func HandleEcho(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ :=getString(args, "message")
	return ok(msg)
}

// HandleCurrentTime returns the current time.
func HandleCurrentTime(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	now := time.Now().Format(time.RFC3339)
	return ok(now)
}

// HandleFetchURL fetches a URL and returns the response body.
func HandleFetchURL(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	if urlStr == "" {
		return err("url is required")
}

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
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

	return ok(string(body))
}

// HandleEncodeBase64 encodes a string to base64.
func HandleEncodeBase64(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	input, _ :=getString(args, "input")
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return ok(encoded)
}

// HandleDecodeBase64 decodes a base64 string.
func HandleDecodeBase64(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	input, _ :=getString(args, "input")
	decoded, decodeErr := base64.StdEncoding.DecodeString(input)
	if decodeErr != nil {
		return err(decodeErr.Error())
}

	return ok(string(decoded))
}