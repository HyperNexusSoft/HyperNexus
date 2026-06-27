package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HandleReadFile reads and returns the contents of a file
func HandleReadFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("path is required")
}

	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(string(content))
}

// HandleWriteFile writes content to a file at the specified path
func HandleWriteFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	content, _ :=getString(args, "content")

	if path == "" {
		return err("path is required")
}

	dir := filepath.Dir(path)
	if dirErr := os.MkdirAll(dir, 0755); dirErr != nil {
		return err(dirErr.Error())
}

	writeErr := os.WriteFile(path, []byte(content), 0644)
	if writeErr != nil {
		return err(writeErr.Error())
}

	return ok(fmt.Sprintf("Successfully wrote to %s", path))
}

// HandleListDirectory lists contents of a directory
func HandleListDirectory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		path = "."
	}

	entries, readErr := os.ReadDir(path)
	if readErr != nil {
		return err(readErr.Error())
}

	var result []map[string]interface{}
	for _, entry := range entries {
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		entryType := "file"
		if entry.IsDir() {
			entryType = "directory"
		} else if info.Mode()&os.ModeSymlink != 0 {
			entryType = "symlink"
		}
		result = append(result, map[string]interface{}{
			"name":    entry.Name(),
			"type":    entryType,
			"size":    info.Size(),
			"modTime": info.ModTime().Format(time.RFC3339),
		})

	jsonBytes, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return err(jsonErr.Error())
}

	return ok(string(jsonBytes))
}

}

// HandleCreateDirectory creates a new directory
func HandleCreateDirectory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	recursive, _ :=getBool(args, "recursive")

	if path == "" {
		return err("path is required")
}

	var mkdirErr error
	if recursive {
		mkdirErr = os.MkdirAll(path, 0755)
	} else {
		mkdirErr = os.Mkdir(path, 0755)

	if mkdirErr != nil {
		return err(mkdirErr.Error())
}

	return ok(fmt.Sprintf("Directory created: %s", path))
}

}

// HandleDeleteFile deletes a file or directory
func HandleDeleteFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	recursive, _ :=getBool(args, "recursive")

	if path == "" {
		return err("path is required")
}

	info, stat

[Switched to Model: gpt-oss-120b | Provider: cerebras via Global Random Retry]

We need to implement a Go-native MCP tool module for "filesystem-with-morph". We need to study original MCP server source code above, but we don't have it in the prompt. However, we can infer typical filesystem tool functions: maybe list directory, read file, write file, delete file, move/rename, etc. "with-morph" maybe includes morphing files? Not sure.

We need to implement 2-6 handlers. Provide simple implementations using stdlib. Must use ok() and err("error") functions defined elsewhere (parity.go). Also getString, getInt, getBool exist.

We need to avoid redeclaring ToolResponse, ok, e, etc.

We need to output GO_FILE with package tools and code, and MANIFEST JSON with filename, server_name, handlers list.

We need to implement handlers with signature: func HandleXxx(ctx context.Context, args map[string]interface{}) (ToolResponse, error)

We need to use getString(args, "key") etc.

We need to import needed packages.

We need to ensure no variable named e is used (shadowing). Use apiErr etc.

Implement maybe:

- HandleReadFile: reads a file path given, returns content as text.

- HandleWriteFile: writes content to file path.

- HandleListDir: lists entries in a directory.

- HandleDeleteFile: deletes a file.

- HandleMoveFile: moves/renames a file.

That's 5 handlers.

Implement each with proper error handling.

Return ok(string) where string is maybe JSON or plain text.

We need to use TextContent? Not sure. The ok function likely expects a TextContent struct? But we can just pass string.

We'll assume ok takes a string.

Implementation details:

- getString(args, "path") for file path.

- For write, also getString(args, "content").

- For list, maybe getString(args, "path") and return JSON array of names.

- For delete, getString(args, "path").

- For move, getString(args, "src") and getString(args, "dst").

Make sure to use context maybe for cancellation? Not needed.

Use os.ReadFile, os.WriteFile, os.Remove, os.Rename, os.ReadDir.

Return errors via err(e.Error()).

Make sure to import context, os, fmt, encoding/json, strings maybe.

Ok.

Now produce manifest: list handlers with tool===GO_FILE===
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

}

// HandleReadFile reads the content of a file and returns it as text.
// Expected args:
//   "path": string – absolute or relative file path.
func HandleListDir(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	path, _ :=getString(args, "path")
	if path == "" {
		return err("missing 'path' argument")
}

	absPath, apiErr := filepath.Abs(path)
	if apiErr != nil {
		return err(fmt.Sprintf("failed to resolve absolute path: %v", apiErr))
}

	entries, readErr := os.ReadDir(absPath)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read directory: %v", readErr))
}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())

	jsonBytes, marshalErr := json.Marshal(names)
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to marshal result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}

}

// HandleDeleteFile removes a file from the filesystem.
// Expected args:
//   "path": string – file path to delete.
func HandleMoveFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	src, _ :=getString(args, "src")
	if src == "" {
		return err("missing 'src' argument")
}

	dst, _ :=getString(args, "dst")
	if dst == "" {
		return err("missing 'dst' argument")
}

	absSrc, srcErr := filepath.Abs(src)
	if srcErr != nil {
		return err(fmt.Sprintf("failed to resolve source path: %v", srcErr))
}

	absDst, dstErr := filepath.Abs(dst)
	if dstErr != nil {
		return err(fmt.Sprintf("failed to resolve destination path: %v", dstErr))
}

	// Ensure destination directory exists
	destDir := filepath.Dir(absDst)
	if mkErr := os.MkdirAll(destDir, 0755); mkErr != nil {
		return err(fmt.Sprintf("failed to create destination directories: %v", mkErr))
}

	renameErr := os.Rename(absSrc, absDst)
	if renameErr != nil {
		return err(fmt.Sprintf("failed to move file: %v", renameErr))
}

	return ok(fmt.Sprintf("moved %s to %s", absSrc, absDst))

}