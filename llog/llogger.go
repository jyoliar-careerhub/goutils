package llog

type LLoger interface {
	Log(llog *LLog) error
}

func Fatal(msg string) error {
	return Level(FATAL).Msg(msg).Log()
}

func Error(msg string) error {
	return Level(ERROR).Msg(msg).Log()
}

func LogErr(err error) error {
	if llogErr, ok := err.(*LLogError); ok {
		return llogErr.Log()
	} else {
		return Error(err.Error())
	}
}

func Warn(msg string) error {
	return Level(WARN).Msg(msg).Log()
}

func Info(msg string) error {
	return Level(INFO).Msg(msg).Log()
}

func Debug(msg string) error {
	return Level(DEBUG).Msg(msg).Log()
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
