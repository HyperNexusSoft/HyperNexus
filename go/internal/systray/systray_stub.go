//go:build !windows

package systray

import (
	"github.com/MDMAtk/TormentNexus/internal/eventbus"
)

// Start is a no-op on non-Windows platforms
func Start(eb *eventbus.EventBus) {
	// Headless mode: no system tray UI
}

// NotifyActivity is a no-op on non-Windows platforms
func NotifyActivity(dir string) {
	// Headless mode: do nothing
}
