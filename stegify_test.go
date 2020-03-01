package main_test

import (
	"bytes"
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building stegify binary: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	_ = os.Remove("stegify")
	os.Exit(exitCode)
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		data       string
		results    []string
		shouldFail bool
	}{
		{
			name:    "Basic encode",
			args:    []string{"encode", "--carrier", "examples/street.jpeg", "--data", "examples/lake.jpeg", "--result", "result.jpeg"},
			data:    "examples/lake.jpeg",
			results: []string{"result.jpeg"},
		},
		{
			name:    "Encode with multiple carriers using --carrier flag and --result flag",
			args:    []string{"encode", "--carrier", "examples/street.jpeg", "--carrier", "examples/lake.jpeg", "--data", "examples/video.mp4", "--result", "result1.jpeg", "--result", "result2.jpeg"},
			data:    "examples/video.mp4",
			results: []string{"result1.jpeg", "result2.jpeg"},
		},
		{
			name:    "Encode with multiple carriers using --carriers flag and --results flag",
			args:    []string{"encode", "--carriers", "examples/street.jpeg examples/lake.jpeg", "--data", "examples/video.mp4", "--results", "result3.jpeg result4.jpeg"},
			data:    "examples/video.mp4",
			results: []string{"result3.jpeg", "result4.jpeg"},
		},
		{
			name:    "Encode with multiple carriers using --carriers flag and --result flag",
			args:    []string{"encode", "--carriers", "examples/street.jpeg examples/lake.jpeg", "--data", "examples/video.mp4", "--result", "result5.jpeg", "--result", "result6.jpeg"},
			data:    "examples/video.mp4",
			results: []string{"result5.jpeg", "result6.jpeg"},
		},
		{
			name:    "Encode with multiple carriers using --carrier flag and --results flag",
			args:    []string{"encode", "--carrier", "examples/street.jpeg", "--carrier", "examples/lake.jpeg", "--data", "examples/video.mp4", "--results", "result7.jpeg result8.jpeg"},
			data:    "examples/video.mp4",
			results: []string{"result7.jpeg", "result8.jpeg"},
		},
		{
			name:    "Encode with multiple carriers using mixed flags",
			args:    []string{"encode", "--carriers", "examples/street.jpeg", "--carrier", "examples/lake.jpeg", "--data", "examples/video.mp4", "--results", "result9.jpeg", "--result", "result10.jpeg"},
			data:    "examples/video.mp4",
			results: []string{"result10.jpeg", "result9.jpeg"}, // --carrier/--result is with priority over --carriers/--results
		},
		{
			name:    "Encode with single carrier should add default result name",
			args:    []string{"encode", "--carrier", "examples/street.jpeg", "--data", "examples/lake.jpeg"},
			data:    "examples/lake.jpeg",
			results: []string{"result0"},
		},
		{
			name:    "Encode with multiple carriers using --carriers should add default result names",
			args:    []string{"encode", "--carriers", "examples/street.jpeg examples/lake.jpeg", "--data", "examples/video.mp4"},
			data:    "examples/video.mp4",
			results: []string{"result0", "result1"},
		},
		{
			name:       "Encode carriers count does not match results count should return an error",
			args:       []string{"encode", "--carriers", "examples/street.jpeg examples/lake.jpeg", "--data", "examples/video.mp4", "--results", "result1.jpeg"},
			shouldFail: true,
		},
		{
			name:       "Encode without data file should fail",
			args:       []string{"encode", "--carrier", "examples/street.jpeg", "--result", "result.jpeg"},
			shouldFail: true,
		},
		{
			name:       "Encode without carrier file should fail",
			args:       []string{"encode", "--data", "examples/lake.jpeg"},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Executing: stegify %s", strings.Join(test.args, " "))
			cmd := exec.Command("./stegify", test.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				if test.shouldFail {
					return
				} else {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			for _, resultName := range test.results {
				defer os.Remove(resultName)
			}

			err = steg.MultiCarrierDecodeByFileNames(test.results, "decode_result")
			if err != nil {
				t.Fatalf("Error decoding results: %v", err)
			}
			defer os.Remove("decode_result")

			assertEqualFiles(t, test.data, "decode_result")
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expected   string
		result     string
		shouldFail bool
	}{
		{
			name:     "Basic decode",
			args:     []string{"decode", "--carrier", "examples/test_decode.jpeg", "--result", "result.jpeg"},
			expected: "examples/lake.jpeg",
			result:   "result.jpeg",
		},
		{
			name:     "Decode from multiple carriers using --carier flag",
			args:     []string{"decode", "--carrier", "examples/test_multi_carrier_decode1.jpeg", "--carrier", "examples/test_multi_carrier_decode2.jpeg", "--result", "result1.mp4"},
			expected: "examples/video.mp4",
			result:   "result1.mp4",
		},
		{
			name:     "Decode from multiple carriers using --cariers flag",
			args:     []string{"decode", "--carriers", "examples/test_multi_carrier_decode1.jpeg examples/test_multi_carrier_decode2.jpeg", "--result", "result2.mp4"},
			expected: "examples/video.mp4",
			result:   "result2.mp4",
		},
		{
			name:     "Decode without result flag should add default",
			args:     []string{"decode", "--carrier", "examples/test_decode.jpeg"},
			expected: "examples/lake.jpeg",
			result:   "result",
		},
		{
			name:       "Decode with multiple results should fail",
			args:       []string{"decode", "--carriers", "examples/test_multi_carrier_decode1.jpeg examples/test_multi_carrier_decode2.jpeg", "--result", "result1.mp4", "--result", "result2.mp4"},
			shouldFail: true,
		},
		{
			name:       "Decode without carrier file should fail",
			args:       []string{"decode", "--result", "result"},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Executing: stegify %s", strings.Join(test.args, " "))
			cmd := exec.Command("./stegify", test.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				if test.shouldFail {
					return
				} else {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			defer os.Remove(test.result)

			assertEqualFiles(t, test.expected, test.result)
		})
	}
}

func assertEqualFiles(t *testing.T, expected string, given string) {
	expectedReader, err := os.Open(expected)
	if err != nil {
		t.Fatalf("Error opening data file %s:%v", expected, err)
	}
	defer expectedReader.Close()

	decodeResult, err := os.Open(given)
	if err != nil {
		t.Fatalf("Error opening decode result file:%v", err)
	}
	defer decodeResult.Close()

	wantedBytes, err := ioutil.ReadAll(expectedReader)
	if err != nil {
		t.Fatalf("Error reading data file: %v", err)
	}

	resultBytes, err := ioutil.ReadAll(decodeResult)
	if err != nil {
		t.Fatalf("Error reading decode result file: %v", err)
	}

	if !bytes.Equal(wantedBytes, resultBytes) {
		t.Error("Assertion failed!")
	}
}
