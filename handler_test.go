package logger

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func disableMetaServer(t *testing.T) {
	t.Helper()

	// Fix issue when logging.Client attempts to detect its
	// env by querying GCE_METADATA_HOST and nothing is there
	// so your test is very slow. This tries to causes the
	// detection to fail faster and not hang your test so long
	curEnv := os.Getenv("GCE_METADATA_HOST")
	t.Cleanup(func() { os.Setenv("GCE_METADATA_HOST", curEnv) })
	_ = os.Setenv("GCE_METADATA_HOST", "localhost")
}

// func TestNewRequestLogger_NewLogger(t *testing.T) {
// 	disableMetaServer(t)

// 	type args struct {
// 		constructor func(*logging.Client, string, ...logging.LoggerOption) func(http.Handler) http.Handler
// 		client      *logging.Client
// 		projectID   string
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantLogAll  bool
// 		wantLogger  interface{}
// 		wantHandler interface{}
// 	}{
// 		{
// 			name: "NewRequestLogger",
// 			args: args{
// 				constructor: NewRequestLogger,
// 				client:      &logging.Client{},
// 				projectID:   "my-Project",
// 			},
// 			wantLogAll:  true,
// 			wantLogger:  &GCPLogger{},
// 			wantHandler: &gcpHandler{},
// 		},
// 		{
// 			name: "NewLogger",
// 			args: args{
// 				constructor: NewLogger,
// 				client:      &logging.Client{},
// 				projectID:   "my-Project",
// 			},
// 			wantLogAll:  false,
// 			wantLogger:  &GCPLogger{},
// 			wantHandler: &gcpHandler{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			handler := tt.args.constructor(tt.args.client, tt.args.projectID)

// 			got, ok := handler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).(*gcpHandler)
// 			if !ok {
// 				t.Errorf("NewRequestLogger() handler = %T, wantHandler %T", got, tt.wantHandler)
// 			}
// 			if got.logAll != tt.wantLogAll {
// 				t.Errorf("NewRequestLogger() handler.logAll = %v, wantLogAll %v", got.logAll, tt.wantLogAll)
// 			}
// 			if got.projectID != tt.args.projectID {
// 				t.Errorf("NewRequestLogger() handler.projectID = %v, want %v", got.projectID, tt.args.projectID)
// 			}
// 			if diff := deep.Equal(got.next, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})); diff != nil {
// 				t.Errorf("NewRequestLogger() handler.next = %v", diff)
// 			}
// 			if got.parentLogger == nil {
// 				t.Errorf("NewRequestLogger() handler.parentLogger = %v", got.parentLogger)
// 			}
// 			if got.childLogger == nil {
// 				t.Errorf("NewRequestLogger() handler.childLogger = %v", got.childLogger)
// 			}
// 		})
// 	}
// }

