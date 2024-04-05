package grpcmw

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SetTraceIdUnaryMW() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := ctxToMetadata(ctx, mw.CtxKeyTraceID, uuid.New().String())

		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func GetTraceIdUnaryMW() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := metadataToCtx(ctx, mw.CtxKeyTraceID, uuid.New().String())

		return handler(newCtx, req)
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

func SetTraceIdStreamMW() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := ctxToMetadata(ctx, mw.CtxKeyTraceID, uuid.New().String())

		return streamer(newCtx, desc, cc, method, opts...)
	}
}

func GetTraceIdStreamMW() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		newCtx := metadataToCtx(ctx, mw.CtxKeyTraceID, uuid.New().String())

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

func ctxToMetadata(ctx context.Context, key string, defauleValue string) context.Context {
	value := ctx.Value(key)

	var valueStr string
	if value == nil {
		valueStr = defauleValue
	} else {
		originValue, ok := value.(string)
		if !ok {
			llog.Msg("context value is not string").Level(llog.WARN).Data("key", key).Data("value", value).Log(ctx)
			valueStr = defauleValue
		} else {
			valueStr = originValue
		}
	}

	return metadata.AppendToOutgoingContext(ctx, key, valueStr)
}

func metadataToCtx(ctx context.Context, key string, defaultValue string) context.Context {
	values := metadata.ValueFromIncomingContext(ctx, key)

	var value string
	if len(values) == 0 {
		value = defaultValue
	} else {
		value = strings.Join(values, ",")
	}

	return mw.SetIfNotExists(ctx, key, value)
}
