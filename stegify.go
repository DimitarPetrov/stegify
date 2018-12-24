package main

import (
	"fmt"
	"stegify/steg"
)

func main() {
	fmt.Println(steg.Encode("test2.png", "README.md", "steg_result"))
	fmt.Println(steg.Decode("steg_result.png", "decode_result.md"))
}
