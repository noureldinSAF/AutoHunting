package zonetransfer

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// lookupNS finds all nameservers for a domain
func lookupNS(domain string, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	nss, err := net.DefaultResolver.LookupNS(ctx, strings.TrimSuffix(domain, "."))
	if err != nil {
		return nil, err
	}

	var out []string
	for _, r := range nss {
		host := r.Host
		if !strings.HasSuffix(host, ".") {
			host += "."
		}
		out = append(out, host)
	}
	return out, nil
}

// resolveAorAAAA resolves a hostname to an IP address (prefers IPv4)
func resolveAorAAAA(host string, timeout time.Duration) (string, error) {
	name := strings.TrimSuffix(host, ".")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", name)
	if err != nil {
		return "", err
	}

	// Prefer IPv4
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	// Fallback to IPv6 if no IPv4
	if len(ips) > 0 {
		return ips[0].String(), nil
	}

	return "", fmt.Errorf("no IP address found for %s", host)
}

// attemptAXFR attempts zone transfer and returns subdomains if successful
func attemptAXFR(domain, nsHost string, port int, timeout time.Duration) ([]string, error) {
	ip, err := resolveAorAAAA(nsHost, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve NS %s: %v", nsHost, err)
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	m := new(dns.Msg)
	m.SetAxfr(domain)

	t := new(dns.Transfer)
	t.DialTimeout = timeout
	t.ReadTimeout = timeout
	t.WriteTimeout = timeout

	ch, err := t.In(m, addr)
	if err != nil {
		return nil, fmt.Errorf("AXFR failed: %v", err)
	}

	var subdomains []string
	count := 0
	maxRecords := 100000 // Safety limit

	for env := range ch {
		if env.Error != nil {
			return nil, fmt.Errorf("AXFR error: %v", env.Error)
		}

		for _, rr := range env.RR {
			count++
			if count > maxRecords {
				return subdomains, fmt.Errorf("safety limit reached: %d records", maxRecords)
			}

			// Extract subdomain from DNS record
			header := rr.Header()
			if header.Rrtype == dns.TypeA || header.Rrtype == dns.TypeAAAA || header.Rrtype == dns.TypeCNAME {
				name := strings.TrimSuffix(header.Name, ".")
				// Only include if it's a subdomain (not the root domain)
				if name != strings.TrimSuffix(domain, ".") && strings.HasSuffix(name, strings.TrimSuffix(domain, ".")) {
					subdomains = append(subdomains, name)
				}
			}
		}
	}

	return subdomains, nil
}

// IsVulnerable checks if a domain is vulnerable to zone transfer
func IsVulnerable(domain string, timeout int) bool {
	domainFQDN := dns.Fqdn(domain)
	timeoutDuration := time.Duration(timeout) * time.Second

	ns, err := lookupNS(domainFQDN, timeoutDuration)
	if err != nil || len(ns) == 0 {
		return false
	}

	// Try zone transfer on each nameserver
	for _, nsHost := range ns {
		_, err := attemptAXFR(domainFQDN, nsHost, 53, timeoutDuration)
		if err == nil {
			return true
		}
	}

	return false
}

// GetSubdomains attempts zone transfer and returns all discovered subdomains
func GetSubdomains(domain string, timeout int) ([]string, error) {
	domainFQDN := dns.Fqdn(domain)
	timeoutDuration := time.Duration(timeout) * time.Second

	ns, err := lookupNS(domainFQDN, timeoutDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup NS: %v", err)
	}
	if len(ns) == 0 {
		return nil, fmt.Errorf("no NS records found")
	}

	// Try each nameserver until one succeeds
	var allSubdomains []string
	seen := make(map[string]bool)

	for _, nsHost := range ns {
		subdomains, err := attemptAXFR(domainFQDN, nsHost, 53, timeoutDuration)
		if err != nil {
			continue // Try next nameserver
		}

		// Deduplicate subdomains
		for _, sub := range subdomains {
			if !seen[sub] {
				seen[sub] = true
				allSubdomains = append(allSubdomains, sub)
			}
		}

		// If we got results from one nameserver, return them
		if len(allSubdomains) > 0 {
			return allSubdomains, nil
		}
	}

	return nil, fmt.Errorf("zone transfer failed for all nameservers")
}

//
/*
key_10v7qLj3aQIkWd9PdJSKtjO.Pqf4WVBV2F7SQ3c67122mwyYNxgYCVziiPFNSYZYpvK

*/