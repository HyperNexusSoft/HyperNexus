# HyperNexus — Universal AI Control Plane

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript)
![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js)
![License](https://img.shields.io/badge/license-Enterprise%2FOSS-orange)

> **HyperNexus** is the ultimate local-first control plane for multi-agent workflows, Model Context Protocol (MCP) tooling, provider routing, session continuity, and operator observability. It gives your AI persistent memory across 38+ tools.

---

## 🚀 Quick Start

### One-Command Install

```bash
# Windows (PowerShell)
irm https://hypernexus.site/install.ps1 | iex

# macOS / Linux
curl -fsSL https://hypernexus.site/install.sh | bash

# npm
npx @hypernexus/install@latest
```

### Download

| Platform | Link |
|----------|------|
| **Windows** | [hypernexus.exe](https://releases.hypernexus.site/latest/hypernexus.exe) |
| **macOS** | [Install Script](https://releases.hypernexus.site/latest/install.sh) |
| **Linux** | [Install Script](https://releases.hypernexus.site/latest/install.sh) |
| **All Downloads** | [hypernexus.site/download](https://hypernexus.site/download.html) |

### Start

```bash
hypernexus serve
```

Open <http://localhost:7779/dashboard>

---

## 💰 Pricing

| Plan | Price | Includes |
|------|-------|----------|
| **Community** | Free | Local binary, open source |
| **Professional** | $50/seat/year | Local license + cloud hosting |

- **Cloud Dashboard**: <https://cloud.hypernexus.site>
- **Product Hunt**: <https://www.producthunt.com/products/hypernexus?launch=hypernexus>

---

## 🤖 Supported AI Clients (38+)

HyperNexus automatically configures MCP for all these tools:

### Tier 1: Major AI Tools

| Client | Config Location |
|--------|-----------------|
| Claude Desktop | `%APPDATA%/Claude/claude_desktop_config.json` |
| Claude Code | `~/.claude/settings.json` |
| Cursor | `%APPDATA%/Cursor/User/mcp.json` |
| VS Code | `%APPDATA%/Code/User/settings.json` |
| Windsurf | `%APPDATA%/Windsurf/User/mcp.json` |
| Gemini | `~/.gemini/settings.json` |
| GitHub Copilot | `%APPDATA%/GitHub Copilot/mcp.json` |

### Tier 2: Popular Coding Assistants

| Client | Config Location |
|--------|-----------------|
| Aider | `~/.aider/mcp.json` |
| Continue | `~/.continue/config.json` |
| Cline | VS Code extension storage |
| Roo | `~/.roo/mcp.json` |
| OpenHands | `~/.openhands/mcp.json` |
| Factory | `~/.factory/mcp.json` |
| CodeWhale | `~/.codewhale/mcp.json` |

### Tier 3: Emerging Tools

| Client | Config Location |
|--------|-----------------|
| Goose | `~/.goose/mcp.json` |
| IFlow | `~/.iflow/mcp.json` |
| OpenCode | `~/.opencode/mcp.json` |
| OpenClaw | `~/.openclaw/mcp.json` |
| Antigravity | `~/.antigravity/mcp.json` |
| Trae | `~/.trae/mcp.json` |
| Zed | `~/.config/zed/settings.json` |
| Kiro | `~/.kiro/mcp.json` |

### Tier 4: Specialized Tools

| Client | Config Location |
|--------|-----------------|
| Codex | `~/.codex/mcp.json` |
| Grok | `~/.grok/mcp.json` |
| Qwen | `~/.qwen/mcp.json` |
| Qwen Code | `~/.qwen-code/mcp.json` |
| Kimi Code | `~/.kimi-code/mcp.json` |
| Moonshot | `~/.moonshot/mcp.json` |
| Pi | `~/.pi/mcp.json` |
| Hermes | `~/.hermes/mcp.json` |

### Tier 5: Enterprise/Team Tools

| Client | Config Location |
|--------|-----------------|
| Omnigent | `~/.omnigent/mcp.json` |
| Citadel | `~/.citadel/mcp.json` |
| Agent Fusion | `~/.agent-fusion/mcp.json` |
| Herdr | `~/.herdr/mcp.json` |
| Claude Squad | `~/.claude-squad/mcp.json` |
| Cliproxyapi | `~/.cliproxyapi/mcp.json` |
| JetBrains | `%APPDATA%/JetBrains/mcp.json` |

---

## 🏢 Corporate Mode

HyperNexus supports corporate mode for commercial deployments:

```bash
# Enable corporate mode
export HN_EDITION=corporate
export HN_CLOUD_ENDPOINT=https://cloud.hypernexus.site

# Or set during install
npx @hypernexus/install@latest  # Auto-detects corporate mode
```

Corporate mode includes:

- **Branding**: HyperNexus Corp branding across all UI
- **Cloud Connection**: Automatic connection to cloud.hypernexus.site
- **Auto-Updates**: Daily update checks via scheduled tasks
- **SSO/OIDC**: Enterprise authentication support
- **Audit Trails**: Full activity logging

---

## 🔄 Auto-Updates

HyperNexus includes automatic update checking:

```bash
# Check for updates
hypernexus-update --check

# Force update
hypernexus-update --force
```

**Windows**: Scheduled task runs daily at 9 AM  
**Linux/macOS**: Cron job runs daily at 9 AM

Updates are pulled from: `https://releases.hypernexus.site/latest/version.json`

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  OPERATOR LAYER                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │  Web Dash   │  │  CLI (TS)   │  │  VS Code    │       │
│  │  Port 7779  │  │  hypernexus │  │  Extension  │       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│         │                │                │                 │
│         └────────────────┴────────────────┘                 │
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  GO SIDECAR (Port 7778) — The Authoritative Kernel    │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │  │ SkillStore│ │ EventBus │ │  Vault  │ │ Healer  │     │
│  │  │ (BM25)  │ │ (Swarm) │ │(sqlite) │ │(Immune) │     │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │  │  Router  │ │ PairOrchestrator│ │ CodeExecutor│ │ MCP Sync │     │
│  │  │(Progressive)│ │(Consensus) │ │(Sandbox) │ │(Registry)│     │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  └──────────────────────────────────────────────────────┘
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  EXTERNAL MODELS & TOOLS                     │        │
│  │  OpenAI · Anthropic · Gemini · OpenRouter · Ollama  │
│  │  600+ MCP Servers · 3,900+ Native Go Tools          │
│  └──────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
```

### Key Ports

| Service | Port | Purpose |
|---------|------|---------|
| Go Sidecar | 7778 | Authoritative native kernel |
| Next.js Dashboard | 7779 | Web observation deck |
| Cloud Server | 7780 | HyperNexus Cloud API |

---

## 🧠 Core Features

### 1. Progressive MCP Tool Routing

Models should never be overwhelmed with a 50,000-token tool dump. HyperNexus employs a multi-layered, progressive disclosure system:

- **Semantic Search**: Local vector embeddings match the active prompt against a global MCP directory of 14,250+ servers.
- **The Router**: Only the top highly relevant tool schemas are injected into the active LLM context.
- **Universal Parity**: Byte-for-byte identical tool signatures for Claude Code, Codex, Gemini CLI, Cursor, and Windsurf.

### 2. Dual-Tier Memory Architecture (L1 / L2)

Context is finite; memory must be infinite.

- **L1 — Session Scratchpad**: Ephemeral, lightning-fast memory tied directly to the active session.
- **L2 — The Vault**: Permanent semantic storage in SQLite with `sqlite-vec`. Saves exact transcripts and LLM-compressed heuristics.
- **Context Harvesting**: Every session autonomously queries the L2 Vault to pull in relevant historical heuristics.
- **Heat Mechanics**: Relevance increases heat, time causes decay — biological memory modeling.

### 3. The Resilient LLM Waterfall

Uptime is non-negotiable. HyperNexus's inference client natively catches 429s (Rate Limits) and 5xx (Server Errors), seamlessly cascading the exact payload down a prioritized chain without crashing:

1. **NVIDIA NIM** / Primary APIs
2. **OpenRouter** (Secondary aggregator fallback)
3. **Local LM Studio / Ollama** (Ultimate offline fallback)

### 4. Multi-Agent Swarm & P2P Mesh

HyperNexus coordinates specialized models inside shared chatrooms via the Agent-to-Agent (A2A) protocol.

- **Role Rotation**: Models take turns acting as Planner, Implementer, Tester, and Critic.
- **Consensus & Debate**: Agents autonomously bid on tasks, share context via a neural transcript, and debate implementations until consensus is reached.
- **PairOrchestrator**: Enforces a strict `Planner → Checker → Implementer → Critic` state machine with weighted consensus scoring.

### 5. Autonomous Immune System

Every failure is an opportunity for diagnosis, remediation, and verification.

- **HealerService**: Multi-turn `Diagnose → Fix → Verify → Retry` loop using the native `CodeExecutor`.
- **L2 Vault Persistence**: All healing events and extracted facts are saved as long-term memory for fleet-wide intelligence sharing.
- **Supervisor Nudge Protocol**: Autonomously maintains development momentum by re-engaging inactive agents through professional, context-aware directives.

---

## 📦 Monorepo Structure

```
HyperNexus/
├── go/                          # Go sidecar (the authoritative kernel)
│   ├── cmd/
│   │   ├── hypernexus/          # Main agent binary
│   │   └── cloud/               # Cloud server binary
│   ├── internal/
│   │   ├── config/              # Branding, discovery, config
│   │   ├── httpapi/             # HTTP API handlers
│   │   ├── mcp/                 # MCP client sync (38+ clients)
│   │   ├── memorystore/         # L2 vault implementation
│   │   └── ...
│   └── ...
├── apps/
│   └── web/                     # Next.js dashboard
├── packages/
│   ├── ui/                      # React components
│   ├── memory/                  # Memory services
│   ├── agents/                  # Agent orchestration
│   └── ...
├── npm/
│   └── @hypernexus/
│       └── install/             # npm installer
├── dist/
│   ├── install.ps1              # Windows installer
│   ├── install.sh               # Linux/macOS installer
│   ├── hypernexus-update        # Linux/macOS updater
│   └── hypernexus-update.ps1    # Windows updater
├── .gitlab-ci.yml               # CI/CD pipeline
└── README.md
```

---

## 🚀 Installation

### Option 1: One-Command Install (Recommended)

```bash
# Windows (PowerShell)
irm https://hypernexus.site/install.ps1 | iex

# macOS / Linux
curl -fsSL https://hypernexus.site/install.sh | bash

# npm
npx @hypernexus/install@latest
```

### Option 2: Download Binary

| Platform | Link |
|----------|------|
| **Windows** | [hypernexus.exe](https://releases.hypernexus.site/latest/hypernexus.exe) |
| **macOS (Intel)** | [hypernexus-darwin-amd64](https://releases.hypernexus.site/latest/hypernexus-darwin-amd64) |
| **macOS (Apple Silicon)** | [hypernexus-darwin-arm64](https://releases.hypernexus.site/latest/hypernexus-darwin-arm64) |
| **Linux (x64)** | [hypernexus-linux-amd64](https://releases.hypernexus.site/latest/hypernexus-linux-amd64) |
| **Linux (ARM64)** | [hypernexus-linux-arm64](https://releases.hypernexus.site/latest/hypernexus-linux-arm64) |

### Option 3: Build from Source

```bash
git clone https://gitlab.com/hypernexus-org/HyperNexus.git
cd HyperNexus/go
go build -o hypernexus ./cmd/hypernexus
```

---

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HN_EDITION` | `hypernexus` | Set to `corporate` for commercial mode |
| `HN_CLOUD_ENDPOINT` | `https://cloud.hypernexus.site` | Cloud API endpoint |
| `HN_CLOUD_AUTH` | (empty) | Cloud authentication token |
| `HN_UPDATE_URL` | `https://releases.hypernexus.site/latest/version.json` | Auto-update URL |

### Config Directory

```
~/.hypernexus/
├── branding.json        # Branding configuration
├── config.json          # Main configuration
├── mcp.json             # MCP server configuration
├── current_version      # Installed version
└── ...
```

### Corporate Branding

```json
{
  "edition": "hypernexus",
  "product_name": "HyperNexus",
  "company_name": "HyperNexus Corp",
  "tray_tooltip": "HyperNexus (Running)",
  "dashboard_title": "HyperNexus Dashboard",
  "cloud_endpoint": "https://cloud.hypernexus.site",
  "config_dir": ".hypernexus",
  "registry_key": "HyperNexus"
}
```

---

## 🌐 Cloud Deployment

### Hetzner Server

| Service | Port | Status |
|---------|------|--------|
| HyperNexus Agent | 7778 | ✅ Running |
| HyperNexus Cloud | 7780 | ✅ Running |
| Marketing Agent | - | ✅ Running |
| Umami Analytics | 3000 | ✅ Running |

### Endpoints

| URL | Purpose |
|-----|---------|
| <https://hypernexus.site> | Landing page |
| <https://cloud.hypernexus.site> | Cloud dashboard |
| <https://releases.hypernexus.site> | Binary releases |
| **Open Source** | <https://tormentnexus.site> |

---

## 🔧 CI/CD Pipeline

GitLab CI/CD automatically builds on every push to `main`:

| Stage | Builds | Platforms |
|-------|--------|-----------|
| **Build** | Go agent, Cloud server | Linux, Windows, macOS (Intel + ARM) |
| **Test** | Go vet + tests | - |
| **Package** | Installers | .exe, .zip, .tar.gz |
| **Deploy** | Hetzner services | Linux |
| **Release** | Versioned releases | All platforms |

### Triggering a Release

```bash
# Update version
echo "1.0.1" > VERSION

# Commit and push
git add VERSION
git commit -m "bump: v1.0.1"
git push origin main

# GitLab CI/CD will:
# 1. Build all platform binaries
# 2. Create versioned release on releases.hypernexus.site
# 3. Update version.json for auto-updater
# 4. Optionally deploy to Hetzner (manual trigger)
```

---

## 📊 Current Capabilities (v1.0.0)

| Capability | Status | Details |
|------------|--------|---------|
| **MCP Registry** | Stable | 14,250+ tracked MCP servers, 11,024+ populated in SQLite catalog, 600+ verified servers, 11,000+ verified tools |
| **Native Go Tools** | Beta | 3,900+ native Go tool implementations replacing external MCP servers |
| **Progressive Tool Routing** | Stable | Semantic vector search + BM25 ranking injects only the most relevant tools into LLM context windows |
| **Dual-Tier Memory** | Stable | L1 (session scratchpad) + L2 (semantic SQLite vault) with heat-score lifecycle and autonomous context harvesting |
| **LLM Waterfall** | Stable | Cascading failover: NVIDIA NIM → OpenRouter → Local LM Studio / Ollama with 429/5xx handling |
| **Multi-Agent Swarm** | Beta | A2A protocol coordination, role rotation (Planner→Implementer→Tester→Critic), consensus engine |
| **Autonomous Healer** | Stable | `Diagnose → Fix → Verify → Retry` loop with native code execution and L2 vault persistence |
| **Browser Automation** | Beta | Native chromedp handlers: navigate, screenshot, evaluate, click, fill, get HTML |
| **Skill Registry** | Stable | 3,229+ assimilated skills from 7 harness ecosystems with Jaccard deduplication |
| **Dashboard** | Stable | Next.js 16 + React 19 + Tailwind CSS 4 with real-time telemetry, knowledge graph, healer view, swarm visualizer |
| **Session Import** | Beta | Automatic ingestion of Claude, Aider, and other harness session artifacts |
| **Enterprise Licensing** | Experimental | Ed25519-signed license token validation with offline verification |
| **Auto-Updates** | Stable | Daily update checks via scheduled tasks, version.json API |
| **Corporate Mode** | Stable | HyperNexus branding, cloud connection, SSO/OIDC support |
| **38+ AI Clients** | Stable | Automatic MCP configuration for all major AI tools |

---

## 🔗 Links

| Resource | URL |
|----------|-----|
| **Website** | <https://hypernexus.site> |
| **Cloud Dashboard** | <https://cloud.hypernexus.site> |
| **Downloads** | <https://hypernexus.site/download.html> |
| **Pricing** | <https://hypernexus.site/pricing.html> |
| **Product Hunt** | <https://www.producthunt.com/products/hypernexus?launch=hypernexus> |
| **GitLab** | <https://gitlab.com/HyperNexusLLC/HyperNexus> |
| **GitHub** | <https://github.com/HyperNexusLLC/HyperNexus> |
| **TormentNexus (OSS)** | <https://tormentnexus.site> |

---

## 📱 Social Media

| Platform | Handle |
|----------|--------|
| **Twitter/X** | [@hypernexus](https://x.com/hypernexus) |
| **Bluesky** | [@hypernexusllc.bsky.social](https://bsky.app/profile/hypernexusllc.bsky.social) |
| **LinkedIn** | [HyperNexus LLC](https://www.linkedin.com/company/hypernexusllc) |
| **Reddit** | [r/HyperNexusLLC](https://reddit.com/r/HyperNexusLLC) |

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a merge request

---

## 📄 License

- **Community Edition**: MIT License (open source)
- **Professional Edition**: Commercial license ($50/seat/year)

---

## 🙏 Acknowledgments

Built with:

- [Go](https://golang.org/) — Backend kernel
- [Next.js](https://nextjs.org/) — Dashboard
- [React](https://react.dev/) — UI components
- [SQLite](https://sqlite.org/) — Local storage
- [MCP](https://modelcontextprotocol.io/) — Tool protocol

---

**HyperNexus** — Your AI finally remembers everything.
