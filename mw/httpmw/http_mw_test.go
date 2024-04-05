package httpmw

import (
	"context"
	"net/http"
	"testing"

	"github.com/jae2274/goutils/mw"
	"github.com/stretchr/testify/require"
)

func TestSetHeaderToCtx(t *testing.T) {
	t.Run("If header is not exists", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		r = setHeaderToCtx(r, mw.CtxKeyTraceID, XRequestId, "123")

		require.Equal(t, "123", r.Context().Value(mw.CtxKeyTraceID))
	})

	t.Run("If header is not exists but already existed in request context", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		r = r.WithContext(context.WithValue(r.Context(), mw.CtxKeyTraceID, "456"))

		r = setHeaderToCtx(r, mw.CtxKeyTraceID, XRequestId, "123")

		require.Equal(t, "456", r.Context().Value(mw.CtxKeyTraceID))
	})

	t.Run("If header is exists", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		r.Header.Add(XRequestId, "456")

		r = setHeaderToCtx(r, mw.CtxKeyTraceID, XRequestId, "123")

		require.Equal(t, "456", r.Context().Value(mw.CtxKeyTraceID))
	})

	t.Run("If header is exists with multiple values", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		r.Header.Add(XRequestId, "456")
		r.Header.Add(XRequestId, "789")

		r = setHeaderToCtx(r, mw.CtxKeyTraceID, XRequestId, "123")

		require.Equal(t, "456,789", r.Context().Value(mw.CtxKeyTraceID))
	})
}
