package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const Delay = 2

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

func getPorts(port string) []string {
	var returnPorts []string

	if strings.Contains(port, "-") {
		ports := strings.Split(port, "-")
		startRange := portStringToInt(ports[0])
		endRange := portStringToInt(ports[1])
		for i := startRange; i < endRange+1; i++ {
			returnPorts = append(returnPorts, fmt.Sprint(i))
		}
	} else if strings.Contains(port, ",") {
		returnPorts = strings.Split(port, ",")
	} else {
		returnPorts = append(returnPorts, port)
	}

	return returnPorts
}

func main() {
	var host string
	var ports string

	//var args = getArgs(os.Args)
	flag.StringVar(&host, "host", "", "Host IP Address")
	flag.StringVar(&ports, "p", "", "port number or range")
	flag.Parse()

	portList := getPorts(ports)

	//TODO(#1): add subnet mask usage
	ip, _ := verifyIP(host)

	if ip == "" {
		log.Fatalf("No ip address specified")
	}

	switch len(portList) {
	case 1:
		// using goroutines with single port would cause inconsistent responses
		testIPAddress(ip, fmt.Sprint(portList[0]))
	default:
		for _, port := range portList {
			go testIPAddress(ip, fmt.Sprint(port))
			time.Sleep(Delay * time.Millisecond)
		}
	}

	if len(IPs) == 0 {
		log.Println("no valid IPs found")
		return
	}
}
