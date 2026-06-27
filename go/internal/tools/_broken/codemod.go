package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

type CodemodPattern struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Find        string   `json:"find"`
    Replace     string   `json:"replace"`
    Extensions  []string `json:"extensions"`
    IsRegex     bool     `json:"is_regex"`
}

type CodemodResult struct {
    File    string `json:"file"`
    Applied bool   `json:"applied"`
    Changes int    `json:"changes"`
    Message string `json:"message"`
}

type PreviewResult struct {
    File         string `json:"file"`
    Original     string `json:"original"`
    Modified     string `json:"modified"`
    ChangesCount int    `json:"changes_count"`
}

func getCodemodDir() string {
    home, homeErr := os.UserHomeDir()
    if homeErr != nil {
        home = "/tmp"
    }
    return filepath.Join(home, ".codemod")
}

func getPatternsPath() string {
    return filepath.Join(getCodemodDir(), "patterns.json")
}

func loadPatterns() ([]CodemodPattern, error) {
    patternsPath := getPatternsPath()
    data, readErr := os.ReadFile(patternsPath)
    if readErr != nil {
        if os.IsNotExist(readErr) {
            return []CodemodPattern{}, nil
        }
        return nil, readErr
    }
    var patterns []CodemodPattern
    parseErr := json.Unmarshal(data, &patterns)
    if parseErr != nil {
        return nil, parseErr
    }
    return patterns, nil
}

func savePatterns(patterns []CodemodPattern) error {
    patternsPath := getPatternsPath()
    dirErr := os.MkdirAll(getCodemodDir(), 0755)
    if dirErr != nil {
        return dirErr
    }
    data, jsonErr := json.MarshalIndent(patterns, "", "  ")
    if jsonErr != nil {
        return jsonErr
    }
    return os.WriteFile(patternsPath, data, 0644)
}

// HandleListCodemods lists all available codemod patterns
func HandleListCodemods(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    patterns, loadErr := loadPatterns()
    if loadErr != nil {
        return err(loadErr.Error())
}

    defaultPatterns := []CodemodPattern{
        {
            Name:        "remove-console-log",
            Description: "Remove all console.log statements",
            Find:        `console\.log\([^)]*\)`,
            Replace:     "",
            Extensions:  []string{".js", ".ts", ".jsx", ".tsx", ".go"},
            IsRegex:     true,
        },
        {
            Name:        "remove-debugger",
            Description: "Remove debugger statements",
            Find:        `debugger\s*;?`,
            Replace:     "",
            Extensions:  []string{".js", ".ts", ".jsx", ".tsx"},
            IsRegex:     true,
        },
        {
            Name:        "remove-todo-comments",
            Description: "Remove TODO and FIXME comments",
            Find:        `(?://|#|/\*)\s*(TODO|FIXME):?\s*.*`,
            Replace:     "",
            Extensions:  []string{".js", ".ts", ".go", ".py", ".rb"},
            IsRegex:     true,
        },
        {
            Name:        "const-to-let",
            Description: "Convert const to let for mutable variables",
            Find:        `\bconst\b`,
            Replace:     "let",
            Extensions:  []string{".js", ".ts", ".jsx", ".tsx"},
            IsRegex:     false,
        },
        {
            Name:        "arrow-function",
            Description: "Convert function expressions to arrow functions",
            Find:        `function\s+\w+\s*\([^)]*\)`,
            Replace:     "",
            Extensions:  []string{".js", ".ts"},
            IsRegex:     true,
        },
    }
    allPatterns := append(defaultPatterns, patterns...)
    data, jsonErr := json.Marshal(allPatterns)
    if jsonErr != nil {
        return err(jsonErr.Error())
}

    return ok(string(data))
}

