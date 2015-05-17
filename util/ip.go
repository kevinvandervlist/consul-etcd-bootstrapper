package util

import (
	"net"
	"os"
)

func GetIP() (*net.IPAddr, error) {
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	ip, err := net.ResolveIPAddr("ip", hn)
	if err != nil {
		return nil, err
	}

	return ip, nil
}
