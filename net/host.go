package net

import (
	"strings"
)

// host or host:port
func SplitHostPort(hostPort string) (host, port string) {
	host = hostPort
	colon := strings.LastIndexByte(host, ':')
	if colon != -1 {
		host, port = hostPort[:colon], hostPort[colon+1:]
	}

	return
}
