package slogr

import (
	"bytes"
	"context"
	"github.com/lillrurre/slogr/level"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

const testFile = "test.log"

func testLogger(lvl level.Level, writer io.Writer) *Logger {
	return NewLogger(&Options{
		Level:            lvl,
		DisableTimeField: true,
		Tags:             map[string]string{"test": "log"},
		Colorful:         false,
	}, writer)
}

func TestDefault(t *testing.T) {
	l := Default()
	if !l.Enabled(context.Background(), slog.LevelDebug) {
		t.Fatalf("expected log level to be debug")
	}
}

func TestNewLogger(t *testing.T) {
	NewLogger(&Options{})
}

func TestLogger_Debug(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Debug, f)

	// Successful log
	{
		expected := []byte(`{"level":"DEBUG","msg":"test debug","test":"log"}`)

		l.Debug("test debug")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %v\ngot:      %v", expected, b)
		}
	}
	_ = os.Remove(testFile)

}

func TestLogger_DebugContext(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Debug, f)

	// Successful log
	{
		expected := []byte(`{"level":"DEBUG","msg":"test debug","test":"log"}`)

		l.DebugContext(context.Background(), "test debug")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_Info(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Info, f)

	// Successful log
	{
		expected := []byte(`{"level":"INFO","msg":"test info","test":"log"}`)
		l.Info("test info")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.Debug("should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_InfoContext(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Info, f)

	// Successful log
	{
		expected := []byte(`{"level":"INFO","msg":"test info","test":"log"}`)

		l.InfoContext(context.Background(), "test info")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.DebugContext(context.Background(), "should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_Warn(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Warn, f)

	// Successful log
	{
		expected := []byte(`{"level":"WARN","msg":"test warn","test":"log"}`)
		l.Warn("test warn")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.Debug("should not write")
		l.Info("should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_WarnContext(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Warn, f)

	// Successful log
	{
		expected := []byte(`{"level":"WARN","msg":"test warn","test":"log"}`)

		l.WarnContext(context.Background(), "test warn")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.DebugContext(context.Background(), "should not write")
		l.InfoContext(context.Background(), "should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_Error(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Error, f)

	// Successful log
	{
		expected := []byte(`{"level":"ERROR","msg":"test error","test":"log"}`)
		l.Error("test error")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.Debug("should not write")
		l.Info("should not write")
		l.Warn("should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_ErrorContext(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	l := testLogger(level.Error, f)

	// Successful log
	{
		expected := []byte(`{"level":"ERROR","msg":"test error","test":"log"}`)

		l.ErrorContext(context.Background(), "test error")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.DebugContext(context.Background(), "should not write")
		l.InfoContext(context.Background(), "should not write")
		l.WarnContext(context.Background(), "should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_Fatal(t *testing.T) {
	// This test starts the test as a separate process with $TEST_FATAL=1
	// This will result in a FATAL log and exit(1)
	// The content should be written to the file and the status code must be 1

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	l := testLogger(level.Fatal, f)

	// The actual test
	if os.Getenv("TEST_FATAL") == "1" {
		l.Fatal("test fatal")
	}

	// Successful log
	{
		// Run it as a separate task
		cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
		cmd.Env = append(os.Environ(), "TEST_FATAL=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); !ok && e.Success() {
			t.Errorf("expected test to result in exit status 1")
		}

		expected := []byte(`{"level":"FATAL","msg":"test fatal","test":"log"}`)
		// Check the output
		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.Debug("should not write")
		l.Info("should not write")
		l.Warn("should not write")
		l.Error("should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	_ = os.Remove(testFile)
}

func TestLogger_FatalContext(t *testing.T) {
	// This test starts the test as a separate process with $TEST_FATAL=1
	// This will result in a FATAL log and exit(1)
	// The content should be written to the file and the status code must be 1

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	l := testLogger(level.Fatal, f)

	// The actual test
	if os.Getenv("TEST_FATAL") == "1" {
		l.FatalContext(context.Background(), "test fatal")
	}

	{
		// Run it as a separate task
		cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
		cmd.Env = append(os.Environ(), "TEST_FATAL=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); !ok && e.Success() {
			t.Errorf("expected test to result in exit status 1")
		}

		expected := []byte(`{"level":"FATAL","msg":"test fatal","test":"log"}`)

		// Check the output
		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	// Too low level logs nothing
	{
		var expected []byte

		if err = os.Truncate(testFile, 0); err != nil {
			t.Fatalf("failed to truncate file: %v", err)
		}
		// lower than the log level
		l.DebugContext(context.Background(), "should not write")
		l.InfoContext(context.Background(), "should not write")
		l.WarnContext(context.Background(), "should not write")
		l.ErrorContext(context.Background(), "should not write")

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}
	_ = os.Remove(testFile)
}

func TestLogger_With(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	l := testLogger(level.Info, f)

	// No args return same logger
	{
		nl := l.With()
		if !reflect.DeepEqual(nl.Handler(), l.Handler()) {
			t.Error("expected loggers to be the same")
		}
	}

	// With args returns a new logger
	{
		nl := l.With("test", "logger")
		if reflect.DeepEqual(nl, l) {
			t.Error("expected loggers to be different")
		}
	}

	// Test output
	{
		expected := []byte(`{"level":"INFO","msg":"test info","test":"log","with":"args","extra":"stuff"}`)
		l.With("with", "args").Info("test info", "extra", "stuff")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	_ = os.Remove(testFile)
}

func TestLogger_WithGroup(t *testing.T) {
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	l := testLogger(level.Info, f)

	// WithGroup returns new logger
	{
		nl := l.WithGroup("group")
		if reflect.DeepEqual(nl, l) {
			t.Error("expected loggers to be different")
		}
	}

	// Test output
	{
		expected := []byte(`{"level":"INFO","msg":"test info","test":"log","group":{"with":"args","extra":"stuff"}}`)
		l.WithGroup("group").With("with", "args").Info("test info", "extra", "stuff")
		_ = f.Close()

		b, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("read error: %v", err)
		}

		if !bytes.Equal(expected, b[:len(b)-1]) {
			t.Errorf("\nexpected: %s\ngot:      %s", expected, string(b))
		}
	}

	_ = os.Remove(testFile)
}
