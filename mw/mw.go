package mw

import "context"

type ctxKeyType string

const (
	CtxKeyTraceID string = "trace_id"
)

func SetIfNotExists(ctx context.Context, key, value any) context.Context {
	if _, ok := ctx.Value(key).(string); !ok {
		return context.WithValue(ctx, key, value)
	}
	return ctx
}
