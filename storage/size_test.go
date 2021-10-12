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

func TestFormat(t *testing.T) {
	x, err := Format(1.01234567890, "GiB")
	if err != nil {
		t.Fatal(err)
	}

	if x.Unit() != "GiB" {
		t.Fatalf("Expected unit: GiB, got: %s", x.Unit())
	}

	size := x.Truncate(5)
	if size != 1.01234 {
		t.Fatalf("Expected size: 1.01234, got: %.5f", size)
	}

	size = x.Round(5)
	if size != 1.01235 {
		t.Fatalf("Expected size: 1.01235, got: %.5f", size)
	}

	_, unit, err := x.Convert("B", 128, true)
	if err != nil {
		t.Fatal(err)
	}

	if unit != "B" {
		t.Fatalf("Expected unit: B, got: %s", unit)
	}

	_, _, err = x.Convert("OiB", 128, false)
	if err == nil {
		t.Fatalf("Unexpected size unit: OiB")
	}
}

func TestCompare(t *testing.T) {
	x, err := Format(2, "GiB")
	if err != nil {
		t.Fatal(err)
	}

	y, err := Format(4, "GiB")
	if err != nil {
		t.Fatal(err)
	}

	if x.Compare(y) != -1 {
		t.Fatalf("Unexpected compare %s > %s", x.Show(), y.Show())
	}

	z, err := Format(y.Sub(x))
	if err != nil {
		t.Fatal(err)
	}

	if x.Compare(z) != 0 {
		t.Fatalf("Unexpected compare %s != %s", x.Show(), z.Show())
	}

	x, err = Format(y.Add(z))
	if err != nil {
		t.Fatal(err)
	}

	if x.Compare(y) != 1 {
		t.Fatalf("Unexpected compare %s < %s", x.Show(), y.Show())
	}
}
