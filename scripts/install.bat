@echo off
echo ============================================================
echo TormentNexus Universal Installer for Windows
echo ============================================================
echo.
echo This will install TormentNexus support for ALL detected
echo AI coding clients on your system (38+ supported).
echo.
echo Press Ctrl+C to cancel, or any key to continue...
pause >nul

echo.
echo [1/3] Installing Python client support...
python "%~dp0install-client-support.py"
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Python not found or script failed.
    echo Please install Python 3.9+ and try again.
    pause
    exit /b 1
)

echo.
echo [2/3] Installing Pi extension...
if exist "%USERPROFILE%\.pi\agent\extensions" (
    copy /Y "%~dp0..\packages\tormentnexus\index.ts" "%USERPROFILE%\.pi\agent\extensions\tormentnexus.ts" 2>nul
    echo   Pi extension installed.
) else (
    echo   Pi not detected - skipping.
)

echo.
echo [3/3] Installing VS Code extension...
if exist "%USERPROFILE%\.vscode\extensions" (
    mkdir "%USERPROFILE%\.vscode\extensions\tormentnexus" 2>nul
    xcopy /E /Y "%~dp0..\apps\vscode\*" "%USERPROFILE%\.vscode\extensions\tormentnexus\" 2>nul
    echo   VS Code extension installed.
) else (
    echo   VS Code not detected - skipping.
)

echo.
echo ============================================================
echo TormentNexus installation complete!
echo.
echo Start the TN Kernel with: tormentnexus serve
echo Open dashboard at: http://localhost:7779
echo Cloud dashboard: https://cloud.hypernexus.site
echo ============================================================
pause
