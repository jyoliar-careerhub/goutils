package httpmw

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/mw"
)

const (
	XRequestId string = "X-Request-Id"
)

func SetTraceIdMW(role string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = setHeaderToCtx(r, mw.CtxKeyTraceID, XRequestId, uuid.New().String())

			next.ServeHTTP(w, r)
		})
	}
}

func setHeaderToCtx(r *http.Request, ctxKey any, headerKey string, defauleValue string) *http.Request {
	originValues, ok := r.Header[headerKey]

	var value string
	if ok {
		value = strings.Join(originValues, ",")
	} else {
		value = defauleValue
	}

	return r.WithContext(mw.SetIfNotExists(r.Context(), ctxKey, value))
}
