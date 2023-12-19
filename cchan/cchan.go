package cchan

import (
	"log"
	"time"
)

func SendResult[T any, QUIT any](result T, err error, resultChan chan<- T, errChan chan<- error, quitChan <-chan QUIT) bool {
	if err != nil {
		ok := SendOrQuit(err, errChan, quitChan)
		return ok
	} else {
		ok := SendOrQuit(result, resultChan, quitChan)
		return ok
	}
}

func SendResults[T any, QUIT any](results []T, err error, resultChan chan<- T, errChan chan<- error, quitChan <-chan QUIT) bool {
	if err != nil {
		ok := SendOrQuit(err, errChan, quitChan)
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

func TooMuchError[QUIT any](periodErrCount uint, limitErrPeriod time.Duration, errChan <-chan error, quitChan chan QUIT) {
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
