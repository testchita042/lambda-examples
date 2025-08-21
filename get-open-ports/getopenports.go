package getopenports

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Returns a list of available addresses in the format of []<ip4>:port
func GetOpenPorts() ([]string, error) {
	var openPorts []string
	var mutex sync.Mutex
	// Check all ports from 1 to 65535
	const maxPort = 65535
	const numWorkers = 1000 // Number of concurrent port scanners

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
			var wg sync.WaitGroup
			sem := make(chan bool, numWorkers)

			// Scan ports from 1 to maxPort
			for port := 1; port <= maxPort; port++ {
				wg.Add(1)
				sem <- true // Acquire semaphore
				go func(p int) {
					defer wg.Done()
					defer func() { <-sem }() // Release semaphore

					address := fmt.Sprintf("%s:%d", ipStr, p)
					conn, err := net.DialTimeout("tcp", address, time.Millisecond*200)
					if err == nil {
						mutex.Lock()
						openPorts = append(openPorts, address)
						mutex.Unlock()
						conn.Close()
					}
				}(port)
			}
			wg.Wait() // Wait for all port scans to complete
		}
	}
	return openPorts, nil
}
