package grpcmw

import (
	"context"

	"github.com/google/uuid"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
	"google.golang.org/grpc"
)

func SetTraceIdMW() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		ctx = mw.SetIfNotExists(ctx, mw.CtxKeyTraceID, uuid.New().String())

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func LogErrMW() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		resp, err := handler(ctx, req)
		if err != nil {
			llog.LogErr(ctx, err)
		}
		return resp, err
	}
}
