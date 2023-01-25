package logger

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// GoogleCloudExporter implements exporting to Google Cloud Logging
type GoogleCloudExporter struct {
	projectID string
	client    *logging.Client
	opts      []logging.LoggerOption
	logAll    bool
}

// NewGoogleCloudExporter returns a configured GoogleCloudExporter
func NewGoogleCloudExporter(client *logging.Client, projectID string, opts ...logging.LoggerOption) *GoogleCloudExporter {
	return &GoogleCloudExporter{
		projectID: projectID,
		client:    client,
		opts:      opts,
	}
}

// LogAll controls if this logger will log all requests, or only requests that contain
// logs written to the request Logger
func (e *GoogleCloudExporter) LogAll(v bool) *GoogleCloudExporter {
	e.logAll = v

	return e
}

// Middleware returns a middleware that exports logs to Google Cloud Logging
func (e *GoogleCloudExporter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		parentLogger := e.client.Logger(
			"request_parent_log",
			e.opts...,
		)

		childLogger := e.client.Logger(
			"request_child_log",
			e.opts...,
		)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()
			traceID := gcpTraceIDFromRequest(r, e.projectID)
			l := newGCPLogger(childLogger, traceID)
			r = r.WithContext(newContext(r.Context(), l))
			sw := &statusWriter{ResponseWriter: w}

			next.ServeHTTP(sw, r)

			l.mu.Lock()
			logCount := l.logCount
			maxSeverity := l.maxSeverity
			l.mu.Unlock()

			if !e.logAll && logCount == 0 {
				return
			}

			// status code should also set the minimum maxSeverity to Error
			if sw.Status() > 399 && maxSeverity < logging.Error {
				maxSeverity = logging.Error
			}

			sc := trace.SpanFromContext(r.Context()).SpanContext()

			parentLogger.Log(logging.Entry{
				Timestamp:    begin,
				Severity:     maxSeverity,
				Trace:        traceID,
				SpanID:       sc.SpanID().String(),
				TraceSampled: sc.IsSampled(),
				Payload: gcpPayload{
					Message: "Parent Log Entry",
				},
				HTTPRequest: &logging.HTTPRequest{
					Request:      r,
					RequestSize:  requestSize(r.Header.Get("Content-Length")),
					Latency:      time.Since(begin),
					Status:       sw.Status(),
					ResponseSize: sw.length,
					RemoteIP:     r.Header.Get("X-Forwarded-For"),
				},
			})
		})
	}
}

// gcpTraceIDFromRequest formats a trace_id value for GCP Stackdriver
func gcpTraceIDFromRequest(r *http.Request, projectID string) string {
	var traceID string
	if sc, ok := new(propagation.HTTPFormat).SpanContextFromRequest(r); ok {
		traceID = sc.TraceID.String()
	} else {
		sc := trace.SpanFromContext(r.Context()).SpanContext()
		if sc.IsValid() {
			traceID = sc.TraceID().String()
		} else {
			_, span := otel.Tracer("").Start(r.Context(), r.URL.String())
			traceID = span.SpanContext().TraceID().String()
		}
	}

	return fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
}

type gcpPayload struct {
	Message interface{} `json:"message"`
}

// logger interface exists for testability
type logger interface {
	Log(e logging.Entry)
}

type gcpLogger struct {
	lg          logger
	traceID     string
	mu          sync.Mutex
	maxSeverity logging.Severity
	logCount    int
}

func newGCPLogger(lg logger, traceID string) *gcpLogger {
	return &gcpLogger{
		lg:      lg,
		traceID: traceID,
	}
}

// Debug logs a debug message.
func (l *gcpLogger) Debug(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Debug, v)
}

// Debugf logs a debug message with format.
func (l *gcpLogger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Debug, fmt.Sprintf(format, v...))
}

// Info logs a info message.
func (l *gcpLogger) Info(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Info, v)
}

// Infof logs a info message with format.
func (l *gcpLogger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Info, fmt.Sprintf(format, v...))
}

// Warn logs a warning message.
func (l *gcpLogger) Warn(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Warning, v)
}

// Warnf logs a warning message with format.
func (l *gcpLogger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Warning, fmt.Sprintf(format, v...))
}

// Error logs an error message.
func (l *gcpLogger) Error(ctx context.Context, v interface{}) {
	l.log(ctx, logging.Error, v)
}

// Errorf logs an error message with format.
func (l *gcpLogger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.log(ctx, logging.Error, fmt.Sprintf(format, v...))
}

func (l *gcpLogger) log(ctx context.Context, severity logging.Severity, p interface{}) {
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
			Payload: gcpPayload{
				Message: p,
			},
			Severity:     severity,
			Trace:        l.traceID,
			SpanID:       span.SpanContext().SpanID().String(),
			TraceSampled: span.SpanContext().IsSampled(),
		},
	)
}
