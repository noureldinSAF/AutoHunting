# Subdomain Enumeration Sources Documentation

This document provides detailed information about all 66 subdomain enumeration sources available in the tool. Sources are organized by whether they require an API key or not.

**Total Sources:** 66  
**Free Sources (No API Key):** 20  
**API Key Required:** 46

---

## Table of Contents

- [Free Sources (No API Key Required)](#free-sources-no-api-key-required)
- [API Key Required Sources](#api-key-required-sources)

---

## Free Sources (No API Key Required)

### 1. abuseipdb

**Description:** AbuseIPDB is a service that provides IP address reputation data. This source extracts subdomains from abuse reports and related data.

**Request Example:**
```bash
curl -X GET "https://www.abuseipdb.com/check/example.com" \
  -H "User-Agent: Mozilla/5.0"
```

**Response Type:** HTML  
**Parsing Method:** Regex extraction from HTML content

**Documentation:** https://www.abuseipdb.com/

---

### 2. alienvault

**Description:** AlienVault OTX (Open Threat Exchange) provides threat intelligence data including subdomain information.

**Request Example:**
```bash
curl -X GET "https://otx.alienvault.com/api/v1/indicators/domain/example.com/passive_dns" \
  -H "X-OTX-API-KEY: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "passive_dns": [
    {
      "hostname": "subdomain.example.com",
      "last": "2024-01-01T00:00:00"
    }
  ]
}
```

**Note:** This source requires an API key.

**Documentation:** https://otx.alienvault.com/api

---

### 3. anubis

**Description:** Anubis is a subdomain enumeration service that provides passive DNS data.

**Request Example:**
```bash
curl -X GET "https://jldc.me/anubis/subdomains/example.com"
```

**Response Type:** JSON/Text  
**Parsing Method:** Direct extraction from response

**Documentation:** https://github.com/jonluca/Anubis

---

### 4. commoncrawl

**Description:** Common Crawl is a massive web crawl dataset. This source searches through Common Crawl indexes to find subdomains.

**Request Example:**
```bash
curl -X GET "https://index.commoncrawl.org/CC-MAIN-2024-10-index?url=*.example.com&output=json"
```

**Response Type:** Text (CDX format)  
**Example Response:**
```
20240101000000 https://subdomain.example.com/path 200 ...
```

**Parsing Method:** Regex extraction from CDX text format

**Documentation:** https://commoncrawl.org/

---

### 5. crtsh

**Description:** crt.sh is a Certificate Transparency log search engine. It queries CT logs to find all certificates issued for a domain and its subdomains.

**Request Example:**
```bash
curl -X GET "https://crt.sh/?q=%25.example.com&output=json"
```

**Response Type:** JSON  
**Example Response:**
```json
[
  {
    "id": 123456,
    "name_value": "example.com\nwww.example.com\napi.example.com"
  }
]
```

**Parsing Method:** JSON parsing, extracts `name_value` field and splits by newlines. Wildcard entries (`*.domain.com`) are filtered out.

**Documentation:** https://crt.sh/

---

### 6. cyfare

**Description:** Cyfare provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://cyfare.com/subdomains/example.com"
```

**Response Type:** HTML/Text  
**Parsing Method:** HTML/text extraction

---

### 7. digitorus

**Description:** Digitorus (Certificate Details) provides certificate information including subdomains found in SSL certificates.

**Request Example:**
```bash
curl -X GET "https://certificatedetails.com/example.com"
```

**Response Type:** HTML  
**Parsing Method:** Regex-based subdomain extraction from HTML content using the utility extractor

**Note:** Even 404 pages may contain subdomain information.

**Documentation:** https://certificatedetails.com/

---

### 8. hackertarget

**Description:** HackerTarget provides various security tools including subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.hackertarget.com/hostsearch/?q=example.com"
```

**Response Type:** Text (CSV-like)  
**Example Response:**
```
subdomain1.example.com,192.168.1.1
subdomain2.example.com,192.168.1.2
```

**Parsing Method:** Line-by-line parsing, extracts subdomain from comma-separated format (subdomain,ip)

**Documentation:** https://hackertarget.com/hostsearch/

---

### 9. hudsonrock

**Description:** Hudson Rock provides threat intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://hudsonrock.com/api/subdomains/example.com"
```

**Response Type:** JSON/Text  
**Parsing Method:** Response parsing

---

### 10. myssl

**Description:** MySSL provides SSL certificate information and subdomain discovery.

**Request Example:**
```bash
curl -X GET "https://myssl.com/api/v1/discover_sub_domains?domain=example.com"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://myssl.com/

---

### 11. racent

**Description:** Racent provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://racent.com/subdomains/example.com"
```

**Response Type:** HTML/JSON  
**Parsing Method:** Response parsing

---

### 12. rapiddns

**Description:** RapidDNS provides DNS and subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://rapiddns.io/subdomain/example.com?page=1&full=1"
```

**Response Type:** HTML  
**Parsing Method:** Regex-based subdomain extraction from HTML pages. Supports pagination.

**Documentation:** https://rapiddns.io/

---

### 13. reconcloud

**Description:** ReconCloud provides reconnaissance data including subdomains.

**Request Example:**
```bash
curl -X GET "https://reconcloud.io/api/subdomains/example.com"
```

**Response Type:** JSON/Text  
**Parsing Method:** Response parsing

---

### 14. riddler

**Description:** Riddler.io provides subdomain enumeration and DNS data.

**Request Example:**
```bash
curl -X GET "https://riddler.io/search?q=pld:example.com&view_type=data_table"
```

**Response Type:** HTML/Text  
**Parsing Method:** Regex-based subdomain extraction from HTML/text content

**Documentation:** https://riddler.io/

---

### 15. shodanx

**Description:** ShodanX is a free alternative to Shodan that provides subdomain enumeration from Shodan's public data.

**Request Example:**
```bash
curl -X GET "https://www.shodan.io/domain/example.com"
```

**Response Type:** HTML  
**Parsing Method:** HTML parsing using golang.org/x/net/html. Extracts subdomains from `<ul id="subdomains">` list items.

**Documentation:** https://www.shodan.io/

---

### 16. shrewdeye

**Description:** ShrewdEye provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://shrewdeye.app/subdomains/example.com"
```

**Response Type:** HTML/JSON  
**Parsing Method:** Response parsing

---

### 17. sitedossier

**Description:** Site Dossier provides subdomain enumeration and website information.

**Request Example:**
```bash
curl -X GET "http://www.sitedossier.com/parentdomain/example.com"
```

**Response Type:** HTML  
**Parsing Method:** Regex-based subdomain extraction from HTML. Supports pagination through recursive enumeration.

**Documentation:** http://www.sitedossier.com/

---

### 18. threatcrowd

**Description:** ThreatCrowd is a free threat intelligence service that provides subdomain enumeration.

**Request Example:**
```bash
curl -X GET "http://ci-www.threatcrowd.org/searchApi/v2/domain/report/?domain=example.com"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "response_code": "1",
  "subdomains": [
    "subdomain1.example.com",
    "subdomain2.example.com"
  ],
  "undercount": "0"
}
```

**Parsing Method:** JSON parsing, extracts `subdomains` array

**Documentation:** https://www.threatcrowd.org/

---

### 19. threatminer

**Description:** ThreatMiner provides threat intelligence data including subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://www.threatminer.org/domain.php?q=example.com&rt=5"
```

**Response Type:** HTML/JSON  
**Parsing Method:** Response parsing

**Documentation:** https://www.threatminer.org/

---

### 20. urlscan

**Description:** urlscan.io provides URL scanning and subdomain discovery services.

**Request Example:**
```bash
curl -X GET "https://urlscan.io/api/v1/search/?q=domain:example.com"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "results": [
    {
      "page": {
        "domain": "subdomain.example.com",
        "url": "https://subdomain.example.com"
      }
    }
  ]
}
```

**Parsing Method:** JSON parsing, extracts domains from search results

**Documentation:** https://urlscan.io/docs/api/

---

### 21. waybackarchive

**Description:** Wayback Archive (Internet Archive) provides historical web data. This source searches through archived URLs to find subdomains.

**Request Example:**
```bash
curl -X GET "http://web.archive.org/cdx/search/cdx?url=*.example.com/*&output=txt&fl=original&collapse=urlkey"
```

**Response Type:** Text (CDX format)  
**Example Response:**
```
20240101000000 https://subdomain.example.com/path
20240101000001 https://www.example.com/page
```

**Parsing Method:** Regex-based subdomain extraction from CDX text format. URL-decodes content before extraction.

**Documentation:** https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server

---

## API Key Required Sources

### 22. alienvault

**Description:** AlienVault OTX (Open Threat Exchange) provides threat intelligence data including subdomain information.

**Request Example:**
```bash
curl -X GET "https://otx.alienvault.com/api/v1/indicators/domain/example.com/passive_dns" \
  -H "X-OTX-API-KEY: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "passive_dns": [
    {
      "hostname": "subdomain.example.com",
      "last": "2024-01-01T00:00:00"
    }
  ]
}
```

**Documentation:** https://otx.alienvault.com/api

---

### 23. arpsyndicate

**Description:** ARP Syndicate provides threat intelligence and subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.arpsyndicate.com/v1/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://arpsyndicate.com/

---

### 24. bevigil

**Description:** BeVigil is a mobile app security platform that provides subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://osint.bevigil.com/api/subdomains/example.com" \
  -H "X-Access-Token: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://bevigil.com/

---

### 25. binaryedge

**Description:** BinaryEdge provides threat intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.binaryedge.io/v2/query/domains/subdomain/example.com" \
  -H "X-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "events": [
    {
      "target": "subdomain.example.com"
    }
  ]
}
```

**Documentation:** https://docs.binaryedge.io/

---

### 26. bufferover

**Description:** BufferOver (Rapid7) provides DNS data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://dns.bufferover.run/dns?q=.example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://github.com/rapid7/sonar/wiki/API

---

### 27. builtwith

**Description:** BuiltWith provides technology profiling and subdomain discovery.

**Request Example:**
```bash
curl -X GET "https://api.builtwith.com/v14/api.json?KEY=YOUR_API_KEY&DOMAIN=example.com"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://builtwith.com/api

