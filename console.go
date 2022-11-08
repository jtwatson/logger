package logger

import (
	"context"
	"log"
	"net/http"
)

type color int

const (
	red    color = 31
	yellow color = 33
	blue   color = 34
	gray   color = 37
)

// NewConsoleLogger logs all output to console
func NewConsoleLogger(r *http.Request) *ConsoleLogger {
	return &ConsoleLogger{r}
}

type ConsoleLogger struct {
	r *http.Request
}

// Debug logs a debug message.
func (l *ConsoleLogger) Debug(ctx context.Context, v interface{}) {
	l.console("DEBUG", gray, v)
}

// Debugf logs a debug message with format.
func (l *ConsoleLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.consolef("DEBUG", gray, format, v...)
}

// Info logs a info message.
func (l *ConsoleLogger) Info(ctx context.Context, v interface{}) {
	l.console("INFO ", blue, v)
}

// Infof logs a info message with format.
func (l *ConsoleLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.consolef("INFO ", blue, format, v...)
}

// Warn logs a warning message.
func (l *ConsoleLogger) Warn(ctx context.Context, v interface{}) {
	l.console("WARN ", yellow, v)
}

// Warnf logs a warning message with format.
func (l *ConsoleLogger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.consolef("WARN ", yellow, format, v...)
}

// Error logs an error message.
func (l *ConsoleLogger) Error(ctx context.Context, v interface{}) {
	l.console("ERROR", red, v)
}

// Errorf logs an error message with format.
func (l *ConsoleLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.consolef("ERROR", red, format, v...)
}

func (l *ConsoleLogger) console(level string, c color, v interface{}) {
	log.Printf(colorPrint(level, c)+": %s %s", l.r.URL.Path, v)
}

func (l *ConsoleLogger) consolef(level string, c color, format string, v ...interface{}) {
	log.Printf(colorPrint(level, c)+": "+l.r.URL.Path+" "+format, v...)
}

func colorPrint(s string, c color) string {
	return string([]byte{0x1b, '[', byte('0' + c/10), byte('0' + c%10), 'm'}) + s + "\x1b[0m"
}