// func Test_gcpHandler_ServeHTTP(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		status int
// 		logs   int
// 		level  logging.Severity
// 	}
// 	type fields struct {
// 		projectID string
// 		logAll    bool
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		wantLevel logging.Severity
// 	}{
// 		{
// 			name: "logAll=true",
// 			fields: fields{
// 				projectID: "my-big-project",
// 				logAll:    true,
// 			},
// 			args: args{
// 				status: http.StatusOK,
// 				logs:   1,
// 				level:  logging.Info,
// 			},
// 			wantLevel: logging.Info,
// 		},
// 		{
// 			name: "logAll=true no logging",
// 			fields: fields{
// 				projectID: "my-big-project",
// 				logAll:    true,
// 			},
// 			args: args{
// 				status: http.StatusOK,
// 			},
// 			wantLevel: logging.Default,
// 		},
// 		{
// 			name: "logAll=false no logging",
// 			fields: fields{
// 				projectID: "my-big-project",
// 			},
// 			args: args{
// 				status: http.StatusOK,
// 			},
// 		},
// 		{
// 			name: "logAll=false with logging",
// 			fields: fields{
// 				projectID: "my-bigger-project",
// 			},
// 			args: args{
// 				status: http.StatusOK,
// 				logs:   1,
// 				level:  logging.Warning,
// 			},
// 			wantLevel: logging.Warning,
// 		},
// 		{
// 			name: "logAll=true no logging",
// 			fields: fields{
// 				projectID: "my-big-project",
// 				logAll:    true,
// 			},
// 			args: args{
// 				status: http.StatusInternalServerError,
// 			},
// 			wantLevel: logging.Error,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			var handlerCalled bool
// 			var traceID string
// 			l := &captureLogger{}
// 			handler := &gcpHandler{
// 				parentLogger: l,
// 				childLogger:  &captureLogger{},
// 				projectID:    tt.fields.projectID,
// 				logAll:       tt.fields.logAll,
// 				next: http.HandlerFunc(
// 					func(w http.ResponseWriter, r *http.Request) {
// 						for i := 0; i < tt.args.logs; i++ {
// 							switch tt.args.level {
// 							case logging.Info:
// 								Info(r.Context(), "some log")
// 							case logging.Warning:
// 								Warn(r.Context(), "some log")
// 							case logging.Error:
// 								Error(r.Context(), "some log")
// 							default:
// 							}
// 						}

// 						if l, ok := FromContext(r.Context()).(*GCPLogger); ok {
// 							traceID = l.traceID
// 						} else {
// 							traceID = "not found in child logger"
// 						}

// 						w.WriteHeader(tt.args.status)
// 						handlerCalled = true
// 					},
// 				),
// 			}

// 			w := httptest.NewRecorder()
// 			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
// 			handler.ServeHTTP(w, r)

// 			if !handlerCalled {
// 				t.Errorf("Failed to call handler")
// 			}
// 			if tt.args.logs == 0 {
// 				return
// 			}
// 			if l.e.Severity != tt.wantLevel {
// 				t.Errorf("Severity = %v, want %v", l.e.Severity, tt.wantLevel)
// 			}
// 			if l.e.Trace != traceID {
// 				t.Errorf("Trace = %v, want %v", l.e.Trace, traceID)
// 			}
// 			if pl, ok := l.e.Payload.(payload); ok {
// 				if m, ok := pl.Message.(string); ok {
// 					if m != "Parent Log Entry" {
// 						t.Errorf("Message = %v, want %v", m, "Parent Log Entry")
// 					}
// 				} else {
// 					t.Fatalf("Message = %T, want %T", pl.Message, "")
// 				}
// 			} else {
// 				t.Fatalf("Payload = %T, want %T", l.e.Payload, payload{})
// 			}
// 			if l.e.HTTPRequest.Status != tt.args.status {
// 				t.Errorf("Status = %v, want %v", l.e.HTTPRequest.Status, tt.args.status)
// 			}
// 		})
// 	}
// }

type captureLogger struct {
	e logging.Entry
}

func (c *captureLogger) Log(e logging.Entry) {
	c.e = e
}

