package llog

import (
	"slices"
	"time"
)

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

func Log(llog *LLog) error {
	if llog.Datas == nil {
		llog.Datas = make(map[string]any)
	}

	if llog.Tags == nil {
		llog.Tags = make([]string, 0)
	}

	if llog.Level == "" {
		llog.Level = INFO
	}

	for k, v := range defaultDatas {
		if _, ok := llog.Datas[k]; !ok {
			llog.Datas[k] = v
		}
	}

	for _, tag := range defaultTags {
		if !slices.Contains(llog.Tags, tag) {
			llog.Tags = append(llog.Tags, tag)
		}
	}

	llog.CreatedAt = LogTime(time.Now())

	return logcfg.lloger.Log(llog)
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

func SetDefaultLLoger(lloger LLoger) {
	logcfg.lloger = lloger
}

var (
	defaultDatas           = make(map[string]any)
	defaultTags            = make([]string, 0)
	defaultContextDataKeys = make([]string, 0)
)

func SetDefaultDatas(datas map[string]any) {
	defaultDatas = datas
}

func SetDefaultTags(tags []string) {
	defaultTags = tags
}

func SetDefaultData(key string, value any) {
	defaultDatas[key] = value
}

func SetDefaultTag(tag string) {
	defaultTags = append(defaultTags, tag)
}

func SetDefaultContextData(key string) {
	defaultContextDataKeys = append(defaultContextDataKeys, key)
}

func SetDefaultContextDatas(keys []string) {
	defaultContextDataKeys = keys
}
