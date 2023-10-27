package slogr

import (
	"context"
	"fmt"
	"github.com/lillrurre/slogr/color"
	"github.com/lillrurre/slogr/level"
	"io"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	w := new(io.Writer)
	h := NewHandler(*w, HandlerOptions{})
	if reflect.DeepEqual(w, h.out) {
		t.Errorf("expected %+v, got %+v", w, h.out)
	}
}

func TestLogHandler_Enabled(t *testing.T) {

	testInputs := []level.Level{level.Level(-100), level.Debug, level.Info, level.Warn, level.Error, level.Fatal}

	testCases := []struct {
		expected []bool
		lvl      level.Level
		h        *Handler
	}{
		{
			expected: []bool{false, true, true, true, true, true},
			lvl:      level.Debug,
		},
		{
			expected: []bool{false, false, true, true, true, true},
			lvl:      level.Info,
		},
		{
			expected: []bool{false, false, false, true, true, true},
			lvl:      level.Warn,
		},
		{
			expected: []bool{false, false, false, false, true, true},
			lvl:      level.Error,
		},
		{
			expected: []bool{false, false, false, false, false, true},
			lvl:      level.Fatal,
		},
	}

	for _, testCase := range testCases {
		h := NewHandler(os.Stdout, HandlerOptions{Level: testCase.lvl})
		for i, in := range testInputs {
			enabled := h.Enabled(context.Background(), in.Level())
			if testCase.expected[i] != enabled {
				t.Errorf("expected %+v, got %+v", testCase.expected[i], enabled)
			}
		}
	}

}

func TestLogHandler_Handle(t *testing.T) {
	// Test time is appended to the buffer
	{
		now := time.Now()
		f, err := os.OpenFile("handler.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			t.Fatalf("unexpected error opening file: %v", err)
		}
		h := NewHandler(f, HandlerOptions{TimeFieldFormat: time.RFC3339Nano})
		r := slog.NewRecord(now, slog.LevelInfo, "lol", uintptr(2))
		if err = h.Handle(context.Background(), r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_ = f.Close()

		expected := fmt.Sprintf(`{"level":"INFO","time":%q,"msg":"lol"}`, now.Format(time.RFC3339Nano))
		b, err := os.ReadFile("handler.log")
		if err != nil {
			t.Errorf("unexpecte read error: %v", err)
			_ = os.Remove("handler.log")
		}

		if expected != string(b[:len(b)-1]) {
			t.Errorf("expected %s, got %s", expected, string(b[:len(b)-1]))
		}
		_ = os.Remove("handler.log")
	}

	// Test time is appended to the buffer
	{
		f, err := os.OpenFile("handler.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			t.Fatalf("unexpected error opening file: %v", err)
		}
		h := NewHandler(f, HandlerOptions{DisableTimeField: true, AddSource: true})
		r := slog.NewRecord(time.Now(), slog.LevelInfo, "lol", uintptr(1))
		if err = h.Handle(context.Background(), r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_ = f.Close()

		expected := `{"level":"INFO","source":":0","msg":"lol"}`
		b, err := os.ReadFile("handler.log")
		if err != nil {
			t.Errorf("unexpecte read error: %v", err)
			_ = os.Remove("handler.log")
		}

		if expected != string(b[:len(b)-1]) {
			t.Errorf("expected %s, got %s", expected, string(b[:len(b)-1]))
		}
		_ = os.Remove("handler.log")
	}

	// Test colorful appends colors
	{
		f, err := os.OpenFile("handler.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			t.Fatalf("unexpected error opening file: %v", err)
		}
		h := NewHandler(f, HandlerOptions{DisableTimeField: true, Colorful: true})
		r := slog.NewRecord(time.Now(), slog.LevelDebug, "lol", uintptr(0))
		if err = h.Handle(context.Background(), r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_ = f.Close()

		expected := string(color.Debug) + `{"level":"DEBUG","msg":"lol"}`
		b, err := os.ReadFile("handler.log")
		if err != nil {
			t.Errorf("unexpecte read error: %v", err)
			_ = os.Remove("handler.log")
		}

		if expected != string(b[:len(b)-1]) {
			t.Errorf("expected %s, got %s", expected, string(b[:len(b)-1]))
		}
		_ = os.Remove("handler.log")
	}
}

func TestLogHandler_WithAttrs(t *testing.T) {

}

func TestLogHandler_WithGroup(t *testing.T) {

}
