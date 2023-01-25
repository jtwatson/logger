// package logger is an HTTP request logger that implements correlated logging to GCP via Logging REST API. Each HTTP request is l
// ogged as the parent with all event logs occurring during the request as child logs. This allows for easy viewing in
// GCP Log Explorer. The logs will also be correlated to Cloud Trace if you instrument your code with tracing.
package logger

import (
	"context"
	"net/http"
)

// Logger implements logging methods for this package
type Logger struct {
	ctx context.Context
	lg  ctxLogger
}

// FromCtx returns the logger from the context. If
// no logger is found, it will write to stderr
func FromCtx(ctx context.Context) *Logger {
	return &Logger{
		ctx: ctx,
		lg:  fromContext(ctx),
	}
}

// FromReq returns the logger from the http request. If
// no logger is found, it will write to stderr
func FromReq(r *http.Request) *Logger {
	return &Logger{
		ctx: r.Context(),
		lg:  fromRequest(r),
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(v interface{}) {
	l.lg.Debug(l.ctx, v)
}

// Debugf logs a debug message with format.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.lg.Debugf(l.ctx, format, v...)
}

// Info logs a info message.
func (l *Logger) Info(v interface{}) {
	l.lg.Info(l.ctx, v)
}

// Infof logs a info message with format.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.lg.Infof(l.ctx, format, v...)
}

// Warn logs a warning message.
func (l *Logger) Warn(v interface{}) {
	l.lg.Warn(l.ctx, v)
}

// Warnf logs a warning message with format.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.lg.Warnf(l.ctx, format, v...)
}

// Error logs an error message.
func (l *Logger) Error(v interface{}) {
	l.lg.Error(l.ctx, v)
}

// Errorf logs an error message with format.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.lg.Errorf(l.ctx, format, v...)
}
