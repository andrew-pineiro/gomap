package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"
)

type Arguments struct {
	IPAddress string
	Mask      string
	Ports     []string
}
type ValidIPs struct {
	IPAddress string
	Port      string
}

func verifyIP(ip string) (string, string) {

	if strings.Contains(ip, "/") {
		_, ipnet, err := net.ParseCIDR(ip)
		if err != nil {
			log.Fatalf("ERROR: not a valid IP address %s", ip)
		}
		return ipnet.IP.String(), ipnet.Mask.String()
	} else {
		ipnet, err := netip.ParseAddr(ip)
		if err != nil {
			log.Fatalf("ERROR: not a valid IP address %s", ip)
		}
		return ipnet.String(), ""
	}
}

func testIPAddress(ip string, port string) bool {
	timeout := time.Millisecond * 10
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func getArgs(args []string) Arguments {
	ip, mask := verifyIP(args[1])
	a := Arguments{
		IPAddress: ip,
		Mask:      mask,
	}
	if len(args) > 2 {
		if strings.Contains(args[2], "-") {
			a.Ports = strings.Split(args[2], "-")
		} else {
			a.Ports = append(a.Ports, args[2])
		}
	}
	return a
}

func main() {
	var args = getArgs(os.Args)

	if args.IPAddress != "" {
		var ips []ValidIPs
		for _, port := range args.Ports {
			if args.Mask != "" {
				//TODO(#1): add subnet mask usage
				log.Fatalf("subnet mask usage not implemented")
			} else {
				valid := testIPAddress(args.IPAddress, port)
				if valid {
					validIP := ValidIPs{
						IPAddress: args.IPAddress,
						Port:      port,
					}
					ips = append(ips, validIP)
				}
			}
		}
		if len(ips) == 0 {
			fmt.Println("no valid IPs found")
		}
		for _, ip := range ips {
			//TODO(#3): match valid ports with known ports
			fmt.Println("valid IP: ", ip.IPAddress, " Valid Port: ", ip.Port)
		}
	}

}
