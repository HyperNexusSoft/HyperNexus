@echo off
cd /d "C:\Users\hyper\workspace\tormentnexus"
setlocal enabledelayedexpansion

echo ========================================
echo  TormentNexus Multi-Agent Installer
echo  Run this as Administrator for services!
echo ========================================
echo.

echo === Step 1: Windows Services ===
echo.
echo Registering Go Sidecar (port 7778)...
sc create "TormentNexusSidecar" binPath="\"C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe\" serve" start=auto displayname="TormentNexus Sidecar"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)

echo Registering Dashboard (port 7779)...
sc create "TormentNexusDashboard" binPath="\"C:\Program Files\nodejs\node.exe\" \"C:\Users\hyper\workspace\tormentnexus\apps\web\node_modules\.bin\next.cmd\" dev -p 7779" start=auto displayname="TormentNexus Dashboard"

echo Registering Watchdog...
sc create "TormentNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\tormentnexus\watchdog.py\"" start=auto displayname="TormentNexus Watchdog"
echo.

echo === Step 2: Pi Coding Agent ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
if exist "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts" (
    echo ℹ️  Pi extension already exists at global path. Skipping copy.
) else (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\packages\tormentnexus\index.ts" "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts"
    if !errorlevel!==0 (echo ✅ Pi extension installed) else (echo ⚠️)
)
echo.



echo === Step 3: Ollama / vLLM (Tool Prediction Engine) ===
echo.
echo TormentNexus uses a local LLM for tool prediction (ConversationalToolInjector).
echo.
echo Choose an option:
echo   [1] Ollama — easiest, auto-start as Windows service (recommended)
echo   [2] vLLM  — faster inference, GPU-accelerated
echo   [S] Skip — tool prediction degrades to keyword matching
echo.
choice /C 12S /N /M "Select [1], [2], or [S]: "
if errorlevel 3 goto :skip_llm
if errorlevel 2 goto :install_vllm
if errorlevel 1 goto :install_ollama

:install_ollama
echo Installing Ollama...
curl -sL -o "%TEMP%\ollama_windows.exe" "https://github.com/ollama/ollama/releases/latest/download/OllamaSetup.exe"
if exist "%TEMP%\ollama_windows.exe" (
    start /wait "" "%TEMP%\ollama_windows.exe" /S
    echo Installing Gemma 4 model (this downloads ~8GB, may take a while)...
    ollama pull gemma4 2>nul || ollama pull gemma3:12b 2>nul
    echo.
    echo Setting up Ollama as auto-start service...
    sc config ollama start=auto >nul 2>nul
    sc start ollama >nul 2>nul
    echo ✅ Ollama + Gemma 4 installed at http://127.0.0.1:11434
) else (
    echo ⚠️  Download failed. Install manually from https://ollama.ai
)
goto :end_llm

:install_vllm
echo vLLM installation requires Python + CUDA.
echo.
echo pip install vllm
echo vllm serve gemma-4 --port 11434 --api-key token-abc123
echo.
echo Set TORMENTNEXUS_OLLAMA_URL=http://127.0.0.1:11434
echo.
echo 📋 Manual setup required — see https://github.com/vllm-project/vllm
goto :end_llm

:skip_llm
echo ⏭️  Skipping LLM install. Tool prediction will use BM25 keyword matching.
goto :end_llm

:end_llm
echo.

echo === Step 4: CodeWhale Plugin ===
echo.
where codewhale >nul 2>nul
if %errorlevel%==0 (
    if not exist "%USERPROFILE%\.codewhale\plugins\tormentnexus" mkdir "%USERPROFILE%\.codewhale\plugins\tormentnexus\skills"
    xcopy /E /I /Y "C:\Users\hyper\workspace\tormentnexus\.codewhale\plugins\tormentnexus" "%USERPROFILE%\.codewhale\plugins\tormentnexus" >nul
    codewhale mcp add "tormentnexus" --command "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" --arg "mcp" >nul 2>nul
    echo ✅ CodeWhale plugin installed
) else (echo ⏭️)
echo.

echo === Step 4: Gemini CLI ===
echo.
where gemini >nul 2>nul
if %errorlevel%==0 (
    if not exist "%USERPROFILE%\.gemini\extensions" mkdir "%USERPROFILE%\.gemini\extensions"
    xcopy /E /I /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\extensions\tormentnexus" "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul
    gemini extensions link "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul 2>nul
    if not exist "%USERPROFILE%\.gemini\skills\tormentnexus" mkdir "%USERPROFILE%\.gemini\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.gemini\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Gemini CLI extension + skill
) else (echo ⏭️)
echo.

echo === Step 5: Claude Desktop ===
echo.
if exist "%APPDATA%\Claude\claude_desktop_config.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\claude-desktop-mcp.json" "%APPDATA%\Claude\claude_desktop_config.json.tn-template" >nul
    echo ✅ Template saved — merge mcpServers block manually
)
echo.

echo === Step 6: Claude Code CLI ===
echo.
where claude >nul 2>nul
if %errorlevel%==0 (
    claude mcp add --transport stdio tormentnexus -- "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" "mcp" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Claude CLI MCP) else (echo ✅ May already exist)
) else (echo ⏭️)
echo.

