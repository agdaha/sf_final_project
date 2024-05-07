package logger

import (
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		level string
		env   string
	}
	tests := []struct {
		name string
		args args
		want *slog.Logger
	}{
		{name: "Text logger INFO",
			args: args{
				level: "INFO",
				env:   "devel",
			},
			want: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		},
		{name: "Text logger DEBUG",
			args: args{
				level: "DEBUG",
				env:   "devel",
			},
			want: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		},
		{name: "JSON logger INFO",
			args: args{
				level: "INFO",
				env:   "production",
			},
			want: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		},
		{name: "JSON logger DEBUG",
			args: args{
				level: "DEBUG",
				env:   "production",
			},
			want: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.level, tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
