package main

import (
	"flag"

	"github.com/cyinnove/logify"

	"github.com/zomaxsec/vhoster/internal/runner"
	"github.com/zomaxsec/vhoster/pkg/utils"
)

var opts = &runner.Options{}

func main() {

	var ports string

	flag.StringVar(&opts.HostsFile, "hosts", "", "File containing hosts to enumerate")
	flag.StringVar(&opts.IPsFile, "ips", "", "File containing IPs to enumerate (optional)")
	flag.StringVar(&ports, "ports", "", "Ports to enumerate (80,443,300,8080)")
	opts.Ports = utils.ParsePorts(ports)
	flag.IntVar(&opts.Timeout, "timeout", 10, "Timeout in seconds")
	flag.IntVar(&opts.Concurrency, "concurrency", 2, "Concurrency level")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&opts.Silent, "silent", false, "Silent output")
	flag.StringVar(&opts.OutputFile, "output", "", "Output file")
	flag.IntVar(&opts.MaxTries, "max-tries", 3, "Max number of times to try (default 3)")
	flag.Parse()

	if opts.HostsFile == "" {
		logify.Fatalf("No hosts file specified which is required")
	}

	if err := runner.Run(opts); err != nil {
		logify.Fatalf("Error: %v", err)
	}
}
