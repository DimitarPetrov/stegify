//Command line tool capable of steganography encoding and decoding any file within given images as carriers
package main

import (
	"flag"
	"fmt"
	"github.com/DimitarPetrov/stegify/steg"
	"os"
	"strings"
)

const encode = "encode"
const decode = "decode"

type sliceFlag []string

func (sf *sliceFlag) String() string {
	return strings.Join(*sf, " ")
}

func (sf *sliceFlag) Set(value string) error {
	*sf = append(*sf, value)
	return nil
}

var carrierFilesSlice sliceFlag
var carrierFiles = flag.String("carriers", "", "carrier files in which the data is encoded (separated by space and surrounded by quotes)")
var dataFile = flag.String("data", "", "data file which is being encoded in carrier")
var resultFilesSlice sliceFlag
var resultFiles = flag.String("results", "", "names of the result files (separated by space and surrounded by quotes)")

func init() {
	flag.StringVar(carrierFiles, "c", "", "carrier files in which the data is encoded (separated by space surrounded by quotes, shorthand for --carriers)")
	flag.Var(&carrierFilesSlice, "carrier", "carrier file in which the data is encoded (could be used multiple times for multiple carriers)")
	flag.StringVar(dataFile, "d", "", "data file which is being encoded in carrier (shorthand for --data)")
	flag.Var(&resultFilesSlice, "result", "name of the result file (could be used multiple times for multiple result file names)")
	flag.StringVar(resultFiles, "r", "", "names of the result files (separated by space and surrounded by quotes, shorthand for --results)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: stegify [encode/decode] [flags...]")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stdout, `NOTE: When multiple carriers are provided with different kinds of flags, the names provided through "carrier" flag are taken first and with "carriers"/"c" flags second. Same goes for the result/results flag.`)
		fmt.Fprintln(os.Stdout, `NOTE: When no results are provided a default values will be used for the names of the results.`)
	}
}

func main() {
	operation := parseOperation()
	flag.Parse()
	carriers := parseCarriers()
	results := parseResults()

	if len(results) == 0 {
		if operation == encode {
			for i := range carriers {
				results = append(results, fmt.Sprintf("result%d", i))
			}
		} else {
			results = append(results, "result")
		}
	}

	if len(results) != len(carriers) && operation == encode {
		fmt.Fprintln(os.Stderr, "Carrier and result files count must be equal when encoding.")
		os.Exit(1)
	}

	if (dataFile == nil || *dataFile == "") && operation == encode {
		fmt.Fprintln(os.Stderr, "Data file must be specified. Use stegify --help for more information.")
		os.Exit(1)
	}

	switch operation {
	case encode:
		err := steg.MultiCarrierEncodeByFileNames(carriers, *dataFile, results)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case decode:
		if len(results) != 1 {
			fmt.Fprintln(os.Stderr, "Only one result file expected.")
			os.Exit(1)
		}
		err := steg.MultiCarrierDecodeByFileNames(carriers, results[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func parseOperation() string {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Operation must be specified [encode/decode]. Use stegify --help for more information.")
		os.Exit(1)
	}
	operation := os.Args[1]
	if operation != encode && operation != decode {
		helpFlags := map[string]bool{
			"--help": true,
			"-help":  true,
			"--h":    true,
			"-h":     true,
		}
		if helpFlags[operation] {
			flag.Parse()
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Unsupported operation: %s. Only [encode/decode] operations are supported.\n Use stegify --help for more information.", operation)
		os.Exit(1)
	}

	os.Args = append(os.Args[:1], os.Args[2:]...) // needed because go flags implementation stop parsing after first non-flag argument
	return operation
}

func parseCarriers() []string {
	carriers := make([]string, 0)
	if len(carrierFilesSlice) != 0 {
		carriers = append(carriers, carrierFilesSlice...)
	}

	if len(*carrierFiles) != 0 {
		carriers = append(carriers, strings.Split(*carrierFiles, " ")...)
	}

	if len(carriers) == 0 {
		fmt.Fprintln(os.Stderr, "Carrier file must be specified. Use stegify --help for more information.")
		os.Exit(1)
	}

	return carriers
}

func parseResults() []string {
	results := make([]string, 0)
	if len(resultFilesSlice) != 0 {
		results = append(results, resultFilesSlice...)
	}

	if len(*resultFiles) != 0 {
		results = append(results, strings.Split(*resultFiles, " ")...)
	}

	return results
}
