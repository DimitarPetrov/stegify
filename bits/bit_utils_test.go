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
