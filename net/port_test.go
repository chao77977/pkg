package net

import (
	"testing"
)

func TestParsePort(t *testing.T) {
	p, err := ParsePort("1234")
	if err != nil {
		t.Fatal(err)
	}

	if p != 1234 {
		t.Fatalf("Expected port number: 1234, got: %d", p)
	}

	_, err = ParsePort("65536")
	if err == nil {
		t.Fatalf("don't get expected error")
	}
}
