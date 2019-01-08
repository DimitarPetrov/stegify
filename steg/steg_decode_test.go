package steg

import (
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Decode("../examples/benchmark_test_decode.jpeg", "benchmark_result")
	}

	os.Remove("benchmark_result")
}
