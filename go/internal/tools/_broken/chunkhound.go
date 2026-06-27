package tools

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

// HandleSearchChunks searches a file for lines matching a regex pattern
func HandleSearchChunks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, _ :=getString(args, "file")
	pattern, _ :=getString(args, "pattern")

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return err(readErr.Error())
}

	re, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		return err(regexErr.Error())
}

	lines := strings.Split(string(content), "\n")
	var matches []string
	for _, line := range lines {
		if re.MatchString(line) {
			matches = append(matches, line)

	}

	result := strings.Join(matches, "\n")
	return ok(result)
}

}

// HandleCountChunks counts the number of lines (chunks) in a file
func HandleCountChunks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, _ :=getString(args, "file")

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return err(readErr.Error())
}

	lines := strings.Split(string(content), "\n")
	count := len(lines)
	return ok(strconv.Itoa(count))
}

// HandleExtractChunks extracts specific line ranges from a file
func HandleExtractChunks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, _ :=getString(args, "file")
	start, _ :=getInt(args, "start")
	end, _ :=getInt(args, "end")

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return err(readErr.Error())
}

	lines := strings.Split(string(content), "\n")
	if start < 0 {
		start = 0
	}
	if end >= len(lines) {
		end = len(lines) - 1
	}
	if start > end {
		return err("start cannot be greater than end")
}

	extracted := lines[start : end+1]
	result := strings.Join(extracted, "\n")
	return ok(result)
}

// HandleFilterChunks filters file lines by inclusion/exclusion patterns
func HandleFilterChunks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, _ :=getString(args, "file")
	include, _ :=getString(args, "include")
	exclude, _ :=getString(args, "exclude")

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return err(readErr.Error())
}

	includeRe, includeErr := regexp.Compile(include)
	if includeErr != nil {
		return err(includeErr.Error())
}

	excludeRe, excludeErr := regexp.Compile(exclude)
	if excludeErr != nil {
		return err(excludeErr.Error())
}

	lines := strings.Split(string(content), "\n")
	var filtered []string
	for _, line := range lines {
		if includeRe.MatchString(line) && !excludeRe.MatchString(line) {
			filtered = append(filtered, line)

	}

	result := strings.Join(filtered, "\n")
	return ok(result)
}
}