---

### 28. c99

**Description:** C99.nl provides various OSINT tools including subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.c99.nl/subdomainfinder?key=YOUR_API_KEY&domain=example.com"
```

**Response Type:** JSON/Text  
**Parsing Method:** Response parsing

**Documentation:** https://api.c99.nl/

---

### 29. censys

**Description:** Censys provides internet-wide scanning data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://search.censys.io/api/v2/hosts/search?q=services.tls.certificates.leaf_data.subject.common_name:example.com" \
  -u "YOUR_API_ID:YOUR_SECRET"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "result": {
    "hits": [
      {
        "services": [
          {
            "certificate": "subdomain.example.com"
          }
        ]
      }
    ]
  }
}
```

**Documentation:** https://search.censys.io/api

---

### 30. certspotter

**Description:** Cert Spotter monitors Certificate Transparency logs to discover subdomains.

**Request Example:**
```bash
curl -X GET "https://api.certspotter.com/v1/issuances?domain=example.com&include_subdomains=true&expand=dns_names" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "dns_names": [
    "example.com",
    "www.example.com",
    "api.example.com"
  ]
}
```

**Documentation:** https://sslmate.com/certspotter/api

---

### 31. chaos

**Description:** ProjectDiscovery Chaos provides subdomain enumeration from various sources.

