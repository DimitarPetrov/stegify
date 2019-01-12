package bits

import (
	"fmt"
	"testing"
)

func TestQuartersOfByte(t *testing.T) {
	var tests = []struct {
		input  byte
		result [4]byte
	}{
		{231, [4]byte{3, 2, 1, 3}},
		{90, [4]byte{1, 1, 2, 2}},
		{239, [4]byte{3, 2, 3, 3}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("QuartersOfByte(%08b)", test.input), func(t *testing.T) {
			if actual := QuartersOfByte(test.input); actual != test.result {
				t.Errorf("Expected %08b (%d) but got %08b (%d)", test.result, test.result, actual, actual)
			}
		})
	}
}

func ExampleQuartersOfByte() {
	fmt.Printf("%08b", QuartersOfByte(231)) //231 is 11100111 in binary
	//Output:
	//[00000011 00000010 00000001 00000011]
}

func TestSetLastTwoBits(t *testing.T) {
	var tests = []struct {
		b, v, result byte
	}{
		{byte(134), byte(1), byte(133)},
		{byte(123), byte(1), byte(121)},
		{byte(234), byte(2), byte(234)},
		{byte(23), byte(2), byte(22)},
		{byte(134), byte(3), byte(135)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("SetLastTwoBits(%08b,%08b)", test.b, test.v), func(t *testing.T) {
			if actual := SetLastTwoBits(test.b, test.v); actual != test.result {
				t.Errorf("Expected %08b (%d) but got %08b (%d)", test.result, test.result, actual, actual)
			}
		})
	}
}

func ExampleSetLastTwoBits() {
	fmt.Printf("%08b", SetLastTwoBits(134, 1)) // 134 is 10000110 and 1 is 00000001 in binary
	//Output:
	//10000101
}

func TestGetLastTwoBits(t *testing.T) {
	var tests = []struct {
		b, result byte
	}{
		{byte(134), byte(2)},
		{byte(123), byte(3)},
		{byte(234), byte(2)},
		{byte(23), byte(3)},
		{byte(133), byte(1)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("GetLastTwoBits(%08b)", test.b), func(t *testing.T) {
			if actual := GetLastTwoBits(test.b); actual != test.result {
				t.Errorf("Expected %08b (%d) but got %08b (%d)", test.result, test.result, actual, actual)
			}
		})
	}
}

func ExampleGetLastTwoBits() {
	fmt.Printf("%08b", GetLastTwoBits(134)) // 134 is 10000110 in binary
	//Output:
	//00000010
}

func TestConstructByteOfQuarters(t *testing.T) {
	var tests = []struct {
		first, second, third, fourth, result byte
	}{
		{byte(3), byte(2), byte(1), byte(3), byte(231)},
		{byte(1), byte(1), byte(2), byte(2), byte(90)},
		{byte(3), byte(2), byte(3), byte(3), byte(239)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ConstructByteOfQuarters(%08b, %08b, %08b, %08b)", test.first, test.second, test.third, test.fourth),
			func(t *testing.T) {
				if actual := ConstructByteOfQuarters(test.first, test.second, test.third, test.fourth); actual != test.result {
					t.Errorf("Expected %08b (%d) but got %08b (%d)", test.result, test.result, actual, actual)
				}
			})
	}
}

func ExampleConstructByteOfQuarters() {
	fmt.Printf("%08b", ConstructByteOfQuarters(3, 2, 1, 3)) // 00000011 00000010 00000001 00000011 in binary
	//Output:
	//11100111
}

func TestConstructByteOfQuartersAsSlice(t *testing.T) {
	var tests = []struct {
		input  []byte
		result byte
	}{
		{[]byte{3, 2, 1, 3}, byte(231)},
		{[]byte{1, 1, 2, 2}, byte(90)},
		{[]byte{3, 2, 3, 3}, byte(239)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ConstructByteOfQuartersAsSlice(%08b)", test.input),
			func(t *testing.T) {
				if actual := ConstructByteOfQuartersAsSlice(test.input); actual != test.result {
					t.Errorf("Expected %08b (%d) but got %08b (%d)", test.result, test.result, actual, actual)
				}
			})
	}
}

func ExampleConstructByteOfQuartersAsSlice() {
	var quarters = []byte{3, 2, 1, 3} // 00000011 00000010 00000001 00000011 in binary
	fmt.Printf("%08b", ConstructByteOfQuartersAsSlice(quarters))
	//Output:
	//11100111
}
