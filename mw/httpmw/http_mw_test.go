package httpmw

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/llog"
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

func TestHttpServer(t *testing.T) {
	url, err := initHttp(t)
	require.NoError(t, err)

	t.Run(fmt.Sprintf("If header \"%s\" is set", XRequestId), func(t *testing.T) {
		xRequestIdValue := "test_ramdom_value"
		traceId, err := getTraceIdWithHeader(url, XRequestId, xRequestIdValue)
		require.NoError(t, err)
		require.Equal(t, xRequestIdValue, traceId)
	})

	t.Run(fmt.Sprintf("If header \"%s\" is not set", XRequestId), func(t *testing.T) {
		traceId, err := getTraceIdWithHeader(url, "OTHER_HEADER", "ItDoesn'tMatter")
		require.NoError(t, err)
		require.NotEmpty(t, traceId)
	})
}

func initHttp(t *testing.T) (string, error) {
	router := mux.NewRouter()
	router.HandleFunc("/get-ctx-id", func(w http.ResponseWriter, r *http.Request) {
		traceIdValue := r.Context().Value(mw.CtxKeyTraceID)
		if traceIdValue == nil {
			llog.Error(r.Context(), "traceId is not set")
			http.Error(w, "traceId is not set", http.StatusInternalServerError)
		}

		traceId, ok := traceIdValue.(string)
		if !ok {
			llog.Level(llog.ERROR).Msg("traceId is not string").Data("traceId", traceIdValue).Log(r.Context())
			http.Error(w, "traceId is not string", http.StatusInternalServerError)
		}

		w.Write([]byte(traceId))
	})
	router.Use(SetTraceIdMW())

	errChan := make(chan error)
	go func() {
		err := http.ListenAndServe(":5245", router)
		errChan <- err
	}()
	time.Sleep(1 * time.Second)
	select {
	case err := <-errChan:
		return "", err
	default:
	}

	url := "http://localhost:5245"
	return url, nil
}

func getTraceIdWithHeader(url string, headerKey, headerValue string) (string, error) {
	req, err := http.NewRequest("GET", url+"/get-ctx-id", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set(headerKey, headerValue)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