**Request Example:**
```bash
curl -X GET "https://chaos.projectdiscovery.io/api/v1/domains/example.com" \
  -H "Authorization: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "subdomains": [
    "subdomain1.example.com",
    "subdomain2.example.com"
  ]
}
```

**Documentation:** https://chaos.projectdiscovery.io/

---

### 32. chinaz

**Description:** Chinaz provides DNS and subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.chinaz.com/subdomain?key=YOUR_API_KEY&domain=example.com"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 33. coderog

**Description:** CodeRog provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.coderog.io/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 34. digitalyama

**Description:** DigitalYama provides subdomain enumeration and DNS data.

**Request Example:**
```bash
curl -X GET "https://api.digitalyama.com/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 35. dnsdb

**Description:** DNSDB (Farsight Security) provides passive DNS data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.dnsdb.info/lookup/rrset/name/*.example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** Text (JSON Lines)  
**Example Response:**
```json
{"rrname":"subdomain.example.com","rrtype":"A","rdata":"192.168.1.1"}
```

**Documentation:** https://docs.dnsdb.info/

---

### 36. dnsdumpster

**Description:** DNS Dumpster provides DNS reconnaissance data including subdomains.

**Request Example:**
```bash
curl -X GET "https://api.dnsdumpster.com/domain/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "a": [
    {"host": "subdomain1.example.com"},
    {"host": "subdomain2.example.com"}
  ],
  "ns": [
    {"host": "ns1.example.com"}
  ]
}
```

**Documentation:** https://dnsdumpster.com/

---

### 37. dnsrepo

**Description:** DNSRepo provides DNS and subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://dnsrepo.noc.org/api/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 38. driftnet

**Description:** Driftnet provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.driftnet.io/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 39. facebook

**Description:** Facebook provides subdomain enumeration through their security research tools.

**Request Example:**
```bash
curl -X GET "https://developers.facebook.com/tools/debug/subdomains?domain=example.com" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://developers.facebook.com/

