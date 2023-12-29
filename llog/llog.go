package llog

import (
	"slices"
	"time"
)

type LLog struct {
	msg       string
	tags      []string
	datas     map[string]string
	createdAt time.Time
}

func Msg(msg string) *LLog {
	return &LLog{msg: msg}
}

func Tag(tag string) *LLog {
	return &LLog{tags: []string{tag}}
}

func Tags(tags ...string) *LLog {
	return &LLog{tags: tags}
}

func Data(key, value string) *LLog {
	return &LLog{datas: map[string]string{key: value}}
}

func Datas(datas map[string]string) *LLog {
	return &LLog{datas: datas}
}

func (l *LLog) Msg(msg string) *LLog {
	l.msg = msg
	return l
}

func (l *LLog) Tag(tag string) *LLog {
	if l.tags == nil {
		l.tags = []string{}
	}

	if !slices.Contains(l.tags, tag) {
		l.tags = append(l.tags, tag)
	}

	return l
}

func (l *LLog) Tags(tags ...string) *LLog {
	l.tags = tags
	return l
}

func (l *LLog) Data(key, value string) *LLog {
	if l.datas == nil {
		l.datas = map[string]string{}
	}
	l.datas[key] = value
	return l
}

func (l *LLog) Datas(datas map[string]string) *LLog {
	l.datas = datas
	return l
}
