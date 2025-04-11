package xslog

import (
	"fmt"
	"io"
	"log/slog"
	"math"
)

type Level string

// All available log levels.
const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelNone  Level = "none"
)

// All available log formats.
const (
	FormatText   = "text"
	FormatGCloud = "gcloud"
)

type Config struct {
	Level  Level
	Format string
}

func New(w io.Writer, config Config) *slog.Logger {
	lvl, err := parseLevel(config.Level)
	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}

	var handler slog.Handler
	switch config.Format {
	case FormatText:
		handler = slog.NewTextHandler(w, opts)
	case FormatGCloud:
		handler = newGoogleCloudHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	l := slog.New(handler)

	return l
}

// newGoogleCloudHandler creates a logger that fits the Google Cloud Logging format.
// More info: https://cloud.google.com/logging/docs/agent/logging/configuration#process-payload
func newGoogleCloudHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.LevelKey:
			a.Key = "severity"
		case slog.MessageKey:
			a.Key = "message"
		}
		return a
	}

	return slog.NewJSONHandler(w, opts)
}

// parseLevel returns the slog.Level based on the provided string.
func parseLevel(lvl Level) (slog.Level, error) {
	switch lvl {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	case LevelNone:
		// Level with math.MaxInt is used to disable logging.
		return math.MaxInt, nil
	default:
		return math.MaxInt, fmt.Errorf("log log_level %s unknown", lvl)
	}
}
