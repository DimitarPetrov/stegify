package steg

import (
	"encoding/binary"
	"fmt"
	"github.com/DimitarPetrov/stegify/bits"
	"image"
	"image/draw"
	_ "image/jpeg" //register jpeg image format
	"image/png"
	"io"
	"os"
)

const dataSizeHeaderReservedBytes = 20 // 20 bytes results in 30 usable bits

//Encode performs steganography encoding of data Reader in carrier
//and writes it to the result Writer encoded as PNG image.
func Encode(carrier io.Reader, data io.Reader, result io.Writer) error {
	RGBAImage, format, err := getImageAsRGBA(carrier)
	if err != nil {
		return fmt.Errorf("error parsing carrier image: %v", err)
	}

	dataBytes := make(chan byte, 128)
	errChan := make(chan error)

	go readData(data, dataBytes, errChan)

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

	switch format {
	case "png", "jpeg":
		return png.Encode(result, RGBAImage)
	default:
		return fmt.Errorf("unsupported carrier format")
	}
}

//EncodeByFileNames performs steganography encoding of data file in carrier file
//and saves the steganography encoded product in new file.
func EncodeByFileNames(carrierFileName, dataFileName, resultFileName string) (err error) {
	carrier, err := os.Open(carrierFileName)
	if err != nil {
		return fmt.Errorf("error opening carrier file: %v", err)
	}
	defer func() {
		closeErr := carrier.Close()
		if err == nil {
			err = closeErr
		}
	}()

	data, err := os.Open(dataFileName)
	if err != nil {
		return fmt.Errorf("error opening data file: %v", err)
	}
	defer func() {
		closeErr := data.Close()
		if err == nil {
			err = closeErr
		}
	}()

	result, err := os.Create(resultFileName)
	if err != nil {
		return fmt.Errorf("error creating result file: %v", err)
	}
	defer func() {
		closeErr := result.Close()
		if err == nil {
			err = closeErr
		}
	}()

	err = Encode(carrier, data, result)
	if err != nil {
		_ = os.Remove(resultFileName)
	}
	return err
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
		}
		*colorSegment = bits.SetLastTwoBits(*colorSegment, byte)
		return true, nil

	case err := <-errChan:
		return false, err

	}
}

func readData(reader io.Reader, bytes chan<- byte, errChan chan<- error) {
	b := make([]byte, 1)
	for {
		if _, err := reader.Read(b); err != nil {
			if err == io.EOF {
				break
			}
			errChan <- fmt.Errorf("error reading data %v", err)
			return
		}
		for _, b := range bits.QuartersOfByte(b[0]) {
			bytes <- b
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
