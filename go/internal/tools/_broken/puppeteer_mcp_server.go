package tools`.
        *   Handler signature: `func HandleXxx(ctx context.Context, args map[string]interface{}) (ToolResponse, error)`.
        *   Success/Error helpers: `ok("text")`, `err("msg")` (from `parity.go`).
        *   Input helpers: `getString`, `getInt`, `getBool` (from `parity.go`).
        *   Imports: ONLY standard library (`context`, `encoding/json`, `fmt`, `io`, `net/http`, `net/url`, `os`, `os/exec`, `path/filepath`, `strconv`, `strings`, `time`, `regexp`, `sort`).
        *   No external packages (no `chromedp`, no `rod`, no `puppeteer` Go bindings). This is the tricky part. The prompt asks to reimplement the *functionality* of the Puppeteer server in Go, but restricts me to the standard library.
        *   *Self-Correction/Refinement:* Since I cannot use a browser automation library (like `chromedp`) due to the "NO external packages" rule, I have to simulate the *interface* of the tools or implement a very basic version using `os/exec` to launch Chrome/Chromium with remote debugging and interact via the CDP (Chrome DevTools Protocol) over WebSockets using `net/http` and `encoding/json`.
        *   However, implementing a full CDP client in pure Go with just stdlib in a single file is extremely complex and error-prone for a "clean compilable" snippet.
        *   *Alternative Interpretation:* The prompt says "Implement a Go-native MCP tool module... reimplement EVERY tool/function". It also says "No external packages
No external packages
No pseudocode, no TODOs. Every function MUST compile."
        *   If I can't do WebSockets in stdlib, I can't really control the browser.
        *   *Pivot:* I will implement the tools to *launch* the browser and *connect* to it (via HTTP to check version), and for the specific actions (screenshot, click), I will implement the logic to construct the CDP commands, but since I can't open a WebSocket in pure stdlib easily, I might have to rely on `os/exec` to run a helper script or just return a message explaining that the specific action requires a WebSocket client which isn't in stdlib.
        *   *Better Pivot:* I will implement the tools that *manage* the browser (launch, connect) and the *navigation* (which can sometimes be done via HTTP or just launching with a URL). For the DOM manipulation, I will simulate the success response or use `os/exec` to invoke a CLI tool if one existed.
        *   *Wait, the prompt asks to reimplement the server.* If I can't do it fully, I fail the "reimplement EVERY tool
Go-native" way to do it without libraries.
        *   For `screenshot`, I can launch `google-chrome --headless --screenshot=... <url>`.
        *   For `click`, `fill`, etc., these require stateful interaction. Without a library, this is impossible in a single file without writing a massive WebSocket implementation.
        *   *However*, the prompt says "You write clean compilable Go code... No external packages."
        *   I will implement the tools using `os/exec` to control Chrome via command-line arguments where possible (Navigate, Screenshot). For the interactive ones (Click, Fill), I will implement a "mock" or "placeholder" response that explains the limitation, OR I will try to implement a very basic HTTP interaction if the browser exposes a REST API (it doesn't really, it's CDP).
        *   *Actually*, maybe I can use `os/exec` to run a Node.js script? No, that defeats the purpose of a Go-native server.
        *   *Let's try to implement a basic CDP HTTP client.* CDP has a `/json` endpoint. It doesn't have a REST API for clicking.
        *   *Okay, I will focus on the tools that can be done via CLI arguments and `os/exec`.*
            *   `puppeteer_navigate`: Launch Chrome with URL.
            *   `puppeteer_screenshot`: Launch Chrome with screenshot flag.
            *   `puppeteer_connect_active_tab`: Check `http://localhost:9222/json/version`.
        *   For the others, I will return an error stating "Interactive DOM manipulation requires a WebSocket client not available in standard library." This is the only honest way to fulfill "No external packages" and "Must compile".
        *   *Wait*, the prompt says "reimplement EVERY tool". Returning an error for half of them is a failure.
        *   *Is there ANY way?* Maybe I can use `net/rpc`? No.
        *   *Let's look at the "No external packages" rule again.* "NO external packages (no github.com imports)
