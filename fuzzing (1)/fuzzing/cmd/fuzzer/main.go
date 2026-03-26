package main

import (
	"flag"
	"fmt"
	"os"

	"fuzzing/internal/fuzzer"
)

var cfg = &fuzzer.Config{}

func main() {
	flag.StringVar(&cfg.Target, "u", "", "Base URL to fuzz (e.g. https://example.com)")
	flag.StringVar(&cfg.Wordlist, "w", "", "Path to wordlist file (one path per line)")
	flag.IntVar(&cfg.DelayMs, "d", 0, "Delay in ms between fuzz requests (0 = no delay)")
	flag.IntVar(&cfg.Timeout, "t", 15, "Timeout in seconds for each request")
	flag.IntVar(&cfg.BaselineDelayMs, "b", 100, "Delay in ms between the 3 baseline requests")
	flag.StringVar(&cfg.Method, "X", "GET", "HTTP method")
	flag.StringVar(&cfg.StatusFilterStr, "status", "", "Only show these status codes (comma-separated, e.g. 200,301,302). Empty = show all")
	flag.Parse()

	if cfg.Target == "" || cfg.Wordlist == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -u <url> -w <wordlist>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}


	if err := fuzzer.Run(*cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
