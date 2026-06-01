import os

memory_path = r"c:\Users\hyper\workspace\borg\MEMORY.md"

new_observation = """

## Multi-Agent Systemic Observation (2026-06-01) - v1.0.0-alpha.86

1. **Rogue Node/TormentNexus DB Locks on Windows**:
   - Multiple background `node.exe` and `tormentnexus.exe` processes can remain active and hold handles on `tormentnexus.db`, blocking Git actions (e.g. unlinking/writing files during rebases or pulls).
   - **Resolution**: Use `taskkill /F /IM node.exe` and `taskkill /F /IM tormentnexus.exe` in Windows environments to fully terminate lock holders before executing Git updates or database migrations.
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
