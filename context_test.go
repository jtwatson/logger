package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestContextLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		v  []interface{}
		v2 interface{}
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
				v:  []interface{}{"Message"},
				v2: "Message",
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
				v:  []interface{}{"Message"},
				v2: errors.New("Message"),
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

			var buf bytes.Buffer
			ctx := NewContext(context.Background(), &GCPLogger{
				lg: &testLogger{
					buf: &buf,
				},
			})

			format := "Formatted %s"

			Debug(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantDebug {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s, tt.wantDebug)
			}
			buf.Reset()

			Debugf(ctx, format, tt.args.v...)
			if s := buf.String(); s != tt.wantDebugf {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s, tt.wantDebugf)
			}
			buf.Reset()

			Info(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantInfo {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s, tt.wantInfo)
			}
			buf.Reset()

			Infof(ctx, format, tt.args.v...)
			if s := buf.String(); s != tt.wantInfof {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s, tt.wantInfof)
			}
			buf.Reset()

			Warn(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantWarn {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s, tt.wantWarn)
			}
			buf.Reset()

			Warnf(ctx, format, tt.args.v...)
			if s := buf.String(); s != tt.wantWarnf {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s, tt.wantWarnf)
			}
			buf.Reset()

			Error(ctx, tt.args.v2)
			if s := buf.String(); s != tt.wantError {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s, tt.wantError)
			}
			buf.Reset()

			Errorf(ctx, format, tt.args.v...)
			if s := buf.String(); s != tt.wantErrorf {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s, tt.wantErrorf)
			}
			buf.Reset()
		})
	}
}
