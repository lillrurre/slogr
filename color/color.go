package color

import "github.com/lillrurre/slogr/level"

var (
	White = []byte("\033[0m")
	Debug = []byte("\033[0;34m")
	Info  = []byte("\033[0m")
	Warn  = []byte("\033[0;33m")
	Error = []byte("\033[0;31m")
	Fatal = []byte("\033[0;35m")
)

func From(l level.Level) []byte {
	switch l {
	case level.DebugLevel:
		return Debug
	case level.InfoLevel:
		return Info
	case level.WarnLevel:
		return Warn
	case level.ErrorLevel:
		return Error
	case level.FatalLevel:
		return Fatal
	default:
		return White
	}
}
