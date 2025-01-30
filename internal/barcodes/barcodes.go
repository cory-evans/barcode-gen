package barcodes

import (
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/qr"
)

// SplitBarcodePrefix takes a string and returns the prefix, number, and number of digits in the number.
func SplitBarcodePrefix(data string) (string, int, int) {

	var prefix string
	var number int
	var nDigits int

	dataLength := len(data)

	// enumerate the string backwards to find the first non-numeric character
	for i := dataLength - 1; i >= 0; i-- {
		if data[i] < '0' || data[i] > '9' {
			nDigits = dataLength - i - 1

			prefix = data[:dataLength-nDigits]
			number, _ = strconv.Atoi(data[dataLength-nDigits:])

			break
		}
	}

	return prefix, number, nDigits
}

func Generate(barcodeType string, data string, width, height int) (barcode.Barcode, error) {
	var bc barcode.Barcode
	var err error

	switch barcodeType {
	case "Code128":
		bc, err = code128.Encode(data)
	case "Datamatrix":
		bc, err = datamatrix.Encode(data)
	case "QR":
		bc, err = qr.Encode(data, qr.M, qr.Auto)
	}

	if err != nil {
		return nil, err
	}

	scaledBC, err := barcode.Scale(bc, width, height)
	if err != nil {
		return nil, err
	}

	return scaledBC, err
}
