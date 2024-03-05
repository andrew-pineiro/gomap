package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
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

func checkPort(port string) string {
	portMatch := map[int]string{
		20:  "ftp data",
		21:  "ftp",
		22:  "ssh",
		23:  "telnet",
		135: "rpc",
		139: "smb",
		445: "smb",
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Panicf("ERROR: cannot convert port to int %d", intPort)
	}
	return portMatch[intPort]
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

		startRange, err := strconv.Atoi(args.Ports[0])
		if err != nil {
			log.Panicf("ERROR: issue converting %s to int. %s", args.Ports[0], err)
		}
		endRange, err := strconv.Atoi(args.Ports[1])
		if err != nil {
			log.Panicf("ERROR: issue converting %s to int. %s", args.Ports[1], err)
		}

		for i := startRange; i < endRange+1; i++ {
			if args.Mask != "" {
				//TODO(#1): add subnet mask usage
				log.Fatalf("subnet mask usage not implemented")
			} else {
				go testIPAddress(args.IPAddress, fmt.Sprint(i))
				time.Sleep(Delay * time.Millisecond)
			}
		}
		if len(IPs) == 0 {
			log.Println("no valid IPs found")
			return
		}
	}

}
