package steg

import (
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Encode("../examples/street.jpeg", "../examples/lake.jpeg", "benchmark_result")
	}

	os.Remove("benchmark_result.jpeg")
}
