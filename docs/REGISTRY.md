# TormentNexus Registry & Catalog Documentation

> Last updated: 2026-07-13 | Version: 1.0.0-b2

## Overview

TormentNexus aggregates MCP tools, skills, and prompts from multiple sources into a unified catalog. This document maps every registry, database, and scraper script in the project.

---

## 1. MCP Tool Registries

### 1.1 Go Native Handlers (`go/internal/mcpimpl/registry.go`)

- **Count:** 5,668 registered tool handlers
- **Source:** Hardcoded Go implementations for each MCP server
- **How to update:** Add new handler entries to `AllHandlers()` function
- **Active tools:** 51 currently registered and serving via kernel API

### 1.2 Published Catalog (`catalog.db` → `published_mcp_servers`)

- **Count:** 8 native Go tools
- **Source:** Auto-registered from Go native handlers at kernel boot
- **Schema:** `uuid, canonical_id, display_name, description, tags, categories, transport, status, created_at, updated_at`

### 1.3 Links Backlog (`catalog.db` → `links_backlog`)

- **Count:** 3,200 entries
- **Sources:**
  - `awesome-mcp-servers` (3,042) — scraped from `punkpeye/awesome-mcp-servers` GitHub repo
  - `github-search` (155) — GitHub topic search for `mcp+server`
  - `glama-mock` (3) — hardcoded fallback presets
- **Schema:** `uuid, url, normalized_url, title, description, tags, source, is_duplicate, duplicate_of, research_status, http_status, synced_at, created_at, updated_at`
- **Scraper:** `scripts/scrape-mcp-servers.py`

### 1.4 Published Skills (`catalog.db` → `published_skills`)

- **Count:** 0 (table exists but not populated)
- **Schema:** TBD — needs to be populated from `.tormentnexus/skills/`

---

## 2. Skill Libraries

### 2.1 Local Skills (`.tormentnexus/skills/`)

- **Count:** 2,987 SKILL.md files
- **Location:** `C:/Users/hyper/workspace/tormentnexus/.tormentnexus/skills/`
- **Format:** Each skill is a directory with a `SKILL.md` file
- **Categories:** Includes 000-999 numbered skills + named skills
- **Examples:**
  - `000_jeremy_content_consistency_validator`
  - `2d_games`, `3d_games`, `3d_web_experience`
  - `6sense_intent`, `a3_problem_solver`
  - `abc_xyz_classifier`, `abstract_domain_library`

### 2.2 Agent Extension Skills

- **Locations:** `.agent/skills/`, `.antigravity/extensions/`, `.claude/plugins/`, `.codewhale/plugins/`, `.codex/skills/`, `.gemini/skills/`, `.kimi-code/skills/`, `.mavis/skills/`
- **Count:** ~10-20 across all locations
- **Format:** Same SKILL.md format

### 2.3 Pi Skills (Installed via npm)

- **Location:** `~/.pi/agent/npm/node_modules/*/skills/`
- **Count:** ~10-15 (context-management, pi-subagents, prompt-template-authoring, etc.)
- **Discovery:** Auto-detected by Pi agent at startup

---

## 3. Prompt Templates

### 3.1 Pi Prompt Templates

- **Location:** `~/.pi/agent/npm/node_modules/*/prompts/`
- **Format:** `.md` files with prompt template syntax
- **Examples:** `/writing-test`, `/code-review`, `/execute-plan`

### 3.2 Project Prompt Templates

- **Location:** `.prompts/` or `prompts/` directories
- **Count:** TBD — check `.prompts/` directory

---

## 4. Database Inventory

### 4.1 Kernel Databases (`/root/.tormentnexus/`)

| Database | Tables | Purpose |
|---|---|---|
| `memory.db` | 56 | L2 vault, L3 cold archive, L4 limbo, GraphRAG, FTS5 |
| `catalog.db` | 3 | MCP tool catalog, skills, links backlog |
| `accounts.db` | 3 | Tenant accounts and billing |
| `l3_cold_archive.db` | 14 | Cold storage for old memories |

### 4.2 Docker Tenant Databases

- Each tenant has its own `tormentnexus.db` in the container volume
- Located at `/var/lib/hypernexus/tenants/{tenant}/`

### 4.3 Other Databases

- `/opt/tormentnexus/tormentnexus.db` — Session import data (2 tables)
- `/opt/tormentnexus/catalog.db` — Duplicate catalog (kernel creates in CWD)

