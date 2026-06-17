
package tools

/**
 * @file registry.go
 * @module go/internal/tools
 *
 * WHAT: Go-native registry for standard library and parity tools.
 * Maps tool names to their native Go implementations.
 */

import (
	"context"
	"fmt"
	"sync"
)

type ToolHandler func(ctx context.Context, args map[string]interface{}) (ToolResponse, error)

type Registry struct {
	mu       sync.RWMutex
	handlers map[string]ToolHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		handlers: make(map[string]ToolHandler),
	}
	r.registerAll()
	return r
}

func (r *Registry) registerAll() {
	// Native Handlers
	r.handlers["read_file"] = HandleRead
	r.handlers["write_file"] = HandleWrite
	r.handlers["edit_file"] = HandleEdit
	r.handlers["str_replace_editor"] = HandleEdit
	r.handlers["grep_search"] = HandleGrep
	r.handlers["search_files"] = HandleGrep
	r.handlers["glob"] = HandleGlob
	r.handlers["find_files"] = HandleGlob
	r.handlers["apply_patch"] = HandleApplyPatch
	r.handlers["multi_edit"] = HandleMultiEdit
	r.handlers["bash"] = HandleBash
	r.handlers["ls"] = HandleLS
	r.handlers["list_directory"] = HandleLS
	r.handlers["web_fetch"] = HandleWebFetch

	// Filesystem MCP Tools

	// Ollama MCP Tools (AI & LLM Integration)

	// TTS MCP Tools (Media & Design)

	// Vercel MCP Tools (Cloud & DevOps)

	// DexPaprika MCP Tools (Finance & Crypto)

	// National Weather Service (NWS) MCP Tools (Weather & Location)

	// ast-grep-mcp Tools (Category 11)

	// PAL Tools (Category 12)

	// Short/alias mappings for PAL tools without prefix

	// Serena Tools (Category 13)

	// Claude Code Aliases
	r.handlers["Read"] = HandleRead
	r.handlers["Write"] = HandleWrite
	r.handlers["Edit"] = HandleEdit
	r.handlers["Bash"] = HandleBash
	r.handlers["LS"] = HandleLS
	r.handlers["WebFetch"] = HandleWebFetch
	r.handlers["Glob"] = HandleGlob
	r.handlers["Grep"] = HandleGrep
	r.handlers["MultiEdit"] = HandleMultiEdit

	// Codex Aliases
	r.handlers["shell"] = HandleBash
	r.handlers["create_file"] = HandleWrite
	r.handlers["view_file"] = HandleRead
	r.handlers["apply_diff"] = HandleApplyPatch
	r.handlers["search_files_codex"] = HandleGrep

	// OpenCode / Pi Aliases
	r.handlers["read"] = HandleRead
	r.handlers["write"] = HandleWrite
	r.handlers["edit"] = HandleEdit
	r.handlers["grep"] = HandleGrep
	r.handlers["ls"] = HandleLS
	r.handlers["glob_pi"] = HandleGlob

	// Thoughtbox Tools (Category 14)

	// Fetch Tool (Assimilated)

	// Tavily Tools (Assimilated)

	// Chrome DevTools Tools (Assimilated)

	// Firecrawl Tools (Assimilated from firecrawl-mcp)

	// Exa Search Tools (Assimilated from SSE exa)

	// arXiv Tools (Assimilated from arxiv-mcp-server)

	// Semantic Scholar Tools (Assimilated from paper_search_server)

	// mem0 Memory Tools (Assimilated from @mem0/mcp-server)

	// Alpaca Trading Tools (Assimilated from alpaca-mcp-server)

	// Alpha Vantage Financial Tools (Assimilated from av-mcp)

	// Hugging Face Hub Tools (Assimilated from SSE 	r.handlers["hf_search_models"] = HandleHFSearchModels

	// Semgrep Security Tools (Assimilated from semgrep + semgrepstream)

	// Octagon Financial Intelligence (Assimilated from octagon + octagon-deep-research)

	// Browser Automation Tools (Assimilated from playwright/browser-use/browsermcp/puppeteer/browserbase)

	// ChromaDB Vector Store Tools (Assimilated from chroma-mcp)

	// Basic Memory Tools (Assimilated from basic-memory)

	// MindsDB ML Database Tools (Assimilated from SSE mindsdb)

	// ═══════════════════════════════════════════════════════════════
	// ASSIMILATED MCP SERVERS — Phase 2: Full Native Reimplementation
	// ═══════════════════════════════════════════════════════════════

	// GitHub Copilot API Tools (Assimilated from github SSE)

	// Supabase Tools (Assimilated from supabase SSE)

	// Desktop Commander Tools (Assimilated from @wonderwhy-er/desktop-commander)

	// Gemini API Tools (Assimilated from gemini-mcp)

	// DBHub Universal Database Tools (Assimilated from @bytebase/dbhub)

	// ConPort Context Portal Tools (Assimilated from context-portal-mcp)

	// ChunkHound Code Search Tools (Assimilated from chunkhound)

	// NotebookLM Tools (Assimilated from @roomi-fields/notebooklm-mcp)

	// Vibe Check Tools (Assimilated from @pv-bhat/vibe-check-mcp)

	// SuperMemory Tools (Assimilated from mcp-supermemory-ai)

	// Probe Code Search Tools (Assimilated from @probelabs/probe)

	// Cipher Memory Aggregator Tools (Assimilated from @byterover/cipher)

	// DeepContext Code Understanding Tools (Assimilated from @wildcard-ai/deepcontext)

	// Windows MCP Tools (Assimilated from windows-mcp)

	// Prism Code Quality Tools (Assimilated from prism-mcp-server)

	// TaskMaster AI Task Management Tools (Assimilated from task-master-ai)

	// ═══════════════════════════════════════════════════════════════
	// SKILL REGISTRY - Database-backed skill management with deduplication
	// ═══════════════════════════════════════════════════════════════

	// Skill Registry Tools

	// OpenMemory — local persistent memory store

	// AutoMem — graph-vector memory for AI agents

	// lsmcp — LSP code manipulation and analysis

	// CodeAlive — semantic code search and context engine

	// Prometheus MCP — monitoring queries

	// Smart-Thinking — graph-based reasoning

	// Mimir — Neo4j-backed persistent memory

	// Sysmon — system monitoring

	// Docker — container management

	// Social — Twitter/X and Reddit

	// Git — repository operations

	// Terraform — infrastructure management

	// Google News — news headlines and search

	// OpenRouter Deep Research

	// Prompt Library — SQLite-backed prompt storage

	// Context Server — SQLite-backed context management

	// WebPeel — web data extraction

	// Omnisearch — universal search

	// Grants — government grants discovery

	// Food Data Central — USDA nutrition database

	// Panther — security monitoring

	// Srclight — code indexing for AI agents

	// Coolify — deployment & infrastructure management

	// Harness Integrations

	// Bobbybookmarks Integration

}

func (r *Registry) Execute(ctx context.Context, name string, args map[string]interface{}) (ToolResponse, error) {
	r.mu.RLock()
	handler, ok := r.handlers[name]
	r.mu.RUnlock()
	if !ok {
		return ToolResponse{}, fmt.Errorf("tool handler not found for: %s", name)
	}
	return handler(ctx, args)
}

func (r *Registry) HasTool(name string) bool {
	r.mu.RLock()
	_, ok := r.handlers[name]
	r.mu.RUnlock()
	return ok
}

// List returns all registered tool names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		result = append(result, name)
	}
	return result
}
