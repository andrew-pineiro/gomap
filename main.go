package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const Delay = 2

type Arguments struct {
	IPAddress string
	Mask      string
	Ports     []string
}
type ValidIPs struct {
	IPAddress string
	Port      string
}

var IPs []ValidIPs

func portStringToInt(port string) int {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("ERROR: could not convert port %s to int: %s", port, err)
	}
	return portInt
}

func checkPort(port string) string {
	portMatch := map[int]string{
		20:  "ftp data",
		21:  "ftp",
		22:  "ssh",
		23:  "telnet",
		25:  "smtp",
		43:  "whois",
		53:  "dns",
		88:  "kerberos",
		109: "pop2",
		110: "pop3",
		123: "ntp",
		135: "rpc",
		137: "netbios",
		139: "smb",
		220: "imap",
		445: "smb-ad",
	}
	intPort := portStringToInt(port)
	val := portMatch[intPort]
	if len(val) <= 0 {
		val = "unknown"
	}
	return val
}
func verifyIP(ip string) (string, string) {
	if strings.Contains(ip, "/") {
		ip, ipnet, err := net.ParseCIDR(ip)
		if err != nil {
			log.Fatalf("ERROR: %s is not a valid IP address", ip)
		}
		return ip.String(), ipnet.Mask.String()
	} else {
		ip := net.ParseIP(ip)
		if ip == nil {
			log.Fatalf("ERROR: %s is not a valid IP address", ip)
		}
		return ip.String(), "32"
	}
}

func testIPAddress(ip string, port string) {
	timeout := time.Millisecond * Delay
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return
	}

	defer conn.Close()
	validIP := ValidIPs{
		IPAddress: ip,
		Port:      port,
	}
	IPs = append(IPs, validIP)
	log.Println("Valid IP:", validIP.IPAddress+";", "Valid Port:", validIP.Port, "("+checkPort(validIP.Port)+")")
}

func getArgs(args []string) Arguments {
	ip, mask := verifyIP(args[1])
	a := Arguments{
		IPAddress: ip,
		Mask:      mask,
	}

	if len(args) > 2 {
		if strings.Contains(args[2], "-") {
			ports := strings.Split(args[2], "-")
			startRange := portStringToInt(ports[0])
			endRange := portStringToInt(ports[1])
			for i := startRange; i < endRange+1; i++ {
				a.Ports = append(a.Ports, fmt.Sprint(i))
			}
		} else if strings.Contains(args[2], ",") {
			a.Ports = strings.Split(args[2], ",")
		} else {
			a.Ports = append(a.Ports, args[2])
		}
	}
	return a
}

func main() {
	var args = getArgs(os.Args)
	//TODO(#1): add subnet mask usage
	if args.IPAddress == "" {
		log.Fatalf("No ip address specified")
	}

	switch len(args.Ports) {
	case 1:
		// using goroutines with single port would cause inconsistent responses
		testIPAddress(args.IPAddress, fmt.Sprint(args.Ports[0]))
	default:
		for _, port := range args.Ports {
			go testIPAddress(args.IPAddress, fmt.Sprint(port))
			time.Sleep(Delay * time.Millisecond)
		}
	}

	if len(IPs) == 0 {
		log.Println("no valid IPs found")
		return
	}
}
