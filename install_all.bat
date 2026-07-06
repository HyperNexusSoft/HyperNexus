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
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)

echo Registering Watchdog...
sc create "TormentNexusWatchdog" binPath="\"C:\Python314\pythonw.exe\" -u \"C:\Users\hyper\workspace\tormentnexus\watchdog.py\"" start=auto displayname="TormentNexus Watchdog"
if %errorlevel%==0 (echo ✅) else (echo ⚠️ may already exist)
echo.

echo === Step 2: Pi Coding Agent Extension ===
echo.
if not exist "%USERPROFILE%\.pi\agent\extensions" mkdir "%USERPROFILE%\.pi\agent\extensions"
copy /Y "C:\Users\hyper\workspace\tormentnexus\.pi\extensions\tormentnexus.ts" "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts"
if %errorlevel%==0 (echo ✅ Pi extension installed) else (echo ⚠️ Pi extension copy failed)
echo.

echo === Step 3: CodeWhale Integration ===
echo.
where codewhale >nul 2>nul
if %errorlevel%==0 (
    echo CodeWhale detected — installing...
    if not exist "%USERPROFILE%\.codewhale\skills\tormentnexus" mkdir "%USERPROFILE%\.codewhale\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.codewhale\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.codewhale\skills\tormentnexus\SKILL.md"
    if %errorlevel%==0 (echo ✅ CodeWhale skill) else (echo ⚠️ Skill copy)
    codewhale mcp add "tormentnexus" --command "C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe" --arg "mcp" >nul 2>nul
    if %errorlevel%==0 (echo ✅ CodeWhale MCP) else (echo ⚠️ MCP may already exist)
) else (echo ⏭️ CodeWhale not installed)
echo.

echo === Step 4: Gemini CLI Extension + Skill ===
echo.
where gemini >nul 2>nul
if %errorlevel%==0 (
    echo Gemini CLI detected — installing...
    if not exist "%USERPROFILE%\.gemini\extensions" mkdir "%USERPROFILE%\.gemini\extensions"
    xcopy /E /I /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\extensions\tormentnexus" "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul
    gemini extensions link "%USERPROFILE%\.gemini\extensions\tormentnexus" >nul 2>nul
    if %errorlevel%==0 (echo ✅ Gemini extension linked) else (echo ⚠️ Extension link)
    if not exist "%USERPROFILE%\.gemini\skills\tormentnexus" mkdir "%USERPROFILE%\.gemini\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.gemini\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.gemini\skills\tormentnexus\SKILL.md"
    if %errorlevel%==0 (echo ✅ Gemini skill) else (echo ⚠️ Skill copy)
) else (echo ⏭️ Gemini CLI not installed)
echo.

echo === Step 5: Claude Desktop MCP ===
echo.
if exist "%APPDATA%\Claude\claude_desktop_config.json" (
    echo Claude Desktop detected — updating MCP config...
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\claude-desktop-mcp.json" "%APPDATA%\Claude\claude_desktop_config.json.tn-template" >nul
    echo ✅ Claude Desktop MCP template saved to %%APPDATA%%\Claude\claude_desktop_config.json.tn-template
    echo ⚠️  Manually merge the tormentnexus block into claude_desktop_config.json
)
echo.

echo === Step 6: Claude Code CLI MCP ===
echo.
where claude >nul 2>nul
if %errorlevel%==0 (
    echo Claude CLI detected — updating MCP config in settings.json...
    if exist "%USERPROFILE%\.claude\settings.json" (
        copy "C:\Users\hyper\workspace\tormentnexus\.editor-configs\claude-desktop-mcp.json" "%USERPROFILE%\.claude\tormentnexus-mcp-merge.json.tn" >nul
        echo ✅ Claude CLI MCP template saved
        echo ⚠️  Manually merge the mcpServers.tormentnexus block
    )
) else (echo ⏭️ Claude CLI not installed)
echo.

echo === Step 7: Cursor MCP ===
echo.
if exist "%USERPROFILE%\.cursor\mcp.json" (
    echo Cursor detected — MCP template available
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\cursor-mcp.json" "%USERPROFILE%\.cursor\mcp.json.tn-template" >nul
    echo ✅ Cursor MCP template saved
) else (echo ⏭️ Cursor not configured)
echo.

echo === Step 8: Windsurf MCP ===
echo.
where windsurf >nul 2>nul
if %errorlevel%==0 (
    echo Windsurf detected — adding MCP via CLI...
    windsurf --add-mcp "{\"name\":\"tormentnexus\",\"command\":\"C:\\\\Users\\\\hyper\\\\workspace\\\\tormentnexus\\\\tormentnexus.exe\",\"args\":[\"mcp\"]}" >nul 2>nul
    if !errorlevel!==0 (echo ✅ Windsurf MCP added) else (echo ⚠️ Windsurf MCP — try manually)
) else (echo ⏭️ Windsurf not installed)
echo.

