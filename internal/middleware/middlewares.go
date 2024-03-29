package middleware

import (
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/auth"
	"github.com/engineerXIII/maiSystemBackend/internal/session"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	sessUC  session.UCSession
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
