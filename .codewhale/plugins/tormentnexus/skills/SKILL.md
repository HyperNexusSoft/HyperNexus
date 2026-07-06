---
name: tormentnexus
description: TormentNexus AI control plane integration — persistent L2 vector memory, semantic tool discovery, session import, skill registry, code search, subagent orchestration
metadata:
  short-description: TormentNexus AI control plane integration
---

# TormentNexus Integration Skill

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry. CodeWhale connects to it via MCP stdio (configured in `~/.codewhale/mcp.json`).

## MCP Server

The tormentnexus MCP server runs as:

- **Command**: `<workspace>/tormentnexus/tormentnexus.exe mcp`
- **Transport**: stdio (JSON-RPC over stdin/stdout)
- **Server name**: `tormentnexus` (tools prefixed `mcp_tormentnexus_*`)
- **Config**: `~/.codewhale/mcp.json`

## Available MCP Tools

The server exposes 49+ tools. Key categories:

### File & System
- `mcp_tormentnexus_read`, `write`, `edit` — file I/O
- `mcp_tormentnexus_grep`, `find`, `ls` — search
- `mcp_tormentnexus_bash` — shell execution
- `mcp_tormentnexus_repomap` — repo map generation

### Window & Automation
- `mcp_tormentnexus_click_chat_button`, `set_chat_input`, `submit_chat_input`, `advance_chat` — UI automation
- `mcp_tormentnexus_detect_chat_surface`, `detect_chat_state` — surface detection
- `mcp_tormentnexus_simulate_input`, `inspect_window_ui`, `list_processes`, `kill_process` — system control

### Memory & Knowledge
- `mcp_tormentnexus_memory_scratchpad_*` — L1 scratchpad (get/set/append)
- `mcp_tormentnexus_memory_extract_relations` — graph extraction
- `mcp_tormentnexus_add_bookmark` — bookmark storage

### MCP Routing
- `mcp_tormentnexus_mcp_list_servers`, `mcp_list_tools`, `mcp_call_tool`, `mcp_server_test`, `mcp_status` — route through TN's Go sidecar to 20+ downstream MCP servers

### Integrations
- `mcp_tormentnexus_jira_create_issue`, `confluence_search`, `cloud_troubleshoot`, `generate_devops_pipeline`, `install_mcp_server` — enterprise tools
- `mcp_tormentnexus_code_interpreter`, `download_llamafile`, `system_status`, `billing_status`, `get_system_stats` — utility tools

## Recommended Usage

### Before significant work
Use `mcp_tormentnexus_memory_scratchpad_get` to check for relevant L1 context, then search via Go sidecar MCP tools for deeper L2 memory.

### During development
Use `mcp_tormentnexus_repomap` for codebase orientation. Use `mcp_tormentnexus_mcp_call_tool` to route through TN's sidecar for tools not directly exposed.

### After decisions
Use `mcp_tormentnexus_memory_scratchpad_set` or `mcp_tormentnexus_add_bookmark` to persist key decisions and patterns for cross-session recall.

## Configuration

The MCP server is already configured if you ran `install_services.bat`. To verify:

```
codewhale mcp list
codewhale mcp connect
codewhale mcp tools | grep tormentnexus
```
