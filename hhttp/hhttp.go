package hhttp

import (
	"context"
	"net/http"

	"github.com/jae2274/goutils/llog"
)

func ErrorHandler(ctx context.Context, w http.ResponseWriter, err error) bool {
	return ErrorHandlerWithMsg(ctx, w, err, "Internal Server Error")
}

func ErrorHandlerWithMsg(ctx context.Context, w http.ResponseWriter, err error, msg string) bool {
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return true
	}
	return false
}
