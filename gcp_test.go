package logger

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"

	"cloud.google.com/go/logging"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		lg      *logging.Logger
		traceID string
	}
	tests := []struct {
		name string
		args args
		want ctxLogger
	}{
		{
			name: "new",
			args: args{
				lg:      &logging.Logger{},
				traceID: "hello",
			},
			want: &gcpLogger{
				lg:      &logging.Logger{},
				traceID: "hello",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := newGCPLogger(tt.args.lg, tt.args.traceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logger(t *testing.T) {
	t.Parallel()

	type args struct {
		format string
		v      []interface{}
		v2     interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantDebug  string
		wantDebugf string
		wantInfo   string
		wantInfof  string
		wantWarn   string
		wantWarnf  string
		wantError  string
		wantErrorf string
	}{
		{
			name: "Strings",
			args: args{
				format: "Formatted %s",
				v:      []interface{}{"Message"},
				v2:     "Message",
			},
			wantDebug:  "Message",
			wantDebugf: "Formatted Message",
			wantInfo:   "Message",
			wantInfof:  "Formatted Message",
			wantWarn:   "Message",
			wantWarnf:  "Formatted Message",
			wantError:  "Message",
			wantErrorf: "Formatted Message",
		},
		{
			name: "String & Error",
			args: args{
				format: "Formatted %s",
				v:      []interface{}{"Message"},
				v2:     errors.New("Message"),
			},
			wantDebug:  "Message",
			wantDebugf: "Formatted Message",
			wantInfo:   "Message",
			wantInfof:  "Formatted Message",
			wantWarn:   "Message",
			wantWarnf:  "Formatted Message",
			wantError:  "Message",
			wantErrorf: "Formatted Message",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			var buf bytes.Buffer

			l := &gcpLogger{
				lg: &testLogger{
					buf: &buf,
				},
			}

			l.Debug(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantDebug {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s, tt.wantDebug)
			}
			buf.Reset()

			l.Debugf(ctx, tt.args.format, tt.args.v...)
			if s := buf.String(); s != tt.wantDebugf {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s, tt.wantDebugf)
			}
			buf.Reset()

			l.Info(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantInfo {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s, tt.wantInfo)
			}
			buf.Reset()

			l.Infof(ctx, tt.args.format, tt.args.v...)
			if s := buf.String(); s != tt.wantInfof {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s, tt.wantInfof)
			}
			buf.Reset()

			l.Warn(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantWarn {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s, tt.wantWarn)
			}
			buf.Reset()

			l.Warnf(ctx, tt.args.format, tt.args.v...)
			if s := buf.String(); s != tt.wantWarnf {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s, tt.wantWarnf)
			}
			buf.Reset()

			l.Error(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantError {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s, tt.wantError)
			}
			buf.Reset()

			l.Errorf(ctx, tt.args.format, tt.args.v...)
			if s := buf.String(); s != tt.wantErrorf {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s, tt.wantErrorf)
			}
			buf.Reset()
		})
	}
}

type testLogger struct {
	buf *bytes.Buffer
}

func (t *testLogger) Log(e logging.Entry) {
	_, _ = t.buf.WriteString(e.Payload.(gcpPayload).Message.(string))
}
