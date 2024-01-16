package llog

import (
	"slices"
	"time"

	"github.com/jae2274/goutils/enum"
)

type LogLevelValues struct{}

type LogLevel = enum.Enum[LogLevelValues]

const (
	FATAL = LogLevel("FATAL")
	ERROR = LogLevel("ERROR")
	WARN  = LogLevel("WARN")
	INFO  = LogLevel("INFO")
	DEBUG = LogLevel("DEBUG")
)

func (LogLevelValues) Values() []string {
	return []string{string(FATAL), string(ERROR), string(WARN), string(INFO), string(DEBUG)}
}

type LogTime time.Time

func (lt LogTime) MarshalText() (text []byte, err error) {

	return []byte(time.Time(lt).Format(time.RFC3339Nano)), nil
}

type LLog struct {
	Level     LogLevel          `json:"level"`
	Msg       string            `json:"msg"`
	Tags      []string          `json:"tags,omitempty"`
	Datas     map[string]string `json:"datas,omitempty"`
	CreatedAt LogTime           `json:"createdAt"`
}

type LLogBuilder struct {
	level     LogLevel
	msg       string
	tags      []string
	datas     map[string]string
	createdAt LogTime
}

func Level(level LogLevel) *LLogBuilder {
	return &LLogBuilder{level: level}
}

func Msg(msg string) *LLogBuilder {
	return &LLogBuilder{msg: msg}
}

func Tag(tag string) *LLogBuilder {
	return &LLogBuilder{tags: []string{tag}}
}

func Tags(tags ...string) *LLogBuilder {
	return &LLogBuilder{tags: tags}
}

func Data(key, value string) *LLogBuilder {
	return &LLogBuilder{datas: map[string]string{key: value}}
}

func Datas(datas map[string]string) *LLogBuilder {
	return &LLogBuilder{datas: datas}
}

func (l *LLogBuilder) Level(level LogLevel) *LLogBuilder {
	l.level = level
	return l
}

func (l *LLogBuilder) Msg(msg string) *LLogBuilder {
	l.msg = msg
	return l
}

func (l *LLogBuilder) Tag(tag string) *LLogBuilder {
	if l.tags == nil {
		l.tags = []string{}
	}

	if !slices.Contains(l.tags, tag) {
		l.tags = append(l.tags, tag)
	}

	return l
}

func (l *LLogBuilder) Tags(tags ...string) *LLogBuilder {
	l.tags = tags
	return l
}

func (l *LLogBuilder) Data(key, value string) *LLogBuilder {
	if l.datas == nil {
		l.datas = map[string]string{}
	}
	l.datas[key] = value
	return l
}

func (l *LLogBuilder) Datas(datas map[string]string) *LLogBuilder {
	l.datas = datas
	return l
}

func (l *LLogBuilder) Build() *LLog {
	return &LLog{
		Level:     l.level,
		Msg:       l.msg,
		Tags:      l.tags,
		Datas:     l.datas,
		CreatedAt: l.createdAt,
	}
}
