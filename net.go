package gotk

import (
	"fmt"
	"net"
	"strconv"
)

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetLocalIP() (addr string, err error) {
	var (
		ok    bool
		ipnet *net.IPNet
		addrs []net.Addr
	)

	if addrs, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}

	for _, v := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok = v.(*net.IPNet); !ok || ipnet.IP.IsLoopback() {
			continue
		}

		if ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("not found")
}

func GetIP(addrs ...string) (ip string, err error) {
	var (
		addr string
		conn net.Conn
	)

	if addr = "8.8.8.8:80"; len(addrs) > 0 {
		addr = addrs[0]
	}

	if conn, err = net.Dial("udp", addr); err != nil {
		return "", err
	}
	defer conn.Close()

	if ipAddr, ok := conn.LocalAddr().(*net.UDPAddr); !ok {
		return "", fmt.Errorf("failed to get ip")
	} else {
		return ipAddr.IP.String(), nil
	}
}

func PortFromAddr(addr string) (port int, err error) {
	var (
		v   int64
		tmp string
	)

	if _, tmp, err = net.SplitHostPort(addr); err != nil {
		return 0, err
	}

	if v, _ = strconv.ParseInt(tmp, 10, 64); v <= 0 {
		return 0, fmt.Errorf("invalid port")
	}

	return int(v), nil
}
