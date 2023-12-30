package ptr

import (
	"testing"
)

func BenchmarkIsNil(b *testing.B) {
	b.Run("basic error", func(b *testing.B) {
		var _ *SampleError = nil
		for i := 0; i < b.N; i++ {
			IsNil(nil)
		}
	})

	b.Run("custom pointer error", func(b *testing.B) {
		var sampleErr *SampleError = nil
		for i := 0; i < b.N; i++ {
			IsNil(sampleErr)
		}
	})

	b.Run("Not nil", func(b *testing.B) {
		var sampleErr *SampleError = &SampleError{}
		for i := 0; i < b.N; i++ {
			IsNil(&sampleErr)
		}
	})
}