echo === Step 7: Codex CLI ===
echo.
where codex >nul 2>nul
if %errorlevel%==0 (
    codex mcp add "tormentnexus" --env TORMENTNEXUS_WORKSPACE_ROOT="C:\Users\hyper\workspace\tormentnexus" -- "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" "mcp" >nul 2>nul
    if not exist "%USERPROFILE%\.codex\skills\tormentnexus" mkdir "%USERPROFILE%\.codex\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.codex\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.codex\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Codex CLI MCP + skill
) else (echo ⏭️)
echo.

echo === Step 8: Cursor ===
echo.
if exist "%USERPROFILE%\.cursor\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\cursor-mcp.json" "%USERPROFILE%\.cursor\mcp.json.tn-template" >nul
    echo ✅ Template saved
) else (echo ⏭️)
echo.

echo === Step 9: Windsurf ===
echo.
where windsurf >nul 2>nul
if %errorlevel%==0 (
    windsurf --add-mcp "{\"name\":\"tormentnexus\",\"command\":\"C:\\\\Users\\\\hyper\\\\workspace\\\\tormentnexus\\\\tormentnexus.exe\",\"args\":[\"mcp\"]}" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Windsurf MCP) else (echo ⚠️)
) else (echo ⏭️)
echo.

echo === Step 10: VS Code ===
echo.
if exist "%USERPROFILE%\.vscode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json.tn-template" >nul
) else (
    mkdir "%USERPROFILE%\.vscode" >nul 2>nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json" >nul
)
echo ✅ VS Code MCP
echo.

echo === Step 11: OpenCode (Fix) ===
echo.
if exist "%USERPROFILE%\.opencode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\opencode-mcp.json" "%USERPROFILE%\.opencode\mcp.json.tn-template" >nul
    echo ✅ Template saved — merge to fix bin/tormentnexus.exe path
) else (echo ⏭️)
echo.

echo === Step 12: Continue (Fix) ===
echo.
if exist "%USERPROFILE%\.continue\config.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\continue-mcp.json" "%USERPROFILE%\.continue\tormentnexus-mcp-merge.json.tn" >nul
    echo ✅ Template saved — merge to fix bin/tormentnexus.exe path
) else (echo ⏭️)
echo.

echo === Step 13: Mavis / MiniMax Code ===
echo.
if exist "%USERPROFILE%\.mavis\mcp\mcp.json" (
    if not exist "%USERPROFILE%\.mavis\skills\tormentnexus" mkdir "%USERPROFILE%\.mavis\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\mcp.json" "%USERPROFILE%\.mavis\mcp\tormentnexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.mavis\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Mavis MCP + skill
) else (echo ⏭️)
echo.

echo === Step 14: Antigravity IDE ===
echo.
if exist "%USERPROFILE%\.antigravity" (
    if not exist "%USERPROFILE%\.antigravity\extensions\tormentnexus" mkdir "%USERPROFILE%\.antigravity\extensions\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\mcp_config.json" "%USERPROFILE%\.antigravity\tormentnexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\extensions\tormentnexus\SKILL.md" "%USERPROFILE%\.antigravity\extensions\tormentnexus\SKILL.md" >nul
    if not exist "%USERPROFILE%\.antigravity\agents\tormentnexus" mkdir "%USERPROFILE%\.antigravity\agents\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\agents\tormentnexus\agent.md" "%USERPROFILE%\.antigravity\agents\tormentnexus\agent.md" >nul
    echo ✅ Antigravity MCP + extension + agent
) else (echo ⏭️)
echo.

echo === Step 15: Kimi Desktop ===
echo.
if exist "%USERPROFILE%\.kimi-code" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.kimi-code\mcp.json" "%USERPROFILE%\.kimi-code\mcp.json.tn-template" >nul
    echo ✅ Template saved — merge mcpServers block
) else (echo ⏭️)
echo.

echo === Step 16: ZCode Desktop ===
echo.
if exist "%USERPROFILE%\.zcode\v2" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.zcode\mcp.json" "%USERPROFILE%\.zcode\v2\mcp.json.tn-template" >nul
    echo ✅ Template saved — merge mcpServers block
) else (echo ⏭️)
echo.

echo === Step 17: Starting Services ===
echo.
sc start TormentNexusSidecar >nul 2>nul
sc start TormentNexusDashboard >nul 2>nul
sc start TormentNexusWatchdog >nul 2>nul
echo.
echo ========================================
echo  TormentNexus Multi-Agent Installer
echo  Complete!
echo.
echo  ✅ Pi Coding Agent        ~\.pi\agent\extensions\
echo  ✅ CodeWhale               ~\.codewhale\plugins\tormentnexus\
echo  ✅ Gemini CLI              ~\.gemini\extensions\ + skills
echo  ✅ Claude CLI (MCP)        claude mcp add
echo  ✅ Codex CLI (MCP+skill)   ~\.codex\skills\tormentnexus\
echo  ✅ Windsurf (MCP)          --add-mcp flag
echo  ✅ VS Code (MCP)           .vscode\mcp.json
echo  ✅ Mavis (MCP+skill)       ~\.mavis\skills\tormentnexus\
echo  ✅ Antigravity IDE         3 files installed
echo  📋 Claude Desktop          template — merge manually
echo  📋 Cursor                  template — merge manually
echo  📋 Kimi Desktop            template — merge manually
echo  📋 ZCode Desktop           template — merge manually
echo  📋 OpenCode                template — merge to fix path
echo  📋 Continue                template — merge to fix path
echo ========================================
echo.
echo  NOTE: Templates saved with .tn-template suffix.
echo  The correct binary path is: tormentnexus.exe
echo  (NOT bin/tormentnexus.exe)
echo.
pause
