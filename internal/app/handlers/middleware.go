package handlers

import "github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"

// Middleware type
type Middleware struct {
	userManager UserManager
	logger *zaplogger.ZapLogger
}

// ------------------------------------------------------------
//
//	Constructor of middleware
//	Input:
//		logger *slog.Logger
//	Output:
//		*Middleware
func NewMiddleware(userManager UserManager, logger *zaplogger.ZapLogger) *Middleware {

	return &Middleware{
		userManager: userManager,
		logger: logger,
	}
}
