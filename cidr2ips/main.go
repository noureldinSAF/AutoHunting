package main

import (
	"log"
	"fmt"
	"net/netip"
	"os"
	"bufio"
)

func CIDRToIPs(cidr string) ([]string, error) {
	var ips []string

	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, err
	}

	addr := prefix.Addr()

    for prefix.Contains(addr) {
		ips = append(ips, addr.String())
		addr = addr.Next()
	}

	return ips, nil
}

// this function read the cidrs from list.txt
func readCIDRsFRomStdin() ([]string, error) {
	var cidrs []string
	
	// Create a scanner to read from Stdin (keyboard or piped file)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			cidrs = append(cidrs, line)
		}
	}

	// Check if there was an error during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cidrs, nil
}

func main() {

	var cidrs []string
	var err error 

	cidrs, err = readCIDRsFRomStdin()
	if err != nil {
		log.Fatalf("Error reading CIDRs from stdin: %v", err)
	}

	if len(cidrs) == 0 {
		log.Fatalf("No CIDR blocks provided")
	}

	for _, cidr := range cidrs {
		ips, err := CIDRToIPs(cidr)
		if err != nil {
			log.Printf("Error processing CIDR %s: %v", cidr, err)
			continue
		}

		for _, ip := range ips {
			fmt.Println(ip)
		}
	}
}




