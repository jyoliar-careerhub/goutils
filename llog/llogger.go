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

func Log(llog *LLogBuilder) {
	llog.createdAt = LogTime(time.Now())
	if llog.level == "" {
		llog.level = INFO
	}

	logcfg.lloger.Log(llog.Build())
}

func SetLLoger(lloger LLoger) {
	logcfg.lloger = lloger
}
