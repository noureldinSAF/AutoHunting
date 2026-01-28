# Simple VHost Enumeration Plan

---

## 0) Inputs

* **Required:** `domains.txt` (subdomains / hostnames)
* **Optional:** `ips.txt` (IP addresses)

**Logic**

* If only domains → resolve IPs via DNS
* If only IPs → use them directly
* If both → merge DNS IPs + provided IPs

---

## 1) Normalize domains

* lowercase
* remove `*.` and duplicates
  Save as `domains.txt`.

---

## 2) DNS probe (auto IP discovery)

* Resolve A / AAAA records for `domains.txt`
* Collect unique IPs

Output:

* `resolved_domains.txt` (domain → IP)
* `final_ips.txt` (DNS IPs + optional `ips.txt`)

> Keep unresolved domains — they may still be valid vhosts.

---

## 3) Find live web targets

* Scan `final_ips.txt` for web ports:

   * 80, 443 (optionally 8080, 8443, etc.)
* Build targets like:

   * `http://IP:PORT`
   * `https://IP:PORT`

Save to `targets.txt`.

---

## 4) Build baseline per target

For each target:

* Send 2–3 requests with invalid `Host` headers
* Record:

   * status code
   * content length

---

## 5) Fuzz Host headers

For each target × domain:

* Send request to IP
* Set:

   * `Host: domain`
   * HTTPS: SNI = domain
* Capture same fields as baseline

---

## 6) Detect vhosts

Mark as **potential vhost** if response differs from baseline:

* status
* content length


Save:

```
domain -> IP:port -> evidence
```
