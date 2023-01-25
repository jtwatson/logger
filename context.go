package logger

import (
	"context"
	"net/http"
)

type key int

const (
	logKey key = iota
)

// fromContext gets the logger out of the context.
// If not logger is stored in the context, a NopLogger is returned.
func fromContext(ctx context.Context) ctxLogger {
	if ctx == nil {
		return &stdErrLogger{}
	}
	l, ok := ctx.Value(logKey).(ctxLogger)
	if !ok {
		return &stdErrLogger{}
	}

	return l
}

// fromRequest gets the logger in the request's context.
// This is a shortcut for xlog.FromContext(r.Context())
func fromRequest(r *http.Request) ctxLogger {
	if r == nil {
		return &stdErrLogger{}
	}

	return fromContext(r.Context())
}

// newContext returns a copy of the parent context and associates it with the provided logger.
func newContext(ctx context.Context, l ctxLogger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

// ctxLogger defines the logging interface with context
type ctxLogger interface {
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
