package barcodes_test

import (
	"testing"

	"github.com/cory-evans/barcode-gen/internal/barcodes"
)

type splitBarcodePrefixTest struct {
	input   string
	prefix  string
	number  int
	nDigits int
}

func TestSplitBarcodePrefix(t *testing.T) {

	var tests []splitBarcodePrefixTest = []splitBarcodePrefixTest{
		{"ABC123", "ABC", 123, 3},
		{"ABC1234", "ABC", 1234, 4},
		{"ABC00A00", "ABC00A", 0, 2},
		{"ABC00A12", "ABC00A", 12, 2},
		{"0000", "", 0, 4},
		{"ABC", "ABC", 0, 0},
	}

	for _, test := range tests {
		prefix, number, nDigits := barcodes.SplitBarcodePrefix(test.input)

		if prefix != test.prefix {
			t.Errorf("Test: %s, Expected prefix to be ABC, got %s", test.input, prefix)
		}

		if number != test.number {
			t.Errorf("Test: %s, Expected number to be 123, got %d", test.input, number)
		}

		if nDigits != test.nDigits {
			t.Errorf("Test: %s, Expected nDigits to be 3, got %d", test.input, nDigits)
		}
	}
}
