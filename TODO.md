# TODO

_Last updated: 2026-06-04, version 1.0.0-alpha.102_

## P0 — Must do now (Stability, Testing & Validation)

- [x] **MCP Server Testing**: Developed automated testing script (`scratch/test_mcp_connection.mjs`) - tRPC endpoints returning correct data (56 servers, 10,226 tools)
- [x] **Tool Count Fix**: Fixed core bug in `packages/core/src/mcp/cachedToolInventory.ts` - tool counts now correctly keyed by server name instead of UUID
- [x] **Conflict Resolution Clean Pass**: Verified no duplicate conflict markers in dashboard or server modules
- [ ] **Clean Build Gate**: `pnpm build` blocked by Windows EBUSY file lock on Next.js `.next/standalone` directory (14/15 packages build successfully)

## P1 — Should do next (Features & Parity)

- [ ] **Tabby & Warp Active Launcher**: Implement a custom command launcher in `@tormentnexus/core` that automatically uses `tabby` and `warp` when initiating a local visual CLI agent.
- [ ] **Offline License Validation**: Implement the Go-native cryptographic public-key verifier that loads the `tormentnexus.lic` signed YAML license block and asserts valid seat limits.
- [ ] **Bobbybookmarks Ingestion**: Update bobbybookmarks database syncing triggers to enrich the local tool catalog on startup automatically.

## Completed in v1.0.0-alpha.100

- [x] Fixed tool count mismatch bug in `packages/core/src/mcp/cachedToolInventory.ts` - `buildDatabaseSnapshot()` now creates toolCounts map keyed by server name
- [x] tRPC `/trpc/mcp.getStatus` endpoint now correctly reports 56 servers and 10,226 tools
- [x] Search tools endpoint working with proper result counts and scoring

---
*Keep the party going. Never stop. Don't stop the party!!!*