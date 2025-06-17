package pipe_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jae2274/goutils/cchan"
	"github.com/jae2274/goutils/cchan/pipe"
	"github.com/jae2274/goutils/ptr"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Inputs          []DivideTarget
	ExpectedOutputs []int
	ExpectedErrs    []error
	TestAction      func(*testing.T, []DivideTarget, []int, []error)
}

func TestPipe(t *testing.T) {
	t.Run("Test cases", func(t *testing.T) {
		TestCases := []TestCase{
			{
				[]DivideTarget{
					{10, 3},
					{20, 5},
					{30, 7},
					{40, 9},
				},
				[]int{16, 16, 36, 64},
				[]error{},
				test,
			},
			{
				[]DivideTarget{
					{10, 3},  // 정상 동작하여 결과 전달
					{-20, 5}, // errNagativeNumber 발생
					{30, 7},  // 정상 동작하여 결과 전달
					{40, 0},  // errDivideByZero 발생
				},
				[]int{16, 36},
				[]error{
					&errNagativeNumber{-20, 5},
					&errDivideByZero{40, 0},
				},
				test,
			},
		}

		for _, testCase := range TestCases {
			testCase.TestAction(t, testCase.Inputs, testCase.ExpectedOutputs, testCase.ExpectedErrs)
		}
	})

	t.Run("No Buffer Test", func(t *testing.T) {
		inputChan := make(chan int)
		ctx := context.Background()

		outputChan, _ := pipe.Transform(ctx, inputChan, nil, square) //outputChan은 버퍼가 없으므로 해당 채널의 데이터를 수신하지 않으면 전송자가 블록된다.
		fmt.Println(outputChan)                                      //컴파일 에러를 방지하기 위해 사용한 코드

		inputChan <- 3
		select {
		case inputChan <- 3: //outputChan의 데이터를 수신하지 않아 Pipe 내부에서 전송자가 블록되었으므로, inputChan의 전송자도 블록된다.
			require.Fail(t, "sumTargetChan is not blocked")
		default:
			log.Default().Println("sumTargetChan is blocked")
		}

		require.Len(t, outputChan, 0)
	})

	t.Run("Buffer Test", func(t *testing.T) {
		inputChan := make(chan int)
		ctx := context.Background()

		outputChan, _ := pipe.Transform(ctx, inputChan, ptr.P(3), square) //outputChan의 버퍼는 3개이다.
		fmt.Println(outputChan)                                           //컴파일 에러를 방지하기 위해 사용한 코드

		for i := 0; i < 4; i++ { //inputChan은 버퍼가 존재하지 않음에 유의한다. 4개까지 전송이 가능한 이유는 outputChan에 3개까지 전송 이후부터 블록되기 때문이다.
			inputChan <- 3 //버퍼가 3개이므로, 3개의 데이터를 전송해도 block되지 않는다.
		}

		select {
		case inputChan <- 3: // 3개의 버퍼가 모두 차있으므로 block된다.
			require.Fail(t, "sumTargetChan is not blocked")
		default:
			log.Default().Println("sumTargetChan is blocked")
		}

		require.Len(t, outputChan, 3)
	})

	t.Run("Quit Test", func(t *testing.T) {
		inputChan := make(chan string)
		ctx, cancelFunc := context.WithCancel(context.Background())

		stepNamesChan := make(chan string, 3)

		makeStep := func(step string) pipe.Step[string, string, error] {
			return pipe.NewStep(nil,
				func(m string) (string, error) {
					time.Sleep(time.Second * 5)

					stepNamesChan <- step

					return m, nil
				})
		}
		resultChan, errChan := pipe.Pipeline3(ctx, inputChan, makeStep("step1"), makeStep("step2"), makeStep("step3"))

		inputChan <- "Hello!"

		cancelFunc()                       // 파이프라인 종료 트리거
		time.Sleep(time.Millisecond * 100) // context 종료 전파 대기
		isClosed, _ := cchan.IsClosed(resultChan)
		require.True(t, isClosed)
		isClosed, err := cchan.IsClosed(errChan)
		require.True(t, ptr.IsNil(err))
		require.True(t, isClosed)

		require.Len(t, stepNamesChan, 0)
		time.Sleep(time.Second * 6) // 모든 step이 종료되기를 기다림

		require.Len(t, stepNamesChan, 1)
		//Action 내부에서 Blocking되어 있었던 step1이 종료되지 않은 상태였음을 알 수 있다.
		require.Equal(t, "step1", <-stepNamesChan)
	})
}

