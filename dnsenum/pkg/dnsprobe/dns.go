package dnsprobe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cyinnove/tldify"
	"github.com/miekg/dns"

	"github.com/zomaxsec/dnsenum/pkg/client"
)

type ProbeResult struct {
	Host    string
	Results []*client.Result
}

func RunProbe(timeout int, host string, records []uint16, resolvers []string, strategy string) *ProbeResult {
	results := &ProbeResult{}

	switch strategy {
	case "fast":
		results = runFastDnsProbe(timeout, host, records, resolvers)

	case "deep":
		results = runDeepProbe(timeout, host, records, resolvers)

	default:
		results = runFastDnsProbe(timeout, host, records, resolvers)
	}

	return results
}

// Fast Strategy

func runFastDnsProbe(timeout int, host string, records []uint16, resolvers []string) *ProbeResult {
	results := []*client.Result{}

	if len(records) == 0 {
		records = []uint16{dns.TypeA}
	}
	if len(resolvers) == 0 {
		resolvers = client.DefaultResolvers
	}

	wildCardType, ok := detectWildcard(host)
	var wildcardHash string
	if ok && wildCardType != nil {
		wildcardHash = getHashForResults(wildCardType)
	}

resolversLabel:
	for _, r := range resolvers {
		foundAny := false
		for _, record := range records {
			result, err := client.Query(host, timeout, record, r)
			if err != nil || result == nil {
				continue
			}

			if result.DNStatus != "NOERROR" || len(result.Values) == 0 {
				continue
			}

			if !ok {
				results = append(results, result)
				foundAny = true
			} else if getHashForResults(result) != wildcardHash {
				results = append(results, result)
				foundAny = true
			}
		}
		if foundAny {
			break resolversLabel
		}
	}

	return &ProbeResult{
		Host:    host,
		Results: results,
	}
}

// Deep Strategy

func runDeepProbe(timeout int, host string, records []uint16, resolvers []string) *ProbeResult {
	if len(records) == 0 {
		records = []uint16{dns.TypeA}
	}
	if len(resolvers) == 0 {
		resolvers = client.DefaultResolvers
	}

	wildCardType, ok := detectWildcard(host)
	var wildcardHash string
	if ok && wildCardType != nil {
		wildcardHash = getHashForResults(wildCardType)
	}

	var results []*client.Result
	recordMap := make(map[string]map[string]bool)
	resolversChecked := 0

	for _, resolver := range resolvers {
		for _, record := range records {
			result, err := client.Query(host, timeout, record, resolver)
			if err != nil {
				resolversChecked++
				continue
			}
			if result == nil {
				resolversChecked++
				continue
			}

			resolversChecked++

			if result.DNStatus == "REFUSED" || result.DNStatus == "SERVFAIL" {
				continue
			}

			if result.DNStatus == "NOERROR" && len(result.Values) > 0 {
				if ok && wildCardType != nil && getHashForResults(result) == wildcardHash {
					continue
				}

				if recordMap[result.RecordName] == nil {
					recordMap[result.RecordName] = make(map[string]bool)
				}
				for _, value := range result.Values {
					recordMap[result.RecordName][value] = true
				}
			}
		}
	}

	for recordName, values := range recordMap {
		var valueList []string
		for value := range values {
			valueList = append(valueList, value)
		}
		results = append(results, &client.Result{
			RecordName: recordName,
			Values:     valueList,
			DNStatus:   "NOERROR",
		})
	}

	return &ProbeResult{
		Host:    host,
		Results: results,
	}
}
func detectWildcard(host string) (*client.Result, bool) {
	randomSub := time.Now().String()

	parsedDomain, err := tldify.Parse(host)
	if err != nil {
		return nil, false
	}

	fullDomain := fmt.Sprintf("%s.%s.%s", randomSub, parsedDomain.Domain, parsedDomain.TLD)

	res, err := client.Query(fullDomain, 10, dns.TypeA, randomResolver())
	if err != nil {
		return nil, false
	}

	if res != nil && res.DNStatus == "NOERROR" && len(res.Values) > 0 {
		return res, true
	}

	return nil, false
}

func randomResolver() string {
	return client.DefaultResolvers[rand.Intn(len(client.DefaultResolvers))]
}

func getHashForResults(result *client.Result) string {
	if result == nil {
		return ""
	}
	hash := sha256.New()
	hash.Write([]byte(result.RecordName))
	hash.Write([]byte(strings.Join(result.Values, ",")))

	return hex.EncodeToString(hash.Sum(nil))
}
