package storage

import (
	"testing"
)

func TestFormatByte(t *testing.T) {
	var b float64 = 1 * 1024 * 1024 * 1024

	f := FormatByte(b)
	if f.Show() != "1.00 GiB" {
		t.Fatalf("Unexpected size %s != 1.00 GiB", f.Show())
	}
}
