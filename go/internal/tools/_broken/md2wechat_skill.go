package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// md2wechat skill: convert Markdown to WeChat Official Account HTML,
// diagnose local readiness, discover/validate advanced layouts, and generate previews.

// HandleMd2wechatConvert converts a Markdown file to WeChat Official Account HTML.
// It shells out to the md2wechat CLI binary.
func HandleMd2wechatConvert(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	file, _ :=getString(args, "file")
	if file == "" {
		return err("file is required")
}

	mode, _ :=getString(args, "mode")
	if mode == "" {
		mode = "api"
	}
	theme, _ :=getString(args, "theme")
	output, _ :=getString(args, "output")
	provider, _ :=getString(args, "provider")

	cmdArgs := []string{"convert", file, "--mode", mode, "--json"}
	if theme != "" {
		cmdArgs = append(cmdArgs, "--theme", theme)

	if output != "" {
		cmdArgs = append(cmdArgs, "--output", output)

	if provider != "" {
		cmdArgs = append(cmdArgs, "--provider", provider)

	cmd := exec.CommandContext(ctx, "md2wechat", cmdArgs...)
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("md2wechat convert failed: %s: %s", cmdErr.Error(), string(out)))
}

	result := strings.TrimSpace(string(out))
	if result == "" {
		return ok("Conversion completed but produced no output.")
}

	return ok(result)
}

}
}
}

// HandleMd2wechatInspect inspects a Markdown article for readiness, metadata,
// and publish checks before converting or uploading.
func HandleMd2wechatInspect(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	file, _ :=getString(args, "file")
	if file == "" {
		return err("file is required")
}

	cmd := exec.CommandContext(ctx, "md2wechat", "inspect", file, "--json")
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("md2wechat inspect failed: %s: %s", cmdErr.Error(), string(out)))
}

	result := strings.TrimSpace(string(out))
	if result == "" {
		return ok("Inspection completed but produced no output.")
}

	return ok(result)
}

// HandleMd2wechatLayoutList discovers advanced layout modules available in
// the current md2wechat installation. Optionally filters by serves goal.
func HandleMd2wechatLayoutList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	serves, _ :=getString(args, "serves")

	cmdArgs := []string{"layout", "list", "--json"}
	if serves != "" {
		cmdArgs = append(cmdArgs, "--serves", serves)

	cmd := exec.CommandContext(ctx, "md2wechat", cmdArgs...)
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("md2wechat layout list failed: %s: %s", cmdErr.Error(), string(out)))
}

	result := strings.TrimSpace(string(out))
	if result == "" {
		return ok("No layout modules found.")
}

	return ok(result)
}

}

// HandleMd2wechatLayoutValidate validates advanced layout syntax (:::block)
// in a Markdown file and returns any errors or warnings.
func HandleMd2wechatLayoutValidate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	file, _ :=getString(args, "file")
	if file == "" {
		return err("file is required")
}

	cmd := exec.CommandContext(ctx, "md2wechat", "layout", "validate", "--file", file, "--json")
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("md2wechat layout validate failed: %s: %s", cmdErr.Error(), string(out)))
}

	result := strings.TrimSpace(string(out))
	if result == "" {
		return ok("Validation completed but produced no output.")
}

	return ok(result)
}

// HandleMd2wechatCapabilities queries the running md2wechat CLI for its
// current capabilities, providers, themes, and prompts — the discovery-first
// source of truth.
func HandleMd2wechatCapabilities(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	client := http.DefaultClient

	apiURL, _ :=getString(args, "api_url")
	if apiURL == "" {
		// Fall back to CLI discovery if no API URL is provided
		cmd := exec.CommandContext(ctx, "md2wechat", "capabilities", "--json")
		out, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			return err(fmt.Sprintf("md2wechat capabilities failed: %s: %s", cmdErr.Error(), string(out)))
}

		result := strings.TrimSpace(string(out))
		if result == "" {
			return ok("Capabilities query returned no output.")
}

		return ok(result)
}

	// Query the local API server for capabilities
	resp, fetchErr := client.Get(apiURL + "/capabilities")
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to reach md2wechat API at %s: %s", apiURL, fetchErr.Error()))
}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read capabilities response: %s", readErr.Error()))
}

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("capabilities API returned status %d: %s", resp.StatusCode, string(body)))
}

	// Validate JSON
	var parsed map[string]interface{}
	if jsonErr := json.Unmarshal(body, &parsed); jsonErr != nil {
		return ok(string(body))
}

	formatted, formatErr := json.MarshalIndent(parsed, "", "  ")
	if formatErr != nil {
		return ok(string(body))
}

	return ok(string(formatted))
}

// HandleMd2wechatPreview generates a read-only HTML confirmation page
// from inspect state for a given Markdown article.
func HandleMd2wechatPreview(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	file, _ :=getString(args, "file")
	if file == "" {
		return err("file is required")
}

	output, _ :=getString(args, "output")

	cmdArgs := []string{"preview", file, "--json"}
	if output != "" {
		cmdArgs = append(cmdArgs, "--output", output)

	cmd := exec.CommandContext(ctx, "md2wechat", cmdArgs...)
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("md2wechat preview failed: %s: %s", cmdErr.Error(), string(out)))
}

	result := strings.TrimSpace(string(out))
	if result == "" {
		return ok("Preview generated but produced no output.")
}

	return ok(result)
}

}

// Helper: check if md2wechat binary exists in PATH
func md2wechatBinaryExists() bool {
	_, lookupErr := exec.LookPath("md2wechat")
	return lookupErr == nil
}

// init ensures the binary is discoverable at load time (informational only)
func init() {
	if !md2wechatBinaryExists() {
		_, _ = fmt.Fprintf(os.Stderr, "warning: md2wechat binary not found in PATH; tools will fail at runtime\n")

}
}