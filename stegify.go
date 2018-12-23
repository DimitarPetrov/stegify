package main

import (
	"fmt"
	"stegify/steg"
)

func main() {
	fmt.Println(steg.Encode("test.jpeg", "test.png"))
}
