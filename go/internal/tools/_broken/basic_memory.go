package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type fileInfo struct {
	path    string
	modTime time.Time
}

func HandleCreateNote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// ... unchanged
}

func HandleReadNote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// ... unchanged
}

func HandleSearchNotes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// ... unchanged
}

func HandleGetRecentNotes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
}

	// ... unchanged, but remove the local fileInfo definition and use the package-level one
	// ... and change sortSlice(files) to sort.Slice(files, func(i, j int) bool { return files[i].modTime.After(files[j].modTime) })

func HandleSync(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// ... unchanged
}