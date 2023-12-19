package pipe_test

import (
	"fmt"
	"log"
	"testing"
	"time"

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

// 테스트 대상
// 각 파이프별로 정상 동작 확인
// 각 파이프별로 에러 발생시 에러가 전달되는지 확인
// bufferSize로 전달된 값이 실제 버퍼 크기에 적용되는지 확인
// quit 채널이 트리거되면 각 파이프가 종료되는지 확인
// 첫번째 파이프라인의 채널이 종료시 이후의 파이프가 연쇄적으로 종료되는지 확인
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
		errChan := make(chan error, 10)
		quitChan := make(chan bool)

		outputChan := pipe.Pipe(inputChan, errChan, quitChan, nil, square) //outputChan은 버퍼가 없으므로 해당 채널의 데이터를 수신하지 않으면 전송자가 블록된다.
		fmt.Println(outputChan)                                            //컴파일 에러를 방지하기 위해 사용한 코드

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
		errChan := make(chan error, 10)
		quitChan := make(chan bool)

		outputChan := pipe.Pipe(inputChan, errChan, quitChan, ptr.P(3), square) //outputChan의 버퍼는 3개이다.
		fmt.Println(outputChan)                                                 //컴파일 에러를 방지하기 위해 사용한 코드

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

}

func test(t *testing.T, inputs []DivideTarget, expectedOutputs []int, errs []error) {
	checkValidChan := make(chan DivideTarget)
	errChan := make(chan error, 10) // 파이프 내부에서 error 전송시 block되지 않도록 충분한 버퍼를 가진 채널을 생성한다.
	quitChan := make(chan bool)

	divideTargetChan := pipe.Pipe(checkValidChan, errChan, quitChan, nil,
		func(target DivideTarget) (*DivideTarget, error) {
			a, b, err := checkPositive(target.denominator, target.numerator)
			if err != nil {
				return nil, err
			}
			return &DivideTarget{a, b}, nil
		},
	)
	sumTargetChan := pipe.Pipe(divideTargetChan, errChan, quitChan, nil, func(dt *DivideTarget) (*sumTarget, error) {
		return divide(dt.denominator, dt.numerator)
	})
	squareTargetChan := pipe.Pipe(sumTargetChan, errChan, quitChan, nil, sum)
	resultChan := pipe.Pipe(squareTargetChan, errChan, quitChan, nil, square)

	for _, input := range inputs {
		checkValidChan <- input
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

	close(checkValidChan)
	time.Sleep(time.Millisecond * 100) // 파이프라인이 종료되기를 기다림
	require.True(t, isClose(divideTargetChan))
	require.True(t, isClose(sumTargetChan))
	require.True(t, isClose(squareTargetChan))
	require.True(t, isClose(resultChan))
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
func isClose[T any](c <-chan T) bool {
	select {
	case _, ok := <-c:
		if ok {
			return false
		}
		return true
	default:
		return false
	}
}
