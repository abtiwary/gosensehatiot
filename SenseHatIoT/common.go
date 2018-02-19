package SenseHatIoT

/*
 * common.go
 * Common functions!
 *
 * Principal author(s) : Abhishek Tiwary
 *                       abhishek.tiwary@dolby.com
 *
 */

import (
	"net"
	"io"
	"os"
)

/*
 * Get the local IP address
 */
func GetLocalIP() string {
	addrs, erradd := net.InterfaceAddrs()
	if erradd != nil {
		io.WriteString(os.Stderr, erradd.Error())
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
