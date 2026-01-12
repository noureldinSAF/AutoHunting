package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	defaultServer = "whois.radb.net:43"
	defaultTimeout       = 20 * time.Second
)

type ASNLookup struct {
	asn string
	server string
	timeout time.Duration
}

// constructor for ASNLookup
func NewASNLookup(asn string, server string, timeout time.Duration) *ASNLookup {
	asn = normalizeASN(asn)
	if server == "" {
		server = defaultServer
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &ASNLookup{
		asn: asn,
		server: server,
		timeout: timeout,
	}
}

func (l *ASNLookup) Query() ([]string, error) {
	resp, err := l.query()
	if err != nil {
		return nil, err
	}

	return parseCIDRs(resp), nil
}

func (l *ASNLookup) query() (string, error) {
	server := l.server

	conn, err := net.DialTimeout("tcp", server, l.timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(l.timeout))

	if _, err := fmt.Fprintf(conn, "-i origin AS%s\n", l.asn); err != nil {
		return "", err
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func parseCIDRs(whoisResp string) []string {
	var cidrs []string 

	for _, line := range strings.Split(whoisResp, "\n") {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "route:") {
			continue
		}

	    cidr := strings.TrimPrefix(line, "route:")

		cidr = strings.TrimSpace(cidr)
		cidrs = append(cidrs, cidr)
	}

	return cidrs
}

func main() {
	asn, server, stdin, timeout := parseFlags()

	asns := collectASNs(asn, stdin)

	if len(asns) == 0 {
		fmt.Fprintln(os.Stderr, "No ASN provided. Use -asn flag or -stdin to read from standard input.")
		os.Exit(1)
	}

	for _, asn := range asns {
		lookup := NewASNLookup(asn, server, timeout)
		cidrs, err := lookup.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error querying ASN %s: %v\n", asn, err)
			continue
		}
		
		for _, cidr := range cidrs {
			fmt.Fprintln(os.Stdout, cidr)
		}
	}	
}



	