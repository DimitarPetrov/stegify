package steg

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		carrier, err := os.Open("../examples/street.jpeg")
		if err != nil {
			b.Fatalf("Error opening carrier file: %v", err)
		}

		data, err := os.Open("../examples/lake.jpeg")
		if err != nil {
			b.Fatalf("Error opening data file: %v", err)
		}

		var result bytes.Buffer
		err = Encode(carrier, data, &result)
		if err != nil {
			b.Fatalf("Error encoding file: %v", err)
		}

		carrier.Close()
		data.Close()
	}
}

func BenchmarkEncodeByFileNames(b *testing.B) {

	for i := 0; i < b.N; i++ {
		err := EncodeByFileNames("../examples/street.jpeg", "../examples/lake.jpeg", "benchmark_result")
		if err != nil {
			b.Fatalf("Error encoding file: %v", err)
		}
	}

	err := os.Remove("benchmark_result.jpeg")
	if err != nil {
		b.Fatalf("Error removing benchmark_result file: %v", err)
	}
}

func TestEncode(t *testing.T) {

	carrier, err := os.Open("../examples/street.jpeg")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier.Close()

	data, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening data file: %v", err)
	}
	defer data.Close()

	var encodeResult bytes.Buffer
	err = Encode(carrier, data, &encodeResult)
	if err != nil {
		t.Fatalf("Error encoding file: %v", err)
	}

	var result bytes.Buffer
	err = Decode(&encodeResult, &result)
	if err != nil {
		t.Fatalf("Error decoding file: %v", err)
	}

	wanted, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening file examples/lake.jpg: %v", err)
	}
	defer wanted.Close()

	wantedBytes, err := ioutil.ReadAll(wanted)
	if err != nil {
		t.Fatalf("Error reading file examples/lake.jpg: %v", err)
	}

	resultBytes, err := ioutil.ReadAll(&result)
	if err != nil {
		t.Fatalf("Error reading result Writer: %v", err)
	}

	if !bytes.Equal(wantedBytes, resultBytes) {
		t.Error("Assertion failed!")
	}
}

func TestEncodeByFileNames(t *testing.T) {

	err := EncodeByFileNames("../examples/street.jpeg", "../examples/lake.jpeg", "encoded_result")
	if err != nil {
		t.Fatalf("Error encoding file: %v", err)
	}

	defer func() {
		err = os.Remove("encoded_result.jpeg")
		if err != nil {
			t.Fatalf("Error removing result file: %v", err)
		}
	}()

	err = DecodeByFileNames("encoded_result.jpeg", "result")
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

	err := EncodeByFileNames("not_existing_file", "../examples/lake.jpeg", "encoded_result")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {

	carrier, err := os.Open("../README.md")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier.Close()

	data, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening data file: %v", err)
	}
	defer data.Close()

	var result bytes.Buffer
	err = Encode(carrier, data, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeByFileNamesShouldReturnErrorWhenDataFileMissing(t *testing.T) {

	err := EncodeByFileNames("../examples/street.jpeg", "not_existing_file", "encoded_result")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)

}

func TestEncodeShouldReturnErrorWhenDataFileTooLarge(t *testing.T) {
	carrier, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier.Close()

	data, err := os.Open("../examples/test_decode.jpeg")
	if err != nil {
		t.Fatalf("Error opening data file: %v", err)
	}
	defer data.Close()

	var result bytes.Buffer
	err = Encode(carrier, data, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}
