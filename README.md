# Automation Tools for Bug Bounty Hunters

Welcome to the Reconnaissance Tool repository – a modular and extensible reconnaissance framework designed for bug bounty hunters and security researchers. It helps automate discovery and mapping of your targets across domains, subdomains, ports, services, and more.

## 🛠️ Features

- Multi-threaded enumeration modules: subdomain discovery, port scanning, service fingerprinting, and HTTP probing.
- API integrations: harness third-party sources such as VirusTotal, Shodan, Censys, etc. for OSINT.
- Unified output: normalizes results from various enumerators into a consistent JSON format for easy parsing.
- Modular architecture: easily add your own enumerators or adjust concurrency settings.
- Extensible configuration: choose which modules to run and configure API keys via a YAML/JSON config file.

## 🚀 Getting started

1. **Clone the repository and build:**
   ```bash
   git clone https://github.com/yourusername/AutoHunting.git
   cd AutoHunting
   go build -o autohunt .
   ```
2. **Configure:**
   Create a config file (e.g., `config.yaml`) with your target domains, API keys, and module settings. See [`configs/example.yaml`](configs/example.yaml) for a template.
3. **Run a scan:**
   ```bash
   ./autohunt -config config.yaml
   ```
   Results will be saved in the `output/` directory.

## 🤖 Notes

This tool is built using patterns described in my [Go Syntax & Notes](../GoLearning/) repository. It demonstrates concurrency with goroutines, HTTP requests, and parsing logic tailored for recon tasks. Feel free to explore the code to learn how these patterns are applied in a real-world reconnaissance tool.

## 📘 Module Usage Details

### asn2cidr
Converts autonomous system numbers (ASNs) into CIDR ranges. Navigate to the `asn2cidr` directory and run:
```bash
./asnmap -asn AS32934   # ASN for Facebook
```

### cidr2ips
Transforms a list of CIDR blocks into individual IP addresses. Place your CIDRs in a `list.txt` file and run:
```bash
cat list.txt | go run .
```

### asn2cidr + cidr2ips
Combine both modules by piping output from `asn2cidr` into `cidr2ips`:
```bash
./asn2cidr/asnmap -asn AS32934 | go run .
```
Use `wc -l` to count the resulting IPs.

### Domain Enumeration (DomEnum)
Enumerate all domains related to a company by name (passive or active). From `DomEnum/cmd/DomEnum`:
```bash
go run . -h
go run . -q Swisscom -o swisscomDomains.txt                # passive
go run . -q Swisscom -o swisscomDomains.txt -active       # passive + active
go run . -q Swisscom -o swisscomDomains.txt -active -t 60 # set timeout in seconds
```
Passive enumeration uses three sources:
1. `crtsh` – no API key required.
2. `whoisfreaks` – free with API key; results are paginated (50 domains per page). To fetch all pages automatically, modify `ro.CurrentPage >= 1` to `ro.TotalPages >= ro.TotalPages` in `whoisfreaks.go`. A paid plan removes rate limits.
3. `whoisxmlapi` – paid.

APIs for all modules are configured in [`DomEnum/internal/config/config.yaml`](https://github.com/noureldinSAF/AutoHunting/tree/main/DomEnum/cmd/DomEnum/internal/config/config.yaml). Provide keys for `whoisfreaks` or `whoisxmlapi` to enable those sources.

### DnsEnum
Checks DNS records (A, AAAA, CNAME, etc.) for a target:
```bash
go run .
# or display help
go run ./main.go -h
```

### SubEnum
Enumerates subdomains. In `SubEnum/cmd/subenum`:
```bash
go run . -h
go run . -active -c 10 -i domains.txt -o subs.txt                      # active enumeration
go run . -active -c 20 -i domains.txt -o subs.txt -e -max-mutations-size 50 # limit mutation size
```
Active enumeration does not perform brute forcing.

### vhost (virtual host enumeration)
Determines which subdomains resolve to which IP addresses:
```bash
go run . -hosts subs.txt -output vhostedSubs
```

### portScanner
Performs TCP port scanning:
```bash
go run . -host-file subs.txt -host-threads 52 -threads 86 -output-file ports.txt -timeout 3
```

### URLEnum (passive and active)
Enumerates URLs from subdomains:
```bash
go run . -i subs.txt -o urls1.txt -pc 20 -ac 50 -timeout 400 -subs -active
```
Notes:
1. The `commoncrawl` source does not work in Codespaces; to use outside Codespaces remove the API key requirement in `commoncrawl.go` (`RequireAPIKey` should return `false`).
2. Specify a reasonable timeout for headless browser enumeration.
3. Increasing concurrency can reduce accuracy.
4. Passive enumeration takes ~2–3 minutes; active enumeration may take up to an hour.
5. Non‑script files (e.g., SVG, JPEG) are ignored automatically.
6. The tool returns unique URLs by default.

### JSAnalyzer
Analyzes JavaScript files and extracts secrets:
```bash
go run . -i js.txt -o output.json -timeout 600 -c 10 -only secrets
```
## 🤝 Contributing

Contributions are welcome! If you'd like to add new enumeration modules, improve performance, or fix bugs, please open an issue or submit a pull request. Be sure to follow Go best practices (`go fmt`) and include tests where appropriate.
