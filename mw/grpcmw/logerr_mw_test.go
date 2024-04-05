package grpcmw

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jae2274/goutils/mw/grpcmw/internal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LogErrGrpcServer struct {
	internal.UnimplementedLogErrServiceServer
}

var errTest = fmt.Errorf("for test")

func (l *LogErrGrpcServer) Error(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errTest
}

func (l *LogErrGrpcServer) ErrorStream(_ *emptypb.Empty, stream internal.LogErrService_ErrorStreamServer) error {
	for i := 0; i < 3; i++ {
		stream.Send(&emptypb.Empty{})
	}
	return errTest
}

func TestLogErrMw(t *testing.T) {
	logErrClient := initLogErrGrpc(t)

	t.Run("Unary", func(t *testing.T) {
		mockStdout, err := os.Create("test.log")
		require.NoError(t, err)
		defer os.Remove("test.log")

		tempStdout := os.Stdout

		os.Stdout = mockStdout
		_, expectedErr := logErrClient.Error(context.Background(), &emptypb.Empty{})
		os.Stdout = tempStdout

		require.Contains(t, expectedErr.Error(), errTest.Error())

		log, _ := os.ReadFile("test.log")
		require.NoError(t, err)
		require.Contains(t, string(log), errTest.Error())
	})

	t.Run("Stream", func(t *testing.T) {
		mockStdout, err := os.Create("test.log")
		require.NoError(t, err)
		defer os.Remove("test.log")

		tempStdout := os.Stdout

		os.Stdout = mockStdout
		stream, err := logErrClient.ErrorStream(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)

		var expectedErr error
		count := 0
		for {
			_, expectedErr = stream.Recv()
			if expectedErr != nil {
				break
			}
			count++
		}

		os.Stdout = tempStdout

		require.Contains(t, expectedErr.Error(), errTest.Error())
		require.Equal(t, 3, count)

		log, err := os.ReadFile("test.log")
		require.NoError(t, err)
		require.Contains(t, string(log), errTest.Error())
	})
}

func initLogErrGrpc(t *testing.T) internal.LogErrServiceClient {

	listener, err := net.Listen("tcp", ":32249")
	require.NoError(t, err)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(LogErrUnaryMW()),
		grpc.ChainStreamInterceptor(LogErrStreamMW()),
	)

	internal.RegisterLogErrServiceServer(grpcServer, &LogErrGrpcServer{})

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
	conn, err := grpc.NewClient("localhost:32249", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return internal.NewLogErrServiceClient(conn)
}
