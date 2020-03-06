package steg_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/DimitarPetrov/stegify/steg"
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
		err = steg.Encode(carrier, data, &result)
		if err != nil {
			b.Fatalf("Error encoding file: %v", err)
		}

		carrier.Close()
		data.Close()
	}
}

func BenchmarkEncodeByFileNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := steg.EncodeByFileNames("../examples/street.jpeg", "../examples/lake.jpeg", "benchmark_result.jpeg")
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
	AssertEncode(t, []string{"../examples/street.jpeg"}, "../examples/lake.jpeg",
		func(readers []io.Reader, reader io.Reader, writer io.Writer) {
			if len(readers) != 1 {
				t.Fatalf("Exactly one reader expected")
			}
			var encodeResult bytes.Buffer
			err := steg.Encode(readers[0], reader, &encodeResult)
			if err != nil {
				t.Fatalf("Error encoding files: %v", err)
			}

			err = steg.Decode(&encodeResult, writer)
			if err != nil {
				t.Fatalf("Error decoding files: %v", err)
			}
		})
}

func TestMultiCarrierEncode(t *testing.T) {
	AssertEncode(t, []string{"../examples/street.jpeg", "../examples/lake.jpeg"}, "../examples/video.mp4",
		func(readers []io.Reader, reader io.Reader, writer io.Writer) {
			var encodeResult1 bytes.Buffer
			var encodeResult2 bytes.Buffer
			err := steg.MultiCarrierEncode(readers, reader, []io.Writer{&encodeResult1, &encodeResult2})
			if err != nil {
				t.Fatalf("Error encoding files: %v", err)
			}

			err = steg.MultiCarrierDecode([]io.Reader{&encodeResult1, &encodeResult2}, writer)
			if err != nil {
				t.Fatalf("Error decoding files: %v", err)
			}
		})
}

func TestEncodeByFileNames(t *testing.T) {
	err := steg.EncodeByFileNames("../examples/street.jpeg", "../examples/lake.jpeg", "encoded_result.jpeg")
	if err != nil {
		t.Fatalf("Error encoding file: %v", err)
	}

	defer func() {
		err = os.Remove("encoded_result.jpeg")
		if err != nil {
			t.Fatalf("Error removing result file: %v", err)
		}
	}()

	err = steg.DecodeByFileNames("encoded_result.jpeg", "result")
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

func TestEncodeShouldReturnErrorWhenCarrierFileNotExists(t *testing.T) {
	err := steg.EncodeByFileNames("not_existing_file", "../examples/lake.jpeg", "encoded_result.jpeg")
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierEncodeShouldReturnErrorWhenCarrierFileNotProvided(t *testing.T) {
	err := steg.MultiCarrierEncodeByFileNames([]string{}, "../examples/lake.jpeg", []string{"encoded_result.jpeg"})
	if err == nil {
		os.Remove("encoded_result.jpeg")
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierEncodeByFileNamesShouldReturnErrorWhenCarrierAndResultFilesCountDoesNotMatch(t *testing.T) {
	err := steg.MultiCarrierEncodeByFileNames([]string{"../examples/street.jpeg", "../examples/lake.jpeg"}, "../examples/lake.jpeg", []string{"encoded_result.jpeg"})
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
	err = steg.Encode(carrier, data, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestEncodeByFileNamesShouldReturnErrorWhenDataFileMissing(t *testing.T) {
	err := steg.EncodeByFileNames("../examples/street.jpeg", "not_existing_file", "encoded_result.jpeg")
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
	err = steg.Encode(carrier, data, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierEncodeShouldReturnErrorWhenDataFileTooLarge(t *testing.T) {
	carrier1, err := os.Open("../examples/lake.jpeg")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier1.Close()

	carrier2, err := os.Open("../examples/street.jpeg")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier2.Close()

	data, err := os.Open("../examples/test_decode.jpeg")
	if err != nil {
		t.Fatalf("Error opening data file: %v", err)
	}
	defer data.Close()

	var result1 bytes.Buffer
	var result2 bytes.Buffer
	err = steg.MultiCarrierEncode([]io.Reader{carrier1, carrier2}, data, []io.Writer{&result1, &result2})
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierEncodeByFileNamesShouldReturnErrorWhenDataFileTooLarge(t *testing.T) {
	err := steg.MultiCarrierEncodeByFileNames([]string{"../examples/lake.jpeg", "../examples/street.jpeg"}, "../examples/test_decode.jpeg", []string{"result1", "result2"})
	if err == nil {
		os.Remove("result1")
		os.Remove("result2")
		t.FailNow()
	}
	t.Log(err)
}

func AssertEncode(t *testing.T, carrierNames []string, dataName string, f func([]io.Reader, io.Reader, io.Writer)) {
	carriers := make([]io.Reader, 0, len(carrierNames))

	for _, name := range carrierNames {
		carrier, err := os.Open(name)
		if err != nil {
			t.Fatalf("Error opening %s file: %v", name, err)
		}
		defer carrier.Close()
		carriers = append(carriers, carrier)
	}

	data, err := os.Open(dataName)
	if err != nil {
		t.Fatalf("Error opening %s file: %v", dataName, err)
	}
	defer data.Close()

	var result bytes.Buffer

	f(carriers, data, &result)

	wanted, err := os.Open(dataName)
	if err != nil {
		t.Fatalf("Error opening file %s: %v", dataName, err)
	}
	defer wanted.Close()

	wantedBytes, err := ioutil.ReadAll(wanted)
	if err != nil {
		t.Fatalf("Error reading file %s: %v", dataName, err)
	}

	resultBytes, err := ioutil.ReadAll(&result)
	if err != nil {
		t.Fatalf("Error reading result Writer: %v", err)
	}

	if !bytes.Equal(wantedBytes, resultBytes) {
		t.Error("Assertion failed!")
	}
}
