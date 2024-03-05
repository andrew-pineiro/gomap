package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"
)

type Arguments struct {
	IPAddress string
	Mask      net.IPMask
	Ports     []string
}
type ValidIPs struct {
	IPAddress string
	Port      string
}

func verifyIP(ip string) bool {
	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		_, err := netip.ParseAddr(ip)
		if err != nil {
			//fmt.Printf("ERROR: %s is not a valid ip address\n", args[1])
			return false
		}
	}
	return true
}

func testIPAddress(ip string, port string) bool {
	//TODO(#1): add subnet mask usage
	//TODO(#2): handle not found ip addresses
	timeout := time.Millisecond * 10
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func getArgs(args []string) Arguments {
	var a Arguments
	if verifyIP(args[1]) {
		a.IPAddress = args[1]
		if len(args) > 2 {
			if strings.Contains(args[2], "-") {
				a.Ports = strings.Split(args[2], "-")
			} else {
				a.Ports = append(a.Ports, args[2])
			}
		}
	} else {
		fmt.Printf("ERROR: %s is not a valid IP address", args[1])
	}
	return a
}

func main() {
	var args = getArgs(os.Args)

	if args.IPAddress != "" {
		var ips []ValidIPs
		for _, port := range args.Ports {
			valid := testIPAddress(args.IPAddress, port)
			if valid {
				validIP := ValidIPs{
					IPAddress: args.IPAddress,
					Port:      port,
				}
				ips = append(ips, validIP)
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
