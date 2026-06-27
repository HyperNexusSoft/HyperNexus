package tools

import (
    "context"
    "fmt"
    "sort"
    "strings"
)

// Global in-memory store for notes (not thread-safe)
var notes = make(map[string]string)

func HandleCreateNote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    title, _ :=getString(args, "title")
    content, _ :=getString(args, "content")
    if title == "" || content == "" {
        return err("missing title or content")
}

    notes[title] = content
    return ok(fmt.Sprintf("Note '%s' created", title))
}

func HandleListNotes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    var list []string
    for title := range notes {
        list = append(list, title)

    sort.Strings(list)
    return ok(strings.Join(list, ", "))
}
}