func Test_traceIDFromRequest(t *testing.T) {
	type args struct {
		mockReq   func(traceStr string) (*http.Request, string)
		projectID string
	}
	tests := []struct {
		name            string
		args            args
		wantTracePrefix string
		wantTraceStr    string
	}{
		// The order these are significant
		{
			// This test relies on the global tracing provider NOT being set
			name: "no trace in request",
			args: args{
				mockReq: func(wantTraceStr string) (*http.Request, string) {
					return &http.Request{URL: &url.URL{}}, wantTraceStr
				},
				projectID: "my-project",
			},
			wantTracePrefix: "projects/my-project/traces/",
			wantTraceStr:    "00000000000000000000000000000000",
		},
		{
			// This test sets the global tracing provider (I don't think this can be un-done)
			name: "with trace in request",
			args: args{
				mockReq: func(_ string) (r *http.Request, traceStr string) {
					otel.SetTracerProvider(sdktrace.NewTracerProvider())
					ctx, span := otel.Tracer("test/examples").Start(context.Background(), "test trace")

					r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
					r = r.WithContext(ctx)

					return r, span.SpanContext().TraceID().String()
				},
				projectID: "my-project",
			},
			wantTracePrefix: "projects/my-project/traces/",
		},
		{
			// With the global tracing provider set, this test shows that
			// trace Propagation is a higher priority then trace in request context
			name: "with propagation span in headers",
			args: args{
				mockReq: func(wantTraceStr string) (r *http.Request, traceStr string) {
					r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
					r.Header.Add("X-Cloud-Trace-Context", wantTraceStr+"/1;o=1")

					return r, wantTraceStr
				},
				projectID: "my-project",
			},
			wantTracePrefix: "projects/my-project/traces/",
			wantTraceStr:    "105445aa7843bc8bf206b12000100000",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r, traceStr := tt.args.mockReq(tt.wantTraceStr)
			want := tt.wantTracePrefix + traceStr
			if got := gcpTraceIDFromRequest(r, tt.args.projectID); got != want {
				t.Errorf("traceIDFromRequest() = %v, want %v", got, want)
			}
		})
	}
}

func Test_requestSize(t *testing.T) {
	t.Parallel()

	type args struct {
		length string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "success",
			args: args{
				length: "20",
			},
			want: 20,
		},
		{
			name: "falure",
			args: args{
				length: "xxx",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := requestSize(tt.args.length); got != tt.want {
				t.Errorf("requestSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusWriter_Status(t *testing.T) {
	t.Parallel()

	type fields struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Status set",
			fields: fields{
				status: http.StatusForbidden,
			},
			want: 403,
		},
		{
			name: "Status not set",
			want: 200,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &statusWriter{
				status: tt.fields.status,
			}
			if got := w.Status(); got != tt.want {
				t.Errorf("statusWriter.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusWriter_WriteHeader(t *testing.T) {
	t.Parallel()

	type fields struct {
		ResponseWriter http.ResponseWriter
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Success",
			fields: fields{
				ResponseWriter: &httptest.ResponseRecorder{},
			},
			args: args{
				status: 201,
			},
			want: 201,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &statusWriter{
				ResponseWriter: tt.fields.ResponseWriter,
			}
			w.WriteHeader(tt.args.status)
			if got := w.Status(); got != tt.want {
				t.Errorf("statusWriter.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusWriter_Write(t *testing.T) {
	t.Parallel()

	type fields struct {
		ResponseWriter http.ResponseWriter
		status         int
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLength int
		wantStatus int
		wantErr    bool
	}{
		{
			name: "No status set",
			fields: fields{
				ResponseWriter: &httptest.ResponseRecorder{},
			},
			args: args{
				b: []byte("0123456789"),
			},
			wantLength: 10,
			wantStatus: 200,
		},
		{
			name: "Status set",
			fields: fields{
				ResponseWriter: &httptest.ResponseRecorder{},
				status:         201,
			},
			args: args{
				b: []byte("01234567891234567890"),
			},
			wantLength: 20,
			wantStatus: 201,
		},
		{
			name: "Write error",
			fields: fields{
				ResponseWriter: &responseRecorder{err: errors.New("Bang")},
				status:         201,
			},
			args: args{
				b: []byte("01234567891234567890"),
			},
			wantLength: 20,
			wantStatus: 201,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &statusWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				status:         tt.fields.status,
			}
			got, err := w.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("statusWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.wantLength {
				t.Errorf("statusWriter.Write() = %v, wantLength %v", got, tt.wantLength)
			}
			if got := w.Status(); got != tt.wantStatus {
				t.Errorf("statusWriter.Status() = %v, wantStatus %v", got, tt.wantStatus)
			}
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	err error
}

func (rw *responseRecorder) Write(buf []byte) (int, error) {
	return len(buf), rw.err
}
