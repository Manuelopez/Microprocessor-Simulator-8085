package util

import (
	"fmt"
	"strconv"
)

const (
	HIGH_BITS = 0
	LOW_BITS  = 1
)

func DecimalToBinary(value int) [8]byte {
	ba := [8]byte{}

	bs := strconv.FormatUint(uint64(value), 2)
	if len([]rune(bs)) < 8 {
		for i := len([]rune(bs)); i < 8; i++ {
			bs = "0" + bs
		}
	}

	for i, c := range bs {
		if c == '0' {
			ba[i] = 0
		} else {
			ba[i] = 1
		}
	}

	return ba
}

func DecimalToBinary16(value int) ([8]byte, [8]byte) {

	hbits := [8]byte{}
	lbits := [8]byte{}
	bs := strconv.FormatUint(uint64(value), 2)
	if len([]rune(bs)) < 16 {
		for i := len([]rune(bs)); i < 8; i++ {
			bs = "0" + bs
		}
	}

	for i, c := range bs {
		if c == '0' {
			if i < 8 {

				hbits[i] = 0
			} else {
				lbits[i-1] = 0
			}
		} else {
			if i < 8 {

				hbits[i] = 1
			} else {
				lbits[i-8] = 1
			}

		}
	}

	return hbits, lbits

}

func BinaryToDecimal(value []byte) int {

	sum := int(0)
	for _, x := range value {
		sum = (sum * 2) + int(x)
	}

	return sum
}

func BinaryToHex(value [8]byte) (string, error) {
	sb := ""

	startAt := false

	for _, b := range value {
		if b == 1 {
			startAt = true
		}
		if startAt {

			sb = strconv.Itoa(int(b)) + sb
		}
	}

	ui, err := strconv.ParseUint(sb, 2, 64)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", ui), nil
}

func dToH(value uint64) string {
	return fmt.Sprintf("%x", value)
}
