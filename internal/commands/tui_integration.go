package commands

import (
	"github.com/denkhaus/agentforge/internal/types"
	"github.com/samber/do"
)

// getTUIModule returns a TUI module instance (placeholder)
func getTUIModule() types.TUIManager {
	// TODO: Implement TUI manager when ready
	return nil
}

// getTUIManagerFromDI returns a TUI manager from the DI container
func getTUIManagerFromDI(injector *do.Injector) (types.TUIManager, error) {
	return do.Invoke[types.TUIManager](injector)
}