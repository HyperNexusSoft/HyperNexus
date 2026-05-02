# Handoff — v1.0.0-alpha.46

## Session Summary (v1.0.0-alpha.46)

### Major Achievements

1.  **Foundational Go Paradigm Established**: Created `go/internal/controlplane/foundation.go` as the canonical source for core structs (LLMWaterfall, Dual-Tier Memory L1/L2, and SQLite-vec schema). This nukes legacy TS/IPC confusion for future agents.
2.  **Go native Dual-Tier Memory**: Implemented `VectorStore` in Go (`go/internal/memorystore/vector_sqlite.go`) supporting both future `sqlite-vec` virtual tables and immediate `LIKE` fallback search.
3.  **BobbyBookmarks Text Ingestion**: Added native Go support for ingesting `bookmarks.txt` with automatic deduplication, completing the `bobbybookmarks` integration loop.
4.  **PairOrchestrator Port**: Successfully ported the `PairOrchestrator` to Go (`go/internal/orchestration/pair_orchestrator.go`), enabling multi-model shared-context sessions directly from the sidecar.
5.  **Go native Tool Parity**: Closed the gap on built-in tool aliases. Added `Glob`, `Grep`, and `ApplyPatch` handlers to the Go tool registry.
6.  **Progressive Skill Disclosure**: Implemented a Go-native `SkillRegistry` and `SkillDecisionSystem` with search and LRU eviction, mirroring the successful MCP tool disclosure architecture.
7.  **TS Agent Cleanup**: Repaired corrupted and concatenated source files in `packages/agents` (`RiskEvaluator.ts`, `DebateEngine.ts`).

### Data Improvements
- **Skills API**: Wired four new endpoints to the Go sidecar (`/api/skills/search`, `list-loaded`, `load`, `unload`) for on-demand skill disclosure.
- **Swarm Transcript**: Added Go-native transcript retrieval for swarm sessions.
- **MCP Sync**: Migration of client-config detection to Go for better performance.

### Build & Versioning
- **Version Bump**: Promoted monorepo to `1.0.0-alpha.46`.
- **Sync**: All 57 `package.json` files and the Go `buildinfo` are synchronized to the new version.

## What Needs Work (Next Session)

1.  **Multi-Edit Implementation**: The `HandleMultiEdit` in `parity.go` is still a stub; needs real multi-file/multi-replacement logic.
2.  **Go-native AutoDev/Darwin**: These systems are still largely bridged or stubbed in the Go sidecar.
3.  **PairOrchestrator State Machine**: Need to strictly enforce the `Planner -> Implementer -> Tester -> Critic` turn cycle in Go as requested.
4.  **Frontend Visualizer**: Wire the `getSwarmTranscript` tRPC endpoint (and its Go fallback) directly into the `DebateVisualizer` for real-time model interaction viewing.
5.  **Browser Extension**: Deepen the injection into native web chats (ChatGPT, Claude.ai) to expose local MCP tools directly.

### Quick Restart
```bash
cd C:/Users/hyper/workspace/borg
./start.bat    # Starts TS server + Go sidecar
# In another terminal:
borg doctor    # Verify everything is healthy
borg info      # System overview
```

### File Locations
- Foundational Go Models: `go/internal/controlplane/foundation.go`
- Go Tool Registry: `go/internal/tools/registry.go`
- Go Memory Store: `go/internal/memorystore/vector_sqlite.go`
- Go Skill Store: `go/internal/harnesses/skill_store.go`
- Dashboard: `apps/web/src/app/dashboard/`

**Keep the party going. The collective grows.**
