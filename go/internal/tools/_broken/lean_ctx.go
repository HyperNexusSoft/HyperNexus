package tools
Then imports:
import (
	"context"
	"fmt"
	"os/exec"
)
Then the runLeanCtx function, same as original, right? Let's check: yes, that's okay.
Then HandleCtxRead: same as original? Wait wait the original HandleCtxRead returns ok(output), which is correct. Wait but wait the user said "Undefined function: ok in HandleCtxRead function" — but that's pre-defined, so no problem, we don't need to define it.
Then HandleCtxShell: fix the if runErr != nil block: add return err(runErr.Error()), then add return ok(output) at the end of the function.
Wait wait let's check the original HandleCtxShell code:
Original:
func HandleCtxShell(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	command, _ :=getString(args, "command")
	if command == "" {
		return err("command is required")
}

	// Using -c flag as per documentation: lean-ctx -c "git status"
	output, runErr := runLeanCtx(ctx, "-c", command)
	if runErr != nil {
	}
}
So we need to fill the if block: return err(runErr.Error()), then after that, return ok(output). That's it.
Wait wait but wait, is there any other error? Let's check: the return types are correct: both Handle functions return (ToolResponse, error), which matches rule 2. Rule 3: return ok("text"), if e != nil { return err(e.Error()) } — which we are doing: in HandleCtxRead, if runErr != nil return err(runErr.Error()), else return ok(output). Same for HandleCtxShell. That's correct.
Wait what about the runLeanCtx function: it returns (string, error), which is correct. It uses exec.CommandContext, which is from os/exec, which is imported. fmt is used for fmt.Errorf, correct.
Wait let's check for unused imports again: after removing os and strings, do we have any unused? context is used in the function parameters, fmt is used in runLeanCtx, os/exec is used for exec.CommandContext. All used. Good.
Wait wait, but wait: the user said "PRE-DEFINED in parity.go (do NOT redeclare): - type ToolResponse, func ok(), func err("error"), func getString(), func getInt(), func getBool(), type TextContent". So we don't need to define any of those, just use them, which we are doing. Correct.
Wait let's write the full code now:
package tools
import (
	"context"
	"fmt"
	"os/exec"
)
// runLeanCtx executes the lean-ctx binary with the given arguments
func runLeanCtx(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "lean-ctx", args...)
	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return "", fmt.Errorf("lean-ctx execution failed: %w, output: %s", execErr, string(output))
}

	return string(output), nil
}
// HandleCtxRead reads a file with compression modes
func HandleCtxRead(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("path is required")
}

	mode, _ :=getString(args, "mode")
	if mode == "" {
		mode = "full"
	}
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "read", path)
	if mode != "" {
		cmdArgs = append(cmdArgs, "--mode", mode)

	output, runErr := runLeanCtx(ctx, cmdArgs...)
	if runErr != nil {
		return err(runErr.Error())
}

	return ok(output)
}

}

// HandleCtxShell executes a shell command with output compression
func runLeanCtx(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "lean-ctx", args...)
	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return "", fmt.Errorf("lean-ctx execution failed: %w, output: %s", execErr, string(output))
}

	return string(output), nil
}
// HandleCtxRead reads a file with compression modes