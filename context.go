package logger

import (
	"context"
	"fmt"
)

// Debug logs a debug message.
func Debug(ctx context.Context, v interface{}) {
	FromContext(ctx).Debug(ctx, v)
}

// Debug logs a debug message with format.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Debugf(ctx, fmt.Sprintf(format, v...))
}

// Info logs a info message.
func Info(ctx context.Context, v interface{}) {
	FromContext(ctx).Info(ctx, v)
}

// Info logs a info message with format.
func Infof(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Infof(ctx, fmt.Sprintf(format, v...))
}

// Warn logs a warning message.
func Warn(ctx context.Context, v interface{}) {
	FromContext(ctx).Warn(ctx, v)
}

// Warn logs a warning message with format.
func Warnf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Warnf(ctx, fmt.Sprintf(format, v...))
}

// Error logs an error message.
func Error(ctx context.Context, v interface{}) {
	FromContext(ctx).Error(ctx, v)
}

// Error logs an error message with format.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Errorf(ctx, fmt.Sprintf(format, v...))
}
