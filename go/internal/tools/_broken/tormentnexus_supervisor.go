package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func HandleSupervisorStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cmd := exec.Command("pgrep", "-f", "torment-nexus")
	output, execErr := cmd.Output()
	if execErr != nil {
		return ok("supervisor: stopped")
}

	if len(output) > 0 {
		return ok("supervisor: running")
}

	return ok("supervisor: stopped")
}

func HandleSupervisorRestart(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cmd := exec.Command("systemctl", "restart", "torment-nexus")
	restartErr := cmd.Run()
	if restartErr == nil {
		return ok("supervisor restart initiated via systemctl")
}

	killCmd := exec.Command("pkill", "-f", "torment-nexus")
	killErr := killCmd.Run()
	if killErr != nil && !strings.Contains(killErr.Error(), "exit status 1") {
		return err(fmt.Sprintf("failed to stop supervisor: %v", killErr))
}

	startCmd := exec.Command("torment-nexus", "--supervisor")
	startErr := startCmd.Start()
	if startErr != nil {
		return err(fmt.Sprintf("failed to start supervisor: %v", startErr))
}

	return ok("supervisor restart initiated")
}