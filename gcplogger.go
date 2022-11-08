package logger

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/trace"
)

// NewGCPLogger manually creates a logger.
//
// This function should only be used outside of a request. Use FromContext in requests.
func NewGCPLogger(lg CloudLogger, traceID string) *GCPLogger {
	return &GCPLogger{
		lg:      lg,
		traceID: traceID,
	}
}

type CloudLogger interface {
	Log(e logging.Entry)
}

type GCPLogger struct {
	lg          CloudLogger
	traceID     string
	mu          sync.Mutex
	maxSeverity logging.Severity
	logCount    int
}

func (l *GCPLogger) log(ctx context.Context, severity logging.Severity, p interface{}) {
	l.mu.Lock()
	if l.maxSeverity < severity {
		l.maxSeverity = severity
	}
	l.logCount++
	l.mu.Unlock()

	if err, ok := p.(error); ok {
		p = err.Error()
	}

	span := trace.SpanFromContext(ctx)

	l.lg.Log(
		logging.Entry{
			Payload: payload{
				Message: p,
			},
			Severity:     severity,
			Trace:        l.traceID,
			SpanID:       span.SpanContext().SpanID().String(),
			TraceSampled: span.SpanContext().IsSampled(),
		},
	)
}

// Debug logs a debug message.
func (l *GCPLogger) Debug(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Debug, v)
}

// Debugf logs a debug message with format.
func (l *GCPLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Debug, fmt.Sprintf(format, v...))
}

// Info logs a info message.
func (l *GCPLogger) Info(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Info, v)
}

// Infof logs a info message with format.
func (l *GCPLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Info, fmt.Sprintf(format, v...))
}

// Warn logs a warning message.
func (l *GCPLogger) Warn(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Warning, v)
}

// Warnf logs a warning message with format.
func (l *GCPLogger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Warning, fmt.Sprintf(format, v...))
}

// Error logs an error message.
func (l *GCPLogger) Error(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Error, v)
}

// Errorf logs an error message with format.
func (l *GCPLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Error, fmt.Sprintf(format, v...))
}

type payload struct {
	Message interface{} `json:"message"`
}
