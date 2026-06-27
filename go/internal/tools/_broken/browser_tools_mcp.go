package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	discoveredHost   = "127.0.0.1"
	discoveredPort   = 3025
	serverDiscovered = false
	client = http.DefaultClient
)

func discoverServer() bool {
	hosts := []string{getDefaultServerHost(), "127.0.0.1", "localhost"}
	ports := []int{getDefaultServerPort()}

	// Add fallback ports
	for p := 3025; p <= 3035; p++ {
		if p != ports[0] {
			ports = append(ports, p)

	}

	for _, host := range hosts {
		for _, port := range ports {
			target := fmt.Sprintf("http://%s:%d/.identity", host, port)
			resp, fetchErr := client.Get(target)
			if fetchErr != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				continue
			}

			var identity struct {
				Signature string `json:"signature"`
			}
			if parseErr := json.NewDecoder(resp.Body).Decode(&identity); parseErr != nil {
				continue
			}

			if identity.Signature == "mcp-browser-connector-24x7" {
				discoveredHost = host
				discoveredPort = port
				serverDiscovered = true
				return true
			}
		}
	}
	return false
}

}

func getDefaultServerPort() int {
	if envPort := os.Getenv("BROWSER_TOOLS_PORT"); envPort != "" {
		if port, convErr := strconv.Atoi(envPort); convErr == nil && port > 0 {
			return port
		}
	}

	portFile := filepath.Join(".", ".port")
	if data, readErr := os.ReadFile(portFile); readErr == nil {
		if port, convErr := strconv.Atoi(strings.TrimSpace(string(data))); convErr == nil && port > 0 {
			return port
		}
	}

	return 3025
}

func getDefaultServerHost() string {
	if host := os.Getenv("BROWSER_TOOLS_HOST"); host != "" {
		return host
	}
	return "127.0.0.1"
}

func withServerConnection(ctx context.Context, apiCall func() (ToolResponse, error)) (ToolResponse, error) {
	if !serverDiscovered {
		if !discoverServer() {
			return err("Failed to discover browser connector server. Please ensure it's running.")

	}

	res, callErr := apiCall()
	if callErr == nil {
		return res, nil
	}

	serverDiscovered = false
	if discoverServer() {
		return apiCall()
}

	return err(fmt.Sprintf("Failed to reconnect to server: %v", callErr))
}

}

func HandleGetConsoleLogs(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/console-logs", discoveredHost, discoveredPort)
		resp, fetchErr := client.Get(target)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err(fmt.Sprintf("Server returned status: %d", resp.StatusCode))
}

		var logs []interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&logs); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(logs))
	})

func HandleGetConsoleErrors(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/console-errors", discoveredHost, discoveredPort)
		resp, fetchErr := client.Get(target)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err(fmt.Sprintf("Server returned status: %d", resp.StatusCode))
}

		var errors []interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errors); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(errors))
	})

func HandleGetNetworkErrors(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/network-errors", discoveredHost, discoveredPort)
		resp, fetchErr := client.Get(target)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err(fmt.Sprintf("Server returned status: %d", resp.StatusCode))
}

		var errors []interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errors); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(errors))
	})

func HandleGetNetworkLogs(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/network-success", discoveredHost, discoveredPort)
		resp, fetchErr := client.Get(target)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err(fmt.Sprintf("Server returned status: %d", resp.StatusCode))
}

		var logs []interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&logs); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(logs))
	})

func HandleTakeScreenshot(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/capture-screenshot", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var result struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
			return err(decodeErr.Error())
}

		if !result.Success {
			return err(result.Error)
}

		return ok("Successfully saved screenshot")
	})

func HandleGetSelectedElement(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/selected-element", discoveredHost, discoveredPort)
		resp, fetchErr := client.Get(target)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err(fmt.Sprintf("Server returned status: %d", resp.StatusCode))
}

		var element struct {
			TagName     string `json:"tagName"`
			TextContent string `json:"textContent"`
			Attributes  []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"attributes"`
			XPath string `json:"xpath"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&element); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(element))
	})

func HandleRunAccessibilityAudit(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-accessibility-audit", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var audit struct {
			Score    float64 `json:"score"`
			Issues   []struct {
				Description string `json:"description"`
				Severity   string `json:"severity"`
				Help       string `json:"help"`
			} `json:"issues"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&audit); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(audit))
	})

func HandleRunPerformanceAudit(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-performance-audit", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var audit struct {
			Score    float64 `json:"score"`
			Metrics  []struct {
				Name  string  `json:"name"`
				Value float64 `json:"value"`
				Unit  string  `json:"unit"`
			} `json:"metrics"`
			Issues []struct {
				Description string `json:"description"`
				Severity   string `json:"severity"`
				Help       string `json:"help"`
			} `json:"issues"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&audit); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(audit))
	})

func HandleRunSEOAudit(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-seo-audit", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var audit struct {
			Score    float64 `json:"score"`
			Issues   []struct {
				Description string `json:"description"`
				Severity   string `json:"severity"`
				Help       string `json:"help"`
			} `json:"issues"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&audit); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(audit))
	})

func HandleRunBestPracticesAudit(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-best-practices-audit", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var audit struct {
			Score    float64 `json:"score"`
			Issues   []struct {
				Description string `json:"description"`
				Severity   string `json:"severity"`
				Help       string `json:"help"`
			} `json:"issues"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&audit); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(audit))
	})

func HandleRunNextJSAudit(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	router, _ :=getString(args, "router")
	if router == "" {
		router = "auto"
	}

	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-nextjs-audit", discoveredHost, discoveredPort)
		data := url.Values{}
		data.Set("router", router)

		resp, fetchErr := client.PostForm(target, data)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var audit struct {
			Score    float64 `json:"score"`
			Issues   []struct {
				Description string `json:"description"`
				Severity   string `json:"severity"`
				Help       string `json:"help"`
			} `json:"issues"`
			Recommendations []string `json:"recommendations"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&audit); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(audit))
	})

func HandleRunAuditMode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-audit-mode", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var result struct {
			Accessibility  interface{} `json:"accessibility"`
			Performance    interface{} `json:"performance"`
			SEO           interface{} `json:"seo"`
			BestPractices interface{} `json:"bestPractices"`
			NextJS        interface{} `json:"nextjs"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(result))
	})

func HandleRunDebuggerMode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return withServerConnection(ctx, func() (ToolResponse, error) {
}
		target := fmt.Sprintf("http://%s:%d/run-debugger-mode", discoveredHost, discoveredPort)
		resp, fetchErr := client.Post(target, "application/json", nil)
		if fetchErr != nil {
			return err(fetchErr.Error())
}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return err(fmt.Sprintf("Server returned status: %d, body: %s", resp.StatusCode, string(body)))
}

		var result struct {
			ConsoleLogs    []interface{} `json:"consoleLogs"`
			NetworkLogs    []interface{} `json:"networkLogs"`
			SelectedElement interface{} `json:"selectedElement"`
			AuditResults   interface{}  `json:"auditResults"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&result); decodeErr != nil {
			return err(decodeErr.Error())
}

		return ok(jsonString(result))
	})

func jsonString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}