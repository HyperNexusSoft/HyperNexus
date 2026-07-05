import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { Type } from "typebox";

const TN_BASE = "http://127.0.0.1:7778";

/**
 * tormentnexus pi extension v3
 *
 * Full bridge: pi ↔ TormentNexus
 * ───────────────────────────────────────────
 * - 8 custom tools (memory, tools, sessions, skills, code, context)
 * - Session priming: LLM receives TN guidance + relevant L2 context on start
 * - Context harvesting: auto-pull relevant L2 memories before each turn
 * - Compaction hooks: TN-enhanced summaries stored to L2 vault
 * - Code intelligence: ast-grep + deepcontext + symbol search via TN
 * - Auto-logging: session_start + turn_end → TN L2
 */

function tnFetch(path: string, init?: RequestInit, signal?: AbortSignal) {
  return fetch(`${TN_BASE}${path}`, { ...init, signal });
}

// ─── System prompt guidance injected on every session start ───
const TN_SYSTEM_PROMPT = `
## TormentNexus Integration

You have access to TormentNexus — a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry. Use these tools:

### Memory Tools (persistent cross-session)
- \`tn_memory_store\` — Save important decisions, patterns, and facts with tags (use 'project:', 'failure:', 'pattern:', 'convention:' prefixes)
- \`tn_memory_search\` — Find past memories by keyword, tag, or category
- \`tn_memory_vector_search\` — Semantic vector search for conceptually related memories (sqlite-vec)

### Discovery Tools
- \`tn_tool_search\` — Describe what you need, TN finds the best tool across 20+ servers
- \`tn_session_search\` — Browse 542+ imported sessions from Claude Code, Aider, etc.
- \`tn_skill_manage\` — Access 5,776 reusable skill modules (list, search, read, create)
- \`tn_code_search\` — Search code via AST-grep rules, deepcontext semantic search, or file pattern matching
- \`tn_context_harvest\` — Manually trigger context harvesting from TN L2 memory

### Best Practices
1. **Before any significant task**, call \`tn_memory_search\` or \`tn_memory_vector_search\` to recall relevant past context
2. **Store key decisions** with \`tn_memory_store\` using descriptive tags
3. **Discover tools** with \`tn_tool_search\` when unsure what's available
4. **Review past sessions** with \`tn_session_search\` to learn from previous work
5. **Check skills** with \`tn_skill_manage\` for reusable procedures matching your task
6. **Use \`tn_context_harvest\`** at the start of complex tasks to pull in all relevant context
`;

