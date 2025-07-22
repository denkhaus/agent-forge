package commands

import (
	"github.com/denkhaus/agentforge/internal/tui"
	"github.com/denkhaus/agentforge/internal/types"
)

// getTUIModule returns a TUI module instance
func getTUIModule() types.TUIManager {
	return tui.NewTUIManager()
}