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

// Pipeline은 여러개의 Step을 연속적으로 연결하여 하나의 채널로 연결한다.
// 각각의 step은 quitChan의 종료 트리거가 별도로 전파되고 종료되므로, 아직 종료되지 않은 step이 존재할 수 있다.
// 모든 step이 종료되기를 기다리지 않으므로, 비정상 종료시에만 quitChan을 트리거하도록 하며, 정상 종료를 의도하는 경우 inputChan을 닫아야 한다.
// Action 내부에서 Blocking되어 있는 동안은 inputChan과 quitChan의 종료 트리거가 전파되지 않는다.
func Pipeline2[INPUT any, M1 any, OUTPUT any, ERROR error, QUIT any](inputChan <-chan INPUT, errChan chan<- ERROR, quitChan <-chan QUIT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, OUTPUT, ERROR],
) <-chan OUTPUT {
	pipeChan := Transform(inputChan, errChan, quitChan, step1.BufferSize, step1.Action)

	return Transform(pipeChan, errChan, quitChan, step2.BufferSize, step2.Action)
}

func Pipeline3[INPUT any, M1 any, M2 any, OUTPUT any, ERROR error, QUIT any](
	inputChan <-chan INPUT, errChan chan<- ERROR, quitChan <-chan QUIT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, OUTPUT, ERROR],
) <-chan OUTPUT {
	pipeChan := Pipeline2(inputChan, errChan, quitChan, step1, step2)
	return Transform(pipeChan, errChan, quitChan, step3.BufferSize, step3.Action)
}

func Pipeline4[INPUT any, M1 any, M2 any, M3 any, OUTPUT any, ERROR error, QUIT any](inputChan <-chan INPUT, errChan chan<- ERROR, quitChan <-chan QUIT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, OUTPUT, ERROR],
) <-chan OUTPUT {
	pipeChan := Pipeline3(inputChan, errChan, quitChan, step1, step2, step3)
	return Transform(pipeChan, errChan, quitChan, step4.BufferSize, step4.Action)
}
