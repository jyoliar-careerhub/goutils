package grpcmw

import (
	"context"

	"github.com/google/uuid"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
	"google.golang.org/grpc"
)

func SetTraceIdUnaryMW() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		ctx = mw.SetIfNotExists(ctx, mw.CtxKeyTraceID, uuid.New().String())

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func LogErrUnaryMW() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		resp, err := handler(ctx, req)
		if err != nil {
			llog.LogErr(ctx, err)
		}
		return resp, err
	}
}

// wrappedStream 구조체는 grpc.ServerStream을 래핑하고,
// 추가 컨텍스트를 포함합니다.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 메서드는 원래 스트림의 컨텍스트 대신
// 래핑된 스트림의 컨텍스트를 반환합니다.
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func SetTraceIdStreamMW() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		newCtx := mw.SetIfNotExists(ctx, mw.CtxKeyTraceID, uuid.New().String())

		wrapped := &wrappedStream{stream, newCtx}
		return handler(srv, wrapped)
	}
}

func LogErrStreamMW() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		if err != nil {
			llog.LogErr(stream.Context(), err)
		}
		return err
	}
}
