package llog

import (
	"context"
)

type LLoger interface {
	Log(llog *LLog) error
}

func Fatal(ctx context.Context, msg string) error {
	return Level(FATAL).Msg(msg).Log(ctx)
}

func Error(ctx context.Context, msg string) error {
	return Level(ERROR).Msg(msg).Log(ctx)
}

func LogErr(ctx context.Context, err error) error {
	if llogErr, ok := err.(*LLogError); ok {
		return llogErr.Log(ctx)
	} else {
		return Error(ctx, err.Error())
	}
}

func Warn(ctx context.Context, msg string) error {
	return Level(WARN).Msg(msg).Log(ctx)
}

func Info(ctx context.Context, msg string) error {
	return Level(INFO).Msg(msg).Log(ctx)
}

func Debug(ctx context.Context, msg string) error {
	return Level(DEBUG).Msg(msg).Log(ctx)
}

func Log(llog *LLog) error {
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
	metadatas              = make(map[string]any)
	defaultTags            = make([]string, 0)
	defaultContextDataKeys = make([]string, 0)
)

func SetMetadatas(datas map[string]any) {
	metadatas = datas
}

func SetDefaultTags(tags []string) {
	defaultTags = tags
}

func SetMetadata(key string, value any) {
	metadatas[key] = value
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
