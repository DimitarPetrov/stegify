package steg

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	"io"
	"os"
	"stegify/bits"
	"strings"
)

const dataSizeReservedBytes = 20

func Encode(carrierFileName string, dataFileName string) error {
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

	dataBytes := make(chan byte)
	errChan := make(chan error)
	go readData(dataFile, dataBytes, errChan)

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	hasMoreBytes := true

	count := 0
	var dataCount uint32 = 0

	for x := 0; x < dx && hasMoreBytes; x++ {
		for y := 0; y < dy && hasMoreBytes; y++ {

			if count >= dataSizeReservedBytes {

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
			}
			count += 4
		}
	}

	select {
		case _, ok := <-dataBytes:
			if ok {
				return fmt.Errorf("data file too large for this carrier")
			}
		default:

	}

	setDataSizeHeader(RGBAImage, bytesOf(dataCount))

	resultFile, err := os.Create("steg_result" + carrierFileName[strings.LastIndex(carrierFileName, "."):])
	defer resultFile.Close()
	if err != nil {
		return fmt.Errorf("error creating result file: %v", err)
	}

	switch format {
		case "png" : png.Encode(resultFile, RGBAImage)
		case "jpeg" : jpeg.Encode(resultFile, RGBAImage, nil)
		case "gif" : gif.Encode(resultFile, RGBAImage, nil)
		default: return fmt.Errorf("unsupported format")
	}

	return nil
}

func bytesOf(counter uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, counter)
	quarters := make([]byte, 16)
	for i := 0; i < 16; i+= 4 {
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

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {

			if count >= (dataSizeReservedBytes / 4) * 3 {
				break
			}

			c := RGBAImage.RGBAAt(x,y)
			c.R = bits.SetLastTwoBits(c.R, dataCountBytes[count])
			c.G = bits.SetLastTwoBits(c.G, dataCountBytes[count+1])
			c.B = bits.SetLastTwoBits(c.B, dataCountBytes[count+2])


			count += 3

		}
	}

}

func setColorSegment(colorSegment *byte, data <-chan byte, errChan <-chan error) (hasMoreBytes bool, err error){

	select {
		case byte, ok := <-data:
			if !ok {
				return false, nil
			} else {
				*colorSegment = bits.SetLastTwoBits(*colorSegment,byte)
				return true, nil
			}

		case err := <-errChan:
			return false, err

	}

}

func readData(reader io.Reader, bytes chan<-byte, errChan chan<-error) {
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

	img, format ,err := image.Decode(reader)
	if err != nil {
		return nil,format, fmt.Errorf("error decoding carrier image: %v", err)
	}

	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	return RGBAImage, format, nil
}

