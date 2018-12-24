package bits

import (
	"fmt"
	"testing"
)

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
