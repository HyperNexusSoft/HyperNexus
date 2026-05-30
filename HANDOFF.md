# Handoff - v1.0.0-alpha.76

## Summary
Upgraded the **Predictive Tool Disclosure (Tool Ads)** engine. It now generates dynamic, context-aware MCP tool recommendations utilizing the core LLMService and ModelSelector, prioritizing free OpenRouter models with cascading fallbacks to local LMStudio endpoints, and keeps a fail-safe backup routing to the Go sidecar.

## Accomplishments

### LLM-Based Tool Predictions (v1.0.0-alpha.76)
- **`getPredictedToolAds` Realization**: Re-implemented prediction mapping inside `MCPServer.ts` to utilize the core LLM execution block.
- **Cascading Fallback Routing**: Integrates with the cheap/free OpenRouter routing strategies, falling back gracefully to alternative models and finally routing to local LMStudio/Ollama instances.
- **Failsafe Backup**: Retains a robust secondary route to the Go sidecar (`http://127.0.0.1:4300/api/mcp/tools/predict`).
- **Prompt Injection integration**: Dynamically merges the predicted tools straight into the worker agent's tool set (`McpWorkerAgent.ts`), embedding details in the system prompt context.

## Current State
- `published_mcp_servers` in `borg.db`: **28,534 rows**
- `published_mcp_config_recipes` in `borg.db`: **27,553 rows**
- VERSION: `1.0.0-alpha.76`
- Package Synchronization: All 27 packages successfully synced to `1.0.0-alpha.76`.

## Next Steps
1. Execute a Swarm Worker task and monitor prompt payloads to verify predicted tool descriptions are successfully injected as system context nudges.
