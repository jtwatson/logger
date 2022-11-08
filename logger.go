// logger is an HTTP request logger that implements correlated logging to GCP via Logging REST API. Each HTTP request is l
// ogged as the parent with all event logs occurring during the request as child logs. This allows for easy viewing in
// GCP Log Explorer. The logs will also be correlated to Cloud Trace if you instrument your code with tracing.
package logger

import (
	"context"
	"net/http"
)

type key int

const (
	logKey key = iota
)

// FromContext gets the logger out of the context.
// If not logger is stored in the context, a NopLogger is returned.
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return StdErrLogger
	}
	l, ok := ctx.Value(logKey).(Logger)
	if !ok {
		return StdErrLogger
	}

	return l
}

// FromRequest gets the logger in the request's context.
// This is a shortcut for xlog.FromContext(r.Context())
func FromRequest(r *http.Request) Logger {
	if r == nil {
		return StdErrLogger
	}

	return FromContext(r.Context())
}

// NewContext returns a copy of the parent context and associates it with the provided logger.
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

// Logger defines the interface for a xlog compatible logger
type Logger interface {
	// Debug logs a debug message.
	Debug(ctx context.Context, v interface{})
	// Debugf logs a debug message with format.
	Debugf(ctx context.Context, format string, v ...interface{})
	// Info logs a info message.
	Info(ctx context.Context, v interface{})
	// Infof logs a info message with format.
	Infof(ctx context.Context, format string, v ...interface{})
	// Warn logs a warning message.
	Warn(ctx context.Context, v interface{})
	// Warnf logs a warning message with format.
	Warnf(ctx context.Context, format string, v ...interface{})
	// Error logs an error message.
	Error(ctx context.Context, v interface{})
	// Errorf logs an error message with format.
	Errorf(ctx context.Context, format string, v ...interface{})
}
