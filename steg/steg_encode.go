package steg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" //register jpeg image format
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"github.com/DimitarPetrov/stegify/bits"
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

//MultiCarrierEncode performs steganography encoding of data Reader in equal pieces in each of the carriers
//and writes it to the result Writers encoded as PNG images.
func MultiCarrierEncode(carriers []io.Reader, data io.Reader, results []io.Writer) error {
	if len(carriers) != len(results) {
		return fmt.Errorf("different number of carriers and results")
	}

	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("error reading data %v", err)
	}

	chunkSize := len(dataBytes) / len(carriers)
	dataChunks := make([]io.Reader, 0, len(carriers))
	chunksCount := 0
	for i := 0; i < len(dataBytes) && chunksCount < len(carriers); i += chunkSize {
		chunksCount++
		if i+chunkSize >= len(dataBytes) || chunksCount == len(carriers) { // last iteration
			dataChunks = append(dataChunks, bytes.NewReader(dataBytes[i:]))
		}
		dataChunks = append(dataChunks, bytes.NewReader(dataBytes[i:i+chunkSize]))
	}

	for i := 0; i < len(carriers); i++ {
		if err := Encode(carriers[i], dataChunks[i], results[i]); err != nil {
			return fmt.Errorf("error encoding chunk with index %d: %v", i, err)
		}
	}
	return nil
}

//EncodeByFileNames performs steganography encoding of data file in carrier file
//and saves the steganography encoded product in new file.
func EncodeByFileNames(carrierFileName, dataFileName, resultFileName string) (err error) {
	return MultiCarrierEncodeByFileNames([]string{carrierFileName}, dataFileName, []string{resultFileName})
}

//MultiCarrierEncodeByFileNames performs steganography encoding of data file in equal pieces in each of the carrier files
//and saves the steganography encoded product in new set of result files.
func MultiCarrierEncodeByFileNames(carrierFileNames []string, dataFileName string, resultFileNames []string) (err error) {
	if len(carrierFileNames) == 0 {
		return fmt.Errorf("missing carriers names")
	}
	if len(carrierFileNames) != len(resultFileNames) {
		return fmt.Errorf("different number of carriers and results")
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

	data, err := os.Open(dataFileName)
	if err != nil {
		return fmt.Errorf("error opening data file %s: %v", dataFileName, err)
	}
	defer func() {
		closeErr := data.Close()
		if err == nil {
			err = closeErr
		}
	}()

	results := make([]io.Writer, 0, len(resultFileNames))
	for _, name := range resultFileNames {
		result, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("error creating result file %s: %v", name, err)
		}
		defer func() {
			closeErr := result.Close()
			if err == nil {
				err = closeErr
			}
		}()
		results = append(results, result)
	}

	err = MultiCarrierEncode(carriers, data, results)
	if err != nil {
		for _, name := range resultFileNames {
			_ = os.Remove(name)
		}
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
