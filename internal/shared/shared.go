package shared

import (
	"fmt"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/container"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
)

type StartupContext struct {
	Injector *do.Injector
	Config   *config.Config
}

func (p *StartupContext) Cleanup() {
	container.Cleanup(p.Injector)
}

func StartupCommon(c *cli.Context) (*StartupContext, error) {

	cfg, err := config.Load(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	injector := container.Setup(cfg)

	return &StartupContext{
		Config:   cfg,
		Injector: injector,
	}, nil
}
