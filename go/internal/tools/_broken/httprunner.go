package tools

import (
	"context"
	"fmt"
	"os/exec"
)

func ok(message string) (ToolResponse, error) {
	return ToolResponse{Content: []TextContent{{Type: "text", Text: message}}}, nil
}

func err(message string) (ToolResponse, error) {
	return ToolResponse{Content: []TextContent{{Type: "text", Text: message}}}, fmt.Errorf("%s", message)
}

func getString(args map[string]interface{}, key string) string {
	if val, found := args[key].(string); found {
		return val
	}
	return ""
}

func getInt(args map[string]interface{}, key string) int {
	if val, found := args[key].(float64); found {
		return int(val)
}

	if val, found := args[key].(int); found {
		return val
	}
	return 0
}

func getBool(args map[string]interface{}, key string) bool {
	if val, found := args[key].(bool); found {
		return val
	}
	return false
}

func HandleBugReport(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	description, _ :=getString(args, "description")
	versions, _ :=getString(args, "versions")
	projectFiles, _ :=getString(args, "project_files")

	report := fmt.Sprintf("Bug Report:\nDescription: %s\nVersions: %s\nProject Files: %s", description, versions, projectFiles)
	return ok(report)
}

func HandleFeatureRequest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	background, _ :=getString(args, "background")
	desiredFeatures, _ :=getString(args, "desired_features")
	exampleDescription, _ :=getString(args, "example_description")

	request := fmt.Sprintf("Feature Request:\nBackground: %s\nDesired Features: %s\nExample: %s", background, desiredFeatures, exampleDescription)
	return ok(request)
}

func HandleCommandExecution(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	command, _ :=getString(args, "command")
	cmd := exec.Command("sh", "-c", command)
	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return err(fmt.Sprintf("Command execution failed: %s", execErr.Error()))
}

	return ok(string(output))
}

func HandleVersionInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	goVersion, _ :=getString(args, "go_version")
	httpRunnerVersion, _ :=getString(args, "http_runner_version")

	versionInfo := fmt.Sprintf("Go Version: %s\nHttpRunner Version: %s", goVersion, httpRunnerVersion)
	return ok(versionInfo)
}