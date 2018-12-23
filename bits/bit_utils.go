//Package bits provides utils for manipulation bits of a single byte
package bits

const (
	firstQuarter  = 192
	secondQuarter = 48
	thirdQuarter  = 12
	fourthQuarter = 3
)

//QuartersOfByte returns a byte's four quarters of bits
func QuartersOfByte(b byte) [4]byte {
	return [4]byte{b & firstQuarter >> 6, b & secondQuarter >> 4, b & thirdQuarter >> 2, b & fourthQuarter}
}

func clearLastTwoBits(b byte) byte {
	return b & byte(252)
}

//SetLastTwoBits modifies last two bits of given byte
func SetLastTwoBits(b byte, value byte) byte {
	return clearLastTwoBits(b) | value
}
