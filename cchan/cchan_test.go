package cchan_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jae2274/goutils/cchan"
	"github.com/stretchr/testify/require"
)

func TestSendResult(t *testing.T) {
	t.Run("context가 종료되면 false를 반환하고 result와 error는 전달되지 않느다.", func(t *testing.T) {
		resultChan, errChan := initChans()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ok := cchan.SendResult(ctx, SampleResult{}, nil, resultChan, errChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResult(ctx, SampleResult{}, errSample, resultChan, errChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, nil, resultChan, errChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, errSample, resultChan, errChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)
	})

	t.Run("context가 종료되지 않고 error가 존재하면 error가 errChan으로 전달되고 result는 전달되지 않는다.", func(t *testing.T) {
		resultChan, errChan := initChans()
		ctx := context.Background()

		ok := cchan.SendResult(ctx, SampleResult{}, errSample, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, errSample, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)
	})
	t.Run("특정 타입 에러 전달시에도 위와 같은 동작을 한다.", func(t *testing.T) {
		resultChan, _ := initChans()
		errChan := make(chan *SampleError, 100)
		ctx := context.Background()

		ok := cchan.SendResult(ctx, SampleResult{}, &SampleError{}, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, &SampleError{}, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)
	})

	t.Run("context가 종료되지 않고 error가 존재하지 않으면 result가 resultChan으로 전달된다.", func(t *testing.T) {
		resultChan, errChan := initChans()
		ctx := context.Background()

		ok := cchan.SendResult(ctx, SampleResult{}, nil, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 1)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, nil, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 2)
		assertLength(t, errChan, 0)
	})

	t.Run("특정 타입의 error가 nil일 경우에도 위와 같은 동작을 한다.", func(t *testing.T) {
		resultChan := make(chan SampleResult, 100)
		errChan := make(chan *SampleError, 100)
		ctx := context.Background()

		var err *SampleError = nil
		ok := cchan.SendResult(ctx, SampleResult{}, err, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 1)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults(ctx, []SampleResult{{}, {}}, nil, resultChan, errChan)
		require.True(t, ok)
		assertLength(t, resultChan, 2)
		assertLength(t, errChan, 0)
	})
}

type SampleError struct{}

func (e SampleError) Error() string { return "Sample Error caused!!!" }

func TestReceiveOrQuit(t *testing.T) {
	t.Run("context가 종료되면 false를 반환하고 data는 nil을 리턴한다.", func(t *testing.T) {
		receiveChan, _ := initChans()
		ctx, cancel := context.WithCancel(context.Background())

		resultC := make(chan ReceivedResult, 1)

		go func() {
			data, ok := cchan.Receive(ctx, receiveChan)
			resultC <- ReceivedResult{data, ok}
		}()
		assertLength(t, resultC, 0)

		cancel() // context 종료
		result := assertLength(t, resultC, 1)[0]
		require.False(t, result.ok)
		require.Nil(t, result.data)

		receiveChan <- SampleResult{}
		data, ok := cchan.Receive(ctx, receiveChan) // context 종료 이후에는 무조건 false, nil을 리턴한다.
		require.False(t, ok)
		require.Nil(t, data)
	})

	t.Run("context가 종료되지 않고 data가 존재하면 data가 전달되고 true를 반환한다.", func(t *testing.T) {
		receiveChan, _ := initChans()
		ctx := context.Background()

		resultC := make(chan ReceivedResult, 1)

		go func() {
			data, ok := cchan.Receive(ctx, receiveChan)
			resultC <- ReceivedResult{data, ok}
		}()
		assertLength(t, resultC, 0)

		receiveChan <- SampleResult{}
		result := assertLength(t, resultC, 1)[0]
		require.True(t, result.ok)
		require.NotNil(t, result.data)

		receiveChan <- SampleResult{}
		data, ok := cchan.Receive(ctx, receiveChan) // receiveChan에 데이터가 존재하면 호출은 바로 리턴된다.
		require.True(t, ok)
		require.NotNil(t, data)
	})
}

type SampleResult struct{}

type ReceivedResult struct {
	data *SampleResult
	ok   bool
}

var errSample = errors.New("sample error")

func initChans() (chan SampleResult, chan error) {
	resultChan := make(chan SampleResult, 100)
	errChan := make(chan error, 100)

	return resultChan, errChan
}

func assertLength[T any](t *testing.T, channel <-chan T, expected int) []T {
	receiveds := getFromChan(channel)
	require.Equal(t, expected, len(receiveds))
	return receiveds
}

func getFromChan[T any](channel <-chan T) []T {
	var result []T

	for {
		select {
		case data := <-channel:
			result = append(result, data)
		case <-moment():
			return result
		}
	}
}

func moment() <-chan time.Time {
	return time.After(time.Millisecond * time.Duration(50)) //context 종료 전파를 위한 대기시간
}

type QuitSignal struct{}
type ProcessedSignal struct{}
