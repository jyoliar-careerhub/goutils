package cchan

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
