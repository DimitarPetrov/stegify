package steg

import (
	"encoding/binary"
	"fmt"
	"github.com/DimitarPetrov/stegify/bits"
	"image"
	"image/draw"
	"image/png"
	_ "image/jpeg" //register jpeg image format
	"io"
	"os"
	"strings"
)

const dataSizeHeaderReservedBytes = 20 // 20 bytes results in 30 usable bits

//Encode performs steganography encoding of data file in carrier file
//and saves the steganography encoded product in new file.
func Encode(carrierFileName string, dataFileName string, newFileName string) error {
	carrier, err := os.Open(carrierFileName)
	defer carrier.Close()
	if err != nil {
		return fmt.Errorf("error opening carrier file: %v", err)
	}

	RGBAImage, format, err := getImageAsRGBA(carrier)
	if err != nil {
		return fmt.Errorf("error parsing carrier image: %v", err)
	}

	dataFile, err := os.Open(dataFileName)
	defer dataFile.Close()
	if err != nil {
		return fmt.Errorf("error opening data file: %v", err)
	}

	dataBytes := make(chan byte, 128)
	errChan := make(chan error)

	go readData(dataFile, dataBytes, errChan)

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	hasMoreBytes := true

	var count int
	var dataCount uint32

	for x := 0; x < dx && hasMoreBytes; x++ {
		for y := 0; y < dy && hasMoreBytes; y++ {

			if count >= dataSizeHeaderReservedBytes {

				c := RGBAImage.RGBAAt(x, y)

				hasMoreBytes, err = setColorSegment(&c.R, dataBytes, errChan)
				if err != nil {
					return err
				}
				if hasMoreBytes {
					dataCount++
				}
				hasMoreBytes, err = setColorSegment(&c.G, dataBytes, errChan)
				if err != nil {
					return err
				}
				if hasMoreBytes {
					dataCount++
				}
				hasMoreBytes, err = setColorSegment(&c.B, dataBytes, errChan)
				if err != nil {
					return err
				}
				if hasMoreBytes {
					dataCount++
				}
				RGBAImage.SetRGBA(x, y, c)
			} else {
				count += 4
			}
		}
	}

	select {
	case _, ok := <-dataBytes: // if there is more data
		if ok {
			return fmt.Errorf("data file too large for this carrier")
		}
	default:

	}

	setDataSizeHeader(RGBAImage, quartersOfBytesOf(dataCount))

	resultFile, err := os.Create(newFileName + carrierFileName[strings.LastIndex(carrierFileName, "."):])
	defer resultFile.Close()
	if err != nil {
		return fmt.Errorf("error creating result file: %v", err)
	}

	switch format {
	case "png", "jpeg":
		png.Encode(resultFile, RGBAImage)
	//case "gif" : gif.Encode(resultFile, RGBAImage, nil)
	default:
		return fmt.Errorf("unsupported carrier format")
	}

	return nil
}

func quartersOfBytesOf(counter uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, counter)
	quarters := make([]byte, 16)
	for i := 0; i < 16; i += 4 {
		quarters[i] = bits.QuartersOfByte(bs[i/4])[0]
		quarters[i+1] = bits.QuartersOfByte(bs[i/4])[1]
		quarters[i+2] = bits.QuartersOfByte(bs[i/4])[2]
		quarters[i+3] = bits.QuartersOfByte(bs[i/4])[3]
	}

	return quarters

}

func setDataSizeHeader(RGBAImage *image.RGBA, dataCountBytes []byte) {

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	count := 0

	for x := 0; x < dx && count < (dataSizeHeaderReservedBytes/4)*3; x++ {
		for y := 0; y < dy && count < (dataSizeHeaderReservedBytes/4)*3; y++ {

			c := RGBAImage.RGBAAt(x, y)
			c.R = bits.SetLastTwoBits(c.R, dataCountBytes[count])
			c.G = bits.SetLastTwoBits(c.G, dataCountBytes[count+1])
			c.B = bits.SetLastTwoBits(c.B, dataCountBytes[count+2])
			RGBAImage.SetRGBA(x, y, c)

			count += 3

		}
	}

}

func setColorSegment(colorSegment *byte, data <-chan byte, errChan <-chan error) (hasMoreBytes bool, err error) {

	select {
	case byte, ok := <-data:
		if !ok {
			return false, nil
		} else {
			*colorSegment = bits.SetLastTwoBits(*colorSegment, byte)
			return true, nil
		}

	case err := <-errChan:
		return false, err

	}

}

func readData(reader io.Reader, bytes chan<- byte, errChan chan<- error) {
	byte := make([]byte, 1)
	for {
		if _, err := reader.Read(byte); err != nil {
			if err == io.EOF {
				break
			}
			errChan <- fmt.Errorf("error reading data %v", err)
			return
		} else {
			for _, byte := range bits.QuartersOfByte(byte[0]) {
				bytes <- byte
			}
		}
	}
	close(bytes)
}

func getImageAsRGBA(reader io.Reader) (*image.RGBA, string, error) {

	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, format, fmt.Errorf("error decoding carrier image: %v", err)
	}

	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	return RGBAImage, format, nil
}
