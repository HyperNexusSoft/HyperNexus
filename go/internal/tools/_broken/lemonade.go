package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Handlers for the lemonade MCP tool module

// HandleListRunners lists the self-hosted runners defined in the repository configuration.
// It reads the .github/runners.json file (simulated via env or hardcoded for this module)
// and returns the list of runner names.
func HandleListRunners(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// In a real scenario, this might read from a file or fetch from GitHub API.
	// Based on the provided .github/runners.json, we return the static list.
	// To make it dynamic, we could check for an env var override, but we stick to the source data.
	runners := []string{
		"jfowers-stx-01",
		"jfowers-stx-02",
		"jfowers-phx-rtx-00",
		"jfowers-stx-rtx-00",
		"sjlab-stx-0",
		"sjlab-stx-halo-06",
		"sjlab-stx-halo-07",
		"sjlab-stx-halo-08",
		"sjlab-stx-halo-09",
		"sjlab-stx-halo-10",
		"sjlab-stx-halo-11",
		"sjlab-stx-halo-12",
		"sjlab-stx-halo-13",
		"sjlab-stx-halo-17",
	}

	// Sort for deterministic output
	sort.Strings(runners)

	result, marshalErr := json.MarshalIndent(runners, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal runners: %v", marshalErr))
}

	return ok(string(result))
}

// HandleAutoLabel simulates the auto-labeling logic described in auto_label.py.
// It takes an issue number and a repo, fetches the issue details (simulated or via gh CLI),
// and returns the suggested labels based on the system prompt logic.
// Since we cannot call the Anthropic API or gh CLI directly in this isolated environment
// without external dependencies or specific environment setup, we implement a heuristic
// classifier based on the rules provided in the Python script's SYSTEM_PROMPT.
func HandleAutoLabel(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	numStr, _ :=getString(args, "number")
	repo, _ :=getString(args, "repo")
	if numStr == "" {
		return err("number argument is required")
}

	num, parseErr := strconv.Atoi(numStr)
	if parseErr != nil {
		return err(fmt.Sprintf("invalid issue number: %v", parseErr))
}

	title, _ :=getString(args, "title")
	body, _ :=getString(args, "body")
	existingLabels, _ :=getString(args, "existing_labels")

	// If we had the gh CLI and API key, we would fetch here.
	// For this implementation, we assume the caller provides title/body or we simulate a fetch.
	// If title is empty, we simulate a fetch attempt (which would fail without gh CLI).
	if title == "" {
		// Attempt to run gh CLI if available
		cmd := exec.CommandContext(ctx, "gh", "issue", "view", strconv.Itoa(num), "--json", "title,body,labels")
		if repo != "" {
			cmd.Args = append(cmd.Args, "--repo", repo)

		output, runErr := cmd.Output()
		if runErr != nil {
			// Fallback: return an error or a generic response if gh is not available
			return err(fmt.Sprintf("could not fetch issue %d (gh CLI not available or auth missing): %v", num, runErr))
}

		var issueData struct {
			Title  string `json:"title"`
			Body   string `json:"body"`
			Labels []struct {
				Name string `json:"name"`
			} `json:"labels"`
		}
		if jsonErr := json.Unmarshal(output, &issueData); jsonErr != nil {
			return err(fmt.Sprintf("failed to parse gh output: %v", jsonErr))
}

		title = issueData.Title
		body = issueData.Body
		existingLabels = ""
		for i, l := range issueData.Labels {
			if i > 0 {
				existingLabels += ","
			}
			existingLabels += l.Name
		}
	}

	// Heuristic Classification Logic (Simulating the LLM prompt rules)
	// This is a simplified version of the logic in auto_label.py's SYSTEM_PROMPT
	labelsToAdd := classifyIssueHeuristic(title, body, existingLabels)

	result, marshalErr := json.Marshal(labelsToAdd)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal labels: %v", marshalErr))
}

	return ok(string(result))
}

}

