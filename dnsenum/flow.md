```mermaid
# DNS Enumeration Tool

Idea is convert this example.com to :
```
- ns.cloudflare -> NS Record
- test.github.io -> CNAME Record
- 1.2.3.4 -> A Record
- 2606:2800:220:1:248:1893:25c8:1946 → AAAA Record
- mail.example.com -> MX Record
- 1-2-3-4.test.net -> PTR Record
- xxxxxxxxxxxxxxxxxxxxx -> TXT Record
```
Normal Usage :

```go
dig example.com CNAME 
dig example.com A 
...
...

```
## 1. Parse User Input 

- User pass list of subdomains via file or stdout
- User can add list of public recursive resolvers to use
- user can specify records that he want to return (CNAME,MX,NS,A,AAAA,...)
- User specify output file (json)

## 2. Prepare the logic of runner
- Handle Runner with concurrency and rate limit and context and logs so on

## 3. DNS Probing 

- Collect list of recursive DNS resolvers
- Create DNS Client using `github.com/miekg/dns`
- Check if it wildcard or not
- Check each host with each existing resolver to avoid false positives
- If records found for an host then return and `continue` this host


## 4. Wildcard Flow

- Before probing subdomains, test for wildcard DNS records
- Generate random subdomain (e.g., `random123456789.example.com`)
- Query multiple random subdomains to detect wildcard patterns
- If wildcard detected:
  - Store wildcard response patterns (A, AAAA, CNAME records)
  - During enumeration, compare each subdomain's response with wildcard pattern
  - Only report subdomains that have unique records different from wildcard
  - Skip subdomains that match wildcard pattern to avoid false positives
- If no wildcard detected:
  - Proceed with normal enumeration
  - All valid DNS responses are considered legitimate

## 5. Record Collection & Validation

- For each subdomain, query all specified record types (A, AAAA, CNAME, MX, NS, TXT, PTR)
- Validate responses from multiple resolvers to ensure consistency
- If multiple resolvers return same result → high confidence
- If resolvers return different results → flag as potential issue
- Collect all unique records per subdomain per record type

| Status            | Exists? | Action in brute force   |
| ----------------- | ------- | ----------------------- |
| NOERROR + records | Yes     | Process fingerprint     |
| NOERROR + empty   | Maybe   | Try other record types  |
| NXDOMAIN          | No      | Safe negative           |
| SERVFAIL          | Unknown | Retry / change resolver |
| REFUSED           | Unknown | Change resolver         |
| FORMERR           | Error   | Fix client              |


## 6. Output Formatting

- Structure output as JSON with following format:
  ```json
  {
    "domain": "example.com",
    "subdomains": [
      {
        "host": "subdomain.example.com",
        "records": {
          "A": ["1.2.3.4", "5.6.7.8"],
          "AAAA": ["2606:2800:220:1:248:1893:25c8:1946"],
          "CNAME": ["target.example.com"],
          "MX": [{"priority": 10, "host": "mail.example.com"}],
          "NS": ["ns1.example.com", "ns2.example.com"],
          "TXT": ["v=spf1 include:_spf.example.com"],
          "PTR": ["1-2-3-4.test.net"]
        },
        "resolvers_checked": 3,
        "wildcard_filtered": false
      }
    ],
    "wildcard_detected": true,
    "wildcard_pattern": {
      "A": ["1.2.3.4"],
      "AAAA": ["2606:2800:220:1:248:1893:25c8:1946"]
    }
  }
  ```
- Write results to specified output file
- Optionally print progress/logs to stdout



## 7. Error Handling & Logging

- Handle DNS query timeouts gracefully
- Log resolver failures and continue with other resolvers
- Track and report statistics:
  - Total subdomains checked
  - Valid subdomains found
  - Wildcard-filtered subdomains
  - Resolver success/failure rates
- Handle context cancellation for graceful shutdown
- Rate limit DNS queries to avoid overwhelming resolvers

## 8. Implementation Checklist

- [ ] CLI argument parsing (subdomains file, resolvers, record types, output file)
- [ ] DNS resolver collection/validation
- [ ] Wildcard detection logic
- [ ] Concurrent DNS query runner with rate limiting
- [ ] Record type query handlers (A, AAAA, CNAME, MX, NS, TXT, PTR)
- [ ] Response validation across multiple resolvers
- [ ] Wildcard filtering logic
- [ ] JSON output formatting
- [ ] Error handling and logging
- [ ] Context support for cancellation
- [ ] Progress reporting
```