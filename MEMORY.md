# Memory: TormentNexus Kernel & TormentNexus Dashboard

## Core Architectural Observations
- **Go Kernel (TormentNexus)**: Port 4300. Ground truth for state, memory (L1/L2), and tool routing. Standardizes on `modernc.org/sqlite` for CGo-free persistence.
- **TypeScript Dashboard (TormentNexus)**: Port 3000. Observation deck using tRPC to communicate with the kernel.
- **MCP Parity**: The system prioritizes Go-native reimplementation of MCP servers to reduce overhead and eliminate runtime dependencies (Node/Python/Docker).
- **Assimilation Strategy**: 
    1. Discovery & Ranking (Top 500)
    2. Native Reimplementation (Go)
    3. Registration in `registry.go`
    4. Submodule/Dependency Removal

## Technical Findings
- **Skill Deduplication**: Using Jaccard similarity (90% threshold) to fold duplicates into canonical entries. 70-89% similarity creates a relational link via `canonical_id`.
- **Progressive Disclosure**: Tools and Skills are loaded in 3 tiers (Manifest -> Summary -> Full Content) to protect LLM context windows.
- **Windows EBUSY Fix**: Next.js build cleanup now renames the `.next` directory before purging to bypass directory locks.

## Skill Assimilation Stats (v1.0.0-alpha.128)
- **Source directories scanned**: 7 (`~/.a5c`, `~/.agent/skills`, `~/.ccs`, `~/.hermes/skills`, `~/.pi`, `~/.agents/skills`, `~/.config/opencode-temp/skills`)
- **Total SKILL.md found**: 3,418
- **Unique skills assimilated**: 3,229 into `~/.tormentnexus/skills/<id>/SKILL.md`
- **Duplicates merged**: 2 (content-hash dedup)
- **Errors**: 0
- **Script**: `data/assimilate_skills.py`
- **Enriched frontmatter**: name, source, category, date, tags
- **Verification**: All skill registry tests pass (`TestSkillSearch`, `TestSkillDecision`, `TestSkillsFallBackToLocalSkillRegistry`)

## Design Preferences
- **Snake Case in DB, Pascal Case in UI**: Maintain idiomatic Go naming in the backend while mapping to dashboard-friendly formats.
- **Autopilot Protocol**: Continuous execution, atomic commits, and automatic handoff documentation.

## Active Assimilation Tracks
- **Track A**: MCP Servers
- **Track B**: Skill Registry — ✅ 3,229 assimilated
- **Track C**: Hermes Addons
- **Track D**: Prompt Library

## FreeLLM A2A Integration Notes
- Skills are `AgentSkill` structs (ID, Name, Description, Tags)
- Declared in `/.well-known/agent-card`
- Resolved by `findAgentForSkill(skillID)` during swarm dispatch
- TormentNexus catalog.db stores scraped MCP servers in partitioned tables:
  - `published_mcp_servers` — raw metadata + status
  - `published_mcp_config_recipes` — execution templates
  - `published_mcp_validation_runs` — verification outcomes
