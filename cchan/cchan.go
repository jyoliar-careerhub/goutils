package cchan

import (
	"log"
	"time"
)

func SendResult[T any, ERROR error, QUIT any](result T, err *ERROR, resultChan chan<- T, errChan chan<- ERROR, quitChan <-chan QUIT) bool {
	if err != nil {
		ok := SendOrQuit(*err, errChan, quitChan)
		return ok
	} else {
		ok := SendOrQuit(result, resultChan, quitChan)
		return ok
	}
}

func SendResults[T any, ERROR error, QUIT any](results []T, err *ERROR, resultChan chan<- T, errChan chan<- ERROR, quitChan <-chan QUIT) bool {
	if err != nil {
		ok := SendOrQuit(*err, errChan, quitChan)
		return ok
	} else {
		for _, result := range results {
			ok := SendOrQuit(result, resultChan, quitChan)
			if !ok {
				return false
			}
		}
		return true
	}
}

func SendOrQuit[T any, QUIT any](data T, sendChan chan<- T, quit <-chan QUIT) bool {
	select {
	case <-quit: // quitChan의 트리거를 우선순위로 둔다.
		return false
	default:
		select {
		case sendChan <- data:
			return true
		case <-quit:
			return false
		}
	}
}

func ReceiveOrQuit[T any, QUIT any](receiveChan <-chan T, quit <-chan QUIT) (*T, bool) {
	select {
	case <-quit: // quitChan의 트리거를 우선순위로 둔다.
		return nil, false
	default:
		select {
		case data, ok := <-receiveChan:
			return &data, ok
		case <-quit:
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

func TooMuchError[ERROR error, QUIT any](periodErrCount uint, limitErrPeriod time.Duration, errChan <-chan ERROR, quitChan chan QUIT) {
	defer func() {
		log.Default().Println("TooMuchError closed")
	}()

	var errCount uint = 0
	errCaughtTimes := make([]time.Time, periodErrCount)

	for {
		err, ok := ReceiveOrQuit(errChan, quitChan)
		if !ok {
			return
		}

		log.Default().Println((*err).Error())
		errCount++
		errCaughtTimes = append(errCaughtTimes, time.Now())
		if len(errCaughtTimes) >= int(periodErrCount) {
			lastErrCaughtTime := errCaughtTimes[len(errCaughtTimes)-1]
			recentErrCaughtTime := errCaughtTimes[len(errCaughtTimes)-10]
			errCaughtPeriod := lastErrCaughtTime.Sub(recentErrCaughtTime)

			if errCaughtPeriod.Abs() < limitErrPeriod.Abs() {
				close(quitChan)
				return
			}
			errCaughtTimes = errCaughtTimes[1:]
		}

	}
}

func Timeout[DATA any, QUIT any](initDuration, duration time.Duration, processedChan <-chan DATA, quitChan chan QUIT) {
	defer func() {
		log.Default().Println("Timeout closed")
	}()
	waitDuration := initDuration

	for {

		select {
		case <-quitChan:
			return
		case <-time.After(waitDuration):
			close(quitChan)
			return
		case _, ok := <-processedChan:
			if !ok {
				return
			}
			waitDuration = duration
		}
	}
}
