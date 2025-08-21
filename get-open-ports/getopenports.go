package getopenports

import (
	"fmt"
	"net"
	"time"
)

// Returns a list of available addresses in the format of []<ip4>:port
func GetOpenPorts() ([]string, error) {
	openPorts := make([]string, 0)
	// Common ports to check
	ports := []int{22, 80, 443, 8080, 3306, 5432}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.To4() == nil {
				continue // skip non-IPv4
			}
			ipStr := ip.String()
			for _, port := range ports {
				address := fmt.Sprintf("%s:%d", ipStr, port)
				conn, err := net.DialTimeout("tcp", address, time.Millisecond*200)
				if err == nil {
					openPorts = append(openPorts, address)
					conn.Close()
				}
			}
		}
	}
	return openPorts, nil
}
