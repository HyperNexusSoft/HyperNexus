package tools

import (
    "context"
    "fmt"
    "strings"
)

func HandleHello(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    // maybe greet
    name, _ :=getString(args, "name")
    if name == "" {
        name = "World"
    }
    return ok(fmt.Sprintf("Hello, %s!", name))
}

func HandleEcho(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    text, _ :=getString(args, "text")
    if text == "" {
        return err("text argument is required")
}

    return ok(text)
}