import os

memory_path = r"c:\Users\hyper\workspace\borg\MEMORY.md"

new_observation = """

## Multi-Agent Systemic Observation (2026-06-04) - v1.0.0-alpha.103

1. **Parallel Batch Validation Throughput & Stability**:
   - Running `WORKERS=6` and `BATCH_SIZE=50` (300 servers per run) completes in approximately 10-15 minutes on Windows.
   - Using 3-second pre-flight checks captures early command/module resolution failures instantly, dramatically reducing timeouts from 60s to 3s for misconfigured servers.
   - Detected several port conflicts (e.g., `EADDRINUSE 0.0.0.0:4100`) from servers that try to bind to the control plane ports or already bound ports.
   - Successfully verified 17 new servers and 295 new tools during the latest runs, scaling the registered registry to 788 servers and 11,066 tools.
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

