package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func New(env string, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					source.File = shortenFilePath(source.File)
					a.Value = slog.AnyValue(source)
				}
			}
			return a
		},
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	return logger
}

func shortenFilePath(fullPath string) string {
	if !strings.Contains(fullPath, "/") {
		return fullPath
	}

	patterns := []string{
		"github.com/",
		"auction-service/",
		"api-gateway/",
		"internal/",
		"pkg/",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(fullPath, pattern); idx != -1 {
			return fullPath[idx:]
		}
	}

	return filepath.Base(fullPath)
}
