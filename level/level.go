package level

import (
	"fmt"
	"log/slog"
)

type Level int

const (
	Debug = Level(slog.LevelDebug)
	Info  = Level(slog.LevelInfo)
	Warn  = Level(slog.LevelWarn)
	Error = Level(slog.LevelError)
	Fatal = Level(slog.LevelError + 4)
)

func (l Level) Level() slog.Level {
	return slog.Level(l)
}

func String(level slog.Level) string {
	switch Level(level) {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL %d", level)
	}
}
