package logger

import (
	"context"
	"log"
)

// StdErrLogger logs all output to stdout
var StdErrLogger Logger = &stdErrLogger{}

type stdErrLogger struct{}

// Debug logs a debug message.
func (l *stdErrLogger) Debug(ctx context.Context, v interface{}) {
	std("DEBUG", v)
}

// Debugf logs a debug message with format.
func (l *stdErrLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	stdf("DEBUG", format, v...)
}

// Info logs a info message.
func (l *stdErrLogger) Info(ctx context.Context, v interface{}) {
	std("INFO ", v)
}

// Infof logs a info message with format.
func (l *stdErrLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	stdf("INFO ", format, v...)
}

// Warn logs a warning message.
func (l *stdErrLogger) Warn(ctx context.Context, v interface{}) {
	std("WARN ", v)
}

// Warnf logs a warning message with format.
func (l *stdErrLogger) Warnf(ctx context.Context, format string, v ...interface{}) {
	stdf("WARN ", format, v...)
}

// Error logs an error message.
func (l *stdErrLogger) Error(ctx context.Context, v interface{}) {
	std("ERROR", v)
}

// Errorf logs an error message with format.
func (l *stdErrLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	stdf("ERROR", format, v...)
}

func std(level string, v ...interface{}) {
	log.Printf(level+": %s", v...)
}

func stdf(level, format string, v ...interface{}) {
	log.Printf(level+": "+format, v...)
}
