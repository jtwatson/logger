package logger

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			name: "logger from ctx",
			args: args{
				NewContext(context.Background(), &GCPLogger{}),
			},
			want: &GCPLogger{},
		},
		{
			name: "StdErrLogger: ctx nil",
			want: &stdErrLogger{},
		},
		{
			name: "StdErrLogger: ctx empty",
			args: args{
				ctx: context.Background(),
			},
			want: &stdErrLogger{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromRequest(t *testing.T) {
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
			name: "nil request",
			want: &stdErrLogger{},
		},
		{
			name: "empty request ctx",
			args: args{
				r: &http.Request{},
			},
			want: &stdErrLogger{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FromRequest(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