echo === Step 9: VS Code MCP ===
echo.
if exist "%USERPROFILE%\.vscode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json.tn-template" >nul
    echo ✅ VS Code MCP template saved
) else (
    mkdir "%USERPROFILE%\.vscode" >nul 2>nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\vscode-mcp.json" "%USERPROFILE%\.vscode\mcp.json" >nul
    if %errorlevel%==0 (echo ✅ VS Code MCP config installed) else (echo ⚠️ VS Code MCP)
)
echo.

echo === Step 10: OpenCode MCP Fix ===
echo.
if exist "%USERPROFILE%\.opencode\mcp.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\opencode-mcp.json" "%USERPROFILE%\.opencode\mcp.json.tn-template" >nul
    echo ✅ OpenCode MCP template saved (replaces bin/tormentnexus.exe path)
) else (echo ⏭️ OpenCode not configured)
echo.

echo === Step 11: Continue MCP Fix ===
echo.
if exist "%USERPROFILE%\.continue\config.json" (
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.editor-configs\continue-mcp.json" "%USERPROFILE%\.continue\tormentnexus-mcp-merge.json.tn" >nul
    echo ✅ Continue MCP template saved (replaces bin/tormentnexus.exe path)
) else (echo ⏭️ Continue not configured)
echo.

echo === Step 12: Mavis / MiniMax Code Integration ===
echo.
if exist "%USERPROFILE%\.mavis\mcp\mcp.json" (
    if not exist "%USERPROFILE%\.mavis\skills\tormentnexus" mkdir "%USERPROFILE%\.mavis\skills\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\mcp.json" "%USERPROFILE%\.mavis\mcp\tormentnexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.mavis\skills\tormentnexus\SKILL.md" "%USERPROFILE%\.mavis\skills\tormentnexus\SKILL.md" >nul
    echo ✅ Mavis MCP + skill installed
) else (echo ⏭️ Mavis not configured)
echo.

echo === Step 13: Antigravity IDE Integration ===
echo.
if exist "%USERPROFILE%\.antigravity" (
    if not exist "%USERPROFILE%\.antigravity\extensions\tormentnexus" mkdir "%USERPROFILE%\.antigravity\extensions\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\mcp_config.json" "%USERPROFILE%\.antigravity\tormentnexus-mcp-merge.json.tn" >nul
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\extensions\tormentnexus\SKILL.md" "%USERPROFILE%\.antigravity\extensions\tormentnexus\SKILL.md" >nul
    if not exist "%USERPROFILE%\.antigravity\agents\tormentnexus" mkdir "%USERPROFILE%\.antigravity\agents\tormentnexus"
    copy /Y "C:\Users\hyper\workspace\tormentnexus\.antigravity\agents\tormentnexus\agent.md" "%USERPROFILE%\.antigravity\agents\tormentnexus\agent.md" >nul
    echo ✅ Antigravity MCP + extension + agent installed
) else (echo ⏭️ Antigravity not configured)
echo.

echo === Step 14: Starting Services ===
echo.
sc start TormentNexusSidecar >nul 2>nul
sc start TormentNexusDashboard >nul 2>nul
sc start TormentNexusWatchdog >nul 2>nul
echo.
echo ========================================
echo  TormentNexus Multi-Agent Installer
echo  Complete!
echo.
echo  Installed for:
echo   ✅ Pi Coding Agent            ~\.pi\agent\extensions\
echo   ✅ CodeWhale                  ~\.codewhale\skills\ + MCP
echo   ✅ Gemini CLI                 ~\.gemini\extensions\ + skills
echo   📋 Claude Desktop             template saved
echo   📋 Claude Code CLI            template saved
echo   📋 Cursor                     template saved
echo   ✅ Windsurf                   --add-mcp via CLI
echo   ✅ VS Code                    .vscode\mcp.json
echo   📋 OpenCode                   template saved (path fix)
echo   📋 Continue                   template saved (path fix)
echo   ✅ Mavis / MiniMax Code       .mavis\skills\ + MCP
echo   ✅ Antigravity IDE            .antigravity\ + extensions
echo ========================================
echo.
echo  NOTE: Some tools need manual merge of the MCP JSON.
echo  Templates saved with .tn-template or .tn suffix.
echo  The correct path is: tormentnexus.exe (not bin/tormentnexus.exe)
echo.
pause
