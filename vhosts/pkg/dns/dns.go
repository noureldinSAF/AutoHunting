package dns

import (
	"time"

	"github.com/miekg/dns"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

func ProbeHost(host string, timeout int, maxTries int) ([]string, error) {
	// Create DNS Resolver with default options
	dsnxOpts := dnsx.Options{
		BaseResolvers:     dnsx.DefaultResolvers,
		Timeout:           time.Duration(timeout) * time.Second,
		MaxRetries:        maxTries,
		TraceMaxRecursion: 255,
		QuestionTypes:     []uint16{dns.TypeA},
	}

	dnsClient, err := dnsx.New(dsnxOpts)
	if err != nil {
		return nil, err
	}

	// DNS A question and returns corresponding IPs
	result, err := dnsClient.Lookup(host)
	if err != nil {
		return nil, err
	}

	return result, err
}
