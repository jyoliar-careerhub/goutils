package llog

import (
	"context"
	"slices"
)

type LLogError struct {
	Err   error
	datas map[string]any
	tags  []string
}

func NewLLogError(err error) *LLogError {
	return &LLogError{Err: err}
}

func (l *LLogError) Data(key string, value any) *LLogError {
	if l.datas == nil {
		l.datas = map[string]any{}
	}
	l.datas[key] = value
	return l
}

func (l *LLogError) Tag(tag string) *LLogError {
	if !slices.Contains(l.tags, tag) {
		l.tags = append(l.tags, tag)
	}
	return l
}

func (l *LLogError) Tags(tags ...string) *LLogError {
	l.tags = tags
	return l
}

func (l *LLogError) Datas(datas map[string]any) *LLogError {
	l.datas = datas
	return l
}

func (l *LLogError) Error() string {
	return l.Err.Error()
}

func (l *LLogError) Log(ctx context.Context) error {
	return Msg(l.Err.Error()).Level(ERROR).Datas(l.datas).Tags(l.tags...).Log(ctx)
}
