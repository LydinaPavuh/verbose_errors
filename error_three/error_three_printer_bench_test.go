package error_three

import (
	"strconv"
	"testing"
)

func BenchmarkPrintErrorThree(b *testing.B) {
	rootErr := Error("Root")
	nextErr := rootErr
	for i := 0; i < 100; i++ {
		nextErr = Error(strconv.Itoa(i), nextErr)
	}
	b.Run("bench", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			PrintErrorThree(rootErr)
		}
	})
}
