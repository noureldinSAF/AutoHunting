package runner

type Options struct {
	Hosts          []string
	IPs            []string
	Ports          []int
	OutputFile     string
	HostsFile      string
	IPsFile        string
	Timeout        int
	Concurrency    int
	MaxTries       int
	Verbose        bool
	Silent         bool
	skipDNSResolve bool
}
