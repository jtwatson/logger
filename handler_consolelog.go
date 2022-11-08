//go:build consolelog

package logger

import (
	"net/http"

	"cloud.google.com/go/logging"
)

// NewRequestLogger is a HTTP handler that logs all requests using a ConsoleLogger. This is intended
// for local development for quicker iterations.
func NewRequestLogger(client *logging.Client, projectID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(NewContext(r.Context(), NewConsoleLogger(r)))
			next.ServeHTTP(w, r)
		})
	}
}

// NewLogger is identical to NewHandler when logging to the console
func NewLogger(client *logging.Client, projectID string) func(http.Handler) http.Handler {
	return NewRequestLogger(client, projectID)
}
