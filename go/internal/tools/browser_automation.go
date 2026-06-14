package tools

/**
 * @file browser_automation.go
 * @module go/internal/tools
 *
 * WHAT: Go-native browser automation handlers using chromedp.
 * Provides navigation, screenshots, HTML extraction, JavaScript evaluation, and form interactions.
 *
 * WHY: Replaces external puppeteer/browser-use MCP servers with a lightweight Go-native implementation.
 */

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// getDefaultAllocator returns a chromedp allocator with sensible defaults.
func getDefaultAllocator() chromedp.ExecAllocatorOption {
	return chromedp.CompositeFlags(
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.WindowSize(1920, 1080),
	)
}

// HandleBrowserNavigate navigates to a URL.
// Args: url (string)
func HandleBrowserNavigate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, ok := args["url"].(string)
	if !ok || urlStr == "" {
		return err("url is required")
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	var title string
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.Title(&title),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return err(fmt.Sprintf("Navigation failed: %v", err))
	}

	return ok(fmt.Sprintf("Navigated to %s (title: %s)", urlStr, title))
}

// HandleBrowserScreenshot captures a screenshot of a page.
// Args: url (string, required), fullPage (bool, optional, default false)
func HandleBrowserScreenshot(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ := args["url"].(string)
	if urlStr == "" {
		return err("url is required")
	}
	fullPage, _ := args["fullPage"].(bool)

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	var buf []byte
	opts := []page.CaptureScreenshotOption{
		page.CaptureScreenshotQuality(80),
	}
	if fullPage {
		opts = append(opts, page.CaptureScreenshotFromSurface(true))
	}

	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, 80),
	); err != nil {
		return err(fmt.Sprintf("Screenshot failed: %v", err))
	}

	base64Img := base64.StdEncoding.EncodeToString(buf)
	return ok(fmt.Sprintf("data:image/png;base64,%s", base64Img))
}

// HandleBrowserGetHTML retrieves the full HTML of a page.
// Args: url (string, required)
func HandleBrowserGetHTML(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ := args["url"].(string)
	if urlStr == "" {
		return err("url is required")
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	var html string
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML(":root", &html, chromedp.NodeQuery),
	); err != nil {
		return err(fmt.Sprintf("HTML retrieval failed: %v", err))
	}

	return ok(html)
}

// HandleBrowserEvaluate executes JavaScript on a page and returns the result.
// Args: url (string, required), script (string, required)
func HandleBrowserEvaluate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ := args["url"].(string)
	script, _ := args["script"].(string)

	if urlStr == "" {
		return err("url is required")
	}
	if script == "" {
		return err("script is required")
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	var result interface{}
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(script, &result),
	); err != nil {
		return err(fmt.Sprintf("Evaluation failed: %v", err))
	}

	resultBytes, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return err(fmt.Sprintf("Failed to marshal result: %v", marshalErr))
	}

	return ok(string(resultBytes))
}

// HandleBrowserClick clicks an element on a page.
// Args: url (string, required), selector (string, required)
func HandleBrowserClick(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ := args["url"].(string)
	selector, _ := args["selector"].(string)

	if urlStr == "" {
		return err("url is required")
	}
	if selector == "" {
		return err("selector is required")
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Click(selector, chromedp.NodeVisible),
	); err != nil {
		return err(fmt.Sprintf("Click failed: %v", err))
	}

	return ok(fmt.Sprintf("Clicked element: %s", selector))
}

// HandleBrowserFillForm fills an input field with a value.
// Args: url (string, required), selector (string, required), value (string, required)
func HandleBrowserFillForm(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ := args["url"].(string)
	selector, _ := args["selector"].(string)
	value, _ := args["value"].(string)

	if urlStr == "" {
		return err("url is required")
	}
	if selector == "" {
		return err("selector is required")
	}
	if value == "" {
		return err("value is required")
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, getDefaultAllocator())
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(func(format string, args ...interface{}) {}))
	defer cancel()

	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitVisible(selector, chromedp.NodeVisible),
		chromedp.SendKeys(selector, value),
	); err != nil {
		return err(fmt.Sprintf("Fill form failed: %v", err))
	}

	return ok(fmt.Sprintf("Filled %s with value: %s", selector, value))
}

// init registers the browser automation handlers when the package is loaded.
func init() {
	// Handlers are registered via the registry's init() chain
	// See registry.go for the registration calls
}