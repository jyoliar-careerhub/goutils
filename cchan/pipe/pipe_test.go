package pipe_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jae2274/goutils/cchan/pipe"
	"github.com/stretchr/testify/require"
)

// 테스트 대상
// 각 파이프별로 정상 동작 확인
// 각 파이프별로 에러 발생시 에러가 전달되는지 확인
// bufferSize로 전달된 값이 실제 버퍼 크기에 적용되는지 확인
// quit 채널이 트리거되면 각 파이프가 종료되는지 확인
// 첫번째 파이프라인의 채널이 종료시 이후의 파이프가 연쇄적으로 종료되는지 확인
func TestPipe(t *testing.T) {
	t.Run("정상 동작 확인", func(t *testing.T) {
		checkValidChan := make(chan DivideTarget)
		errChan := make(chan error)
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
		inputs := []DivideTarget{
			{10, 3},
			{20, 5},
			{30, 7},
			{40, 9},
		}

		expectedOutputs := []int{16, 16, 36, 64}

		for _, input := range inputs {
			checkValidChan <- input
		}

		for _, expectedOutput := range expectedOutputs {
			output, ok := <-resultChan
			require.True(t, ok, "resultChan is closed")
			require.Equal(t, expectedOutput, output)
		}

		close(checkValidChan)
		time.Sleep(time.Millisecond * 100) // 파이프라인이 종료되기를 기다림
		require.True(t, isClose(divideTargetChan))
		require.True(t, isClose(sumTargetChan))
		require.True(t, isClose(squareTargetChan))
		require.True(t, isClose(resultChan))

	})
}

func checkPositive(a, b int) (int, int, error) {
	if a <= 0 || b <= 0 {
		return a, b, errors.New("negative number")
	}

	return a, b, nil
}

type DivideTarget struct {
	denominator, numerator int
}

func divide(a, b int) (*sumTarget, error) {
	if b == 0 {
		return nil, errors.New("divide by zero")
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
