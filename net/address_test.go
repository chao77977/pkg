package net

import (
	"testing"
)

func TestValidIPAddress(t *testing.T) {
	addr := "0.0.0.255"
	if ValidIPAddress(addr) {
		t.Fatalf("Expected unvalid ip address: %s", addr)
	}
}