func test(t *testing.T, inputs []DivideTarget, expectedOutputs []int, errs []error) {
	inputChan := make(chan DivideTarget)
	ctx := context.Background()

	step1 := pipe.NewStep(nil,
		func(target DivideTarget) (*DivideTarget, error) {
			a, b, err := checkPositive(target.denominator, target.numerator)
			if err != nil {
				return nil, err
			}
			return &DivideTarget{a, b}, nil
		})

	step2 := pipe.NewStep(nil,
		func(dt *DivideTarget) (*sumTarget, error) {
			return divide(dt.denominator, dt.numerator)
		})

	step3 := pipe.NewStep(nil, sum)
	step4 := pipe.NewStep(nil, square)

	resultChan, errChan := pipe.Pipeline4(ctx, inputChan, step1, step2, step3, step4)

	for _, input := range inputs {
		inputChan <- input
	}

	for _, expectedOutput := range expectedOutputs {
		output, ok := <-resultChan
		require.True(t, ok, "resultChan is closed")
		require.Equal(t, expectedOutput, output)
	}

	for _, err := range errs {
		output, ok := <-errChan
		require.True(t, ok, "errChan is closed")
		require.Equal(t, err, output)
	}

	require.Len(t, resultChan, 0)
	require.Len(t, errChan, 0)

	close(inputChan)
	time.Sleep(time.Millisecond * 100) // 파이프라인이 종료되기를 기다림
	isClosed, _ := cchan.IsClosed(resultChan)
	require.True(t, isClosed)
	isClosed, _ = cchan.IsClosed(errChan)
	require.True(t, isClosed)
	isClosed, _ = cchan.IsClosed(ctx.Done())
	require.False(t, isClosed) //resultChan의 종료가 context종료에 의해 트리거되지 않았음을 알 수 있다.
}

func TestPassThrough(t *testing.T) {

	t.Run("정상 동작", func(t *testing.T) {
		inputChan := make(chan int)
		ctx := context.Background()

		sideResults := make([]int, 0)
		outputChan := pipe.PassThrough(ctx, inputChan, func(number int) {
			sideResults = append(sideResults, number*10)
		})

		go func() {
			inputChan <- 1
			inputChan <- 2
			inputChan <- 3
			close(inputChan)
		}()

		require.Equal(t, 1, <-outputChan)
		require.Equal(t, 2, <-outputChan)
		require.Equal(t, 3, <-outputChan)

		time.Sleep(time.Millisecond * 100)
		isClosed, _ := cchan.IsClosed(outputChan)
		require.True(t, isClosed)
		isClosed, _ = cchan.IsClosed(ctx.Done())
		require.False(t, isClosed) //resultChan의 종료가 context 종료에 의해 트리거되지 않았음을 알 수 있다.

		require.Len(t, sideResults, 3)
		require.Equal(t, []int{10, 20, 30}, sideResults)
	})

	t.Run("context 종료 발생", func(t *testing.T) {
		inputChan := make(chan int)
		ctx, cancelFunc := context.WithCancel(context.Background())

		sideResults := make([]int, 0)
		outputChan := pipe.PassThrough(ctx, inputChan, func(number int) {
			sideResults = append(sideResults, number*10)
		})

		go func() {
			inputChan <- 1
			inputChan <- 2
			time.Sleep(time.Millisecond * 100) // context 종료와 inputChan 전송이 동시에 발생하지 않도록 대기
			cancelFunc()
			inputChan <- 3
		}()

		require.Equal(t, 1, <-outputChan)
		require.Equal(t, 2, <-outputChan)

		time.Sleep(time.Millisecond * 200) // context 종료 전파 대기
		isClosed, _ := cchan.IsClosed(outputChan)
		require.True(t, isClosed)

		require.Len(t, sideResults, 2, sideResults)
	})

}

func TestAsyncAwaitSteps(t *testing.T) {
	t.Run("AsyncAwaitSteps", func(t *testing.T) {
		// asyncTest(t, nil, 2, 6)
		// asyncTest(t, nil, 3, 4)
		asyncTest(t, nil, 4, 3)
		asyncTest(t, nil, 5, 3)
		asyncTest(t, nil, 6, 2)
		asyncTest(t, nil, 7, 2)
		asyncTest(t, nil, 11, 2)
		asyncTest(t, nil, 12, 1)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		inputChan := make(chan byte)
		go func() {
			inputChan <- 'a'
			time.Sleep(time.Millisecond * 1100) // 첫 번째 데이터는 처리된 후 취소
			cancel()                            // context 취소
			inputChan <- 'b'
			inputChan <- 'c' // 이 데이터는 처리되지 않음
			close(inputChan)
		}()
		asyncContext(t, ctx, inputChan, 100, 1)
	})

	testAsyncTimeout := func(t *testing.T, concurrencySize int, expectedSeconds int) {
		inputChan := make(chan byte)

		go func() {
			inputChan <- 'a'
			inputChan <- 'b'
			inputChan <- 'c'
			inputChan <- 'd'
			close(inputChan)
		}()

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*1100) //첫 사이클이 처리된 후 timeout
		asyncContext(t, ctx, inputChan, 3, 3)
	}
	t.Run("context timeout", func(t *testing.T) {
		testAsyncTimeout(t, 1, 1)
		testAsyncTimeout(t, 2, 2)
		testAsyncTimeout(t, 3, 3)
		testAsyncTimeout(t, 4, 4)
	})

}

