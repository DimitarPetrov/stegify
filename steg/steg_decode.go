package steg

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"
	"stegify/bits"
)

func Decode(carrierFileName string, newFileName string) error {
	carrier, err := os.Open(carrierFileName)
	defer carrier.Close()
	if err != nil {
		return fmt.Errorf("error opening carrier file: %v", err)
	}

	RGBAImage, _, err := getImageAsRGBA(carrier)
	if err != nil {
		return fmt.Errorf("error parsing carrier image: %v", err)
	}

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	dataBytes := make([]byte,0 , 2048)
	dataCount := extractDataCount(RGBAImage)


	return nil
}

func extractDataCount(RGBAImage *image.RGBA) uint32 {
	dataCountBytes := make([]byte, 0, 16)

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	count := 0
	hasMoreBytes := true

	for x := 0; x < dx && hasMoreBytes; x++ {
		for y := 0; y < dy && hasMoreBytes; y++ {

			c := RGBAImage.RGBAAt(x,y)

			if count < dataSizeReservedBytes {
				dataCountBytes = append(dataCountBytes, bits.GetLastTwoBits(c.R), bits.GetLastTwoBits(c.G), bits.GetLastTwoBits(c.B))
			} else {
				hasMoreBytes = false
			}

			count += 4
		}
	}

	dataCountBytes = append(dataCountBytes, byte(0))

	var bs = []byte{bits.ConstructByteOfQuartersAsSlice(dataCountBytes[:4]),
						   bits.ConstructByteOfQuartersAsSlice(dataCountBytes[4:8]),
						   bits.ConstructByteOfQuartersAsSlice(dataCountBytes[8:12]),
						   bits.ConstructByteOfQuartersAsSlice(dataCountBytes[12:])}

	return binary.LittleEndian.Uint32(bs)
}
