package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

type Arguments struct {
	IPAddress string
	Mask      net.IPMask
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
func testPort(host string, ports []string) {
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
		}
		if conn != nil {
			defer conn.Close()
			fmt.Println("Opened", net.JoinHostPort(host, port))
		}
	}
}

func sendPing(ip string) {
	//TODO(#1): add subnet mask usage
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	//TODO(#2): handle not found ip addresses
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getArgs(args []string) Arguments {
	var a Arguments
	if verifyIP(args[1]) {
		a.IPAddress = args[1]
		sendPing(a.IPAddress)
	} else {
		fmt.Printf("ERROR: %s is not a valid IP address", args[1])
	}
	return a
}

func main() {
	var args = getArgs(os.Args)

	if args.IPAddress != "" {
		sendPing(args.IPAddress, args.Mask)
	}

}
