package apiactor_test

import (
	"testing"
	"time"

	"github.com/jae2274/goutils/apiactor"
	"github.com/stretchr/testify/require"
)

type ProcessedSignal struct{}
type QuitSignal struct{}

// 각 Call 메소드의 호출은 이전 Call 메소드 호출이 끝난 후 최소 1초 이후에 시작되어야 한다.
func TestApiActorDelay(t *testing.T) {
	var delay int64 = 1000
	apiActor := apiactor.NewApiActor(delay)
	apiactor.Run(apiActor, make(<-chan QuitSignal))

	start := time.Now().UnixMilli()

	_, err := apiActor.Call(apiactor.NewRequest("GET", "https://google.com")) // 첫 호출은 바로 시작
	require.NoError(t, err)
	_, err = apiActor.Call(apiactor.NewRequest("GET", "https://google.com")) // 이후의 호출은 이전 호출과의 간격이 1초 이하일 경우, 남은 시간만큼 대기하여 1초 이상의 간격이 되도록 한다.
	require.NoError(t, err)
	_, err = apiActor.Call(apiactor.NewRequest("GET", "https://google.com"))
	require.NoError(t, err)
	_, err = apiActor.Call(apiactor.NewRequest("GET", "https://google.com"))
	require.NoError(t, err)
	_, err = apiActor.Call(apiactor.NewRequest("GET", "https://google.com"))
	require.NoError(t, err)

	end := time.Now().UnixMilli()

	//경과시간은 4초 이상이어야 한다.
	//이전 네트워크의 응답 지연시간이 1초 이하일 경우, 설정한 delay보다 남은 시간만큼 대기하므로 각 호출의 지연시간이 크지 않을 경우 약 4~5초로 예상된다.
	//예) 호출 경과 0.2초, delay 1초 -> 바로 이후 호출은 0.8초 대기
	require.Greater(t, end-start, delay*4)
	require.Less(t, end-start, delay*6)
}
