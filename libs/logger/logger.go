package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gitlab.com/subrotokumar/glitchr/pkg/core"
)

type Env string

const (
	EnvDevelopment Env = "dev"
	EnvStaging     Env = "stage"
	EnvProduction  Env = "prod"
)

func New(env core.Env, service string, level *string) *slog.Logger {

	slogLevel := slog.LevelInfo
	if level != nil {
		slogLevel = parseLevel(*level)
	}
	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: env != core.EnvDevelopment,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "time"
				a.Value = slog.StringValue(
					time.Now().Local().Format("2006/01/02 15:04:05"),
				)

			case slog.LevelKey:
				a.Key = "level"

			case slog.SourceKey:
				a.Key = "caller"
			}
			return a
		},
	}

	var handler slog.Handler

	switch env {
	case core.EnvDevelopment:
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)

	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	svc := fmt.Sprintf("%s:%s", service, env)
	return slog.New(handler).With(
		slog.String("app", svc),
	)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
