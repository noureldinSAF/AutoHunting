package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cyinnove/logify"

	"github.com/zomaxsec/dnsenum/internal/runner"
)

func main() {
	var opts = &runner.Options{}

	flag.StringVar(&opts.TargetsFile, "l", "", "File containing list of targets (subdomains)")
	flag.StringVar(&opts.TargetsFile, "list", "", "File containing list of targets (subdomains)")
	flag.StringVar(&opts.CustomResolversFile, "r", "", "File containing custom DNS resolvers")
	flag.StringVar(&opts.CustomResolversFile, "resolvers", "", "File containing custom DNS resolvers")
	flag.StringVar(&opts.OutputFile, "o", "", "Output file for JSON results")
	flag.StringVar(&opts.OutputFile, "output", "", "Output file for JSON results")
	flag.StringVar(&opts.Strategy, "s", "fast", "Strategy: fast or deep")
	flag.StringVar(&opts.Strategy, "strategy", "fast", "Strategy: fast or deep")
	flag.IntVar(&opts.Timeout, "timeout", 5, "DNS query timeout in seconds")
	flag.IntVar(&opts.Concurrency, "c", 10, "Number of concurrent workers")
	flag.IntVar(&opts.Concurrency, "concurrency", 10, "Number of concurrent workers")
	flag.BoolVar(&opts.Silent, "silent", false, "Silent mode")

	var recordsFlag string
	flag.StringVar(&recordsFlag, "t", "A", "Comma-separated record types (A,AAAA,CNAME,MX,NS,TXT,PTR,SRV,SOA)")
	flag.StringVar(&recordsFlag, "type", "A", "Comma-separated record types (A,AAAA,CNAME,MX,NS,TXT,PTR,SRV,SOA)")

	flag.StringVar(&opts.OutputFormat, "of", "json", "Output format: json or lines")
    flag.StringVar(&opts.OutputFormat, "output-format", "json", "Output format: json or lines")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "DNS Enumeration Tool\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -l targets.txt -o results.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  echo 'subdomain.example.com' | %s -t A,AAAA -s deep\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -l targets.txt -r resolvers.txt -o output.json -c 20\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -l targets.txt -o results.json -of json\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "  %s -l targets.txt -o targets_out.txt -of lines\n", os.Args[0])

	}

	flag.Parse()	

	if opts.Silent {
		logify.MaxLevel = logify.Silent
	}

	if recordsFlag != "" {
		records := []string{}
		for _, r := range splitRecords(recordsFlag) {
			if r != "" {
				records = append(records, r)
			}
		}
		opts.Records = records
	}

	if opts.Strategy != "fast" && opts.Strategy != "deep" {
		logify.Fatalf("Invalid strategy: %s (must be 'fast' or 'deep')", opts.Strategy)
	}

	if opts.Concurrency < 1 {
		logify.Fatalf("Concurrency must be at least 1")
	}

	if opts.Timeout < 1 {
		logify.Fatalf("Timeout must be at least 1 second")
	}

	if opts.OutputFormat != "json" && opts.OutputFormat != "lines" {
    logify.Fatalf("Invalid output format: %s (must be 'json' or 'lines')", opts.OutputFormat)
     }

	if err := runner.Run(opts); err != nil {
		logify.Fatalf("Error: %v", err)
	}

	
}

func splitRecords(s string) []string {
	var result []string
	start := 0
	for i, char := range s {
		if char == ',' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
