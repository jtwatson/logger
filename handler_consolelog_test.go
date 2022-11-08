//go:build consolelog

package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequestLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "success",
			want: &ConsoleLogger{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewRequestLogger(nil, "")

			var handlerCalled bool
			w := httptest.NewRecorder()
			handler(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					got := FromContext(r.Context())
					if _, ok := got.(*ConsoleLogger); !ok {
						t.Errorf("NewRequestLogger() context logger = %T, want %T", got, tt.want)
					}
					handlerCalled = true
				},
			)).ServeHTTP(w, &http.Request{})

			if !handlerCalled {
				t.Errorf("Failed to call handler")
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "success",
			want: &ConsoleLogger{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := NewLogger(nil, "")

			var handlerCalled bool
			w := httptest.NewRecorder()
			handler(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					got := FromContext(r.Context())
					if _, ok := got.(*ConsoleLogger); !ok {
						t.Errorf("NewLogger() context logger = %T, want %T", got, tt.want)
					}
					handlerCalled = true
				},
			)).ServeHTTP(w, &http.Request{})

			if !handlerCalled {
				t.Errorf("Failed to call handler")
			}
		})
	}
}
