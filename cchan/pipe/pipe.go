package pipe

import (
	"context"

	"github.com/jae2274/goutils/cchan"
	"github.com/jae2274/goutils/cchan/async"
	"github.com/jae2274/goutils/ppool"
	"github.com/jae2274/goutils/ptr"
)

func Transform[INPUT any, OUTPUT any, ERROR error](ctx context.Context, inputChan <-chan INPUT, bufferSize *int, action func(INPUT) (OUTPUT, ERROR)) (<-chan OUTPUT, <-chan ERROR) {
	bfs := 0
	if bufferSize != nil {
		bfs = *bufferSize
	}

	outputChan := make(chan OUTPUT, bfs)
	errChan := make(chan ERROR, bfs)

	go func() {
		defer close(errChan)
		defer close(outputChan)

		for {
			received, ok := cchan.Receive(ctx, inputChan)
			if !ok {
				return
			}

			output, err := action(*received)

			ok = cchan.SendResult(ctx, output, err, outputChan, errChan)
			if !ok {
				return
			}
		}
	}()

	return outputChan, errChan
}

type Step[INPUT, OUTPUT any, ERROR error] struct {
	BufferSize *int
	Action     func(INPUT) (OUTPUT, ERROR)
}

func NewStep[INPUT, OUTPUT any, ERROR error](bufferSize *int, action func(INPUT) (OUTPUT, ERROR)) Step[INPUT, OUTPUT, ERROR] {
	return Step[INPUT, OUTPUT, ERROR]{bufferSize, action}
}

func NewAsyncAwaitSteps[INPUT, OUTPUT any](
	ctx context.Context,
	bufferSize *int,
	concurrencySize int,
	action func(context.Context, INPUT) (OUTPUT, error),
) (Step[INPUT, <-chan async.Result[OUTPUT], error], Step[<-chan async.Result[OUTPUT], OUTPUT, error]) {
	p := ppool.NewPool(concurrencySize, func() (chan async.Result[OUTPUT], error) {
		return make(chan async.Result[OUTPUT], 1), nil
	})

	asyncStep := NewStep(ptr.P(concurrencySize), func(input INPUT) (<-chan async.Result[OUTPUT], error) {
		ch, err := p.Acquire(ctx)
		if err != nil {
			return nil, err
		}

		go func() {
			defer p.Release(ch)
			value, err := action(ctx, input)
			ch <- async.Result[OUTPUT]{Value: value, Err: err}
		}()

		return ch, nil
	})

	awaitStep := NewStep(bufferSize, func(asyncResult <-chan async.Result[OUTPUT]) (OUTPUT, error) {
		result := <-asyncResult

		return result.Value, result.Err
	})

	return asyncStep, awaitStep
}

// type token struct{}

// func NewUnorderedAsyncStep[INPUT any, OUTPUT any](ctx context.Context, bufferSize *int, concurrencySize int, action func(INPUT) (OUTPUT, error)) Step[INPUT, OUTPUT, error] {
// 	p := ppool.NewPool(concurrencySize, func() (token, error) {
// 		return token{}, nil
// 	})

// 	resultChan := make(chan async.Result[OUTPUT], concurrencySize)

// 	return NewStep(bufferSize, func(input INPUT) (OUTPUT, error) {
// 		tk, err := p.Acquire(ctx)
// 		if err != nil {
// 			return *new(OUTPUT), err
// 		}

// 		go func(tk token, input INPUT) {
// 			defer p.Release(tk)

// 			v, err := action(input)
// 			resultChan <- async.Result[OUTPUT]{Value: v, Err: err}
// 		}(tk, input)

// 		result := <-resultChan
// 		return result.Value, result.Err
// 	})
// }

// Pipeline은 여러개의 Step을 연속적으로 연결하여 하나의 채널로 연결한다.
// 각각의 step은 context의 종료 트리거가 별도로 전파되고 종료되므로, 아직 종료되지 않은 step이 존재할 수 있다.
// 모든 step이 종료되기를 기다리지 않으므로, 비정상 종료시에만 context를 종료 트리거하도록 하며, 정상 종료를 의도하는 경우 inputChan을 닫아야 한다.
// Action 내부에서 Blocking되어 있는 동안은 inputChan과 context의 종료 트리거가 전파되지 않는다.
func Pipeline2[INPUT any, M1 any, OUTPUT any, ERROR error](ctx context.Context, inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err)

	return step2Pipe, errChan
}

