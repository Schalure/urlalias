package handlers

import "github.com/Schalure/urlalias/internal/app/aliasmaker"

// Middleware type
type Middleware struct {
	logger aliasmaker.Loggerer
}

// ------------------------------------------------------------
//
//	Constructor of middleware
//	Input:
//		logger *slog.Logger
//	Output:
//		*Middleware
func NewMiddleware(logger aliasmaker.Loggerer) *Middleware {

	return &Middleware{
		logger: logger,
	}
}
