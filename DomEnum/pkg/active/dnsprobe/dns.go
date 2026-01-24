// dnsprobe/dnsprobe.go
package dnsprobe

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var resolvers = []string{
	"1.1.1.1:53",
	"8.8.8.8:53",
	"8.8.4.4:53",
}

func randomResolver() string {
	return resolvers[rand.Intn(len(resolvers))]
}

func CheckDomainsWithConcurrency(domains []string, concurrency int) []string {
	// de-dup to avoid repeat queries
	seen := make(map[string]struct{}, len(domains))
	unique := make([]string, 0, len(domains))
	for _, d := range domains {
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = struct{}{}
		unique = append(unique, d)
	}

	results := make([]string, 0)
	var resultsMu sync.Mutex

	jobs := make(chan string, concurrency*2)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for d := range jobs {
			if isAliveFast(d) {
				resultsMu.Lock()
				results = append(results, d)
				resultsMu.Unlock()
			}
		}
	}

	// spin up workers
	if concurrency < 1 {
		concurrency = 1
	}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	for _, d := range unique {
		jobs <- d
	}
	close(jobs)

	wg.Wait()
	return results
}

func isAliveFast(host string) bool {
	// quick reject of obviously invalid
	if host == "" {
		return false
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), dns.TypeA)
	m.RecursionDesired = true

	// fast UDP client
	udp := &dns.Client{
		Net:     "udp",
		Timeout: 1500 * time.Millisecond, // MUCH lower than 10s
	}

	resp, _, err := udp.Exchange(m, randomResolver())
	if err != nil || resp == nil {
		return false
	}

	// If truncated, retry over TCP once
	if resp.Truncated {
		tcp := &dns.Client{
			Net:     "tcp",
			Timeout: 2500 * time.Millisecond,
		}
		resp, _, err = tcp.Exchange(m, randomResolver())
		if err != nil || resp == nil {
			return false
		}
	}

	if resp.Rcode != dns.RcodeSuccess || len(resp.Answer) == 0 {
		return false
	}

	// Optional: accept CNAME-only answers too (some resolvers return it in Answer)
	_ = net.IP{} // keeps import if you later validate A records
	return true
}
