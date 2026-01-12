package main

import (
	"strings"
	"flag"
	"time"
	"fmt"
)

func normalizeASN(asn string) string {
	asn = strings.TrimSpace(asn)
	asn = strings.ToUpper(asn)
	if strings.HasPrefix(asn, "AS") {
		asn = strings.TrimPrefix(asn, "AS")
	}
	return asn
}


func parseFlags() (asn, server string, stdin bool, timeout time.Duration) {
	flag.StringVar(&asn, "asn", "", "Autonomous System Number (e.g., AS12345)")
	flag.StringVar(&server, "server", defaultServer, "Whois server to query")
	flag.BoolVar(&stdin, "stdin", false, "Read ASN from standard input")
	flag.DurationVar(&timeout, "timeout", timeout, "Timeout for queries")
	flag.Parse()
	return
}

func collectASNs(asn string, stdin bool) []string {
	var asns []string

	if asn != "" {
		asns = append(asns, normalizeASN(asn))
	}

	if stdin {
		var input string
		for {
			_, err := fmt.Scanln(&input)
			if err != nil {
				break
			}
			asns = append(asns, normalizeASN(input))
		}
	}

	return asns
}	