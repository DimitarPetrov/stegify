package steg

import (
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Encode("../test.png", "../test.jpeg", "benchmark_result")
	}

	os.Remove("benchmark_result.png")
}