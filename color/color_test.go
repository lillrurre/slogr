package color

import (
	"github.com/lillrurre/slogr/level"
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	testCases := []struct {
		expected []byte
		lvl      level.Level
	}{
		{
			expected: White,
			lvl:      100,
		},
		{
			expected: Debug,
			lvl:      level.Debug,
		},
		{
			expected: Info,
			lvl:      level.Info,
		},
		{
			expected: Warn,
			lvl:      level.Warn,
		},
		{
			expected: Error,
			lvl:      level.Error,
		},
		{
			expected: Fatal,
			lvl:      level.Fatal,
		},
	}

	for _, testCase := range testCases {
		if !reflect.DeepEqual(testCase.expected, From(testCase.lvl)) {
			t.Errorf("expected %+v, got %+v", testCase.expected, From(testCase.lvl))
		}
	}
}
