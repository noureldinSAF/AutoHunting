# NPM Scanner – How It Works

Simple step-by-step flow of how the scanner finds secrets in npm packages.

---

## 1. Entry point

- You create a scan with **query** (package name or npm username), **scan type** (package or user), and optional **version/tag**.
- Calling **Run()** starts the flow and returns a **FinalResult** (list of unique secrets and where they were found).

---

## 2. Choose scan mode

- **NormalScan (package)**  
  Scan one package by name.  
  If a version/tag is given (e.g. `"latest"` or `"1.0.0"`), only that version is scanned; otherwise **all versions** of the package are scanned.

- **UserScan**  
  Look up all packages maintained by that npm user, then scan **every version** of each of those packages.

---

## 3. Get package metadata

- **Package scan:**  
  `GET https://registry.npmjs.org/<package>`  
  Response includes all **versions** and **dist-tags** (e.g. `latest` → version number).

- **User scan:**  
  `GET https://registry.npmjs.com/-/v1/search?text=maintainer:<user>&size=1000`  
  Response is the list of packages for that user.  
  Then for each package we fetch its metadata as above.

---

## 4. Decide which versions to scan

- If a **version or tag** was specified: resolve tag to version if needed, then scan only that version.
- If **no version** was specified: scan every version listed in the package metadata.

---

## 5. For each version: download the tarball

- From metadata we have a **tarball URL** per version.
- We **download** that `.tgz` into a temp file (no local `npm install`).

---

## 6. Unpack and walk the archive

- The `.tgz` is **gzip + tar**. We read it with a gzip reader and a tar reader.
- We iterate over each entry (file or directory) in the archive.
- **Skip** if the entry’s **full path** contains any blacklisted segment, e.g.:
  - `node_modules`
  - `package.json`
  - `package-lock.json`
- So we never scan inside `node_modules/` or those files.

---

## 7. Scan each non-skipped file for secrets

- Read the file content from the tar stream.
- Run **regex** on the content (`secretPatterns`) to find patterns that look like secrets (e.g. variable assignments with `key`, `password`, `token`, `apiKey`, etc.).
- All regex matches are collected for that file.

---

## 8. Filter matches (per file)

- Drop matches that are too short (&lt; 25 chars), too long (&gt; 1500), or empty.
- Run **FalsePositiveFilter** regex to drop common false positives (e.g. function calls, array access, URLs, placeholder text).
- Only the remaining matches are kept and attached to that file path.

---

## 9. Keep only versions that had secrets

- For each version we get a **ScanResult**: package name, version, and list of `{ path, secrets }`.
- We **only append** that version’s result if it has at least one secret (no “empty” results).

---

## 10. Aggregate by secret (deduplicate and add locations)

- We convert the list of **ScanResults** (per package/version/file) into a list of **SecretFindings** (per unique secret).
- For each **secret text** we collect every **location** where it was found:  
  `package@version/path/to/file`.
- Same secret in multiple files or versions appears **once**, with **all** locations in the **found** array.
- Output structure:  
  `results: [ { "secret": "...", "found": ["pkg@1.0.0/file.js", ...] }, ... ]`

---

## 11. Output

- **FinalResult** is JSON-serialized (e.g. via **GetJSONOutput()**).
- You get a single JSON object with a **results** array: each item is one unique secret and the list of places it was found.

---

## Flow summary

1. Create scan (query + type + optional version).
2. Run → choose NormalScan or UserScan.
3. Fetch metadata (package or user’s packages).
4. Determine versions to scan (one vs all).
5. For each version: download `.tgz` → unpack in memory.
6. For each file: skip if path is blacklisted → run secret regex → filter by length and false-positive regex.
7. Keep only versions that have at least one secret.
8. Group by secret text, collect all locations (package@version/path).
9. Return JSON: `results: [ { secret, found } ]`.
