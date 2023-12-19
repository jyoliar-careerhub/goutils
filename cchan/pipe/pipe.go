package pipe

import "github.com/jae2274/goutils/cchan"

func Pipe[INPUT any, OUTPUT any, QUIT any](inputChan <-chan INPUT, errChan chan<- error, quitChan <-chan QUIT, bufferSize *int, action func(INPUT) (OUTPUT, error)) <-chan OUTPUT {
	bfs := 0
	if bufferSize != nil {
		bfs = *bufferSize
	}

	outputChan := make(chan OUTPUT, bfs)

	go func() {
		defer close(outputChan)
		for {
			received, ok := cchan.ReceiveOrQuit(inputChan, quitChan)
			if !ok {
				return
			}

			output, err := action(*received)

			ok = cchan.SendResult(output, err, outputChan, errChan, quitChan)
			if !ok {
				return
			}
		}
	}()

	return outputChan
}
