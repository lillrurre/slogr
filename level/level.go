package level

import (
	"fmt"
	"log/slog"
)

type Level int

const (
	DebugLevel = Level(slog.LevelDebug)
	InfoLevel  = Level(slog.LevelInfo)
	WarnLevel  = Level(slog.LevelWarn)
	ErrorLevel = Level(slog.LevelError)
	FatalLevel = Level(slog.LevelError + 4)
)

func (l Level) Level() slog.Level {
	return slog.Level(l)
}

func String(level slog.Level) string {
	switch Level(level) {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL %d", level)
	}
}
