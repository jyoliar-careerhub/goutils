package pipe

import "github.com/jae2274/goutils/cchan"

func Transform[INPUT any, OUTPUT any, QUIT any](inputChan <-chan INPUT, errChan chan<- error, quitChan <-chan QUIT, bufferSize *int, action func(INPUT) (OUTPUT, error)) <-chan OUTPUT {
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

type Step[INPUT, OUTPUT any] struct {
	BufferSize *int
	Action     func(INPUT) (OUTPUT, error)
}

func Pipeline2[INPUT any, M1 any, OUTPUT any, QUIT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1],
	step2 Step[M1, OUTPUT],
) (<-chan OUTPUT, chan error) {
	errChan := make(chan error, errBufferSize)
	pipeChan := Transform(inputChan, errChan, quitChan, step1.BufferSize, step1.Action)

	return Transform(pipeChan, errChan, quitChan, step2.BufferSize, step2.Action), errChan
}

func Pipeline3[QUIT any, INPUT any, M1 any, M2 any, OUTPUT any](
	inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1],
	step2 Step[M1, M2],
	step3 Step[M2, OUTPUT],
) (<-chan OUTPUT, chan error) {
	pipeChan, errChan := Pipeline2(inputChan, quitChan, errBufferSize, step1, step2)
	return Transform(pipeChan, errChan, quitChan, step3.BufferSize, step3.Action), errChan
}

func Pipeline4[QUIT any, INPUT any, M1 any, M2 any, M3 any, OUTPUT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1],
	step2 Step[M1, M2],
	step3 Step[M2, M3],
	step4 Step[M3, OUTPUT],
) (<-chan OUTPUT, chan error) {
	pipeChan, errChan := Pipeline3(inputChan, quitChan, errBufferSize, step1, step2, step3)
	return Transform(pipeChan, errChan, quitChan, step4.BufferSize, step4.Action), errChan
}
