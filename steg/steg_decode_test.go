package steg_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/DimitarPetrov/stegify/steg"
)

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		carrier, err := os.Open("../examples/test_decode.jpeg")
		if err != nil {
			b.Fatalf("Error opening carrier file: %v", err)
		}

		var result bytes.Buffer

		err = steg.Decode(carrier, &result)
		if err != nil {
			b.Fatalf("Error decoding file: %v", err)
		}

		carrier.Close()
	}
}

func BenchmarkDecodeByFileNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := steg.DecodeByFileNames("../examples/test_decode.jpeg", "benchmark_result")
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
	AssertDecodedDataMatchesOriginal(t, []string{"../examples/test_decode.jpeg"}, "../examples/lake.jpeg",
		func(readers []io.Reader, writer io.Writer) {
			if len(readers) != 1 {
				t.Fatalf("Exactly one reader expected")
			}
			err := steg.Decode(readers[0], writer)
			if err != nil {
				t.Fatalf("Error decoding file: %v", err)
			}
		})
}

func TestMultiCarrierDecode(t *testing.T) {
	AssertDecodedDataMatchesOriginal(t, []string{"../examples/test_multi_carrier_decode1.jpeg", "../examples/test_multi_carrier_decode2.jpeg"}, "../examples/video.mp4",
		func(readers []io.Reader, writer io.Writer) {
			err := steg.MultiCarrierDecode(readers, writer)
			if err != nil {
				t.Fatalf("Error decoding file: %v", err)
			}
		})
}

func TestMultiCarrierDecodeOrderMatters(t *testing.T) {
	AssertDecodedDataDoesNotMatchOriginal(t, []string{"../examples/test_multi_carrier_decode2.jpeg", "../examples/test_multi_carrier_decode1.jpeg"}, "../examples/video.mp4",
		func(readers []io.Reader, writer io.Writer) {
			err := steg.MultiCarrierDecode(readers, writer)
			if err != nil {
				t.Fatalf("Error decoding file: %v", err)
			}
		})
}

func TestDecodeByFileNames(t *testing.T) {
	AssertDecodeByFileNames(t, []string{"../examples/test_decode.jpeg"}, "../examples/lake.jpeg", func(strings []string, s string) {
		if len(strings) != 1 {
			t.Fatalf("Exactly one carrier expected")
		}
		err := steg.DecodeByFileNames(strings[0], s)
		if err != nil {
			t.Fatalf("Error decoding file: %v", err)
		}
	})
}

func TestMultiCarrierDecodeByFileNames(t *testing.T) {
	AssertDecodeByFileNames(t, []string{"../examples/test_multi_carrier_decode1.jpeg", "../examples/test_multi_carrier_decode2.jpeg"}, "../examples/video.mp4", func(strings []string, s string) {
		err := steg.MultiCarrierDecodeByFileNames(strings, s)
		if err != nil {
			t.Fatalf("Error decoding file: %v", err)
		}
	})
}

func TestDecodeByFileNamesShouldReturnErrorWhenCarrierFileMissing(t *testing.T) {
	err := steg.DecodeByFileNames("not_existing_file", "result")
	if err == nil {
		os.Remove("result")
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierDecodeByFileNamesShouldReturnErrorWhenCarrierFileMissing(t *testing.T) {
	err := steg.MultiCarrierDecodeByFileNames([]string{"not_existing_file1", "non_existing_file2"}, "result")
	if err == nil {
		os.Remove("result")
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierDecodeByFileNamesShouldReturnErrorWhenNoCarrierFilesProvided(t *testing.T) {
	err := steg.MultiCarrierDecodeByFileNames([]string{}, "result")
	if err == nil {
		os.Remove("result")
		t.FailNow()
	}
	t.Log(err)
}

func TestDecodeShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {
	carrier, err := os.Open("../README.md")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier.Close()

	var result bytes.Buffer
	err = steg.Decode(carrier, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierDecodeShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {
	carrier, err := os.Open("../README.md")
	if err != nil {
		t.Fatalf("Error opening carrier file: %v", err)
	}
	defer carrier.Close()

	var result bytes.Buffer
	err = steg.MultiCarrierDecode([]io.Reader{carrier}, &result)
	if err == nil {
		t.FailNow()
	}
	t.Log(err)
}

func TestMultiCarrierDecodeByFileNamesShouldReturnErrorWhenCarrierFileIsNotImage(t *testing.T) {
	err := steg.MultiCarrierDecodeByFileNames([]string{"../README.md"}, "result")
	if err == nil {
		os.Remove("result")
		t.FailNow()
	}
	t.Log(err)
}

func AssertDecodedDataMatchesOriginal(t *testing.T, carrierNames []string, expected string, f func([]io.Reader, io.Writer)) {
	AssertDecode(t, carrierNames, expected, f, true)
}

func AssertDecodedDataDoesNotMatchOriginal(t *testing.T, carrierNames []string, expected string, f func([]io.Reader, io.Writer)) {
	AssertDecode(t, carrierNames, expected, f, false)
}

func AssertDecode(t *testing.T, carrierNames []string, expected string, f func([]io.Reader, io.Writer), shouldEqual bool) {
	carriers := make([]io.Reader, 0, len(carrierNames))

	for _, name := range carrierNames {
		carrier, err := os.Open(name)
		if err != nil {
			t.Fatalf("Error opening %s file: %v", name, err)
		}
		defer carrier.Close()
		carriers = append(carriers, carrier)
	}

	var result bytes.Buffer

	f(carriers, &result)

	wanted, err := os.Open(expected)
	if err != nil {
		t.Fatalf("Error opening file %s: %v", expected, err)
	}
	defer wanted.Close()

	wantedBytes, err := ioutil.ReadAll(wanted)
	if err != nil {
		t.Fatalf("Error reading file %s: %v", expected, err)
	}

	resultBytes, err := ioutil.ReadAll(&result)
	if err != nil {
		t.Fatalf("Error reading result Writer: %v", err)
	}

	if bytes.Equal(wantedBytes, resultBytes) != shouldEqual {
		t.Error("Assertion failed!")
	}
}

func AssertDecodeByFileNames(t *testing.T, carrierNames []string, expected string, f func([]string, string)) {
	f(carrierNames, "result")

	defer os.Remove("result")

	wanted, err := os.Open(expected)
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
