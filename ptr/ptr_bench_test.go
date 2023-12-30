package ptr

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkIsNil(b *testing.B) {
	old := os.Stdout

	// 표준 출력을 /dev/null로 변경
	file, _ := os.Create("/dev/null")
	os.Stdout = file

	b.Run("Just check nil via if", func(b *testing.B) {
		var sampleErr *SampleError = nil
		for i := 0; i < b.N; i++ {
			fmt.Print(sampleErr == nil)
		}
	})

	b.Run("basic error", func(b *testing.B) {
		var _ *SampleError = nil
		for i := 0; i < b.N; i++ {
			fmt.Print(IsNil(nil))

		}
	})

	b.Run("custom pointer error", func(b *testing.B) {
		var sampleErr *SampleError = nil
		for i := 0; i < b.N; i++ {
			fmt.Print(IsNil(sampleErr))
		}
	})

	b.Run("Not nil", func(b *testing.B) {
		var sampleErr *SampleError = &SampleError{}
		for i := 0; i < b.N; i++ {
			fmt.Print(IsNil(&sampleErr))
		}
	})

	// 표준 출력 원복
	os.Stdout = old
}
