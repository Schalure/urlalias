package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	//	Middleware type
	Middleware struct {
		logger *zap.SugaredLogger
	}

	//	Date from response
	responseData struct {
		status int
		size   int
		data   string
	}

	//	Response writer with login
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// ------------------------------------------------------------
//
//	Override Write() method by http.ResponseWriter
//	Receiver:
//		r *loggingResponseWriter
//	Input:
//		b []byte
//	Output:
//		int - count of write bytes
//		err
func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	size, err := r.ResponseWriter.Write(b)

	r.responseData.data = string(b)
	r.responseData.size += size
	return size, err
}

// ------------------------------------------------------------
//
//	Override WriteHeader() method by http.ResponseWriter
//	Receiver:
//		r *loggingResponseWriter
//	Input:
//		statusCode int
func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
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

// ------------------------------------------------------------
//
//	WhithLogging middleware - method of Middleware type
//	Logs request and response data
//	Receiver:
//		m* Middleware
//	Input:
//		h http.Handler
//	Output:
//		http.Handler
func (m *Middleware) WhithLogging(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := new(responseData)
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)

		m.logger.Infow(
			"Information about request and response",
			"Request URI", r.RequestURI,
			"Request method", r.Method,

			"Response status", responseData.status,
			"Response headers", lw.ResponseWriter.Header(),
			"Response data", responseData.data,

			"duration", duration,
		)

	})
}
