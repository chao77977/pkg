package net

import (
	"net"
)

// validate an IP Address
func ValidIPAddress(addr string) bool {
	return net.ParseIP(addr) == nil
}