No external packages" implies stdlib only.
        *   *Let's assume I can use `os/exec` to drive Chrome.*
        *   I will implement a state manager that keeps track of the `exec.Cmd`.
        *   For `click`, `fill`, etc., I will construct a JavaScript snippet and try to execute it? How? I can't inject JS without a connection.
        *   *Okay, I will implement the tools as best as possible using `os/exec` for launching/navigating/screenshots, and for the others, I will implement a basic HTTP-based CDP client if possible, or just stub them.*
        *   *Actually, let's look at the "puppeteer_connect_active_tab" logic.* It connects to a running Chrome.
        *   If I connect to a running Chrome, I can use the `/json/protocol` descriptor to understand the domain, but I still need a WebSocket.
        *   *Decision:* I will implement the `puppeteer_navigate` and `puppeteer_screenshot` using `os/exec` flags. This is robust and stdlib-only.
        *   For `puppeteer_connect_active_tab`, I will implement the HTTP check.
        *   For the rest (`click`, `fill`, `select`, `hover`, `evaluate`), I will implement a handler that attempts to use a hypothetical WebSocket connection, but since I can't implement a full WS client in this snippet, I will have to return a specific error or a "simulated" success message if the prompt allows "simulated" behavior.
        *   *Strict interpretation:* "reimplement EVERY tool".
        *   I will implement a minimal WebSocket client in the file. It's about 100 lines of code. It's doable.
        *   *Plan:*
            1.  Define a simple `WebSocket` struct wrapping `net.Conn`.
            2.  Implement the Handshake.
            3.  Implement Send/Receive (framing).
            4.  Use this to talk to CDP.
            5.  Implement the tools using this CDP client.

    *   *WebSocket Implementation in Stdlib:*
        *   It's tricky but possible. I need to handle the HTTP upgrade request, masking, and framing.
        *   Given the complexity, this might be too much code for the "clean compilable" constraint.
        *   *Alternative:* The prompt says "No external packages". It doesn't say
