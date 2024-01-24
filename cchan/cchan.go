package cchan

import (
	"context"
	"time"

	"github.com/jae2274/goutils/ptr"
)

func SendResult[T any, ERROR error](ctx context.Context, result T, err ERROR, resultChan chan<- T, errChan chan<- ERROR) bool {
	if !ptr.IsNil(err) {
		ok := Send(ctx, errChan, err)
		return ok
	} else {
		ok := Send(ctx, resultChan, result)
		return ok
	}
}

func SendResults[T any, ERROR error](ctx context.Context, results []T, err ERROR, resultChan chan<- T, errChan chan<- ERROR) bool {
	if !ptr.IsNil(err) {
		ok := Send(ctx, errChan, err)
		return ok
	} else {
		for _, result := range results {
			ok := Send(ctx, resultChan, result)
			if !ok {
				return false
			}
		}
		return true
	}
}

func Send[T any](ctx context.Context, sendChan chan<- T, data T) bool {
	select {
	case <-ctx.Done(): // context의 종료 트리거를 우선순위로 둔다.
		return false
	default:
		select {
		case sendChan <- data:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

func Receive[T any](ctx context.Context, receiveChan <-chan T) (*T, bool) {
	select {
	case <-ctx.Done(): // context의 종료 트리거를 우선순위로 둔다.
		return nil, false
	default:
		select {
		case data, ok := <-receiveChan:
			return &data, ok
		case <-ctx.Done():
			return nil, false
		}
	}
}

// 채널 내부의 데이터를 소진하게 될 수 있으므로 주의한다.
func WaitClosed[T any](ch <-chan T) []T {
	items := make([]T, 0)
	for {
		item, ok := <-ch
		if !ok {
			return items
		}
		items = append(items, item)
	}
}

// 채널 내부의 데이터를 소진하게 될 수 있으므로 주의한다.
func IsClosed[T any](ch <-chan T) (bool, *T) {
	select {
	case item, ok := <-ch:
		if ok {
			return false, &item
		}
		return true, nil
	default:
		return false, nil
	}
}

func SafeClose[T any](ch chan T) {
	select {
	case _, ok := <-ch:
		if ok {
			close(ch)
		}
	default:
		close(ch)
	}
}

func TooMuchError[ERROR error](periodErrCount uint, limitErrPeriod time.Duration, errChan <-chan ERROR, tooMuchErrFunc func(), closedFunc func()) {
	go func() {
		var errCount uint = 0
		errCaughtTimes := make([]time.Time, periodErrCount)

		for {
			_, ok := <-errChan
			if !ok {
				closedFunc()
				return
			}

			errCount++
			errCaughtTimes = append(errCaughtTimes, time.Now())

			if len(errCaughtTimes) >= int(periodErrCount) {
				lastErrCaughtTime := errCaughtTimes[len(errCaughtTimes)-1]
				recentErrCaughtTime := errCaughtTimes[len(errCaughtTimes)-10]
				errCaughtPeriod := lastErrCaughtTime.Sub(recentErrCaughtTime)

				if errCaughtPeriod.Abs() < limitErrPeriod.Abs() {
					tooMuchErrFunc()
				}
				errCaughtTimes = errCaughtTimes[1:]
			}
		}
	}()
}

func Timeout[DATA any](initDuration, duration time.Duration, processedChan <-chan DATA, timeoutFunc func(), closedFunc func()) {
	go func() {
		waitDuration := initDuration

		for {

			select {
			case <-time.After(waitDuration):
				timeoutFunc()
			case _, ok := <-processedChan:
				if !ok {
					closedFunc()
					return
				}
				waitDuration = duration
			}
		}
	}()
}