---

## 5. External Registry Sources

### 5.1 Active Sources (Working)

| Source | URL | Status | Last Scrape |
|---|---|---|---|
| awesome-mcp-servers | `github.com/punkpeye/awesome-mcp-servers` | ✅ Working | 2026-07-13 |
| GitHub Topics | `api.github.com/search/repositories?q=mcp+server` | ✅ Working | 2026-07-13 |
| Official MCP Servers | `github.com/modelcontextprotocol/servers` | ✅ Working (7 dirs) | Manual |

### 5.2 Broken Sources (Need Fix)

| Source | URL | Issue |
|---|---|---|
| Glama.ai | `glama.ai/api/v1/mcp/servers` | Returns HTML instead of JSON (API changed) |
| Smithery | `registry.smery.ai` | Empty response |
| mcp.run | `mcp.run/api/catalog` | 301 redirect (API changed) |

### 5.3 Sources Not Yet Integrated

| Source | URL | Notes |
|---|---|---|
| npm MCP packages | `registry.npmjs.org` | Search for `mcp-server` keyword |
| PyPI MCP packages | `pypi.org` | Search for `mcp-server` keyword |
| MCP Hub | `mcp-hub.com` | Community directory |
| Toolhouse | `toolhouse.ai` | Agent tool marketplace |
| Composio | `composio.dev` | Integration platform |

---

## 6. Scraper Scripts

### 6.1 `scripts/scrape-mcp-servers.py`

- **Purpose:** Bulk scrape MCP servers from awesome-mcp-servers and GitHub topics
- **Output:** Inserts into `catalog.db` → `links_backlog` table
- **Run:** `python3 scripts/scrape-mcp-servers.py`
- **Last run:** 2026-07-13 (3,200 entries)

### 6.2 `scripts/gen-blog.py`

- **Purpose:** Generate blog post HTML files for tormentnexus.site
- **Output:** HTML files in `/var/www/tormentnexus.site/blog/tormentnexus/`
- **Run:** `python3 scripts/gen-blog.py`

### 6.3 Go Catalog Sync (`go/internal/hsync/glama.go`)

- **Purpose:** Sync from Glama.ai API (currently broken — API changed)
- **Run:** Automatic at kernel boot via `CatalogSync` goroutine
- **Fallback:** Uses 3 hardcoded presets when API fails

### 6.4 Go Catalog Ingest (`go/internal/mcp/catalog_ingest.go`)

- **Purpose:** Ingest from multiple sources via adapter pattern
- **Adapters:** `GlamaAiAdapter` (only one currently)
- **Run:** Called from kernel initialization

---

## 7. How to Add a New Registry Source

### Option A: Python Scraper (Quick)

1. Add a new function to `scripts/scrape-mcp-servers.py`
2. Fetch from the API/URL
3. Parse and normalize entries
4. Insert into `catalog.db` → `links_backlog`
5. Run the script

### Option B: Go Adapter (Production)

1. Create a new adapter in `go/internal/mcp/catalog_ingest.go`
2. Implement `CatalogSourceAdapter` interface (`Name()` + `Ingest()`)
3. Add to `adapters` slice in `IngestPublishedCatalog()`
4. Rebuild kernel

---

## 8. Sync Schedule

| Task | Frequency | Script |
|---|---|---|
| MCP server scrape | Manual / on-demand | `scripts/scrape-mcp-servers.py` |
| Glama.ai sync | Every 60min (kernel) | `go/internal/hsync/glama.go` |
| Catalog ingest | At kernel boot | `go/internal/mcp/catalog_ingest.go` |
| DB backup | Daily 3am | `/etc/cron.d/tn-backup` |
| Memory maintenance | Every 60min (kernel) | Go kernel goroutine |

---

## 9. Quick Reference Commands

```bash
# Check catalog counts
ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT source, count(*) FROM links_backlog GROUP BY source;"'

# Run MCP scraper
python3 scripts/scrape-mcp-servers.py

# Check active tools
curl -s http://127.0.0.1:8090/api/runtime/status | jq '.data.cli.toolCount'

# List skills
find .tormentnexus/skills/ -name "SKILL.md" | wc -l

# Check Stripe config
ssh hetzner 'grep STRIPE /opt/tormentnexus/.env'
```
