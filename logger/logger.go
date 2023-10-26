package logger

import (
	"context"
	"github.com/lillrurre/slogr/handler"
	"github.com/lillrurre/slogr/level"
	"io"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
}

type Options struct {
	// Level is an extension of slog.Level, by introducing FatalLevel.
	Level level.Level
	// DisableTimeField disables the time form log entries
	DisableTimeField bool
	// TimeFieldFormat. Defaults to time.RFC3339Nano.
	TimeFieldFormat string
	// AddSource adds the source of the log statement to every log entry
	AddSource bool
	// Tags are appended to the base logger.
	// map[string]string{"version": "0.1.2"} would output "version": "0.1.2" in every log entry.
	Tags     map[string]string
	Colorful bool
}

func NewLogger(opts *Options, writers ...io.Writer) *Logger {

	// use stdout if not writers are specified.
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}

	if opts.TimeFieldFormat == "" {
		opts.TimeFieldFormat = time.RFC3339Nano
	}

	handlerOpts := handler.Options{
		DisableTimeField: opts.DisableTimeField,
		Colorful:         opts.Colorful,
		TimeFieldFormat:  opts.TimeFieldFormat,
		Level:            opts.Level,
		AddSource:        opts.AddSource,
		ReplaceAttr:      nil,
	}

	h := handler.NewHandler(io.MultiWriter(writers...), handlerOpts)
	logger := slog.New(h)

	for key, val := range opts.Tags {
		logger = logger.With(key, val)
	}

	return &Logger{logger}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.DebugContext(context.Background(), msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.InfoContext(context.Background(), msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.WarnContext(context.Background(), msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.ErrorContext(context.Background(), msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.FatalContext(context.Background(), msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.Level(level.FatalLevel), msg, args...)
	os.Exit(1)
}

func (l *Logger) With(args ...any) *Logger {
	ll := l.Logger.With(args...)
	return &Logger{ll}
}

func (l *Logger) WithGroup(name string) *Logger {
	ll := l.Logger.WithGroup(name)
	return &Logger{ll}
}
