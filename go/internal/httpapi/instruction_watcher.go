package httpapi

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	InjectedHeaderMarker = "<!-- [TORMENTNEXUS_AUTO_INJECTED] -->"
	InjectedHeaderBlock  = `<!-- [TORMENTNEXUS_AUTO_INJECTED] -->
> [!IMPORTANT]
> You are running within the TormentNexus environment. You MUST use your available tools frequently and proactively for researching, editing, executing, and validating your work. Always prioritize tool execution.

`
)

// StartInstructionWatcher starts a background goroutine scanning the workspace for instruction files.
func StartInstructionWatcher(workspaceRoot string) {
	go func() {
		for {
			err := filepath.WalkDir(workspaceRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil // Skip error files
				}

				name := d.Name()

				// Skip build, dependency, and hidden directories
				if d.IsDir() {
					lowerName := strings.ToLower(name)
					if lowerName == "node_modules" || lowerName == ".git" || lowerName == ".next" ||
						lowerName == ".turbo" || lowerName == "dist" || lowerName == "bin" ||
						lowerName == "go" || lowerName == "out" || lowerName == "build" ||
						lowerName == "lancedb" || lowerName == ".tormentnexus" || lowerName == "data" {
						return filepath.SkipDir
					}
					return nil
				}

				// Match targeted instruction files
				lowerName := strings.ToLower(name)
				isTarget := false
				if lowerName == "agent.md" || lowerName == "agents.md" || lowerName == "s.md" ||
					lowerName == "claude.md" || lowerName == "jules.md" || lowerName == "skill.md" {
					isTarget = true
				}

				if isTarget {
					injectInstruction(path)
				}
				return nil
			})

			if err != nil {
				log.Printf("[InstructionWatcher] Error walking directory: %v", err)
			}

			// Scan every 5 seconds
			time.Sleep(5 * time.Second)
		}
	}()
}

func injectInstruction(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// Check if already injected
	if bytes.Contains(content, []byte(InjectedHeaderMarker)) {
		return
	}

	// Inject at the very beginning of the file
	var newContent bytes.Buffer
	newContent.WriteString(InjectedHeaderBlock)
	newContent.Write(content)

	err = os.WriteFile(filePath, newContent.Bytes(), 0644)
	if err != nil {
		log.Printf("[InstructionWatcher] FAILED to inject instructions into %s: %v", filePath, err)
		return
	}

	log.Printf("[InstructionWatcher] Successfully injected tool-use instructions into %s", filePath)
}