. But usually, "No external packages" means stdlib.
        *   *Let's try the `os/exec` approach for everything.*
        *   Can I execute JS via CLI? `chrome --headless --dump-dom <url>` dumps the DOM. `chrome --headless --virtual-time-budget=1000 --repl` allows interaction.
        *   *Idea:* I can launch Chrome in `--headless --remote-debugging-port=0` (random port) and then connect to it.
        *   *Okay, I will implement a minimal CDP WebSocket client.* It is the only way to fulfill the requirements faithfully.
        *   *Wait*, the prompt says "You write clean compilable Go code... No pseudocode... No external packages."
        *   If I write a 200-line WebSocket implementation, it might be rejected for being too complex or "bloat".
        *   *Let's look at the provided source code again.* It uses `puppeteer`.
        *   *Maybe I can just implement the HTTP parts and return "Not implemented in stdlib" for the WS parts?*
        *   "reimplement EVERY tool".
        *   *Okay, I will implement the tools using `os/exec` to launch Chrome with specific arguments for the actions that support it (Navigate, Screenshot).*
        *   For the actions that don't (Click, Fill), I will implement a "Simulated" response that logs what *would* happen, or I will try to use the `--repl` mode of Chrome to pipe commands.
        *   *Actually, Chrome's `--repl` is for debugging, not automation.*
        *   *Let's go with the `os/exec` approach for Navigate and Screenshot, and a "Mock" implementation for the rest, explaining the limitation.*
        *   *Wait, the prompt says "reimplement EVERY tool".*
        *   I will implement the tools to *launch* the browser with the specific URL.
        *   `puppeteer_navigate`: Launches browser with URL.
        *   `puppeteer_screenshot`: Launches browser with `--screenshot`.
        *   `puppeteer_click`: Returns error "Requires WebSocket client".
        *   This seems weak.

    *   *Let's try to implement the WebSocket client.* It's the only way to be "correct".
    *   I'll write a `cdpClient` struct.
    *   It will connect to `ws://localhost:9222/devtools/page/...`.
    *   It will send JSON-RPC messages.
    *   I'll limit the implementation to the specific domains needed: `Page`, `Runtime`, `Input`, `DOM`.
    *   This is a big task for a single file, but I can condense it.

    *   *Refined Plan:*

    *   *WebSocket Implementation Details:*
        *   Connect via `net.Dial`.
        *   Send HTTP Upgrade request.
        *   Read response.
        *   Loop: Read frame header (2 bytes), parse length, read payload, unmask.
        *   Send frame: Construct header, mask payload, send.
        *   This is feasible in ~100 lines.

    *   *Wait*, `crypto/rand` is needed for the WebSocket mask key.
    *   *Is `crypto` allowed?* Yes, it's stdlib.

    *   *Let's refine the "No external packages" constraint.* It usually implies "don't make me `go get` anything". `crypto` is fine.

    *   *Code Structure:*
        *   `wsDial`: Connects and upgrades.
        *   `wsRead`: Reads a message.
        *   `wsWrite`: Writes a message.
        *   `cdpSend`: Sends a command.
        *   `cdpWait`: Waits for a result.
        *   Global `currentCDP` struct.

    *   *Handling "Active Tab":*
        *   The Node code connects to an existing browser.
        *   I need to support that.
        *   I also need to support launching a new one if not connected?
        *   The Node code has `ensureBrowser` which launches a new one.
        *   I'll implement a `ensureBrowser` that launches `chrome --headless --remote-debugging-port=9222` if one isn't found/connected.

    *   *Let's simplify.*
        *   I will implement the handlers assuming the browser is running on port 9222 (default).
        *   I will add a helper to launch it if needed.
        *   I will implement the WebSocket logic.

    *   *Wait, the prompt says "No pseudocode".*
    *   Writing a full WebSocket client from scratch is risky for a single prompt response (bugs, edge cases).
    *   *Is there a simpler way?*
    *   What if I just use `os/exec` to run `chrome` with `--dump-dom` or `--print-to-pdf`?
    *   No, that doesn't support "Click".
    *   *Okay, I will implement the WebSocket client.* It's the only way to satisfy "reimplement EVERY tool".

    *   *Drafting the WebSocket Client (Mental Check):*
        *   `key := base64(rand(16))`
        *   `req := fmt.Sprintf("GET %s HTTP/1.1\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path, key)`
        *   `conn.Write([]byte(req))`
        *   `reader.ReadUntil(\r\n\r\n)`
        *   Check `101 Switching Protocols`.
        *   *Read Frame:*
            *   Byte 0: FIN (bit 7), RSV (bits 6-4), Opcode (bits 3-0).
            *   Byte 1: Mask (bit 7), Payload Len (bits 6-0).
            *   Extended Len (2 or 8 bytes) if Payload Len is 126 or 127.
            *   Masking Key (4 bytes) if Mask is 1.
            *   Payload.
        *   *Write Frame:*
            *   Byte 0: FIN=1, Opcode=1 (Text).
            *   Byte 1: Mask=1, Len.
            *   Mask Key.
            *   Payload XOR Mask.

    *   *This is doable.*

    *   *Tool Logic:*
        *   `HandlePuppeteerConnectActiveTab`:
            *   Fetch `http://localhost:9222/json`.
            *   Parse JSON to find a tab.
            *   Connect WS to `webSocketDebuggerUrl`.
            *   Enable domains (`Page.enable`, `Runtime.enable`, `DOM.enable`).
        *   `HandlePuppeteerNavigate`:
            *   Send `{"id":1,"method":"Page.navigate","params":{"url":"..."}}`.
            *   Wait for response.
        *   `HandlePuppeteerScreenshot`:
            *   Send `{"id":2,"method":"Page.captureScreenshot"}`.
            *   Decode Base64 data.
            *   Save to file? The Node code saves to memory or returns data.