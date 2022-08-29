package net

import (
	"errors"
	"net"
	"strconv"
)

// Port represents network port
type Port uint16

func (p Port) String() string {
	return strconv.FormatInt(int64(p), 10)
}

// ReceiveFreePort receives free port
func ReceiveFreePort() (Port, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()

	return Port(l.Addr().(*net.TCPAddr).Port), nil
}

func ParsePort(s string) (Port, error) {
	p, err := strconv.Atoi(s)
	if err != nil {
		return Port(p), err
	}

	if p < 0 || p > 65535 {
		return Port(p), errors.New("port must be between 0 to 65535")
	}

	return Port(p), nil
}
