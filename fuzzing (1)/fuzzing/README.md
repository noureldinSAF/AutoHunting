# HTTP Fuzzer

Basic intelligent HTTP path fuzzer (Go): takes a baseline fingerprint, then flags responses that differ from it.

## Structure (3 dirs)

```
fuzzing/
├── cmd/
│   └── fuzzer/          # CLI entrypoint
│       └── main.go
├── internal/
│   ├── fingerprint/     # Baseline capture & comparison
│   │   └── fingerprint.go
│   └── fuzzer/          # Wordlist, requests, compare, report
│       └── fuzzer.go
├── go.mod
├── flow.md
├── wordlist.txt         # Example wordlist
└── README.md
```

- **cmd/fuzzer** – Parses flags (`-u`, `-w`, `-d`, `-X`) and runs the fuzzer.
- **internal/fingerprint** – Sends a request to a known non-existent path, records status/length/headers, and exposes `Equal()` for comparison.
- **internal/fuzzer** – Loads wordlist, captures baseline, sends one request per path, compares to baseline, prints interesting findings.

## Build & run

```bash
go build -o fuzzer ./cmd/fuzzer/

./fuzzer -u https://example.com -w wordlist.txt
./fuzzer -u https://example.com -w wordlist.txt -d 100 -X GET
```

## Flags

| Flag | Description |
|------|-------------|
| `-u` | Base URL to fuzz (required) |
| `-w` | Path to wordlist file, one path per line (required) |
| `-d` | Delay in ms between requests (default 0) |
| `-X` | HTTP method (default GET) |

## Example output

```
[*] Baseline: 404, 1234 bytes
[+] 200 5678 1 /admin
[+] 301 0 2 /login

[*] Done. 2 interesting, 9 total

--- Matchs ---
  /admin -> 200 (5678 bytes)
  /login -> 301 (0 bytes)
```
