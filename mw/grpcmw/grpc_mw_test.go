package grpcmw

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/jae2274/goutils/mw"
	"github.com/jae2274/goutils/mw/grpcmw/internal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TestServer struct {
	internal.UnimplementedTestServiceServer
}

func (t *TestServer) GetCtxId(ctx context.Context, _ *emptypb.Empty) (*internal.ContextId, error) {
	return &internal.ContextId{Id: fmt.Sprintf("%v", ctx.Value(mw.CtxKeyTraceID))}, nil
}

func (t *TestServer) GetCtxIdStream(_ *emptypb.Empty, stream internal.TestService_GetCtxIdStreamServer) error {
	traceId := fmt.Sprintf("%v", stream.Context().Value(mw.CtxKeyTraceID))

	for i := 0; i < 3; i++ {
		err := stream.Send(&internal.ContextId{Id: fmt.Sprintf("%v", traceId)})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestGrpcMW(t *testing.T) {
	client := initGrpc(t)

	t.Run("GetCtxId", func(t *testing.T) {
		t.Run("If traceId is set", func(t *testing.T) {
			traceId := "test-trace-id"
			ctx := context.WithValue(context.Background(), mw.CtxKeyTraceID, traceId)
			res, err := client.GetCtxId(ctx, &emptypb.Empty{})
			require.NoError(t, err)
			require.Equal(t, traceId, res.Id)
		})

		t.Run("If traceId is not set, it should be set by SetTraceIdUnaryMW", func(t *testing.T) {
			res, err := client.GetCtxId(context.Background(), &emptypb.Empty{})
			require.NoError(t, err)
			require.NotEmpty(t, res.Id)
		})
	})

	t.Run("If traceId is set", func(t *testing.T) {
		t.Run("If traceId is set", func(t *testing.T) {
			traceId := "test-trace-id"
			ctx := context.WithValue(context.Background(), mw.CtxKeyTraceID, traceId)
			clientStream, err := client.GetCtxIdStream(ctx, &emptypb.Empty{})
			require.NoError(t, err)

			for {
				res, err := clientStream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					require.NoError(t, err)
				}

				require.Equal(t, traceId, res.Id)
			}
		})

		t.Run("If traceId is not set, it should be set by SetTraceIdStreamMW", func(t *testing.T) {
			clientStream, err := client.GetCtxIdStream(context.Background(), &emptypb.Empty{})
			require.NoError(t, err)

			for {
				res, err := clientStream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					require.NoError(t, err)
				}

				require.NotEmpty(t, res.Id)
			}
		})
	})
}

func initGrpc(t *testing.T) internal.TestServiceClient {

	listener, err := net.Listen("tcp", ":32249")
	require.NoError(t, err)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(GetTraceIdUnaryMW()),
		grpc.StreamInterceptor(GetTraceIdStreamMW()),
	)

	internal.RegisterTestServiceServer(grpcServer, &TestServer{})

	errChan := make(chan error)
	go func() {
		err := grpcServer.Serve(listener)
		errChan <- err
	}()
	time.Sleep(1 * time.Second)
	select {
	case err := <-errChan:
		require.FailNow(t, err.Error())
	default:
	}
	conn, err := grpc.NewClient("localhost:32249", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(SetTraceIdUnaryMW()),
		grpc.WithChainStreamInterceptor(SetTraceIdStreamMW()),
	)
	require.NoError(t, err)
	return internal.NewTestServiceClient(conn)
}
