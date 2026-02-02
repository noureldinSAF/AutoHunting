## Example Flow Walkthrough

**Input:** `https://example.com`

1. **Initial Request**
   - `OnRequest` fires → Prints `https://example.com`
   - Not a static asset → Request continues
   - HTML fetched and parsed

2. **Link Discovery**
   - `OnHTML` finds: `<a href="/page">`, `<img src="/logo.png">`, `<link href="/style.css">`
   - All three URLs printed via `printUnique()`

3. **Crawling Decision**
   - `/page` → HTML → Will be crawled
   - `/logo.png` → Static → Request aborted
   - `/style.css` → Static → Request aborted

4. **Next Level**
   - `/page` is crawled
   - Links from `/page` discovered (depth 2)
   - Process repeats until max depth reached

5. **Completion**
   - All discovered URLs printed
   - Only HTML pages actually crawled
   - Static assets discovered but not downloaded

## Output Characteristics

✅ **What Gets Printed:**
- All discovered URLs (HTML pages + static assets)
- URLs with query parameters (preferred over versions without)
- Unique URLs only (deduplicated by path)

❌ **What Doesn't Get Printed:**
- Duplicate paths (unless params differ)
- URLs without params when a version with params exists

🚫 **What Doesn't Get Crawled:**
- Static assets (CSS, JS, images, fonts, etc.)
- External domains
- URLs beyond max depth (2 levels)
