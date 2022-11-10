//go:build !consolelog

package logger

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/go-playground/errors/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// NewRequestLogger is a HTTP handler that logs all requests as a parent
// log and links it to all children event logs occurring further
// down the handler call chain
//
// If not configured, the output is set to stderr by default.
func NewRequestLogger(client *logging.Client, projectID string, opts ...logging.LoggerOption) func(http.Handler) http.Handler {
	return handler(client, projectID, true, opts...)
}

// NewLogger is the same as NewHandler, but only logs the parent request
// if a child log is generated during the request. This should be used
// for high frequency request paths to prevent logging every request
func NewLogger(client *logging.Client, projectID string, opts ...logging.LoggerOption) func(http.Handler) http.Handler {
	return handler(client, projectID, false, opts...)
}

func handler(client *logging.Client, projectID string, logAll bool, opts ...logging.LoggerOption) func(http.Handler) http.Handler {
	parentLogger := client.Logger(
		"request_parent_log",
		opts...,
	)

	childLogger := client.Logger(
		"request_child_log",
		opts...,
	)

	return func(next http.Handler) http.Handler {
		return &gcpHandler{
			next:         next,
			parentLogger: parentLogger,
			childLogger:  childLogger,
			projectID:    projectID,
			logAll:       logAll,
		}
	}
}

type gcpHandler struct {
	next         http.Handler
	parentLogger CloudLogger
	childLogger  CloudLogger
	projectID    string
	logAll       bool
}

func (g *gcpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	traceID := traceIDFromRequest(r, g.projectID)
	l := NewGCPLogger(g.childLogger, traceID)
	r = r.WithContext(NewContext(r.Context(), l))
	sw := &statusWriter{ResponseWriter: w}

	g.next.ServeHTTP(sw, r)

	l.mu.Lock()
	logCount := l.logCount
	maxSeverity := l.maxSeverity
	l.mu.Unlock()

	if !g.logAll && logCount == 0 {
		return
	}

	// status code should also set the minimum maxSeverity to Error
	if sw.Status() > 399 && maxSeverity < logging.Error {
		maxSeverity = logging.Error
	}

	sc := trace.SpanFromContext(r.Context()).SpanContext()

	g.parentLogger.Log(logging.Entry{
		Timestamp:    begin,
		Severity:     maxSeverity,
		Trace:        traceID,
		SpanID:       sc.SpanID().String(),
		TraceSampled: sc.IsSampled(),
		Payload: payload{
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
}

// traceIDFromRequest formats a trace_id value for GCP Stackdriver
func traceIDFromRequest(r *http.Request, projectID string) string {
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

func requestSize(length string) int64 {
	l, err := strconv.Atoi(length)
	if err != nil {
		return 0
	}

	return int64(l)
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int64
}

func (w *statusWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}

	return w.status
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}

	n, err := w.ResponseWriter.Write(b)
	w.length += int64(n)
	if err != nil {
		return n, errors.Wrap(err, "http.ResponseWriter.Write()")
	}

	return n, nil
}
