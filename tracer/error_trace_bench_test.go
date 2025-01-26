package tracer

import (
	"fmt"
	"testing"
)

func addFrames(b *testing.B, depth int, message string) {
	if depth <= 1 {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			New(message)
		}
		return
	}
	addFrames(b, depth-1, message)
}

func BenchmarkNew(b *testing.B) {
	for _, frames := range []int{5, 10, 20, 40} {
		suffix := fmt.Sprintf("%d", frames)
		b.Run(suffix, func(b *testing.B) {
			addFrames(b, frames, "test error")
		})
	}
}
