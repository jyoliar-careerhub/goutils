package pipe

import "github.com/jae2274/goutils/cchan"

func Transform[INPUT any, OUTPUT any, ERROR error, QUIT any](inputChan <-chan INPUT, errChan chan<- ERROR, quitChan <-chan QUIT, bufferSize *int, action func(INPUT) (OUTPUT, ERROR)) <-chan OUTPUT {
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

type Step[INPUT, OUTPUT any, ERROR error] struct {
	BufferSize *int
	Action     func(INPUT) (OUTPUT, ERROR)
}

func NewStep[INPUT, OUTPUT any, ERROR error](bufferSize *int, action func(INPUT) (OUTPUT, ERROR)) Step[INPUT, OUTPUT, ERROR] {
	return Step[INPUT, OUTPUT, ERROR]{bufferSize, action}
}

func Pipeline2[INPUT any, M1 any, OUTPUT any, ERROR error, QUIT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, OUTPUT, ERROR],
) (<-chan OUTPUT, chan ERROR) {
	errChan := make(chan ERROR, errBufferSize)
	pipeChan := Transform(inputChan, errChan, quitChan, step1.BufferSize, step1.Action)

	return Transform(pipeChan, errChan, quitChan, step2.BufferSize, step2.Action), errChan
}

func Pipeline3[INPUT any, M1 any, M2 any, OUTPUT any, ERROR error, QUIT any](
	inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, OUTPUT, ERROR],
) (<-chan OUTPUT, chan ERROR) {
	pipeChan, errChan := Pipeline2(inputChan, quitChan, errBufferSize, step1, step2)
	return Transform(pipeChan, errChan, quitChan, step3.BufferSize, step3.Action), errChan
}

func Pipeline4[INPUT any, M1 any, M2 any, M3 any, OUTPUT any, ERROR error, QUIT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, OUTPUT, ERROR],
) (<-chan OUTPUT, chan ERROR) {
	pipeChan, errChan := Pipeline3(inputChan, quitChan, errBufferSize, step1, step2, step3)
	return Transform(pipeChan, errChan, quitChan, step4.BufferSize, step4.Action), errChan
}
