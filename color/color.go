package color

import (
	"github.com/lillrurre/slogr/level"
)

var (
	White = []byte("\033[0m")
	Debug = []byte("\033[0;34m")
	Info  = White
	Warn  = []byte("\033[0;33m")
	Error = []byte("\033[0;31m")
	Fatal = []byte("\033[0;35m")
)

func From(l level.Level) []byte {
	switch l {
	case level.Debug:
		return Debug
	case level.Info:
		return Info
	case level.Warn:
		return Warn
	case level.Error:
		return Error
	case level.Fatal:
		return Fatal
	default:
		return White
	}
}