// HandleApplyCodemod applies a codemod pattern to a file
func HandleApplyCodemod(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    patternName, _ :=getString(args, "pattern_name")
    filePath, _ :=getString(args, "file_path")
    customFind, _ :=getString(args, "find")
    customReplace, _ :=getString(args, "replace")
    isRegex, _ :=getBool(args, "is_regex")
    if isRegex && customFind == "" {
        isRegex = false
    }
    if patternName == "" && customFind == "" {
        return err("either pattern_name or find must be provided")
}

    if filePath == "" {
        return err("file_path is required")
}

    fileInfo, statErr := os.Stat(filePath)
    if statErr != nil {
        return err("file not found: " + filePath)
}

    if fileInfo.IsDir() {
        return err("file_path must be a file, not a directory")
}

    content, readErr := os.ReadFile(filePath)
    if readErr != nil {
        return err("failed to read file: " + readErr.Error())
}

    original := string(content)
    modified := original
    var pattern CodemodPattern
    if customFind != "" {
        pattern = CodemodPattern{
            Name:    "custom",
            Find:    customFind,
            Replace: customReplace,
            IsRegex: isRegex,
        }
    } else {
        patterns, loadErr := loadPatterns()
        if loadErr != nil {
            return err("failed to load patterns: " + loadErr.Error())
}

        found := false
        for _, p := range patterns {
            if p.Name == patternName {
                pattern = p
                found = true
                break
            }
        }
        if !found {
            defaultPatterns := map[string]CodemodPattern{
                "remove-console-log": {Find: `console\.log\([^)]*\)`, Replace: "", IsRegex: true},
                "remove-debugger":    {Find: `debugger\s*;?`, Replace: "", IsRegex: true},
                "remove-todo-comments": {Find: `(?://|#|/\*)\s*(TODO|FIXME):?\s*.*`, Replace: "", IsRegex: true},
                "const-to-let":       {Find: "const", Replace: "let", IsRegex: false},
            }
            if dp, found := defaultPatterns[patternName]; found {
                pattern = dp
                found = true
            }
        }
        if !found {
            return err("pattern not found: " + patternName)

    }
    changes := 0
    if pattern.IsRegex {
        re, compileErr := regexp.Compile(pattern.Find)
        if compileErr != nil {
            return err("invalid regex pattern: " + compileErr.Error())
}

        modified = re.ReplaceAllString(modified, pattern.Replace)
        matches := re.FindAllStringIndex(original, -1)
        changes = len(matches)
    } else {
        modified = strings.ReplaceAll(modified, pattern.Find, pattern.Replace)
        occurrences := strings.Count(original, pattern.Find)
        changes = occurrences
    }
    if changes == 0 {
        result := CodemodResult{
            File:    filePath,
            Applied: false,
            Changes: 0,
            Message: "no changes made - pattern not found in file",
        }
        data, _ := json.Marshal(result)
        return ok(string(data))
}

    writeErr := os.WriteFile(filePath, []byte(modified), 0644)
    if writeErr != nil {
        return err("failed to write file: " + writeErr.Error())
}

    result := CodemodResult{
        File:    filePath,
        Applied: true,
        Changes: changes,
        Message: fmt.Sprintf("applied %d changes to %s", changes, filePath),
    }
    data, jsonErr := json.Marshal(result)
    if jsonErr != nil {
        return err(jsonErr.Error())
}

    return ok(string(data))
}

}

// HandlePreviewCodemod shows what changes would be made without applying them
func HandlePreviewCodemod(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    patternName, _ :=getString(args, "pattern_name")
    filePath, _ :=getString(args, "file_path")
    customFind, _ :=getString(args, "find")
    customReplace, _ :=getString(args, "replace")
    isRegex, _ :=getBool(args, "is_regex")
    if isRegex && customFind == "" {
        isRegex = false
    }
    if patternName == "" && customFind == "" {
        return err("either pattern_name or find must be provided")
}

    if filePath == "" {
        return err("file_path is required")
}

    content, readErr := os.ReadFile(filePath)
    if readErr != nil {
        return err("failed to read file: " + readErr.Error())
}

    original := string(content)
    modified := original
    var pattern CodemodPattern
    if customFind != "" {
        pattern = CodemodPattern{
            Name:    "custom",
            Find:    customFind,
            Replace: customReplace,
            IsRegex: isRegex,
        }
    } else {
        patterns, loadErr := loadPatterns()
        if loadErr != nil {
            return err("failed to load patterns: " + loadErr.Error())
}

        found := false
        for _, p := range patterns {
            if p.Name == patternName {
                pattern = p
                found = true
                break
            }
        }
        if !found {
            defaultPatterns := map[string]CodemodPattern{
                "remove-console-log": {Find: `console\.log\([^)]*\)`, Replace: "", IsRegex: true},
                "remove-debugger":    {Find: `debugger\s*;?`, Replace: "", IsRegex: true},
                "remove-todo-comments": {Find: `(?://|#|/\*)\s*(TODO|FIXME):?\s*.*`, Replace: "", IsRegex: true},
                "const-to-let":       {Find: "const", Replace: "let", IsRegex: false},
            }
            if dp, found := defaultPatterns[patternName]; found {
                pattern = dp
                found = true
            }
        }
        if !found {
            return err("pattern not found: " + patternName)

    }
    changes := 0
    if pattern.IsRegex {
        re, compileErr := regexp.Compile(pattern.Find)
        if compileErr != nil {
            return err("invalid regex pattern: " + compileErr.Error())
}

        modified = re.ReplaceAllString(modified, pattern.Replace)
        matches := re.FindAllStringIndex(original, -1)
        changes = len(matches)
    } else {
        modified = strings.ReplaceAll(modified, pattern.Find, pattern.Replace)
        occurrences := strings.Count(original, pattern.Find)
        changes = occurrences
    }
    result := PreviewResult{
        File:         filePath,
        Original:     original,
        Modified:     modified,
        ChangesCount: changes,
    }
    data, jsonErr := json.Marshal(result)
    if jsonErr != nil {
        return err(jsonErr.Error())
}

    return ok(string(data))
}

}

// HandleCreateCodemod creates a new codemod pattern
func HandleCreateCodemod(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    name, _ :=getString(args, "name")
    description, _ :=getString(args, "description")
    find, _ :=getString(args, "find")
    replace, _ :=getString(args, "replace")
    extensionsStr, _ :=getString(args, "extensions")
    isRegex, _ :=getBool(args, "is_regex")
    if name == "" {
        return err("name is required")
}

    if find == "" {
        return err("find pattern is required")
}

    if isRegex {
        _, compileErr := regexp.Compile(find)
        if compileErr != nil {
            return err("invalid regex pattern: " + compileErr.Error())

    }
    var extensions []string
    if extensionsStr != "" {
        extensions = strings.Split(extensionsStr, ",")
        for i := range extensions {
            extensions[i] = strings.TrimSpace(extensions[i])

    } else {
	}
}
}
}