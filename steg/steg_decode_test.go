package steg

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		err := Decode("../examples/test_decode.jpeg", "benchmark_result")
		if err != nil {
			b.Fatalf("Error decoding file: %v", err)
		}
	}

	err := os.Remove("benchmark_result")
	if err != nil {
		b.Fatalf("Error removing benchmark_result file: %v", err)
	}
}

func TestDecode(t *testing.T) {

	err := Decode("../examples/test_decode.jpeg", "result")
	if err != nil {
		t.Fatalf("Error decoding file: %v", err)
	}

	defer func() {
		err = os.Remove("result")
		if err != nil {
			t.Fatalf("Error removing result file: %v", err)
		}
	}()

	wanted, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening file examples/lake.jpg: %v", err)
	}
	defer wanted.Close()

	result, err := os.Open("result")
	if err != nil {
		t.Fatalf("Error opening file result: %v", err)
	}
	defer result.Close()

	wantedBytes, err := ioutil.ReadAll(wanted)
	if err != nil {
		t.Fatalf("Error reading file examples/lake.jpg: %v", err)
	}

	resultBytes, err := ioutil.ReadAll(result)
	if err != nil {
		t.Fatalf("Error reading file result: %v", err)
	}

	if !bytes.Equal(wantedBytes, resultBytes) {
		t.Error("Assertion failed!")
	}

}

func TestDecodeShouldReturnErrorWhenCarrierFileMissing(t *testing.T) {

	err := Decode("not_existing_file", "result")
	if err == nil {
		os.Remove("result")
		t.FailNow()
	}
	t.Log(err)

}

func TestDecodeShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {

	err := Decode("../README.md", "result")
	if err == nil {
		t.FailNow()
	}
	t.Log(err)

}