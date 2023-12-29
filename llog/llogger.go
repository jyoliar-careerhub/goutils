package llog

import (
	"time"
)

type LLoger interface {
	Log(llog *LLog) error
}

func Default(msg string) {
	Log(Msg(msg))
}

type logConfig struct {
	lloger LLoger
}

func newLogConfig(lloger LLoger) *logConfig {
	return &logConfig{
		lloger: lloger,
	}
}

var logcfg *logConfig = newLogConfig(&StdoutLLogger{})

func Log(llog *LLog) {
	llog.createdAt = time.Now()

	logcfg.lloger.Log(llog)
}

func SetLLoger(lloger LLoger) {
	logcfg.lloger = lloger
}
