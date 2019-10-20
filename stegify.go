//Command line tool capable of steganography encoding and decoding any file in given image as carrier
package main

import (
	"flag"
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
	"os"
)

var carrierFile = flag.String("carrier", "", "carrier file in which the data is encoded")
var dataFile = flag.String("data", "", "data file which is being encoded in carrier")
var resultFile = flag.String("result", "result", "name of the result file")

func init() {
	flag.StringVar(carrierFile, "c", "", "carrier file in which the data is encoded")
	flag.StringVar(dataFile, "d", "", "data file which is being encoded in carrier")
	flag.StringVar(resultFile, "r", "result", "name of the result file")
}

const encode = "encode"
const decode = "decode"

func main() {
	operation := parseOperation()
	flag.Parse()

	if carrierFile == nil || *carrierFile == "" {
		fmt.Fprintln(os.Stderr, "Carrier file must be specified")
		os.Exit(1)
	}

	if (dataFile == nil || *dataFile == "") && operation == encode {
		fmt.Fprintln(os.Stderr, "Data file must be specified")
		os.Exit(1)
	}

	switch operation {
	case encode:
		err := steg.EncodeByFileNames(*carrierFile, *dataFile, *resultFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case decode:
		err := steg.DecodeByFileNames(*carrierFile, *resultFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func parseOperation() string {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Operation must be specified")
		os.Exit(1)
	}
	operation := os.Args[1]
	if operation != encode && operation != decode {
		if operation == "--help" {
			flag.Parse()
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Unsupported operation: %s. Only encode/decode operations supported", operation)
		os.Exit(1)
	}

	os.Args = append(os.Args[:1], os.Args[2:]...) // needed because go flags implementation stop parsing after first non-flag argument
	return operation
}
