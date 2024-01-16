package llog

type LLoger interface {
	Log(llog *LLog) error
}

func Fatal(msg string) {
	Level(FATAL).Msg(msg).Log()
}

func Error(msg string) {
	Level(ERROR).Msg(msg).Log()
}

func Warn(msg string) {
	Level(WARN).Msg(msg).Log()
}

func Info(msg string) {
	Level(INFO).Msg(msg).Log()
}

func Debug(msg string) {
	Level(DEBUG).Msg(msg).Log()
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

func SetLLoger(lloger LLoger) {
	logcfg.lloger = lloger
}
