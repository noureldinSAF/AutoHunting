# Reconnaissance Tool for Bug Bounty Hunters

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

## 🤝 Contributing

Contributions are welcome! If you'd like to add new enumeration modules, improve performance, or fix bugs, please open an issue or submit a pull request. Be sure to follow Go best practices (`go fmt`) and include tests where appropriate.

## 💄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
