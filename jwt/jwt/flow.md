

## 1. Flow for JWT Secret Brute force

A **JWT MAC brute-forcer**: a tool that takes a JWT and a list of possible secrets, tries each secret to verify the token’s signature, and reports the secret if one matches.

**Goal for learners:** Understand how symmetric JWT signing works and why weak secrets make tokens easy to crack.

---

## 2. JWT structure (reminder)

Every JWT has three base64url-encoded parts, separated by dots:

```
header.payload.signature
```
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NzExMzQzMDksIm5hbWUiOiJ6b21hc2VjIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.87P0XJieN4OsGaLi-jsg_-Pi_qqU-NWoUowiNDdhE70
```
- **Header** — e.g. `{"alg":"HS256","typ":"JWT"}` (algorithm + type).
- **Payload** — claims (e.g. `sub`, `exp`, custom data).
- **Signature** — protects header + payload from tampering.

The **signature** is what we attack when we brute-force the secret.

---

## 3. Symmetric signing (HS256 / HS384 / HS512)

**Idea:** One **shared secret** is used both to **sign** and to **verify** the token.  
Algorithm family: **HMAC** (Hash-based Message Authentication Code).

| Algorithm | Construction   |
|----------|------------------|
| HS256    | HMAC + SHA-256   |
| HS384    | HMAC + SHA-384   |
| HS512    | HMAC + SHA-512   |

**How the signature is produced:**

1. Signing input:  
   `message = base64url(header) + "." + base64url(payload)`
2. Signature:  
   `signature = HMAC(secret, message)`  
   then base64url-encoded.

**How verification works (same secret):**

1. Compute:  
   `expected = HMAC(secret, message)`
2. Compare:  
   `expected == signature` → token is valid.

**Takeaway:** Anyone who knows the secret can create or verify tokens. If the secret is weak or guessable, an attacker can brute-force it and forge tokens.

---

## 4. Asymmetric JWTs (RS256, ECDSA) — why we don’t brute-force those

**Algorithms:** e.g. RS256 (RSA + SHA-256), ES256 (ECDSA + SHA-256).

**Idea:** Two different keys:

- **Private key** — used only to **sign**; must stay secret.
- **Public key** — used only to **verify**; can be shared (e.g. in JWKS).

**Verification:**  
`Verify(public_key, signature, message)` → valid or not.  
The verifier never needs the private key.

**Why no brute-force:**  
The secret we’d need is the **private key** (large random number / key pair). Guessing it is not feasible, so a wordlist or “weak secret” scan doesn’t apply. Our scanner is built for **symmetric (HMAC)** JWTs only.



## JWT Secrets Wordlist

[jwt-secrets](https://github.com/wallarm/jwt-secrets/blob/master/jwt.secrets.list)