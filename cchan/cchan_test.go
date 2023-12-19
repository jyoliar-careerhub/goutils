package cchan_test

import (
	"errors"
	"goutils/cchan"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSendResult(t *testing.T) {
	t.Run("quitChan이 트리거되면 false를 반환하고 result와 error는 전달되지 않느다.", func(t *testing.T) {
		resultChan, errChan, quitChan := initChans[QuitSignal]()
		close(quitChan)

		ok := cchan.SendResult(SampleResult{}, nil, resultChan, errChan, quitChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResult(SampleResult{}, errSample, resultChan, errChan, quitChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults([]SampleResult{{}, {}}, nil, resultChan, errChan, quitChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults([]SampleResult{{}, {}}, errSample, resultChan, errChan, quitChan)
		require.False(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 0)
	})

	t.Run("quitChan이 트리거되지 않고 error가 존재하면 error가 errChan으로 전달되고 result는 전달되지 않는다.", func(t *testing.T) {
		resultChan, errChan, quitChan := initChans[QuitSignal]()

		ok := cchan.SendResult(SampleResult{}, errSample, resultChan, errChan, quitChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)

		ok = cchan.SendResults([]SampleResult{{}, {}}, errSample, resultChan, errChan, quitChan)
		require.True(t, ok)
		assertLength(t, resultChan, 0)
		assertLength(t, errChan, 1)
	})

	t.Run("quitChan이 트리거되지 않고 error가 존재하지 않으면 result가 resultChan으로 전달된다.", func(t *testing.T) {
		resultChan, errChan, quitChan := initChans[QuitSignal]()

		ok := cchan.SendResult(SampleResult{}, nil, resultChan, errChan, quitChan)
		require.True(t, ok)
		assertLength(t, resultChan, 1)
		assertLength(t, errChan, 0)

		ok = cchan.SendResults([]SampleResult{{}, {}}, nil, resultChan, errChan, quitChan)
		require.True(t, ok)
		assertLength(t, resultChan, 2)
		assertLength(t, errChan, 0)
	})
}

func TestReceiveOrQuit(t *testing.T) {
	t.Run("quitChan이 트리거되면 false를 반환하고 data는 nil을 리턴한다.", func(t *testing.T) {
		receiveChan, _, quitChan := initChans[QuitSignal]()
		resultC := make(chan ReceivedResult, 1)

		go func() {
			data, ok := cchan.ReceiveOrQuit(receiveChan, quitChan) // receiveChan, quitChan이 트리거되지 않는다면 무한 대기
			resultC <- ReceivedResult{data, ok}
		}()
		assertLength(t, resultC, 0)

		close(quitChan)
		result := assertLength(t, resultC, 1)[0]
		require.False(t, result.ok)
		require.Nil(t, result.data)

		receiveChan <- SampleResult{}                          // 이 데이터는 quitChan 트리거의 우선순위로 인해 무시된다.
		data, ok := cchan.ReceiveOrQuit(receiveChan, quitChan) // close(quitChan) 이후에는 무조건 data는 nil, ok는 false
		require.False(t, ok)
		require.Nil(t, data)
	})

	t.Run("quitChan이 트리거되지 않고 data가 존재하면 data가 전달되고 true를 반환한다.", func(t *testing.T) {
		receiveChan, _, quitChan := initChans[QuitSignal]()
		resultC := make(chan ReceivedResult, 1)

		go func() {
			data, ok := cchan.ReceiveOrQuit(receiveChan, quitChan) // receiveChan, quitChan이 트리거되지 않는다면 무한 대기
			resultC <- ReceivedResult{data, ok}
		}()
		assertLength(t, resultC, 0)

		receiveChan <- SampleResult{}
		result := assertLength(t, resultC, 1)[0]
		require.True(t, result.ok)
		require.NotNil(t, result.data)

		receiveChan <- SampleResult{}
		data, ok := cchan.ReceiveOrQuit(receiveChan, quitChan) // receiveChan에 데이터가 존재하면 호출은 바로 리턴된다.
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

func initChans[QUIT any]() (chan SampleResult, chan error, chan QUIT) {
	resultChan := make(chan SampleResult, 100)
	errChan := make(chan error, 100)
	quitChan := make(chan QUIT)
	return resultChan, errChan, quitChan
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
	return time.After(time.Millisecond * time.Duration(50)) // quitChan에 데이터가 입력될 때까지의 최소 대기 시간
}

type QuitSignal struct{}
type ProcessedSignal struct{}
