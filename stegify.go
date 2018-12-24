package main

import (
	"fmt"
	"stegify/steg"
)

func main() {
	fmt.Println(steg.Encode("test.png", "test.jpeg", "steg_result"))
	//fmt.Println(steg.Decode("benchmark_test_decode.png", "decode_result"))
}
