package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ARIS project structure initialization and research wiki management

func HandleInitializeProject(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	projectDir, _ :=getString(args, "project_dir")
	if projectDir == "" {
		return err("project_dir is required")
}

	// Create project directory
	if mkErr := os.MkdirAll(projectDir, 0755); mkErr != nil {
		return err(fmt.Sprintf("failed to create project dir: %v", mkErr))
}

	// Initialize git repository
	gitCmd := exec.CommandContext(ctx, "git", "init")
	gitCmd.Dir = projectDir
	if output, gitErr := gitCmd.CombinedOutput(); gitErr != nil {
		return err(fmt.Sprintf("git init failed: %s: %v", string(output), gitErr))
}

	// Create CLAUDE.md
	claudeMd := filepath.Join(projectDir, "CLAUDE.md")
	claudeContent := "# ARIS Research Project\n\nThis project uses Auto-claude-code-research-in-sleep (ARIS) for research workflows.\n"
	if writeErr := os.WriteFile(claudeMd, []byte(claudeContent), 0644); writeErr != nil {
		return err(fmt.Sprintf("failed to create CLAUDE.md: %v", writeErr))
}

	// Create .claude directory structure
	claudeDir := filepath.Join(projectDir, ".claude")
	skillsDir := filepath.Join(claudeDir, "skills")
	if mkdirErr := os.MkdirAll(skillsDir, 0755); mkdirErr != nil {
		return err(fmt.Sprintf("failed to create .claude/skills: %v", mkdirErr))
}

	// Create .aris directory
	arisDir := filepath.Join(projectDir, ".aris")
	if mkdirErr := os.MkdirAll(arisDir, 0755); mkdirErr != nil {
		return err(fmt.Sprintf("failed to create .aris: %v", mkdirErr))
}

	return ok(fmt.Sprintf("Initialized ARIS project at %s with git, CLAUDE.md, .claude/skills/, .aris/", projectDir))
}

func HandleInitializeResearchWiki(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	projectDir, _ :=getString(args, "project_dir")
	if projectDir == "" {
		return err("project_dir is required")
}

	wikiDir := filepath.Join(projectDir, "research-wiki")

	// Create directory structure
	dirs := []string{
		wikiDir,
		filepath.Join(wikiDir, "papers"),
		filepath.Join(wikiDir, "ideas"),
		filepath.Join(wikiDir, "experiments"),
		filepath.Join(wikiDir, "claims"),
		filepath.Join(wikiDir, "graph"),
	}

	for _, dir := range dirs {
		if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
			return err(fmt.Sprintf("failed to create %s: %v", dir, mkdirErr))

	}

	// Create initial files
	files := map[string]string{
		filepath.Join(wikiDir, "README.md"):          "# Research Wiki\n\nThis wiki contains research papers, ideas, experiments, and claims.\n",
		filepath.Join(wikiDir, "papers", ".gitkeep"): "",
		filepath.Join(wikiDir, "ideas", ".gitkeep"): "",
		filepath.Join(wikiDir, "experiments", ".gitkeep"): "",
		filepath.Join(wikiDir, "claims", ".gitkeep"): "",
		filepath.Join(wikiDir, "graph", ".gitkeep"): "",
	}

	for path, content := range files {
		if writeErr := os.WriteFile(path, []byte(content), 0644); writeErr != nil {
			return err(fmt.Sprintf("failed to create %s: %v", path, writeErr))

	}

	return ok(fmt.Sprintf("Initialized research wiki at %s", wikiDir))
}
}
}