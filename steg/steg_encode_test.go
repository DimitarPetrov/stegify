package steg

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Encode("../examples/street.jpeg", "../examples/lake.jpeg", "benchmark_result")
	}

	os.Remove("benchmark_result.jpeg")
}

func TestEncode(t *testing.T) {

	err := Encode("../examples/street.jpeg", "../examples/lake.jpeg", "encoded_result")
	if err != nil {
		t.Fatalf("Error encoding file: %v", err)
	}

	defer func() {
		err = os.Remove("encoded_result.jpeg")
		if err != nil {
			t.Fatalf("Error removing result file: %v", err)
		}
	}()

	err = Decode("encoded_result.jpeg", "result")
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

func TestEncodeShouldReturnErrorWhenCarrierFileMissing(t *testing.T) {

	err := Encode("not_existing_file", "../examples/lake.jpeg", "encoded_result")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {

	err := Encode("../README.md", "../examples/lake.jpeg", "encoded_result")
	if err == nil {
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeShouldReturnErrorWhenDataFileMissing(t *testing.T) {

	err := Encode("../examples/street.jpeg", "not_existing_file", "encoded_result")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeShouldReturnErrorWhenDataFileTooLarge(t *testing.T) {

	err := Encode("../examples/lake.jpeg", "../examples/test_decode.jpeg", "encoded_result")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)
}
