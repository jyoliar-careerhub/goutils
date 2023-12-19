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

type Chain[INPUT, OUTPUT any] struct {
	BufferSize *int
	Action     func(INPUT) (OUTPUT, error)
}

func Chain2[INPUT any, M1 any, OUTPUT any, QUIT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	chain1 Chain[INPUT, M1],
	chain2 Chain[M1, OUTPUT],
) (<-chan OUTPUT, chan error) {
	errChan := make(chan error, errBufferSize)
	pipeChan := Pipe(inputChan, errChan, quitChan, chain1.BufferSize, chain1.Action)

	return Pipe(pipeChan, errChan, quitChan, chain2.BufferSize, chain2.Action), errChan
}

func Chain3[QUIT any, INPUT any, M1 any, M2 any, OUTPUT any](
	inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	chain1 Chain[INPUT, M1],
	chain2 Chain[M1, M2],
	chain3 Chain[M2, OUTPUT],
) (<-chan OUTPUT, chan error) {
	pipeChan, errChan := Chain2(inputChan, quitChan, errBufferSize, chain1, chain2)
	return Pipe(pipeChan, errChan, quitChan, chain3.BufferSize, chain3.Action), errChan
}

func Chain4[QUIT any, INPUT any, M1 any, M2 any, M3 any, OUTPUT any](inputChan <-chan INPUT, quitChan <-chan QUIT,
	errBufferSize int,
	chain1 Chain[INPUT, M1],
	chain2 Chain[M1, M2],
	chain3 Chain[M2, M3],
	chain4 Chain[M3, OUTPUT],
) (<-chan OUTPUT, chan error) {
	pipeChan, errChan := Chain3(inputChan, quitChan, errBufferSize, chain1, chain2, chain3)
	return Pipe(pipeChan, errChan, quitChan, chain4.BufferSize, chain4.Action), errChan
}
