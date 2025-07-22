package services

import (
	"github.com/denkhaus/agentforge/internal/logger"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log = logger.WithPackage("services")
}