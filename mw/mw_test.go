package mw

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetIfNotExists(t *testing.T) {
	t.Run("If value is not exists", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), CtxKeyTraceID, "123")

		ctx = SetIfNotExists(ctx, CtxKeyTraceID, "456")

		require.Equal(t, "123", ctx.Value(CtxKeyTraceID))
	})

	t.Run("If value is exists", func(t *testing.T) {
		ctx := context.Background()

		ctx = SetIfNotExists(ctx, CtxKeyTraceID, "456")

		require.Equal(t, "456", ctx.Value(CtxKeyTraceID))
	})
}