// classifyIssueHeuristic implements a simplified version of the LLM classification logic.
func classifyIssueHeuristic(title, body, existingLabels string) []string {
	lowerTitle := strings.ToLower(title)
	lowerBody := strings.ToLower(body)
	fullText := lowerTitle + " " + lowerBody
	existingSet := make(map[string]bool)
	for _, l := range strings.Split(existingLabels, ",") {
		l = strings.TrimSpace(l)
		if l != "" {
			existingSet[l] = true
		}
	}

	var result []string

	// Helper to add if not exists
	addIfNew := func(label string) {
		if !existingSet[label] {
			result = append(result, label)
			existingSet[label] = true
		}
	}

	// Engine Detection (At most one)
	if strings.Contains(fullText, "llama") || strings.Contains(fullText, "llamacpp") {
		addIfNew("engine::llamacpp")
	} else if strings.Contains(fullText, "fastflow") || strings.Contains(fullText, "flm") {
		addIfNew("engine::flm")
	} else if strings.Contains(fullText, "ryzen") || strings.Contains(fullText, "ryzenai") {
		addIfNew("engine::ryzenai")
	} else if strings.Contains(fullText, "vllm") {
		addIfNew("engine::vllm")
	} else if strings.Contains(fullText, "whisper") {
		addIfNew("engine::whispercpp")
	} else if strings.Contains(fullText, "stable-diffusion") || strings.Contains(fullText, "sd-cpp") {
		addIfNew("engine::sd")
	} else if strings.Contains(fullText, "kokoro") {
		addIfNew("engine::kokoro")
	} else if strings.Contains(fullText, "moonshine") {
		addIfNew("engine::moonshine")

	// Area Detection (At most one)
	if strings.Contains(fullText, "cli") && !strings.Contains(fullText, "api") {
		addIfNew("area::cli")
	} else if strings.Contains(fullText, "msi") || strings.Contains(fullText, "dmg") || strings.Contains(fullText, "deb") || strings.Contains(fullText, "rpm") || strings.Contains(fullText, "installer") {
		addIfNew("area::installer")
	} else if strings.Contains(fullText, "api") || strings.Contains(fullText, "rest") || strings.Contains(fullText, "endpoint") {
		addIfNew("area::api")
	} else if strings.Contains(fullText, "tray") || strings.Contains(fullText, "system tray") {
		addIfNew("area::tray")
	} else if strings.Contains(fullText, "ci") || strings.Contains(fullText, "github action") || strings.Contains(fullText, "workflow") {
		addIfNew("area::ci")

	// Runtime Detection (At most one)
	if strings.Contains(fullText, "vulkan") {
		addIfNew("runtime::vulkan")
	} else if strings.Contains(fullText, "rocm") {
		addIfNew("runtime::rocm")
	} else if strings.Contains(fullText, "cuda") {
		addIfNew("runtime::cuda")
	} else if strings.Contains(fullText, "metal") {
		addIfNew("runtime::metal")
	} else if strings.Contains(fullText, "cpu only") || strings.Contains(fullText, "cpu fallback") {
		addIfNew("runtime::cpu")

	// Component Labels
	if strings.Contains(fullText, "tauri") || strings.Contains(fullText, "desktop app") {
		addIfNew("app")
	} else if strings.Contains(fullText, "web app") || strings.Contains(fullText, "web ui") {
		addIfNew("web ui")
	} else if strings.Contains(fullText, "audio") || strings.Contains(fullText, "transcription") || strings.Contains(fullText, "tts") {
		addIfNew("audio")
	} else if strings.Contains(fullText, "c++") && !strings.Contains(fullText, "area::api") && !strings.Contains(fullText, "area::cli") {
		addIfNew("cpp")

	// Type Labels
	if strings.Contains(fullText, "bug") || strings.Contains(fullText, "crash") || strings.Contains(fullText, "error") || strings.Contains(fullText, "fix") {
		addIfNew("bug")
	} else if strings.Contains(fullText, "feature") || strings.Contains(fullText, "enhancement") || strings.Contains(fullText, "improve") {
		addIfNew("enhancement")
	} else if strings.Contains(fullText, "readme") || strings.Contains(fullText, "doc") || strings.Contains(fullText, "guide") {
		addIfNew("documentation")
	} else if strings.Contains(fullText, "question") || strings.Contains(fullText, "how to") {
		addIfNew("question")

	return result
}

}
}
}
}
}

