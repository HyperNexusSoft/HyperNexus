#!/usr/bin/env node
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

const HOME = os.homedir();
const OPENHANDS_DIR = path.join(HOME, ".openhands");

console.log("Installing TormentNexus for OpenHands...");

fs.mkdirSync(path.join(OPENHANDS_DIR, "microagents"), { recursive: true });

fs.copyFileSync(
  path.join(__dirname, "config.toml"),
  path.join(OPENHANDS_DIR, "config.toml")
);

fs.copyFileSync(
  path.join(__dirname, "microagent.md"),
  path.join(OPENHANDS_DIR, "microagents", "tormentnexus.md")
);

const mcpConfig = {
  mcpServers: {
    tormentnexus: {
      command: process.platform === "win32" ? "tormentnexus.exe" : "tormentnexus",
      args: ["mcp"],
      env: { TORMENTNEXUS_WORKSPACE_ROOT: process.cwd() },
    },
  },
};
fs.writeFileSync(
  path.join(OPENHANDS_DIR, "mcp.json"),
  JSON.stringify(mcpConfig, null, 2)
);

console.log("TormentNexus installed for OpenHands!");
console.log("  Config: ~/.openhands/config.toml");
console.log("  Microagent: ~/.openhands/microagents/tormentnexus.md");
console.log("  MCP: ~/.openhands/mcp.json");
