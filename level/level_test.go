package level

import (
	"log/slog"
	"testing"
)

func TestLevel_Level(t *testing.T) {
	testCases := []struct {
		in  Level
		out slog.Level
	}{
		{
			in:  Level(100),
			out: slog.Level(100),
		},
		{
			in:  Debug,
			out: slog.LevelDebug,
		},
		{
			in:  Info,
			out: slog.LevelInfo,
		},
		{
			in:  Warn,
			out: slog.LevelWarn,
		},
		{
			in:  Error,
			out: slog.LevelError,
		},
		{
			in:  Fatal,
			out: slog.Level(Error + 4),
		},
	}

	for _, testCase := range testCases {
		if testCase.out != testCase.in.Level() {
			t.Errorf("expected %+v, got %+v", testCase.out, testCase.in.Level())
		}
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		lvl      slog.Level
		expected string
	}{
		{
			lvl:      100,
			expected: "LEVEL 100",
		},
		{
			lvl:      slog.Level(Debug),
			expected: "DEBUG",
		},
		{
			lvl:      slog.Level(Info),
			expected: "INFO",
		},
		{
			lvl:      slog.Level(Warn),
			expected: "WARN",
		},
		{
			lvl:      slog.Level(Error),
			expected: "ERROR",
		},
		{
			lvl:      slog.Level(Fatal),
			expected: "FATAL",
		},
	}

	for _, testCase := range testCases {
		if testCase.expected != String(testCase.lvl) {
			t.Errorf("expected %s, got %s", testCase.expected, String(testCase.lvl))
		}
	}
}
