package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type Config struct {
	Mode string
}

func New(cfg Config) *slog.Logger {
	return NewWithWriter(cfg, os.Stdout)
}

func NewWithWriter(cfg Config, w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if cfg.Mode == "debug" {
		return slog.New(tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}))
	}

	return slog.New(slog.NewJSONHandler(w, opts))
}

func SetDefault(l *slog.Logger) {
	if l == nil {
		return
	}
	slog.SetDefault(l)
}
