package slogr

import (
	"context"
	"fmt"
	"github.com/lillrurre/slogr/color"
	"github.com/lillrurre/slogr/level"
	"io"
	"log/slog"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type Handler struct {
	opts           HandlerOptions
	preformatted   []byte   // data from WithGroup and WithAttrs
	unopenedGroups []string // groups from WithGroup that haven't been opened
	braces         int      // amount of braces to append at the end
	mu             *sync.Mutex
	out            io.Writer
}

type HandlerOptions struct {
	DisableTimeField bool
	Colorful         bool
	TimeFieldFormat  string
	Level            level.Level
	AddSource        bool
	ReplaceAttr      func(_ []string, attr slog.Attr) slog.Attr
}

func NewHandler(writer io.Writer, opts HandlerOptions) *Handler {
	return &Handler{
		opts: opts,
		mu:   new(sync.Mutex),
		out:  writer,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.opts.Level.Level() <= level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Add log level. This is the first attr, so add start brace and no comma before.
	buf = fmt.Appendf(buf, "{%q:%q,", slog.LevelKey, level.String(r.Level))

	// Add time field
	if !h.opts.DisableTimeField && !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time))
	}

	// Add source
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)))
	}

	// Add message
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message))

	// Insert preformatted attributes just after built-in ones.
	buf = append(buf, h.preformatted...)
	if r.NumAttrs() > 0 {
		buf = h.appendUnopenedGroups(buf)
		r.Attrs(func(a slog.Attr) bool {
			buf = h.appendAttr(buf, a)
			return true
		})
	}

	// Split the last comma and append braces + new line
	buf = fmt.Append(buf[:len(buf)-1], strings.Repeat("}", h.braces+1), "\n")

	if h.opts.Colorful {
		buf = append(color.From(level.Level(r.Level)), buf...)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h

	// Force an append to copy the underlying array and add all groups from WithGroup.
	h2.preformatted = h2.appendUnopenedGroups(slices.Clip(h.preformatted))

	// Now all groups have been opened.
	h2.unopenedGroups = nil

	for _, a := range attrs {
		h2.preformatted = h2.appendAttr(h2.preformatted, a)
	}
	return &h2
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	// Add an unopened group to h2 without modifying h.
	h2.unopenedGroups = make([]string, len(h.unopenedGroups)+1)
	copy(h2.unopenedGroups, h.unopenedGroups)
	h2.unopenedGroups[len(h2.unopenedGroups)-1] = name
	return &h2
}

func (h *Handler) appendAttr(buf []byte, a slog.Attr) []byte {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()

	// Ignore empty
	if a.Equal(slog.Attr{}) {
		return buf
	}

	switch a.Value.Kind() {
	case slog.KindString:
		buf = fmt.Appendf(buf, "%q:%q,", a.Key, a.Value.String())
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%q:%q,", a.Key, a.Value.Time().Format(h.opts.TimeFieldFormat))
	case slog.KindGroup:
		attrs := a.Value.Group()
		// Ignore empty groups.
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			buf = fmt.Appendf(buf, "%q:{", a.Key)
			h.braces++
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga)
		}
	default:
		buf = fmt.Appendf(buf, "%q:%q,", a.Key, a.Value)
	}
	return buf
}

func (h *Handler) appendUnopenedGroups(buf []byte) []byte {
	for _, group := range h.unopenedGroups {
		buf = fmt.Appendf(buf, "%q:{", group)
		h.braces++ // increment the amount of braces to append at the end
	}
	return buf
}

func (h *Handler) SetColorful(colorful bool) {
	h.opts = HandlerOptions{
		DisableTimeField: h.opts.DisableTimeField,
		Colorful:         colorful,
		TimeFieldFormat:  h.opts.TimeFieldFormat,
		Level:            h.opts.Level,
		AddSource:        h.opts.AddSource,
		ReplaceAttr:      h.opts.ReplaceAttr,
	}
}

func (h *Handler) WithColors() slog.Handler {
	h2 := *h
	h2.SetColorful(true)
	return &h2
}