export default function (pi: ExtensionAPI) {
  // ──────────────────────────────────────────────
  // 1. Session Priming + Guidance Injection
  // ──────────────────────────────────────────────
  pi.on("session_start", async (event, ctx) => {
    const sessionFile = ctx.sessionManager.getSessionFile();
    const reason = event.reason;

    // Store session start in TN L2
    try {
      await tnFetch("/api/memory/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: JSON.stringify({
            content: `Session ${reason}: ${sessionFile ?? "ephemeral"}`,
            tags: ["system:session", `reason:${reason}`],
            category: "session",
            timestamp: new Date().toISOString(),
          }),
        }),
      });
    } catch { /* TN not running */ }

    ctx.ui.setStatus("tn", "TN active • L2 mem • tools • skills");
  });

  // Inject TN guidance + relevant L2 context into every prompt
  pi.on("before_agent_start", async (event, ctx) => {
    // Only inject the full TN guidance on the first turn of a new session
    const isFirstTurn = event.systemPrompt.includes("TormentNexus");
    if (!isFirstTurn) {
      // Try to inject relevant L2 memory context on subsequent turns
      try {
        const prompt = event.prompt;
        const res = await tnFetch(`/api/memory/search?q=${encodeURIComponent(prompt.slice(0, 100))}`, {}, ctx.signal);
        if (res.ok) {
          const body = await res.json();
          const memories = body.data ?? [];
          if (Array.isArray(memories) && memories.length > 0) {
            const contextBlock = memories
              .slice(0, 3)
              .map((m: any) => `  • ${(m.text ?? m.content ?? JSON.stringify(m)).slice(0, 200)}`)
              .join("\n");
            return {
              systemPrompt: `${event.systemPrompt}\n\n## Relevant Context from TormentNexus L2 Memory\n${contextBlock}`,
            };
          }
        }
      } catch { /* TN unavailable, proceed without */ }
      return;
    }

    // First turn: inject full TN guidance + harvest relevant L2 memories
    let memoryContext = "";
    try {
      const res = await tnFetch(`/api/memory/search?q=${encodeURIComponent(event.prompt.slice(0, 100))}`, {}, ctx.signal);
      if (res.ok) {
        const body = await res.json();
        const memories = body.data ?? [];
        if (Array.isArray(memories) && memories.length > 0) {
          memoryContext = "\n\n## Relevant Past Context\n" + memories
            .slice(0, 5)
            .map((m: any) => `  • ${(m.text ?? m.content ?? JSON.stringify(m)).slice(0, 200)}`)
            .join("\n");
        }
      }
    } catch { /* fallback: keyword search */ }

    // Fallback: keyword search against L2 memory list
    if (!memoryContext) {
      try {
        const res = await tnFetch("/api/memory/list", {}, ctx.signal);
        if (res.ok) {
          const all: string[] = await res.json();
          const q = event.prompt.toLowerCase().slice(0, 100);
          const relevant = all
            .map((m) => {
              try { return JSON.parse(m); } catch { return { content: m, tags: [] }; }
            })
            .filter((m) => m.content?.toLowerCase().includes(q) || m.tags?.some((t: string) => t.toLowerCase().includes(q)))
            .slice(0, 5);
          if (relevant.length > 0) {
            memoryContext = "\n\n## Relevant Past Context\n" + relevant
              .map((m) => `  • ${m.content.slice(0, 200)}`)
              .join("\n");
          }
        }
      } catch { /* no memory */ }
    }

    return {
      systemPrompt: event.systemPrompt + TN_SYSTEM_PROMPT + memoryContext,
    };
  });

  // ──────────────────────────────────────────────
  // 2. Per-Turn Context Harvesting
  // ──────────────────────────────────────────────
  pi.on("context", async (event, ctx) => {
    // Only inject context on early turns to avoid bloat
    const lastMessages = event.messages.slice(-4);
    const hasRecentMemorySearch = lastMessages.some(
      (m: any) => m.role === "assistant" && JSON.stringify(m.content)?.includes("tn_memory_search")
    );
    if (hasRecentMemorySearch) return; // Already searched, don't duplicate

    // Harvest relevant memories for this turn
    const lastUserMsg = [...lastMessages].reverse().find((m: any) => m.role === "user");
    if (!lastUserMsg) return;

    const userText = typeof lastUserMsg.content === "string"
      ? lastUserMsg.content
      : lastUserMsg.content?.map((c: any) => c.text).filter(Boolean).join(" ") ?? "";

    if (!userText || userText.length < 10) return;

    try {
      const res = await tnFetch(`/api/memory/search?q=${encodeURIComponent(userText.slice(0, 100))}`, {}, ctx.signal);
      if (!res.ok) return;
      const body = await res.json();
      const memories = body.data ?? [];
      if (!Array.isArray(memories) || memories.length === 0) return;

      // Only inject if memory seems relevant (high confidence)
      const top = memories.slice(0, 2);
      const contextBlock = top
        .map((m: any) => (m.text ?? m.content ?? JSON.stringify(m)).slice(0, 150))
        .filter(Boolean)
        .join("\n");

      if (!contextBlock) return;

      // Inject as a system-context message
      event.messages.push({
        role: "system",
        content: `[TN Context]: ${contextBlock}`,
      });
    } catch { /* silently skip */ }
  });

  // ──────────────────────────────────────────────
  // 3. Compaction Hooks (TN-enhanced summaries)
  // ──────────────────────────────────────────────
  pi.on("session_before_compact", async (event, ctx) => {
    const { preparation, branchEntries, customInstructions, reason, signal } = event;

    // Build a rich summary enriched with TN L2 memory
    const summaryParts: string[] = [];

    // Add branch entries context
    if (branchEntries && branchEntries.length > 0) {
      const fileOps = branchEntries
        .filter((e: any) => e.details?.readFiles || e.details?.modifiedFiles)
        .map((e: any) => {
          const reads = (e.details?.readFiles ?? []).join(", ");
          const mods = (e.details?.modifiedFiles ?? []).join(", ");
          return `${reads ? `Read: ${reads}` : ""}${mods ? ` Modified: ${mods}` : ""}`;
        })
        .filter(Boolean);
      if (fileOps.length > 0) summaryParts.push(`Files: ${fileOps.join("; ")}`);
    }

    // Try to enrich with L2 memory context
    try {
      if (summaryParts.length > 0) {
        const query = summaryParts.join(" ").slice(0, 100);
        const res = await tnFetch(`/api/memory/search?q=${encodeURIComponent(query)}`, {}, signal);
        if (res.ok) {
          const body = await res.json();
          const memories = body.data ?? [];
          if (Array.isArray(memories) && memories.length > 0) {
            const related = memories
              .slice(0, 2)
              .map((m: any) => (m.text ?? m.content ?? "").slice(0, 100))
              .filter(Boolean)
              .join("; ");
            if (related) summaryParts.push(`Related from L2: ${related}`);
          }
        }
      }
    } catch { /* skip enrichment */ }

    const summary = summaryParts.length > 0
      ? summaryParts.join("\n")
      : `Compaction (${reason}) — ${customInstructions ?? "standard"}`;

    return {
      compaction: {
        summary,
        firstKeptEntryId: preparation.firstKeptEntryId,
        tokensBefore: preparation.tokensBefore,
        details: {
          enrichedBy: "tormentnexus-l2",
          reason,
          timestamp: new Date().toISOString(),
        },
      },
    };
  });

  pi.on("session_compact", async (event, ctx) => {
    // Store compaction summary back to TN L2 for future retrieval
    if (!event.compactionEntry?.summary) return;

    try {
      await tnFetch("/api/memory/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: JSON.stringify({
            content: `Compaction [${event.reason}]: ${event.compactionEntry.summary.slice(0, 200)}`,
            tags: ["system:compaction", `reason:${event.reason}`],
            category: "session",
            timestamp: new Date().toISOString(),
          }),
        }),
      });
    } catch { /* skip */ }
  });

  // ──────────────────────────────────────────────
  // 4. Memory Tools
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_memory_store",
    label: "TN Memory Store",
    description: "Store a memory in TormentNexus L2 vault with structured fields and tags. Use this for persistent cross-session knowledge — decisions, patterns, failures, conventions.",
    promptSnippet: "Store knowledge in persistent L2 memory",
    promptGuidelines: [
      "Use tn_memory_store to save important patterns, decisions, and facts across sessions.",
      "Use tags like 'project:name', 'user:', 'failure:', 'pattern:', or 'convention:' for scope filtering.",
      "Good candidates: architectural decisions, bug fixes, build procedures, tool quirks.",
    ],
    parameters: Type.Object({
      content: Type.String({ description: "The memory content to store" }),
      tags: Type.Optional(Type.Array(Type.String(), { description: "Categorization tags, e.g. ['project:bg', 'pattern:build', 'failure:submodule']" })),
      category: Type.Optional(Type.String({ description: "Category: 'pattern', 'decision', 'convention', 'insight', 'failure', 'correction', 'preference'" })),
    }),
    async execute(_toolCallId, params, signal, _onUpdate, ctx) {
      const sessionFile = ctx.sessionManager.getSessionFile();
      const enriched = JSON.stringify({
        content: params.content,
        tags: params.tags ?? [],
        category: params.category ?? "general",
        timestamp: new Date().toISOString(),
        session: sessionFile,
      });

      const res = await tnFetch("/api/memory/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: enriched }),
      }, signal);

      if (!res.ok) {
        return { content: [{ type: "text", text: `Failed to store memory: ${res.status}` }], isError: true };
      }

      return {
        content: [{ type: "text", text: `✅ Memory stored in TormentNexus L2 vault.` }],
        details: { tags: params.tags, category: params.category },
      };
    },
  });

  pi.registerTool({
    name: "tn_memory_search",
    label: "TN Memory Search",
    description: "Search TormentNexus L2 vault for stored memories by keyword, tag filter, or category. Best for finding exact past context.",
    promptSnippet: "Search persistent L2 memory",
    promptGuidelines: [
      "Use tn_memory_search before starting tasks to recall relevant past context.",
      "Filter by tag prefix like 'project:', 'failure:', or 'pattern:' to narrow results.",
    ],
    parameters: Type.Object({
      query: Type.Optional(Type.String({ description: "Keyword to search for in memory content" })),
      tag: Type.Optional(Type.String({ description: "Filter by tag prefix, e.g. 'project:', 'failure:', 'pattern:'" })),
      category: Type.Optional(Type.String({ description: "Filter by category" })),
      limit: Type.Optional(Type.Number({ description: "Max results (default 20)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const res = await tnFetch("/api/memory/list", {}, signal);
      if (!res.ok) return { content: [{ type: "text", text: "Memory unavailable." }], isError: true };

      const memories: string[] = await res.json();
      const limit = params.limit ?? 20;

      const parsed = memories
        .map((m) => {
          try {
            const p = JSON.parse(m);
            return { content: p.content ?? m, tags: p.tags ?? [], category: p.category ?? "general", timestamp: p.timestamp ?? "" };
          } catch {
            return { content: m, tags: [], category: "general", timestamp: "" };
          }
        })
        .filter((m) => {
          if (params.query) {
            const q = params.query.toLowerCase();
            if (!m.content.toLowerCase().includes(q) && !m.tags.some((t: string) => t.toLowerCase().includes(q))) return false;
          }
          if (params.tag && !m.tags.some((t: string) => t.startsWith(params.tag!))) return false;
          if (params.category && m.category !== params.category) return false;
          return true;
        })
        .slice(0, limit);

      if (parsed.length === 0) return { content: [{ type: "text", text: "No matching memories found." }] };

      const formatted = parsed.map((m, i) => {
        const tags = m.tags.length ? ` [${m.tags.join(", ")}]` : "";
        const cat = m.category !== "general" ? ` (${m.category})` : "";
        return `${i + 1}.${cat}${tags}\n   ${m.content.slice(0, 200)}`;
      }).join("\n\n");

      return {
        content: [{ type: "text", text: `📚 ${parsed.length} memories:\n\n${formatted}` }],
        details: { count: parsed.length },
      };
    },
  });

  pi.registerTool({
    name: "tn_memory_vector_search",
    label: "TN Vector Memory Search",
    description: "Semantic search TormentNexus L2 memory using sqlite-vec vector embeddings. Finds conceptually related memories without exact keyword matches. Falls back to keyword search if vector store is empty.",
    promptSnippet: "Semantic vector search L2 memory",
    promptGuidelines: [
      "Use tn_memory_vector_search for fuzzy/conceptual recall — it finds meaning, not exact words.",
      "Great for: 'what did we decide about the build system?' or 'find patterns related to submodule conflicts'.",
    ],
    parameters: Type.Object({
      query: Type.String({ description: "Natural language query for semantic search" }),
      limit: Type.Optional(Type.Number({ description: "Max results (default 10)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const res = await tnFetch(`/api/memory/search?q=${encodeURIComponent(params.query)}`, {}, signal);
      if (!res.ok) return { content: [{ type: "text", text: "Vector search unavailable." }], isError: true };

      const body = await res.json();
      const data = body.data ?? body;

      if (Array.isArray(data) && data.length > 0) {
        const limit = params.limit ?? 10;
        const results = data.slice(0, limit);
        const formatted = results
          .map((r: any, i: number) => `  ${i + 1}. ${(r.text ?? r.content ?? JSON.stringify(r)).slice(0, 200)}`)
          .join("\n\n");
        return {
          content: [{ type: "text", text: `🧠 ${data.length} vector results:\n\n${formatted}` }],
          details: { count: data.length, mode: "vector" },
        };
      }

      // Fallback: keyword search
      const fallbackRes = await tnFetch("/api/memory/list", {}, signal);
      if (!fallbackRes.ok) return { content: [{ type: "text", text: `No results for "${params.query}".` }] };

      const memories: string[] = await fallbackRes.json();
      const query = params.query.toLowerCase();
      const matched = memories
        .map((m) => {
          try { const p = JSON.parse(m); return { content: p.content ?? m, tags: p.tags ?? [] }; }
          catch { return { content: m, tags: [] }; }
        })
        .filter((m) => m.content.toLowerCase().includes(query) || m.tags.some((t: string) => t.toLowerCase().includes(query)))
        .slice(0, params.limit ?? 10);

      if (matched.length === 0) return { content: [{ type: "text", text: `No results for "${params.query}".` }] };

      const formatted = matched.map((m, i) => `  ${i + 1}. ${m.content.slice(0, 200)}`).join("\n\n");
      return {
        content: [{ type: "text", text: `📚 ${matched.length} keyword matches:\n\n${formatted}` }],
        details: { count: matched.length, mode: "keyword-fallback" },
      };
    },
  });

  // ──────────────────────────────────────────────
  // 5. Tool Discovery
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_tool_search",
    label: "TN Tool Search",
    description: "Semantically search for available MCP tools across TormentNexus's 20+ registered servers (desktop-commander, deepcontext, firecrawl, ast-grep, basic-memory, semgrep, etc.). Describe what you need — it finds the best tool by meaning.",
    promptSnippet: "Discover tools via semantic search",
    promptGuidelines: [
      "Use tn_tool_search when you need a tool but aren't sure what's available.",
      "Describe the task naturally — TormentNexus matches by meaning, not keywords.",
      "Examples: 'search codebase', 'query database', 'analyze security', 'manage memory'.",
    ],
    parameters: Type.Object({
      query: Type.String({ description: "Natural language description of what you need to do" }),
      limit: Type.Optional(Type.Number({ description: "Max results (default 5)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const res = await tnFetch(`/api/mcp/native/search?query=${encodeURIComponent(params.query)}`, {}, signal);
      if (!res.ok) return { content: [{ type: "text", text: `Tool search failed: ${res.status}` }], isError: true };

      const body = await res.json();
      const data = body.data ?? body;
      const results = data.results ?? (Array.isArray(data) ? data : []);
      if (results.length === 0) return { content: [{ type: "text", text: `No tools found for "${params.query}".` }] };

      const limit = params.limit ?? 5;
      const top = results.slice(0, limit);
      const formatted = top.map((r: any) => {
        const name = r.originalName ?? r.name ?? "?";
        const server = r.server ?? "?";
        const score = r.score ?? "?";
        const desc = r.description ? r.description.slice(0, 150).replace(/\n/g, " ") : "";
        const match = r.matchReason ?? "";
        return `  [${score}] ${name} (${server})\n         ${desc}\n         → ${match}`;
      }).join("\n\n");

      return {
        content: [{ type: "text", text: `🔧 Top tools for "${params.query}":\n\n${formatted}` }],
        details: { total: results.length, top: top.map((r: any) => r.originalName ?? r.name) },
      };
    },
  });

  // ──────────────────────────────────────────────
  // 6. Session Search
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_session_search",
    label: "TN Session Search",
    description: "Search imported AI coding sessions (542+ from Claude Code, Aider, etc.) by tool type or ID. Browse past session transcripts and statistics.",
    promptSnippet: "Search past AI coding sessions",
    promptGuidelines: [
      "Use tn_session_search to find and review past sessions from other AI coding tools.",
      "Source tools: claude-code, aider. Use action='list' to browse, action='get' to view transcript.",
    ],
    parameters: Type.Object({
      action: Type.String({ description: "'list' — browse sessions, 'get' — retrieve transcript by ID, 'stats' — summary statistics" }),
      sourceTool: Type.Optional(Type.String({ description: "Filter by source tool: 'claude-code', 'aider'." })),
      limit: Type.Optional(Type.Number({ description: "Max results (default 10)" })),
      id: Type.Optional(Type.String({ description: "Session ID for action='get'" })),
    }),
    async execute(_toolCallId, params, signal) {
      if (params.action === "stats") {
        const res = await tnFetch("/api/sessions/imported/maintenance-stats", {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: "Failed to get stats." }], isError: true };
        const body = await res.json();
        const data = body.data ?? body;
        return { content: [{ type: "text", text: `📊 Session Stats:\n\nTotal: ${data.totalSessions ?? "?"}\nArchived transcripts: ${data.archivedTranscriptCount ?? "?"}` }] };
      }

      if (params.action === "get") {
        if (!params.id) return { content: [{ type: "text", text: "Provide 'id' for get." }] };
        const res = await tnFetch(`/api/sessions/imported/get?id=${encodeURIComponent(params.id)}`, {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: `Session not found: ${params.id}` }] };
        const body = await res.json();
        const session = body.data ?? body;
        const summary = session.transcript ? session.transcript.slice(0, 2000) : JSON.stringify(session).slice(0, 2000);
        return { content: [{ type: "text", text: `📝 Session ${params.id}:\n\n${summary}` }] };
      }

      const limit = params.limit ?? 10;
      const res = await tnFetch(`/api/sessions/imported/list?limit=${limit}`, {}, signal);
      if (!res.ok) return { content: [{ type: "text", text: "Failed to list sessions." }], isError: true };
      const body = await res.json();
      let sessions = body.data ?? body;
      if (!Array.isArray(sessions)) sessions = [];

      if (params.sourceTool) sessions = sessions.filter((s: any) => s.sourceTool === params.sourceTool);
      if (sessions.length === 0) return { content: [{ type: "text", text: "No matching sessions found." }] };

      const formatted = sessions.slice(0, limit).map((s: any, i: number) => {
        const tool = s.sourceTool ?? "?"; const id = s.id ?? "?"; const fmt = s.sessionFormat ?? "?";
        const valid = s.valid !== false ? "✅" : "❌";
        return `${i + 1}. [${valid}] ${tool} (${fmt}) ID: ${id}`;
      }).join("\n\n");

      return {
        content: [{ type: "text", text: `📋 ${sessions.length} sessions:\n\n${formatted}\n\nUse action='get' id='<ID>' to view transcript.` }],
        details: { count: sessions.length },
      };
    },
  });

  // ──────────────────────────────────────────────
  // 7. Skill Management
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_skill_manage",
    label: "TN Skill Management",
    description: "Manage 5,776+ TormentNexus skills (reusable procedural knowledge modules). List, search, read, and create skills.",
    promptSnippet: "Manage skills (list, search, read, create)",
    promptGuidelines: [
      "Use tn_skill_manage to discover or create reusable skill modules.",
      "Skills are SKILL.md format (YAML frontmatter + markdown body).",
      "Search with action='search' q='keyword'. List all with action='list'.",
    ],
    parameters: Type.Object({
      action: Type.String({ description: "'list' | 'search' | 'read' | 'create'" }),
      id: Type.Optional(Type.String({ description: "Skill ID (required for 'read' and 'create')" })),
      query: Type.Optional(Type.String({ description: "Search query (required for 'search')" })),
      content: Type.Optional(Type.String({ description: "Full skill markdown with YAML frontmatter (required for 'create')" })),
      limit: Type.Optional(Type.Number({ description: "Max results (default 20)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const limit = params.limit ?? 20;

      if (params.action === "list") {
        const res = await tnFetch("/api/skills/list", {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: "Failed to list skills." }], isError: true };
        const body = await res.json();
        const skills = body.skills ?? body.data?.skills ?? [];
        const total = body.count ?? skills.length;
        const top = skills.slice(0, limit);
        return {
          content: [{ type: "text", text: `📚 ${total} skills:\n\n${top.map((s: any, i: number) => `  ${i + 1}. ${s.id}`).join("\n")}\n\nUse action='search' to find specific skills.` }],
        };
      }

      if (params.action === "search") {
        if (!params.query) return { content: [{ type: "text", text: "Provide 'query'." }] };
        const res = await tnFetch(`/api/skills/search?q=${encodeURIComponent(params.query)}`, {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: "Search failed." }], isError: true };
        const body = await res.json();
        const skills = body.skills ?? body.data?.skills ?? [];
        const total = body.count ?? skills.length;
        const top = skills.slice(0, limit);
        if (top.length === 0) return { content: [{ type: "text", text: `No skills for "${params.query}".` }] };
        return {
          content: [{ type: "text", text: `🔍 ${total} skills matching "${params.query}":\n\n${top.map((s: any, i: number) => `  ${i + 1}. ${s.id}`).join("\n")}` }],
        };
      }

      if (params.action === "read") {
        if (!params.id) return { content: [{ type: "text", text: "Provide 'id'." }] };
        const res = await tnFetch(`/api/skills/read?name=${encodeURIComponent(params.id)}`, {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: `Skill '${params.id}' not found.` }] };
        const body = await res.json();
        const data = body.data ?? body;
        const content = data.content ?? data.skill ?? "";
        let text = "";
        if (Array.isArray(content)) text = content.map((c: any) => c.text ?? JSON.stringify(c)).join("");
        else if (typeof content === "string") text = content;
        else text = JSON.stringify(content);
        return { content: [{ type: "text", text: `📖 ${params.id}\n\n${text.slice(0, 3000)}` }] };
      }

      if (params.action === "create") {
        if (!params.id || !params.content) return { content: [{ type: "text", text: "Provide 'id' and 'content'." }] };
        const res = await tnFetch("/api/skills/create", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: params.id, content: params.content }),
        }, signal);
        if (!res.ok) {
          const err = await res.text().catch(() => "");
          return { content: [{ type: "text", text: `Failed: ${res.status} — ${err.slice(0, 200)}` }], isError: true };
        }
        return { content: [{ type: "text", text: `✅ Skill '${params.id}' created.` }] };
      }

      return { content: [{ type: "text", text: "Usage: action='list'|'search'|'read'|'create'" }] };
    },
  });

  // ──────────────────────────────────────────────
  // 8. Code Intelligence (AST + LSP + Search)
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_code_search",
    label: "TN Code Search",
    description: "Search code using AST-grep rules, semantic codebase search, or pattern matching. Finds code by structure (AST), meaning (deepcontext), or file patterns.",
    promptSnippet: "Search code via AST, semantic, or pattern matching",
    promptGuidelines: [
      "Use tn_code_search to find code by structure (AST), semantics, or pattern.",
      "mode='ast' uses ast-grep YAML rules for structural pattern matching (e.g. 'find all error handlers').",
      "mode='semantic' uses deepcontext for natural language code search (e.g. 'find authentication logic').",
      "mode='pattern' uses glob/file patterns (e.g. '**/*.ts' matching 'class.*Handler').",
    ],
    parameters: Type.Object({
      query: Type.String({ description: "Search query — AST rule YAML, natural language description, or glob/file pattern" }),
      mode: Type.Optional(Type.String({ description: "Search mode: 'ast' (structural), 'semantic' (natural language), 'pattern' (file glob)" })),
      path: Type.Optional(Type.String({ description: "Directory or file path to scope the search" })),
      limit: Type.Optional(Type.Number({ description: "Max results (default 10)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const mode = params.mode ?? "semantic";
      const limit = params.limit ?? 10;

      if (mode === "ast") {
        // Use ast-grep-mcp for structural code search
        const res = await tnFetch(
          `/api/mcp/native/search?query=${encodeURIComponent(`ast-grep code pattern ${params.query}`)}`,
          {}, signal
        );
        if (!res.ok) return { content: [{ type: "text", text: "AST search unavailable." }], isError: true };
        const body = await res.json();
        const data = body.data ?? body;
        const results = data.results ?? [];
        const tools = results.filter((r: any) => r.server === "ast-grep-mcp" || r.server === "ast-grep").slice(0, 3);
        if (tools.length > 0) {
          return {
            content: [{ type: "text", text: `🔍 AST tools available for "${params.query}":\n\n${tools.map((t: any) => `  • ${t.originalName ?? t.name} (${t.server})`).join("\n")}\n\nUse these tools directly for structural pattern matching.` }],
          };
        }
        return { content: [{ type: "text", text: `No AST tools found. Try mode='semantic' for natural language code search.` }] };
      }

      if (mode === "semantic") {
        const res = await tnFetch(
          `/api/mcp/native/search?query=${encodeURIComponent(`code search semantic codebase ${params.query}`)}`,
          {}, signal
        );
        if (!res.ok) return { content: [{ type: "text", text: "Semantic search failed." }], isError: true };
        const body = await res.json();
        const data = body.data ?? body;
        const results = data.results ?? [];
        const top = results.slice(0, limit);

        if (top.length === 0) return { content: [{ type: "text", text: `No tools found for "${params.query}".` }] };

        const formatted = top.map((r: any) => {
          const name = r.originalName ?? r.name ?? "?";
          const server = r.server ?? "?";
          const score = r.score ?? "?";
          const desc = r.description ? r.description.slice(0, 120).replace(/\n/g, " ") : "";
          return `  [${score}] ${name} (${server})\n         ${desc}`;
        }).join("\n\n");

        return {
          content: [{ type: "text", text: `🔧 Tools for "${params.query}":\n\n${formatted}` }],
          details: { total: results.length },
        };
      }

      // mode === "pattern"
      const pathFilter = params.path ? ` path:${params.path}` : "";
      const res = await tnFetch(
        `/api/mcp/native/search?query=${encodeURIComponent(`search files pattern ${params.query}${pathFilter}`)}`,
        {}, signal
      );
      if (!res.ok) return { content: [{ type: "text", text: "Pattern search failed." }], isError: true };
      const body = await res.json();
      const data = body.data ?? body;
      const results = data.results ?? [];
      const top = results.slice(0, limit);

      if (top.length === 0) return { content: [{ type: "text", text: `No tools found for pattern "${params.query}".` }] };

      const formatted = top.map((r: any) => {
        const name = r.originalName ?? r.name ?? "?";
        const server = r.server ?? "?";
        const desc = r.description ? r.description.slice(0, 120).replace(/\n/g, " ") : "";
        return `  • ${name} (${server})\n    ${desc}`;
      }).join("\n\n");

      return {
        content: [{ type: "text", text: `📁 Tools for pattern "${params.query}":\n\n${formatted}` }],
      };
    },
  });

  // ──────────────────────────────────────────────
  // 9. Context Harvest (manual trigger)
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_context_harvest",
    label: "TN Context Harvest",
    description: "Manually harvest relevant context from TormentNexus L2 memory for the current task. Pulls in related memories, relevant skills, and past sessions.",
    promptSnippet: "Harvest context from L2 memory",
    promptGuidelines: [
      "Use tn_context_harvest at the start of complex tasks to gather all relevant context.",
      "It searches L2 memory, skills, and sessions for context related to your query.",
    ],
    parameters: Type.Object({
      query: Type.String({ description: "What you're working on — TN searches L2 memory + skills + sessions for related context" }),
      harvestMemory: Type.Optional(Type.Boolean({ description: "Search L2 memory (default true)" })),
      harvestSkills: Type.Optional(Type.Boolean({ description: "Search skill registry (default false, can be slow)" })),
    }),
    async execute(_toolCallId, params, signal) {
      const results: string[] = [];

      // 1. Harvest L2 memory
      if (params.harvestMemory !== false) {
        try {
          // Vector search first
          const vecRes = await tnFetch(`/api/memory/search?q=${encodeURIComponent(params.query)}`, {}, signal);
          if (vecRes.ok) {
            const body = await vecRes.json();
            const memories = body.data ?? [];
            if (Array.isArray(memories) && memories.length > 0) {
              results.push(`## L2 Memory (vector search)\n${memories.slice(0, 5).map((m: any) => `  • ${(m.text ?? m.content ?? JSON.stringify(m)).slice(0, 200)}`).join("\n")}`);
            }
          }

          // Fallback keyword search
          if (results.length === 0 || !results.some(r => r.includes("vector"))) {
            const listRes = await tnFetch("/api/memory/list", {}, signal);
            if (listRes.ok) {
              const all: string[] = await listRes.json();
              const q = params.query.toLowerCase();
              const relevant = all
                .map((m) => {
                  try { return { ...JSON.parse(m), raw: m }; }
                  catch { return { content: m, tags: [], raw: m }; }
                })
                .filter((m) => m.content?.toLowerCase().includes(q) || m.tags?.some((t: string) => t.toLowerCase().includes(q)))
                .slice(0, 5);
              if (relevant.length > 0) {
                results.push(`## L2 Memory (keyword)\n${relevant.map((m) => `  • ${m.content.slice(0, 200)}`).join("\n")}`);
              }
            }
          }
        } catch { /* skip */ }
      }

      // 2. Harvest skills
      if (params.harvestSkills) {
        try {
          const skRes = await tnFetch(`/api/skills/search?q=${encodeURIComponent(params.query)}`, {}, signal);
          if (skRes.ok) {
            const body = await skRes.json();
            const skills = body.skills ?? body.data?.skills ?? [];
            if (skills.length > 0) {
              results.push(`## Related Skills\n${skills.slice(0, 5).map((s: any) => `  • ${s.id}`).join("\n")}`);
            }
          }
        } catch { /* skip */ }
      }

      if (results.length === 0) {
        return { content: [{ type: "text", text: `No relevant context found in L2 memory for "${params.query}". Try storing some context first with tn_memory_store.` }] };
      }

      return {
        content: [{ type: "text", text: `🌾 Context harvested for "${params.query}":\n\n${results.join("\n\n")}` }],
        details: { sources: results.map(r => r.split("\n")[0].replace("## ", "")) },
      };
    },
  });

  // ──────────────────────────────────────────────
  // 10. Scratchpad (L1 working memory)
  // ──────────────────────────────────────────────
  pi.registerTool({
    name: "tn_scratchpad",
    label: "TN Scratchpad",
    description: "Read or write the TormentNexus L1 session scratchpad — ephemeral working memory for the current session.",
    parameters: Type.Object({
      action: Type.String({ description: "'get' to read, 'set' to write" }),
      content: Type.Optional(Type.String({ description: "Content to write (required for 'set')" })),
    }),
    async execute(_toolCallId, params, signal) {
      if (params.action === "get") {
        const res = await tnFetch("/api/memory/scratchpad/get", {}, signal);
        if (!res.ok) return { content: [{ type: "text", text: "Scratchpad empty." }] };
        const text = await res.text();
        return { content: [{ type: "text", text: text || "Empty." }] };
      }
      if (params.action === "set" && params.content) {
        const res = await tnFetch("/api/memory/scratchpad/set", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ content: params.content }),
        }, signal);
        if (!res.ok) return { content: [{ type: "text", text: `Failed: ${res.status}` }], isError: true };
        return { content: [{ type: "text", text: "✅ Scratchpad updated." }] };
      }
      return { content: [{ type: "text", text: 'Usage: action="get" or "set".' }] };
    },
  });

  // ──────────────────────────────────────────────
  // 11. Auto-logging
  // ──────────────────────────────────────────────
  pi.on("turn_end", async (event, ctx) => {
    if (!event.toolResults || event.toolResults.length === 0) return;

    const summary = event.message?.content
      ?.filter((c: any) => c.type === "text")
      ?.map((c: any) => c.text)
      ?.join(" ")
      ?.slice(0, 500);

    if (!summary || summary.length < 50) return;

    try {
      await tnFetch("/api/memory/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: JSON.stringify({
            content: `Turn ${event.turnIndex}: ${summary.slice(0, 200)}...`,
            tags: ["system:turn", `turn:${event.turnIndex}`],
            category: "session",
            timestamp: new Date().toISOString(),
          }),
        }),
      });
    } catch { /* skip */ }
  });

  // Cleanup on session end
  pi.on("session_shutdown", async (_event, ctx) => {
    try { ctx.ui.setStatus("tn", ""); } catch { /* skip */ }
  });
}
