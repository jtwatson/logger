package logger

import (
	"context"
	"log"
)

type stdErrLogger struct{}

// Debug logs a debug message.
func (l *stdErrLogger) Debug(_ context.Context, v interface{}) {
	std("DEBUG", v)
}

// Debugf logs a debug message with format.
func (l *stdErrLogger) Debugf(_ context.Context, format string, v ...interface{}) {
	stdf("DEBUG", format, v...)
}

// Info logs a info message.
func (l *stdErrLogger) Info(_ context.Context, v interface{}) {
	std("INFO ", v)
}

// Infof logs a info message with format.
func (l *stdErrLogger) Infof(_ context.Context, format string, v ...interface{}) {
	stdf("INFO ", format, v...)
}

// Warn logs a warning message.
func (l *stdErrLogger) Warn(_ context.Context, v interface{}) {
	std("WARN ", v)
}

// Warnf logs a warning message with format.
func (l *stdErrLogger) Warnf(_ context.Context, format string, v ...interface{}) {
	stdf("WARN ", format, v...)
}

// Error logs an error message.
func (l *stdErrLogger) Error(_ context.Context, v interface{}) {
	std("ERROR", v)
}

// Errorf logs an error message with format.
func (l *stdErrLogger) Errorf(_ context.Context, format string, v ...interface{}) {
	stdf("ERROR", format, v...)
}

func std(level string, v ...interface{}) {
	log.Printf(level+": %s", v...)
}

func stdf(level, format string, v ...interface{}) {
	log.Printf(level+": "+format, v...)
}
