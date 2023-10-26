package color

import (
	"github.com/lillrurre/slogr"
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	testCases := []struct {
		expected []byte
		lvl      slogr.Level
	}{
		{
			expected: White,
			lvl:      100,
		},
		{
			expected: Debug,
			lvl:      slogr.Debug,
		},
		{
			expected: Info,
			lvl:      slogr.Info,
		},
		{
			expected: Warn,
			lvl:      slogr.Warn,
		},
		{
			expected: Error,
			lvl:      slogr.Error,
		},
		{
			expected: Fatal,
			lvl:      slogr.Fatal,
		},
	}

	for _, testCase := range testCases {
		if !reflect.DeepEqual(testCase.expected, From(testCase.lvl)) {
			t.Errorf("expected %+v, got %+v", testCase.expected, From(testCase.lvl))
		}
	}
}
