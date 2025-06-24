package pipe

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkAsyncAwaitSteps(b *testing.B) {

	action := func(ctx context.Context, input int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(time.Second):
			return input * 2, nil
		}
	}
	// s := rand.NewSource(time.Now().UnixNano())
	// r := rand.New(s)
	// randAction := func(ctx context.Context, input int) (int, error) {

	// 	randTime := r.Intn(3) + 1

	// 	select {
	// 	case <-ctx.Done():
	// 		return 0, ctx.Err()
	// 	case <-time.After(time.Second * time.Duration(randTime)):
	// 		return input * 2, nil
	// 	}
	// }

	b.Run("AsyncAwaitSteps", func(b *testing.B) {
		inputChan := make(chan int, 10)

		go func() {
			for i := 0; i < 30000; i++ {
				inputChan <- i
			}
			close(inputChan)
		}()

		ctx := context.Background()
		asyncStep, outputStep := NewAsyncAwaitSteps(ctx, nil, 10000, action)

		resultChan, errChan := Pipeline2(ctx, inputChan, asyncStep, outputStep)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for result := range resultChan {
				//그저 데이터 소모 코드
				_ = result * 2
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for err := range errChan {
				b.Error(err)
			}
		}()

		wg.Wait()
	})

	b.Run("AsyncAwaitStepsWithRunGC", func(b *testing.B) {
		inputChan := make(chan int, 10)

		go func() {
			for i := 0; i < 30000; i++ {
				inputChan <- i
			}
			close(inputChan)
		}()

		ctx := context.Background()
		asyncStep, outputStep := NewAsyncAwaitSteps(ctx, nil, 10000, action)

		resultChan, errChan := Pipeline2(ctx, inputChan, asyncStep, outputStep)

		go func() {
			time.Sleep(2500 * time.Millisecond)
			runtime.GC()
		}()
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for result := range resultChan {
				//그저 데이터 소모 코드
				_ = result * 2
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for err := range errChan {
				b.Error(err)
			}
		}()

		wg.Wait()
	})
}
