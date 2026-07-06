# tormentnexus

TormentNexus pi extension — connect pi to the TormentNexus local AI control plane.

```
pi install npm:tormentnexus
```

## What It Gives You

### 9 Custom Tools
| Tool | What It Does |
|------|-------------|
| `tn_memory_store` | Save decisions, patterns, facts to L2 vault |
| `tn_memory_search` | Find memories by keyword, tag, or category |
| `tn_memory_vector_search` | Semantic vector search across L2 |
| `tn_tool_search` | Find MCP tools across 20+ servers |
| `tn_session_search` | Browse imported sessions |
| `tn_skill_manage` | Access 5,776+ skill modules |
| `tn_code_search` | Search code by AST, semantics, or pattern |
| `tn_context_harvest` | Pull relevant L2 context |
| `tn_scratchpad` | Read/write L1 session scratchpad |

### 6 Slash Commands
`/tn-store`, `/tn-search`, `/tn-status`, `/tn-plan`, `/tn-purge`, `/tn-summary`

### Keyboard Shortcuts
`Ctrl+Shift+M` — Memory search · `Ctrl+Shift+T` — Tool search · `Ctrl+Shift+P` — Status

### Automatic
- Session priming with L2 context · Per-turn harvesting · Enterprise RBAC
- Auto-logging · Compaction enrichment · Live stats widget

## Requirements
TormentNexus Go sidecar running on `http://127.0.0.1:7778`

## License MIT
