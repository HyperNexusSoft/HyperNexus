package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runMd2wechatCLI runs the md2wechat CLI with the given arguments and returns stdout.
// Returns an error if the binary is not found or the command fails.
func runMd2wechatCLI(arg ...string) (string, error) {
	cmd := exec.Command("md2wechat", arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	execErr := cmd.Run()
	if execErr != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("md2wechat error: %v: %s", execErr, strings.TrimSpace(stderr.String()))
}

		return "", fmt.Errorf("md2wechat error: %v", execErr)
}

	if stderr.Len() > 0 {
		// some output may still be valid; we include stderr in returned string for diagnostic
		return strings.TrimSpace(out.String()) + "\n" + strings.TrimSpace(stderr.String()), nil
	}
	return strings.TrimSpace(out.String()), nil
}

// HandleMd2wechatConvert converts a Markdown file to WeChat Official Account HTML.
// Arguments:
//   - file (string, required): path to the Markdown file.
//   - output (string, optional): path for the output HTML file. Default: /tmp/md2wechat-output.html
//   - mode (string, optional): conversion mode "api" or "ai". Default: "api"
func HandleMd2wechatConvert(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	mdFile, _ :=getString(args, "file")
	if mdFile == "" {
		return err("missing required argument: 'file'")
}

	output, _ :=getString(args, "output")
	if output == "" {
		output = "/tmp/md2wechat-output.html"
	}

	mode, _ :=getString(args, "mode")
	if mode == "" {
		mode = "api"
	}

	cliArgs := []string{"convert", mdFile, "--mode", mode, "--output", output}
	rawResult, execErr := runMd2wechatCLI(cliArgs...)
	if execErr != nil {
		return err(execErr.Error())
}

	// If the CLI output is a JSON envelope (common for md2wechat), try to parse and re-wrap.
	var envelope map[string]interface{}
	if json.Unmarshal([]byte(rawResult), &envelope) == nil {
		// It is already JSON; forward as-is
		return ok(rawResult)
}

	// Otherwise treat rawResult as a plain success message
	result := map[string]interface{}{
		"message": fmt.Sprintf("Converted %s → %s (mode: %s)", mdFile, output, mode),
		"stdout":  rawResult,
		"output":  output,
	}
	jsonBytes, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		return ok(fmt.Sprintf("Conversion complete. Output: %s", output))
}

	return ok(string(jsonBytes))
}

// HandleMd2wechatCapabilities returns the JSON output of `md2wechat capabilities --json`.
// No arguments required.
func HandleMd2wechatCapabilities(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	rawResult, execErr := runMd2wechatCLI("capabilities", "--json")
	if execErr != nil {
		return err(execErr.Error())
}

	// Validate that it is valid JSON
	if json.Valid([]byte(rawResult)) {
		return ok(rawResult)
}

	return ok(fmt.Sprintf("Capabilities output:\n%s", rawResult))
}

// HandleMd2wechatLayoutList lists available advanced layout modules.
// Arguments:
//   - serves (string, optional): filter by goal, one of "attention", "readability", "memorability", "conversion".
func HandleMd2wechatLayoutList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cliArgs := []string{"layout", "list", "--json"}
	serves, _ :=getString(args, "serves")
	if serves != "" {
		cliArgs = append(cliArgs, "--serves", serves)

	rawResult, execErr := runMd2wechatCLI(cliArgs...)
	if execErr != nil {
		return err(execErr.Error())
}

	if json.Valid([]byte(rawResult)) {
		return ok(rawResult)
}

	return ok(fmt.Sprintf("Layout list:\n%s", rawResult))
}

}

// HandleMd2wechatLayoutValidate validates advanced layout syntax in a Markdown file.
// Arguments:
//   - file (string, required): path to the Markdown file to validate.
func HandleMd2wechatLayoutValidate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	mdFile, _ :=getString(args, "file")
	if mdFile == "" {
		return err("missing required argument: 'file'")
}

	cliArgs := []string{"layout", "validate", "--file", mdFile, "--json"}
	rawResult, execErr := runMd2wechatCLI(cliArgs...)
	if execErr != nil {
		return err(execErr.Error())
}

	if json.Valid([]byte(rawResult)) {
		return ok(rawResult)
}

	return ok(fmt.Sprintf("Validation output:\n%s", rawResult))
}

// HandleMd2wechatLayoutRender renders an advanced layout syntax block for a given module.
// Arguments:
//   - name (string, required): name of the layout module (e.g., "hero", "toc").
//   - vars (string, optional): a JSON object of variables to substitute (e.g., {"title":"Hello"}).
//     Alternatively, you can provide a string like "KEY=VALUE,KEY2=VALUE2".
func HandleMd2wechatLayoutRender(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	modName, _ :=getString(args, "name")
	if modName == "" {
		return err("missing required argument: 'name'")
}

	cliArgs := []string{"layout", "render", modName}

	varsRaw, _ :=getString(args, "vars")
	if varsRaw != "" {
		// Attempt to parse as JSON object
		var varsMap map[string]string
		if json.Unmarshal([]byte(varsRaw), &varsMap) == nil {
			for k, v := range varsMap {
				cliArgs = append(cliArgs, "--var", fmt.Sprintf("%s=%s", k, v))

		} else {
			// Fallback: treat as comma-separated KEY=VALUE pairs
			pairs := strings.Split(varsRaw, ",")
			for _, pair := range pairs {
				pair = strings.TrimSpace(pair)
				if pair == "" {
					continue
				}
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 2 {
					cliArgs = append(cliArgs, "--var", fmt.Sprintf("%s=%s", strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])))

			}
		}
	}

	rawResult, execErr := runMd2wechatCLI(cliArgs...)
	if execErr != nil {
		return err(execErr.Error())
}

	// The output may be rendered HTML/plain text, return as is.
	return ok(fmt.Sprintf("Rendered layout module '%s':\n%s", modName, rawResult))
}

}
}

// HandleMd2wechatPreview generates a preview HTML page from a Markdown file.
// This wraps `md2wechat convert` with preview-like arguments (read-only preview).
// Arguments:
//   - file (string, required): path to the Markdown file.
//   - output (string, optional): path for the preview HTML file. Default: /tmp/md2wechat-preview.html
func HandleMd2wechatPreview(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	mdFile, _ :=getString(args, "file")
	if mdFile == "" {
		return err("missing required argument: 'file'")
}

	output, _ :=getString(args, "output")
	if output == "" {
		output = "/tmp/md2wechat-preview.html"
	}

	// Use the same convert command; the HTML output serves as preview.
	cliArgs := []string{"convert", mdFile, "--mode", "api", "--output", output}
}