package handlers

import "github.com/Schalure/urlalias/internal/app/aliasmaker"

// Middleware type
type Middleware struct {
	service *aliasmaker.AliasMakerServise
}

// ------------------------------------------------------------
//
//	Constructor of middleware
//	Input:
//		logger *slog.Logger
//	Output:
//		*Middleware
func NewMiddleware(service *aliasmaker.AliasMakerServise) *Middleware {

	return &Middleware{
		service: service,
	}
}
