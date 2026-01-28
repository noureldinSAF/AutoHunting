package runner

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/cyinnove/logify"

	"github.com/zomaxsec/vhoster/pkg/dns"
	"github.com/zomaxsec/vhoster/pkg/http"
	"github.com/zomaxsec/vhoster/pkg/utils"
)

func Run(opts *Options) error {

	if opts.HostsFile != "" {
		var err error
		opts.Hosts, err = utils.ReadLines(opts.HostsFile)
		if err != nil {
			return err
		}

	}

	if opts.IPsFile != "" {
		var err error
		opts.IPs, err = utils.ReadLines(opts.IPsFile)
		if err != nil {
			return err
		}
		opts.skipDNSResolve = true
	}

	// map[IP][]Hosts - protected by mutex for concurrent access
	resultMap := map[string][]string{}
	var mu sync.Mutex

	// Create semaphore for concurrency control
	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = 10 // Default concurrency
	}
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	if !opts.skipDNSResolve {
		// Concurrent DNS resolution
		for _, host := range opts.Hosts {
			wg.Add(1)
			sem <- struct{}{} // Acquire semaphore

			go func(h string) {
				defer func() {
					<-sem // Release semaphore
					wg.Done()
				}()

				ips, err := dns.ProbeHost(h, opts.Timeout, opts.MaxTries)
				if err != nil {
					return
				}

				// Thread-safe append to resultMap
				mu.Lock()
				for _, ip := range ips {
					resultMap[ip] = append(resultMap[ip], h)
				}
				mu.Unlock()
			}(host)
		}

		wg.Wait()

	} else {
		// Concurrent vhost enumeration
		for _, ip := range opts.IPs {
			wg.Add(1)
			sem <- struct{}{} // Acquire semaphore

			go func(ipAddr string) {
				defer func() {
					<-sem // Release semaphore
					wg.Done()
				}()

				// Generate valid HTTP URL for this IP
				validURL := http.ProbeHTTP(ipAddr, opts.Timeout)
				if validURL == "" {
					return
				}

				invalidHost := fmt.Sprintf("%d-invalid.invalid", rand.Intn(1000000))

				baselineFingerprint, err := http.GetResponse(opts.Timeout, invalidHost, validURL)
				if err != nil || baselineFingerprint == nil {
					return
				}

				// Fuzz Host headers with actual domains (concurrent per IP)
				var hostWg sync.WaitGroup
				hostSem := make(chan struct{}, concurrency)

				for _, host := range opts.Hosts {
					hostWg.Add(1)
					hostSem <- struct{}{}

					go func(h string) {
						defer func() {
							<-hostSem
							hostWg.Done()
						}()

						// Send request with Host header set to the domain
						vhostResp, err := http.GetResponse(opts.Timeout, h, validURL)
						if err != nil {
							return
						}

						// Detect vhost: compare response to baseline
						if isVhost(baselineFingerprint, vhostResp) {
							mu.Lock()
							resultMap[ipAddr] = append(resultMap[ipAddr], h)
							mu.Unlock()
						}
					}(host)
				}

				hostWg.Wait()
			}(ip)
		}

		wg.Wait()
	}

	// Output results
	if err := outputResults(resultMap, opts); err != nil {
		return err
	}

	return nil
}

// isVhost returns true if status code or content length differ significantly
func isVhost(baseline, vhostResp *http.Response) bool {

	if baseline.StatusCode != vhostResp.StatusCode {
		return true
	}

	if baseline.ContentLength != vhostResp.ContentLength {
		return true
	}

	return false
}

// outputResults handles outputting results in JSON format to file and CLI format to console
func outputResults(resultMap map[string][]string, opts *Options) error {
	// Write JSON output to file
	if opts.OutputFile != "" {
		jsonData, err := json.MarshalIndent(resultMap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		outputFile := opts.OutputFile
		if !strings.HasSuffix(outputFile, ".json") {
			outputFile += ".json"
		}

		if err := os.WriteFile(outputFile, jsonData, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	if !opts.Silent {
		printCLIResults(resultMap)
	}

	return nil
}

func printCLIResults(resultMap map[string][]string) {
	if len(resultMap) == 0 {
		fmt.Println("No results found.")
		return
	}

	for ip, hosts := range resultMap {
		if len(hosts) > 0 {
			logify.Infof(ip)
			for _, host := range hosts {
				fmt.Printf("  - %s\n", host)
			}
			fmt.Println()
		}
	}
}
