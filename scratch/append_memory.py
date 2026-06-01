import os

memory_path = r"c:\Users\hyper\workspace\borg\MEMORY.md"

new_observation = """

## Multi-Agent Systemic Observation (2026-06-01) - v1.0.0-alpha.87

1. **Interactive OAuth Handshakes & Connection Timeouts**:
   - Many remote/serverless MCP servers require interactive browser-based OAuth authorization, which hangs indefinitely in automated pipelines.
   - **Resolution**: Enforce a strict 60-second connection deadline using `Promise.race()` during initial connection handshake. This allows interactive authorizations to fail gracefully as timeouts and lets the sequential bulk validator proceed down the backlog queue without manual intervention.

2. **SQLite Unique Name Constraints**:
   - Backlog entries can share identical display names (e.g., `"Developer Utilities"`), which triggers `UNIQUE constraint failed: mcp_servers.name, mcp_servers.user_id` database write failures during automated registration.
   - **Resolution**: Catch uniqueness constraints gracefully in database registration helper blocks, logging the collision without crashing the batch validation loop.
"""

# Read existing content in utf-16-le
with open(memory_path, "r", encoding="utf-16le", errors="replace") as f:
    existing = f.read()

# Append new observation
updated = existing + new_observation

# Write back in utf-16-le
with open(memory_path, "w", encoding="utf-16le") as f:
    f.write(updated)

print("Successfully appended new systemic observations to MEMORY.md!")
