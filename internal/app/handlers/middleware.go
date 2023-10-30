package handlers

import "go.uber.org/zap"

// Middleware type
type Middleware struct {
	logger *zap.SugaredLogger
}

// ------------------------------------------------------------
//
//	Constructor of middleware
//	Input:
//		logger *slog.Logger
//	Output:
//		*Middleware
func NewMiddleware(logger *zap.SugaredLogger) *Middleware {

	return &Middleware{
		logger: logger,
	}
}
