package steg

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/DimitarPetrov/stegify/bits"
)

//Decode performs steganography decoding of Reader with previously encoded data by the Encode function and writes to result Writer.
func Decode(carrier io.Reader, result io.Writer) error {
	RGBAImage, _, err := getImageAsRGBA(carrier)
	if err != nil {
		return fmt.Errorf("error parsing carrier image: %v", err)
	}

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	dataBytes := make([]byte, 0, 2048)
	resultBytes := make([]byte, 0, 2048)

	dataCount := extractDataCount(RGBAImage)

	var count int

	for x := 0; x < dx && dataCount > 0; x++ {
		for y := 0; y < dy && dataCount > 0; y++ {
			if count >= dataSizeHeaderReservedBytes {
				c := RGBAImage.RGBAAt(x, y)
				dataBytes = append(dataBytes, bits.GetLastTwoBits(c.R), bits.GetLastTwoBits(c.G), bits.GetLastTwoBits(c.B))
				dataCount -= 3
			} else {
				count += 4
			}
		}
	}

	if dataCount < 0 {
		dataBytes = dataBytes[:len(dataBytes)+dataCount] //remove bytes that are not part of data and mistakenly added
	}

	dataBytes = align(dataBytes) // len(dataBytes) must be aliquot of 4

	for i := 0; i < len(dataBytes); i += 4 {
		resultBytes = append(resultBytes, bits.ConstructByteOfQuartersAsSlice(dataBytes[i:i+4]))
	}

	if _, err = result.Write(resultBytes); err != nil {
		return err
	}

	return nil
}

//MultiCarrierDecode performs steganography decoding of Readers with previously encoded data chunks by the MultiCarrierEncode function and writes to result Writer.
//NOTE: The order of the carriers MUST be the same as the one when encoding.
func MultiCarrierDecode(carriers []io.Reader, result io.Writer) error {
	for i := 0; i < len(carriers); i++ {
		if err := Decode(carriers[i], result); err != nil {
			return fmt.Errorf("error decoding chunk with index %d: %v", i, err)
		}
	}
	return nil
}

//DecodeByFileNames performs steganography decoding of data previously encoded by the Encode function.
//The data is decoded from file carrier and it is saved in separate new file
func DecodeByFileNames(carrierFileName string, resultName string) (err error) {
	return MultiCarrierDecodeByFileNames([]string{carrierFileName}, resultName)
}

//MultiCarrierDecodeByFileNames performs steganography decoding of data previously encoded by the MultiCarrierEncode function.
//The data is decoded from carrier files and it is saved in separate new file
//NOTE: The order of the carriers MUST be the same as the one when encoding.
func MultiCarrierDecodeByFileNames(carrierFileNames []string, resultName string) (err error) {
	if len(carrierFileNames) == 0 {
		return fmt.Errorf("missing carriers names")
	}

	carriers := make([]io.Reader, 0, len(carrierFileNames))
	for _, name := range carrierFileNames {
		carrier, err := os.Open(name)
		if err != nil {
			return fmt.Errorf("error opening carrier file %s: %v", name, err)
		}
		defer func() {
			closeErr := carrier.Close()
			if err == nil {
				err = closeErr
			}
		}()
		carriers = append(carriers, carrier)
	}

	result, err := os.Create(resultName)
	if err != nil {
		return fmt.Errorf("error creating result file: %v", err)
	}
	defer func() {
		closeErr := result.Close()
		if err == nil {
			err = closeErr
		}
	}()

	err = MultiCarrierDecode(carriers, result)
	if err != nil {
		_ = os.Remove(resultName)
	}
	return err
}

func align(dataBytes []byte) []byte {
	switch len(dataBytes) % 4 {
	case 1:
		dataBytes = append(dataBytes, byte(0), byte(0), byte(0))
	case 2:
		dataBytes = append(dataBytes, byte(0), byte(0))
	case 3:
		dataBytes = append(dataBytes, byte(0))
	}
	return dataBytes
}

func extractDataCount(RGBAImage *image.RGBA) int {
	dataCountBytes := make([]byte, 0, 16)

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	count := 0

	for x := 0; x < dx && count < dataSizeHeaderReservedBytes; x++ {
		for y := 0; y < dy && count < dataSizeHeaderReservedBytes; y++ {
			c := RGBAImage.RGBAAt(x, y)
			dataCountBytes = append(dataCountBytes, bits.GetLastTwoBits(c.R), bits.GetLastTwoBits(c.G), bits.GetLastTwoBits(c.B))
			count += 4
		}
	}

	dataCountBytes = append(dataCountBytes, byte(0))

	var bs = []byte{bits.ConstructByteOfQuartersAsSlice(dataCountBytes[:4]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[4:8]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[8:12]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[12:])}

	return int(binary.LittleEndian.Uint32(bs))
}
