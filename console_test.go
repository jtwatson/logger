package logger

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func Test_consoleLogger(t *testing.T) {
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
			name: "Test 1", args: args{v: []interface{}{"Message"}, v2: "Message"},
			wantDebug: "\x1b[37mDEBUG\x1b[0m: /path Message\n", wantDebugf: "\x1b[37mDEBUG\x1b[0m: /path Formatted Message\n",
			wantInfo: "\x1b[34mINFO \x1b[0m: /path Message\n", wantInfof: "\x1b[34mINFO \x1b[0m: /path Formatted Message\n",
			wantWarn: "\x1b[33mWARN \x1b[0m: /path Message\n", wantWarnf: "\x1b[33mWARN \x1b[0m: /path Formatted Message\n",
			wantError: "\x1b[31mERROR\x1b[0m: /path Message\n", wantErrorf: "\x1b[31mERROR\x1b[0m: /path Formatted Message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := context.Background()
			log.SetOutput(&buf)
			t.Cleanup(func() { log.SetOutput(os.Stderr) })

			u, _ := url.Parse("http://some.domain.com/path")
			l := &ConsoleLogger{r: &http.Request{URL: u}}
			format := "Formatted %s"

			l.Debug(ctx, tt.args.v2)
			if s := buf.String(); s[20:] != tt.wantDebug {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s[20:], tt.wantDebug)
			}
			buf.Reset()

			l.Debugf(ctx, format, tt.args.v...)
			if s := buf.String(); s[20:] != tt.wantDebugf {
				t.Errorf("stdErrLogger.Debug() value = %v, wantValue %v", s[20:], tt.wantDebugf)
			}
			buf.Reset()

			l.Info(ctx, tt.args.v2)
			if s := buf.String(); s[20:] != tt.wantInfo {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s[20:], tt.wantInfo)
			}
			buf.Reset()

			l.Infof(ctx, format, tt.args.v...)
			if s := buf.String(); s[20:] != tt.wantInfof {
				t.Errorf("stdErrLogger.Info() value = %v, wantValue %v", s[20:], tt.wantInfof)
			}
			buf.Reset()

			l.Warn(ctx, tt.args.v2)
			if s := buf.String(); s[20:] != tt.wantWarn {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s[20:], tt.wantWarn)
			}
			buf.Reset()

			l.Warnf(ctx, format, tt.args.v...)
			if s := buf.String(); s[20:] != tt.wantWarnf {
				t.Errorf("stdErrLogger.Warn() value = %v, wantValue %v", s[20:], tt.wantWarnf)
			}
			buf.Reset()

			l.Error(ctx, tt.args.v2)
			if s := buf.String(); s[20:] != tt.wantError {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s[20:], tt.wantError)
			}
			buf.Reset()

			l.Errorf(ctx, format, tt.args.v...)
			if s := buf.String(); s[20:] != tt.wantErrorf {
				t.Errorf("stdErrLogger.Error() value = %v, wantValue %v", s[20:], tt.wantErrorf)
			}
			buf.Reset()
		})
	}
}

func TestNewConsoleLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			name: "some request",
			args: args{
				r: &http.Request{},
			},
			want: &ConsoleLogger{r: &http.Request{}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewConsoleLogger(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsoleLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