---

### 40. fofa

**Description:** FOFA is a cyberspace search engine that provides subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://fofa.info/api/v1/search/all?email=YOUR_EMAIL&key=YOUR_API_KEY&qbase64=ZG9tYWluPWV4YW1wbGUuY29t" \
  -H "Content-Type: application/json"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "results": [
    ["subdomain.example.com", "443", "https"]
  ]
}
```

**Documentation:** https://fofa.info/api

---

### 41. fullhunt

**Description:** FullHunt provides attack surface management and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://fullhunt.io/api/v1/domains/example.com/subdomains" \
  -H "X-API-KEY: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "hosts": [
    {
      "host": "subdomain.example.com"
    }
  ]
}
```

**Documentation:** https://api-docs.fullhunt.io/

---

### 42. github

**Description:** GitHub Code Search allows searching through public repositories for subdomains.

**Request Example:**
```bash
curl -X GET "https://api.github.com/search/code?q=example.com&per_page=100" \
  -H "Accept: application/vnd.github.v3.text-match+json" \
  -H "Authorization: token YOUR_GITHUB_TOKEN"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "total_count": 100,
  "items": [
    {
      "name": "config.js",
      "html_url": "https://github.com/user/repo/blob/main/config.js",
      "text_matches": [
        {
          "fragment": "api.example.com"
        }
      ]
    }
  ]
}
```

**Parsing Method:** 
1. Searches GitHub code for domain references
2. Downloads raw file content
3. Extracts subdomains using regex patterns

**Documentation:** https://docs.github.com/en/rest/search#search-code

---

### 43. gitlab

**Description:** GitLab provides code search capabilities for subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://gitlab.com/api/v4/search?scope=blobs&search=example.com" \
  -H "PRIVATE-TOKEN: YOUR_GITLAB_TOKEN"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing and content extraction

**Documentation:** https://docs.gitlab.com/ee/api/search.html

---

### 44. google

**Description:** Google Custom Search API can be used to find subdomains through web search.

**Request Example:**
```bash
curl -X GET "https://www.googleapis.com/customsearch/v1?key=YOUR_API_KEY&cx=YOUR_SEARCH_ENGINE_ID&q=site:example.com"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "items": [
    {
      "link": "https://subdomain.example.com/page"
    }
  ]
}
```

**Documentation:** https://developers.google.com/custom-search/v1/overview

---

### 45. huntermap

**Description:** HunterMap provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.huntermap.io/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 46. intelx

**Description:** IntelX provides threat intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://2.intelx.io/domain/search?k=YOUR_API_KEY&domain=example.com"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://intelx.io/

---

### 47. leakix

**Description:** LeakIX provides data leak intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://leakix.net/api/subdomains/example.com" \
  -H "api-key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://leakix.net/

---

### 48. merklemap

**Description:** MerkleMap provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.merklemap.com/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 49. netlas

**Description:** Netlas provides internet-wide scanning data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://app.netlas.io/api/domains/?q=example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "items": [
    {
      "data": {
        "domain": "subdomain.example.com"
      }
    }
  ]
}
```

**Documentation:** https://netlas.io/

---

### 50. odin

**Description:** Odin provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.odin.com/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 51. onyphe

**Description:** Onyphe provides cyber threat intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://www.onyphe.io/api/v2/simple/domain/example.com" \
  -H "Authorization: apikey YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "results": [
    {
      "subdomain": "subdomain.example.com"
    }
  ]
}
```

**Documentation:** https://www.onyphe.io/

---

### 52. pugrecon

**Description:** PugRecon provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.pugrecon.com/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 53. quake

**Description:** Quake (360) provides cyberspace search and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://quake.360.cn/api/v3/search/domain?query=example.com" \
  -H "X-QuakeToken: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://quake.360.cn/

---

### 54. rapidapi