// HandleRenderSdCppPRBody simulates the logic of render_sdcpp_pr_body.py.
// It takes validation data (simulated via JSON string) and generates a PR body.
func HandleRenderSdCppPRBody(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	jsonData, _ :=getString(args, "validation_data")
	baseRelease, _ :=getString(args, "base_release")
	cudaRelease, _ :=getString(args, "cuda_release")
	prompt, _ :=getString(args, "prompt")
	seed, _ :=getString(args, "seed")
	steps, _ :=getString(args, "steps")
	models, _ :=getString(args, "models")
	sizes, _ :=getString(args, "sizes")

	if jsonData == "" {
		return err("validation_data is required")
}

	var records []map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(jsonData), &records); unmarshalErr != nil {
		return err(fmt.Sprintf("invalid validation_data JSON: %v", unmarshalErr))
}

	// Process records to find passed images and calculate stats
	var passedRecords []map[string]interface{}
	var totalElapsed float64
	var backendLabels []string
	seenBackends := make(map[string]bool)

	for _, row := range records {
		pass, _ := row["pass"].(bool)
		if pass {
			passedRecords = append(passedRecords, row)
			if elapsed, found := row["elapsed_s"].(float64); found {
				totalElapsed += elapsed
			}
			if label, found := row["label"].(string); ok && !seenBackends[label] {
				backendLabels = append(backendLabels, label)
				seenBackends[label] = true
			}
		}
	}

	sort.Strings(backendLabels)

	// Construct PR Body
	var body strings.Builder
	body.WriteString("This updates Lemonade's existing `sd-cpp` backend pins using the appropriate release stream for each backend.\n\n")
	body.WriteString(fmt.Sprintf("- Base sd-cpp release (`leejet/stable-diffusion.cpp`): `%s`\n", baseRelease))
	body.WriteString(fmt.Sprintf("- CUDA sd-cpp release (`lemonade-sdk/stable-diffusion.cpp`): `%s`\n", cudaRelease))
	body.WriteString(fmt.Sprintf("- Validated backends: `%s`\n", strings.Join(backendLabels, ", ")))
	body.WriteString("- Timing: `request wall time`\n\n")
	body.WriteString("Validated with the same prompt and seed on each available runner:\n\n")
	body.WriteString(fmt.Sprintf("- Prompt: `%s`\n", prompt))
	body.WriteString(fmt.Sprintf("- Seed / steps: `%s` / `%s`\n", seed, steps))
	body.WriteString(fmt.Sprintf("- Models: `%s`\n", models))
	body.WriteString(fmt.Sprintf("- Sizes: `%s`\n\n", sizes))

	body.WriteString("| Backend | Result | Time (s) |\n")
	body.WriteString("|---|---|---|\n")

	for _, row := range passedRecords {
		label, _ :=getString(row, "label")
		elapsed := 0.0
		if e, found := row["elapsed_s"].(float64); found {
			elapsed = e
		}
		body.WriteString(fmt.Sprintf("| %s | ✅ | %.2f |\n", label, elapsed))

	return ok(body.String())
}

}

// HandleSanitizeJSON simulates the script_safe_json function from test_triage_dashboard.py.
// It takes a JSON string and ensures it is safe for embedding in HTML/JS by escaping
// dangerous characters like <, >, &, and line separators.
func HandleSanitizeJSON(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	inputJSON, _ :=getString(args, "json_input")
	if inputJSON == "" {
		return err("json_input is required")
}

	// Parse to ensure it's valid JSON first
	var data interface{}
	if unmarshalErr := json.Unmarshal([]byte(inputJSON), &data); unmarshalErr != nil {
		return err(fmt.Sprintf("invalid JSON input: %v", unmarshalErr))
}

	// Re-marshal with standard escaping, then manually escape < > & and line separators
	// to ensure safety in HTML/JS contexts.
	// Note: json.Marshal already escapes < > & as \u003c, \u003e, \u0026 in Go 1.15+
	// but we explicitly ensure line separators are handled if the input contained them raw.
	// The Python script uses a custom encoder. We will use standard json.Marshal which
	// handles the critical characters for HTML/JS safety.
	
	// However, to be strictly compliant with the Python script's behavior regarding
	// U+2028 and U+2029 which standard JSON libraries sometimes miss in JS contexts:
	// We will re-encode the string manually if needed, but json.Marshal in Go
	// usually handles the critical HTML chars.
	// Let's do a manual pass to be safe and explicit about U+2028/U+2029.
	
	outputBytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to re-marshal: %v", marshalErr))
}

	outputStr := string(outputBytes)
	
	// Replace U+2028 (Line Separator) and U+2029 (Paragraph Separator)
	// These are valid in JSON but break JS string literals.
	outputStr = strings.ReplaceAll(outputStr, "\u2028", "\\u2028")
	outputStr = strings.ReplaceAll(outputStr, "\u2029", "\\u2029")

	return ok(outputStr)
}

// Helper to get string from args
func getString(args map[string]interface{}, key string) string {
	if val, found := args[key]; found {
		if s, found := val.(string); found {
			return s
		}
		return fmt.Sprintf("%v", val)
}

	return ""
}

// Helper to get int from args
func getInt(args map[string]interface{}, key string) int {
	if val, found := args[key]; found {
		switch v := val.(type) {
		case int:
			return v
}
		case float64:
			return int(v)
}
		case string:
			i, _ := strconv.Atoi(v)
			return i
		}
	}
	return 0
}

// Helper to get bool from args
func getBool(args map[string]interface{}, key string) bool {
	if val, found := args[key]; found {
		switch v := val.(bool) {
		case true:
			return true
}
		case false:
			return false
		}
		if s, found := val.(string); found {
			return strings.ToLower(s) == "true"
		}
	}
	return false
}

// Helper to get string from map[string]interface{} (for nested access if needed)
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, found := m[key]; found {
		if s, found := val.(string); found {
			return s
		}
		return fmt.Sprintf("%v", val)

	return ""
}