package handlers

// Middleware type
type Middleware struct {
	logger Loggerer
}

// ------------------------------------------------------------
//
//	Constructor of middleware
//	Input:
//		logger *slog.Logger
//	Output:
//		*Middleware
func NewMiddleware(logger Loggerer) *Middleware {

	return &Middleware{
		logger: logger,
	}
}
