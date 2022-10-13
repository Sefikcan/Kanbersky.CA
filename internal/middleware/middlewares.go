package middleware

import (
	"github.com/sefikcan/kanbersky.ca/pkg/config"
	"github.com/sefikcan/kanbersky.ca/pkg/logger"
)

type MiddlewareManager struct {
	cfg *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg: cfg,
		logger: logger,
	}
}