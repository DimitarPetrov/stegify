//Command line tool capable of steganography encoding and decoding any file in given image as carrier
package main

import (
	"flag"
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
	"os"
)

var operation = flag.String("op", "", "operation (one of the following: encode, decode)")
var carrierFile = flag.String("carrier", "", "carrier file in which the data is encoded")
var dataFile = flag.String("data", "", "data file which is being encoded in carrier")
var resultFile = flag.String("result", "result", "result file of operation (with carrier's file extension when encoding)")

func main() {
	flag.Parse()

	if operation == nil || *operation == "" {
		fmt.Fprintf(os.Stderr, "Operation must be specified")
		return
	}

	if carrierFile == nil || *carrierFile == "" {
		fmt.Fprintf(os.Stderr, "Carrier file must be specified")
		return
	}

	if (dataFile == nil || *dataFile == "") && *operation == "encode" {
		fmt.Fprintf(os.Stderr, "Data file must be specified")
		return
	}

	switch *operation {
	case "encode":
		err := steg.Encode(*carrierFile, *dataFile, *resultFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "decode":
		err := steg.Decode(*carrierFile, *resultFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported operation: %q", *operation)
	}

}