**Description:** RapidAPI provides various APIs including subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://rapidapi.com/api/subdomains/example.com" \
  -H "X-RapidAPI-Key: YOUR_API_KEY" \
  -H "X-RapidAPI-Host: api-host"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://rapidapi.com/

---

### 55. rapidfinder

**Description:** RapidFinder provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.rapidfinder.io/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 56. rapidscan

**Description:** RapidScan provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.rapidscan.io/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 57. redhuntlabs

**Description:** RedHunt Labs provides attack surface management and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.redhuntlabs.com/v1/domains/example.com/subdomains" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://redhuntlabs.com/

---

### 58. robtex

**Description:** Robtex provides DNS and IP information including subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://freeapi.robtex.com/pdns/forward/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "rrname": "subdomain.example.com",
  "rrtype": "A",
  "rdata": "192.168.1.1"
}
```

**Documentation:** https://www.robtex.com/api/

---

### 59. rsecloud

**Description:** RSE Cloud provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.rsecloud.com/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 60. securitytrails

**Description:** SecurityTrails provides historical DNS data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.securitytrails.com/v1/domain/example.com/subdomains" \
  -H "APIKEY: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "subdomains": [
    "subdomain1.example.com",
    "subdomain2.example.com"
  ]
}
```

**Documentation:** https://securitytrails.com/corp/api

---

### 61. shodan

**Description:** Shodan is a search engine for internet-connected devices. It provides DNS data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.shodan.io/dns/domain/example.com?key=YOUR_API_KEY&page=1"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "domain": "example.com",
  "subdomains": [
    "subdomain1",
    "subdomain2"
  ],
  "result": 1,
  "more": false
}
```

**Parsing Method:** JSON parsing, constructs full subdomain names by appending to domain. Supports pagination.

**Documentation:** https://developer.shodan.io/api

---

### 62. threatbook

**Description:** ThreatBook provides threat intelligence and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.threatbook.cn/v3/domain/sub_domain?apikey=YOUR_API_KEY&resource=example.com"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://threatbook.com/

---

### 63. trickest

**Description:** Trickest provides workflow automation for security research including subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.trickest.com/v1/subdomains/example.com" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://trickest.com/

---

### 64. virustotal

**Description:** VirusTotal provides threat intelligence and subdomain enumeration through their API.

**Request Example:**
```bash
curl -X GET "https://www.virustotal.com/vtapi/v2/domain/report?apikey=YOUR_API_KEY&domain=example.com"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "subdomains": [
    "subdomain1.example.com",
    "subdomain2.example.com"
  ]
}
```

**Documentation:** https://developers.virustotal.com/reference

---

### 65. whoisxmlapi

**Description:** WhoisXML API provides WHOIS data and subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=YOUR_API_KEY&domainName=example.com&outputFormat=JSON"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

**Documentation:** https://whoisxmlapi.com/

---

### 66. windvane

**Description:** WindVane provides subdomain enumeration services.

**Request Example:**
```bash
curl -X GET "https://api.windvane.com/subdomains/example.com" \
  -H "X-API-Key: YOUR_API_KEY"
```

**Response Type:** JSON  
**Parsing Method:** JSON parsing

---

### 67. zoomeyeapi

**Description:** ZoomEye is a cyberspace search engine that provides subdomain enumeration.

**Request Example:**
```bash
curl -X GET "https://api.zoomeye.org/domain/search?q=example.com&page=1" \
  -H "API-KEY: YOUR_API_KEY"
```

**Response Type:** JSON  
**Example Response:**
```json
{
  "list": [
    {
      "name": "subdomain.example.com"
    }
  ]
}
```

**Documentation:** https://www.zoomeye.org/doc

---

## Response Type Summary

- **JSON:** 45 sources
- **HTML:** 12 sources  
- **Text/CSV:** 9 sources

## Common Parsing Methods

1. **JSON Parsing:** Direct extraction from JSON response fields
2. **Regex Extraction:** Pattern matching using regular expressions (via `utils.NewSubdomainExtractor`)
3. **HTML Parsing:** DOM traversal using `golang.org/x/net/html` or regex extraction
4. **Text Parsing:** Line-by-line or field-based extraction from text responses

## Notes

- All sources implement the `scraper.Source` interface
- Sources requiring API keys are only loaded if valid keys are provided in the configuration
- Free sources are always available
- The tool uses a unified subdomain extractor utility for HTML/text sources to ensure consistent parsing
- Rate limiting and error handling are implemented per source
- Some sources support pagination for large result sets