func Pipeline3[INPUT any, M1 any, M2 any, OUTPUT any, ERROR error](
	ctx context.Context,
	inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)
	step3Pipe, step3Err := Transform(ctx, step2Pipe, step3.BufferSize, step3.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err, step3Err)

	return step3Pipe, errChan
}

func Pipeline4[INPUT any, M1 any, M2 any, M3 any, OUTPUT any, ERROR error](
	ctx context.Context,
	inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)
	step3Pipe, step3Err := Transform(ctx, step2Pipe, step3.BufferSize, step3.Action)
	step4Pipe, step4Err := Transform(ctx, step3Pipe, step4.BufferSize, step4.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err, step3Err, step4Err)

	return step4Pipe, errChan
}

func Pipeline5[INPUT any, M1 any, M2 any, M3 any, M4 any, OUTPUT any, ERROR error](
	ctx context.Context,
	inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, M4, ERROR],
	step5 Step[M4, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)
	step3Pipe, step3Err := Transform(ctx, step2Pipe, step3.BufferSize, step3.Action)
	step4Pipe, step4Err := Transform(ctx, step3Pipe, step4.BufferSize, step4.Action)
	step5Pipe, step5Err := Transform(ctx, step4Pipe, step5.BufferSize, step5.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err, step3Err, step4Err, step5Err)

	return step5Pipe, errChan
}

func Pipeline6[INPUT any, M1 any, M2 any, M3 any, M4 any, M5 any, OUTPUT any, ERROR error](
	ctx context.Context,
	inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, M4, ERROR],
	step5 Step[M4, M5, ERROR],
	step6 Step[M5, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)
	step3Pipe, step3Err := Transform(ctx, step2Pipe, step3.BufferSize, step3.Action)
	step4Pipe, step4Err := Transform(ctx, step3Pipe, step4.BufferSize, step4.Action)
	step5Pipe, step5Err := Transform(ctx, step4Pipe, step5.BufferSize, step5.Action)
	step6Pipe, step6Err := Transform(ctx, step5Pipe, step6.BufferSize, step6.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err, step3Err, step4Err, step5Err, step6Err)

	return step6Pipe, errChan
}

func Pipeline7[INPUT any, M1 any, M2 any, M3 any, M4 any, M5 any, M6 any, OUTPUT any, ERROR error](
	ctx context.Context,
	inputChan <-chan INPUT,
	step1 Step[INPUT, M1, ERROR],
	step2 Step[M1, M2, ERROR],
	step3 Step[M2, M3, ERROR],
	step4 Step[M3, M4, ERROR],
	step5 Step[M4, M5, ERROR],
	step6 Step[M5, M6, ERROR],
	step7 Step[M6, OUTPUT, ERROR],
) (<-chan OUTPUT, <-chan ERROR) {
	step1Pipe, step1Err := Transform(ctx, inputChan, step1.BufferSize, step1.Action)
	step2Pipe, step2Err := Transform(ctx, step1Pipe, step2.BufferSize, step2.Action)
	step3Pipe, step3Err := Transform(ctx, step2Pipe, step3.BufferSize, step3.Action)
	step4Pipe, step4Err := Transform(ctx, step3Pipe, step4.BufferSize, step4.Action)
	step5Pipe, step5Err := Transform(ctx, step4Pipe, step5.BufferSize, step5.Action)
	step6Pipe, step6Err := Transform(ctx, step5Pipe, step6.BufferSize, step6.Action)
	step7Pipe, step7Err := Transform(ctx, step6Pipe, step7.BufferSize, step7.Action)

	errChan := cchan.Merge(ctx, step1Err, step2Err, step3Err, step4Err, step5Err,
		step6Err,
		step7Err)

	return step7Pipe,
		errChan
}

func PassThrough[TARGET any](ctx context.Context, inputChan <-chan TARGET, action func(TARGET)) <-chan TARGET {
	outputChan := make(chan TARGET)

	go func() {
		defer close(outputChan)
		for {
			received, ok := cchan.Receive(ctx, inputChan)
			if !ok {
				return
			}

			action(*received)

			ok = cchan.Send(ctx, outputChan, *received)
			if !ok {
				return
			}
		}
	}()

	return outputChan
}
