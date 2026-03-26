# S3 Bucket Permission Discovery Flow

Simple steps to discover misconfigured S3 bucket permissions using unauthenticated HTTP requests.

## Overview

Make HTTP GET requests to S3 bucket URLs with specific query parameters. Status code `200` = publicly accessible (misconfigured).

---

## 1. List Objects (`LIST`)

**Endpoint:** `https://bucket.s3.amazonaws.com/?list-type=2`

**What it checks:** Can anonymous users list objects in the bucket?

**Status codes:**
- `200` ✅ **VULNERABLE** - Bucket allows public listing
- `403` ❌ Access denied (expected/secure)
- `404` ❌ Bucket doesn't exist

**Impact:** Attacker can enumerate all files/keys in the bucket.

**Demo**
https://s3.amazonaws.com/tripdata

---

## 2. Get ACL (`ACL`)

**Endpoint:** `https://bucket.s3.amazonaws.com/?acl`

**What it checks:** Can anonymous users read the bucket's Access Control List?

**Status codes:**
- `200` ✅ **VULNERABLE** - ACL is publicly readable
- `403` ❌ Access denied (expected/secure)

**Impact:** Attacker can see who has permissions on the bucket and potentially find more attack vectors.

** Demo **
https://unode1.s3.us-east-1.amazonaws.com
---

## 3. Get Policy (`POLICY`)

**Endpoint:** `https://bucket.s3.amazonaws.com/?policy`

**What it checks:** Can anonymous users read the bucket policy document?

**Status codes:**
- `200` ✅ **VULNERABLE** - Policy is publicly readable
- `403` ❌ Access denied (expected/secure)
- `404` ❌ No policy configured

**Impact:** Attacker learns exact permissions, conditions, and potential privilege escalation paths.

---

## 4. Get Public Access Block (`PUBLIC_ACCESS_BLOCK`)

**Endpoint:** `https://bucket.s3.amazonaws.com/?publicAccessBlock`

**What it checks:** Can anonymous users read the Public Access Block configuration?

**Status codes:**
- `200` ✅ **VULNERABLE** - Config is publicly readable
- `403` ❌ Access denied (expected/secure)

**Impact:** Attacker knows if bucket has public access restrictions enabled/disabled.

---

## Quick Summary

| Check | Query Parameter | Vulnerable Status | What Leaks |
|-------|----------------|-------------------|------------|
| LIST  | `?list-type=2` | 200 | All object keys |
| ACL   | `?acl` | 200 | Permission grants |
| POLICY | `?policy` | 200 | Full policy rules |
| PUBLIC_ACCESS_BLOCK | `?publicAccessBlock` | 200 | Block settings |

**Note:** All checks require `200 OK` response to confirm misconfiguration. Any `403 Forbidden` means the bucket is properly secured against anonymous access.
