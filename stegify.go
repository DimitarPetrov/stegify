//Command line tool capable of steganography encoding and decoding any file within given images as carriers
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DimitarPetrov/stegify/steg"
	"github.com/posener/cmd"
)

var (
	root         = cmd.New()
	carrierFiles = root.String("carrier", "", "carrier files in which the data is encoded (comma separated)")
	resultFiles  = root.String("result", "", "names of the result files (comma separated)")

	// Encode subcommand and flags.
	encode   = root.SubCommand("encode", "Encode data into carrier(s)")
	dataFile = encode.String("data", "", "Data file which is being encoded in the carrier")

	// Decode subcommand and flags.
	decode = root.SubCommand("decode", "Decode data from carrier(s)")
)

func main() {
	_ = root.Parse()

	carriers := splitFlag(*carrierFiles)
	results := splitFlag(*resultFiles)

	if len(carriers) == 0 {
		fail("Carrier file must be specified. Use stegify --help for more information.")
	}

	switch {
	case encode.Parsed():
		if len(results) == 0 { // if no results provided use defaults
			for i := range carriers {
				results = append(results, fmt.Sprintf("result%d", i))
			}
		}
		if len(results) != len(carriers) {
			fail("Carrier and result files count must be equal when encoding.")
		}
		if *dataFile == "" {
			fail("Data file must be specified. Use stegify --help for more information.")
		}

		err := steg.MultiCarrierEncodeByFileNames(carriers, *dataFile, results)
		if err != nil {
			fail(err.Error())
		}
	case decode.Parsed():
		if len(results) == 0 { // if no result provided use default
			results = append(results, "result")
		}
		if len(results) != 1 {
			fail("Only one result file expected.")
		}
		err := steg.MultiCarrierDecodeByFileNames(carriers, results[0])
		if err != nil {
			fail(err.Error())
		}
	default:
		fail("x")
	}
}

// splitFlag splits a comma separated value and omits empty values.
func splitFlag(value string) []string {
	var values []string
	for _, val := range strings.Split(value, ",") {
		if val != "" {
			values = append(values, val)
		}
	}
	return values
}

// fail prints the formatted message to stderr and exits.
func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
