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
	"strings"
	"time"
)

func HandleInitLLMAdapter(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	provider, _ :=getString(args, "provider")
	name, _ :=getString(args, "name")
	description, _ :=getString(args, "description")
	autoLogo, _ :=getBool(args, "auto_logo")

	// Validate required fields
	if provider == "" {
		return err("provider is required")
}

	if name == "" {
		return err("name is required")
}

	// Build command
	cmd := exec.CommandContext(ctx, "python", ".claude/skills/unstract-adapter-extension/scripts/init_llm_adapter.py",
		"--provider", provider,
		"--name", name,
		"--description", description,
	)

	if autoLogo {
		cmd.Args = append(cmd.Args, "--auto-logo")

	// Execute command
	output, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to initialize LLM adapter: %v\nOutput: %s", cmdErr, string(output)))
}

	return ok(fmt.Sprintf("Successfully initialized LLM adapter for %s\nOutput: %s", provider, string(output)))
}

}

func HandleInitEmbeddingAdapter(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	provider, _ :=getString(args, "provider")
	name, _ :=getString(args, "name")
	description, _ :=getString(args, "description")
	autoLogo, _ :=getBool(args, "auto_logo")

	// Validate required fields
	if provider == "" {
		return err("provider is required")
}

	if name == "" {
		return err("name is required")
}

	// Build command
	cmd := exec.CommandContext(ctx, "python", ".claude/skills/unstract-adapter-extension/scripts/init_embedding_adapter.py",
		"--provider", provider,
		"--name", name,
		"--description", description,
	)

	if autoLogo {
		cmd.Args = append(cmd.Args, "--auto-logo")

	// Execute command
	output, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to initialize embedding adapter: %v\nOutput: %s", cmdErr, string(output)))
}

	return ok(fmt.Sprintf("Successfully initialized embedding adapter for %s\nOutput: %s", provider, string(output)))
}

}

func HandleManageModels(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	adapterType, _ :=getString(args, "adapter")
	provider, _ :=getString(args, "provider")
	action, _ :=getString(args, "action")
	models, _ :=getString(args, "models")

	// Validate required fields
	if adapterType == "" {
		return err("adapter type is required")
}

	if provider == "" {
		return err("provider is required")
}

	if action == "" {
		return err("action is required")
}

	// Build command
	cmd := exec.CommandContext(ctx, "python", ".claude/skills/unstract-adapter-extension/scripts/manage_models.py",
		"--adapter", adapterType,
		"--provider", provider,
		"--action", action,
	)

	if models != "" {
		cmd.Args = append(cmd.Args, "--models", models)

	// Execute command
	output, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to manage models: %v\nOutput: %s", cmdErr, string(output)))
}

	return ok(fmt.Sprintf("Successfully managed models for %s %s\nOutput: %s", adapterType, provider, string(output)))
}

}

func HandleCheckAdapterUpdates(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Execute command
	cmd := exec.CommandContext(ctx, "python", ".claude/skills/unstract-adapter-extension/scripts/check_adapter_updates.py")

	output, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to check adapter updates: %v\nOutput: %s", cmdErr, string(output)))
}

	return ok(fmt.Sprintf("Adapter update check completed\nOutput: %s", string(output)))
}

func HandleGetAdaptersList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Execute command to get adapters list
	cmd := exec.CommandContext(ctx, "python", "-c", `
from unstract.sdk1.adapters.adapterkit import Adapterkit
kit = Adapterkit()
adapters = kit.get_adapters_list()
import json
print(json.dumps(adapters))
`)

	output, cmdErr := cmd.Output()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to get adapters list: %v", cmdErr))
}

	// Parse JSON output
	var adaptersList map[string]interface{}
	parseErr := json.Unmarshal(output, &adaptersList)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse adapters list: %v", parseErr))
}

	return ok(fmt.Sprintf("Adapters list:\n%s", string(output)))
}

func HandleValidateModel(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	provider, _ :=getString(args, "provider")
	model, _ :=getString(args, "model")

	// Validate required fields
	if provider == "" {
		return err("provider is required")
}

	if model == "" {
		return err("model is required")
}

	// Execute command to validate model
	cmd := exec.CommandContext(ctx, "python", "-c", fmt.Sprintf(`
from unstract.sdk1.adapters.base1 import %sParameters
print(%sParameters.validate_model({"model": "%s"}))
`, provider, provider, model))

	output, cmdErr := cmd.Output()
	if cmdErr != nil {
		return err(fmt.Sprintf("failed to validate model: %v", cmdErr))
}

	validatedModel := strings.TrimSpace(string(output))
	return ok(fmt.Sprintf("Validated model: %s", validatedModel))
}