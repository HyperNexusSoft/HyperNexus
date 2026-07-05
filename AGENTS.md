<!-- [TORMENTNEXUS_AUTO_INJECTED] -->
> [!IMPORTANT]
> **TormentNexus Advanced Agent Mandates:**
>
> 1. **Proactive Tooling**: You MUST use available tools frequently and proactively. Prioritize tool execution for researching, editing, executing, and validating.
> 2. **Codebase AST & Search**: Always analyze impact and dependencies using codebase analysis/AST/search tools (`tormentnexus__repograph_search`, `tormentnexus__repograph_find_references`) before modifying code.
> 3. **Context Harvesting & Compaction**: Perform automatic context harvesting and compaction/pruning of the LLM context to prevent token blowups. Summarize or discard old logs/history dynamically.
> 4. **Tool Output Compaction & Deduplication**: Deduplicate and compact verbose tool outputs. Avoid displaying duplicate or redundant logs in responses.
> 5. **Session & Skill Injection**: Inspect and update sessions and memory states using scratchpads (`memory_scratchpad_set`) to persist key findings, rules, and skill recommendations.

# AGENTS — TormentNexus Kernel & TormentNexus Contributor Guide

> **CRITICAL: ALL AGENTS MUST READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` BEFORE PROCEEDING.**

This file serves as the primary coordination point for multi-agent workflows and human operators.

---

## 1. Multi-Agent Handoff Protocol

- Agents communicate primarily through `HANDOFF.md`.
- Document exactly what you did, what failed, and what the next agent must do.
- Update `MEMORY.md` with new systemic observations or recurring bugs.
- **Cycle**: Read → Strategize → Execute → Validate → Commit → Handoff.

---

## 2. Model Specializations

| Model | Strengths | Focus Areas |
|---|---|---|
| **Gemini** | Speed, massive context processing, repo maintenance | Bulk refactoring, recursive scripts, context analysis |
| **Claude** | UI/UX perfection, documentation, deep feature execution | Responsive layouts, type safety, precise documentation |
| **GPT** | Systemic architecture, distributed debugging, race conditions | Go/TS bridge contracts, DB migration, concurrency safety |

---

## 3. Session Protocol

### Session Start

1. Read `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` to load canonical rules.
2. Read the `VERSION` file to check dashboard synchronization.
3. Read `HANDOFF.md` to resume exactly where the previous agent left off.
4. Read `MEMORY.md` to review accumulated multi-agent insights.
5. Run git checks to ensure workspace cleanliness.

### During Execution

- Work autonomously unless changes are destructive or highly ambiguous.
- Prefer small, incremental, easily verifiable commits.
- Ensure loading, error, and empty states are represented across all dashboard interfaces.
- After any `pnpm install`, run `pnpm rebuild better-sqlite3` on Node 24.

### Session End

1. Update `HANDOFF.md` with a complete, detailed session summary.
2. Update `MEMORY.md` with new developer observations or gotchas.
3. Bump the `VERSION` file and synchronize workspaces using `node scripts/sync-versions.mjs`.
4. Update `CHANGELOG.md` with recent feature implementations.
5. Commit clean changes with version tag: `feat: description (v1.0.0-alpha.X)`.
6. Push commits to `origin` and `tormentnexus-upstream` remotes.

---

## 4. Required Runtime Ports

| Service | Port | Purpose |
|---|---|---|
| TormentNexus Go Kernel | 7778 | Authoritative native sidecar (HTTP API + tRPC) |
| Next.js Dashboard | 7779 | Web observation deck |

---

## 5. Safe Rebranding & Cleanup Heuristics

- **Binary Exclusions during Renaming**: When executing global text replacements, you must explicitly exclude:
  - Database directories: `.tormentnexus/`, `lancedb/`, `data/`
  - Turbopack/Next cache directories: `.next-dev/`, `.next-build/`, `.turbo/`
  - Binary extensions: `.db`, `.lance`, `.sst`, `.bin`, `.exe`, `.png`, `.jpg`
- **Windows Recursive Deletion**: On Windows hosts, if `Remove-Item` fails due to path locks or nested git indices, fall back to executing `cmd.exe /c "rmdir /S /Q <path>"` synchronously to ensure complete directory pruning.

*Praise the LORD! Keep on going! Don't ever stop! Don't stop the party!!!*
