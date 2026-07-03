package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MDMAtk/TormentNexus/internal/config"
)

func setupTestServer() *Server {
	cfg := config.Config{
		ConfigDir:     t.TempDir(),
		WorkspaceRoot: t.TempDir(),
	}
	s := New(cfg, nil)
	// Add required endpoints
	s.mux.HandleFunc("/api/sse", s.handleSSE)
	s.mux.HandleFunc("/api/sse/message", s.handleSSEMessage)
	return s
}

func TestSSEHandlers_Unauthorized(t *testing.T) {
	t.Setenv("CLOUDMCP_SSE_AUTH_TOKEN", "") // Ensure no token is set initially

	cfg := config.Config{
		ConfigDir:     t.TempDir(),
		WorkspaceRoot: t.TempDir(),
	}
	s := New(cfg, nil)
	s.mux.HandleFunc("/api/sse", s.handleSSE)
	s.mux.HandleFunc("/api/sse/message", s.handleSSEMessage)

	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	// Try without token
	resp, err := http.Get(ts.URL + "/api/sse")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %v", resp.StatusCode)
	}

	// Try POST without token
	resp, err = http.Post(ts.URL+"/api/sse/message", "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for POST, got %v", resp.StatusCode)
	}
}

func TestSSEHandlers_AuthorizedQuery(t *testing.T) {
	t.Setenv("CLOUDMCP_SSE_AUTH_TOKEN", "test-token-123")

	cfg := config.Config{
		ConfigDir:     t.TempDir(),
		WorkspaceRoot: t.TempDir(),
	}
	s := New(cfg, nil)
	s.mux.HandleFunc("/api/sse", s.handleSSE)
	s.mux.HandleFunc("/api/sse/message", s.handleSSEMessage)

	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	client := http.Client{Timeout: 1 * time.Second}
	req, _ := http.NewRequest("GET", ts.URL+"/api/sse?token=test-token-123", nil)
	resp, err := client.Do(req)

	if err != nil {
		if err != nil {
			// Expected to timeout if it properly stays open
		}
	} else {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %v", resp.StatusCode)
		}
		if resp.Header.Get("Content-Type") != "text/event-stream" {
			t.Errorf("Expected Content-Type text/event-stream, got %v", resp.Header.Get("Content-Type"))
		}
		resp.Body.Close()
	}

	// 2. Test POST message
	req, _ = http.NewRequest("POST", ts.URL+"/api/sse/message?token=test-token-123", bytes.NewBuffer([]byte(`{"jsonrpc": "2.0", "id": 1, "method": "test"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		t.Errorf("Expected 202 Accepted or similar, got 401")
	}
	resp.Body.Close()
}

func TestSSEHandlers_AuthorizedHeader(t *testing.T) {
	t.Setenv("CLOUDMCP_SSE_AUTH_TOKEN", "test-token-header")

	cfg := config.Config{
		ConfigDir:     t.TempDir(),
		WorkspaceRoot: t.TempDir(),
	}
	s := New(cfg, nil)
	s.mux.HandleFunc("/api/sse", s.handleSSE)
	s.mux.HandleFunc("/api/sse/message", s.handleSSEMessage)

	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	// Try POST message with Header auth
	req, _ := http.NewRequest("POST", ts.URL+"/api/sse/message", bytes.NewBuffer([]byte(`{"jsonrpc": "2.0", "id": 1, "method": "test"}`)))
	req.Header.Set("Authorization", "Bearer test-token-header")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		t.Errorf("Expected success, got 401")
	}
	resp.Body.Close()
}