func asyncTest(t *testing.T, bufferSize *int, concurrencySize int, expectedSeconds int) {
	ctx := context.Background()
	inputChan := make(chan byte)

	step1 := pipe.NewStep(nil, func(num byte) (int, error) {
		return int(num) - 96, nil
	})

	asyncStep, awaitStep := pipe.NewAsyncAwaitSteps(
		ctx,
		bufferSize, // asyncStep 버퍼 사이즈는 nil로 설정하여 기본값(0) 사용
		concurrencySize,
		func(ctx context.Context, num int) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(time.Second):
				// 1초 후에 결과를 반환
				return num * 2, nil
			}
		},
	)

	step2 := pipe.NewStep(nil, func(num int) (string, error) {
		return fmt.Sprintf("%d", num), nil
	})

	asyncResultChan, errChan := pipe.Pipeline4(ctx, inputChan, step1, asyncStep, awaitStep, step2)

	start := time.Now()
	go func() {
		inputChan <- 'a'
		inputChan <- 'b'
		inputChan <- 'c'
		inputChan <- 'd'
		inputChan <- 'e'
		inputChan <- 'f'
		inputChan <- 'g'
		inputChan <- 'h'
		inputChan <- 'i'
		inputChan <- 'j'
		inputChan <- 'k'
		inputChan <- 'l'
		close(inputChan)
	}()
	exepectedResults := []string{"2", "4", "6", "8", "10", "12", "14", "16", "18", "20", "22", "24"}
	index := 0
	for result := range asyncResultChan {
		require.Equal(t, exepectedResults[index], result)
		index++
	}

	require.Len(t, errChan, 0)
	for err := range errChan {
		require.Fail(t, "errChan should be empty", err)
	}

	end := time.Now()

	expectedEnd := start.Add(time.Duration(expectedSeconds) * time.Second)
	require.WithinDurationf(t, expectedEnd, end, 100*time.Millisecond, "Expected processing to complete within %d seconds, but took %d", expectedEnd.Sub(start).Milliseconds(), end.Sub(start).Milliseconds())

}

func asyncContext(t *testing.T, ctx context.Context, inputChan chan byte, concurrencySize int, expectedProcessedCount int) {

	step1 := pipe.NewStep(nil, func(num byte) (int, error) {
		return int(num) - 96, nil
	})

	asyncStep, awaitStep := pipe.NewAsyncAwaitSteps(
		ctx,
		nil,
		concurrencySize,
		func(ctx context.Context, num int) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(time.Second):
				return num * 2, nil
			}
		},
	)

	asyncResultChan, errChan := pipe.Pipeline3(ctx, inputChan, step1, asyncStep, awaitStep)

	exepectedResults := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24}
	index := 0
	for result := range asyncResultChan {
		require.Equal(t, exepectedResults[index], result)
		index++
	}
	require.Equal(t, expectedProcessedCount, index)
	index = 0
	for err := range errChan {
		require.Fail(t, "errChan should be empty", err)
		index++
	}
	require.Equal(t, index, 0) //context 취소 이후에는 데이터와 에러는 더 이상 전송되지 않아야 함

}

type errNagativeNumber struct {
	a, b int
}

func (e *errNagativeNumber) Error() string {
	return fmt.Sprintf("negative number: %d, %d", e.a, e.b)
}

func checkPositive(a, b int) (int, int, error) {
	if a < 0 || b < 0 {
		return a, b, &errNagativeNumber{a, b}
	}

	return a, b, nil
}

type errDivideByZero struct {
	a, b int
}

func (e *errDivideByZero) Error() string {
	return fmt.Sprintf("divide by zero: %d, %d", e.a, e.b)
}

type DivideTarget struct {
	denominator, numerator int
}

func divide(a, b int) (*sumTarget, error) {
	if b == 0 {
		return nil, &errDivideByZero{a, b}
	}

	return &sumTarget{
		a: a / b,
		b: a % b,
	}, nil
}

type sumTarget struct {
	a, b int
}

func sum(target *sumTarget) (int, error) {
	return target.a + target.b, nil
}

func square(a int) (int, error) {
	return a * a, nil
}
