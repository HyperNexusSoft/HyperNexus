package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func (r *ToolResponse) ok(v string) {
	r.Ok = v
}

func (r *ToolResponse) err(v string) {
	r.Err = v
}

func (r *ToolResponse) getString() string {
	if r.Ok != "" {
		return r.Ok
	}
	return r.Err
}

func (r *ToolResponse) getInt() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (r *ToolResponse) getBool() (bool, error) {
	return false, fmt.Errorf("not implemented")
}

type TextContent string

func HandleAstGrepSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pattern, _ :=getString(args, "pattern")
	if pattern == "" {
		return err("pattern is required")
}

	rule, _ :=getString(args, "rule")
	path, _ :=getString(args, "path")
	if path == "" {
		path = "."
	}

	if _, statErr := os.Stat(path); statErr != nil {
		return err(fmt.Sprintf("invalid path: %v", statErr))
}

	cmdArgs := []string{"--pattern", pattern}
	if rule != "" {
		cmdArgs = append(cmdArgs, "--rule", rule)

	cmdArgs = append(cmdArgs, path)

	cmd := exec.Command("sg", cmdArgs...)
	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return err(fmt.Sprintf("ast-grep failed: %v - output: %s", execErr, string(output)))
}

	return ok(string(output))
}

// 他の関数も同様に修正
}