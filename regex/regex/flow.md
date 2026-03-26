# Go `regexp` Library Cheat Sheet

> **ZomaSec** × **Cyber3rb | سايبر عرب**  
> Security Automation Course  
> 🌐 [zomasec.me](https://zomasec.me) · 🌐 [cyber3rb.com](https://cyber3rb.com)

---

## 1 · Import & Setup

```go
import "regexp"
```

There are **two ways** to compile a pattern:

```go
// Option 1 — returns an error (use in functions)
re, err := regexp.Compile(`\d+`)
if err != nil {
    log.Fatal(err)
}

// Option 2 — panics if pattern is invalid (use for package-level vars)
re := regexp.MustCompile(`\d+`)
```

> 💡 **Rule:** Compile your pattern **once** at the top of your file, not inside a loop.

```go
// ✅ Good — compiled once
var re = regexp.MustCompile(`\d+`)

func findNumbers(s string) []string {
    return re.FindAllString(s, -1)
}

// ❌ Bad — compiled on every call
func findNumbers(s string) []string {
    re := regexp.MustCompile(`\d+`)   // wasteful!
    return re.FindAllString(s, -1)
}
```

---

## 2 · Check if a Pattern Matches

```go
re := regexp.MustCompile(`go+gle`)

// Does the string contain a match?
matched := re.MatchString("I love google")
fmt.Println(matched) // true

// Quick one-liner (no pre-compile needed for simple checks)
matched, _ = regexp.MatchString(`go+gle`, "I love google")
fmt.Println(matched) // true
```

---

## 3 · Find the First Match

```go
re := regexp.MustCompile(`\d+`)

// Returns the first match as a string ("" if none)
result := re.FindString("order 42 and item 99")
fmt.Println(result) // "42"

// Returns the start and end index of the first match (nil if none)
loc := re.FindStringIndex("order 42 and item 99")
fmt.Println(loc) // [6 8]
```

---

## 4 · Find All Matches

```go
re := regexp.MustCompile(`\d+`)
input := "order 42 and item 99 and ref 7"

// Second argument = max results, use -1 for all
all := re.FindAllString(input, -1)
fmt.Println(all) // ["42" "99" "7"]

// Limit to first 2 results
two := re.FindAllString(input, 2)
fmt.Println(two) // ["42" "99"]
```

---

## 5 · Capture Groups ( )

Groups let you extract **parts** of a match.

```go
re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
input := "user@example.com"

// Returns slice: [full_match, group1, group2, group3]
match := re.FindStringSubmatch(input)
fmt.Println(match[0]) // "user@example.com"  ← full match
fmt.Println(match[1]) // "user"              ← group 1
fmt.Println(match[2]) // "example"           ← group 2
fmt.Println(match[3]) // "com"               ← group 3
```

### Find All Matches with Groups

```go
re := regexp.MustCompile(`(\w+)=(\w+)`)
input := "user=admin role=editor status=active"

all := re.FindAllStringSubmatch(input, -1)
for _, m := range all {
    fmt.Printf("key=%s  value=%s\n", m[1], m[2])
}
// key=user    value=admin
// key=role    value=editor
// key=status  value=active
```

---

## 6 · Named Capture Groups `(?P<name>...)`

Named groups make your code more readable.

```go
re := regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
input := "date: 2024-07-15"

match := re.FindStringSubmatch(input)

// Map group names to their values
result := map[string]string{}
for i, name := range re.SubexpNames() {
    if name != "" {
        result[name] = match[i]
    }
}

fmt.Println(result["year"])  // "2024"
fmt.Println(result["month"]) // "07"
fmt.Println(result["day"])   // "15"
```

---

## 7 · Replace Matches

```go
re := regexp.MustCompile(`\d+`)

// Replace ALL matches with a fixed string
result := re.ReplaceAllString("id=42 ref=99", "NUM")
fmt.Println(result) // "id=NUM ref=NUM"

// Use $1, $2 to reference capture groups in the replacement
re2 := regexp.MustCompile(`(\w+)=(\w+)`)
result2 := re2.ReplaceAllString("user=admin", "$1=[REDACTED]")
fmt.Println(result2) // "user=[REDACTED]"
```

### Replace with a Function

```go
re := regexp.MustCompile(`\d+`)

result := re.ReplaceAllStringFunc("id=42 ref=99", func(match string) string {
    n, _ := strconv.Atoi(match)
    return strconv.Itoa(n * 2) // double every number
})
fmt.Println(result) // "id=84 ref=198"
```

---

## 8 · Split a String

```go
re := regexp.MustCompile(`[\s,;]+`) // split on spaces, commas, or semicolons

parts := re.Split("one two,three;four  five", -1)
fmt.Println(parts) // ["one" "two" "three" "four" "five"]

// Limit number of pieces
parts2 := re.Split("one two,three", 2)
fmt.Println(parts2) // ["one" "two,three"]
```

---

## 9 · Common RE2 Syntax

| Pattern | What it matches |
|---|---|
| `.` | Any character except newline |
| `\d` | A digit `[0-9]` |
| `\D` | Non-digit |
| `\w` | Word character `[a-zA-Z0-9_]` |
| `\W` | Non-word character |
| `\s` | Whitespace (space, tab, newline) |
| `\S` | Non-whitespace |
| `[abc]` | One of: `a`, `b`, `c` |
| `[^abc]` | Anything except `a`, `b`, `c` |
| `[a-z]` | Any lowercase letter |
| `^` | Start of string |
| `$` | End of string |
| `*` | 0 or more (greedy) |
| `+` | 1 or more (greedy) |
| `?` | 0 or 1 |
| `{n}` | Exactly n times |
| `{n,m}` | Between n and m times |
| `*?` / `+?` | Non-greedy (match as little as possible) |
| `(abc)` | Capture group |
| `(?:abc)` | Group without capturing |
| `(?P<n>…)` | Named capture group |
| `a\|b` | Either `a` or `b` |
| `(?i)` | Case-insensitive mode |

---

## 10 · Useful Methods Summary

| Method | Description |
|---|---|
| `regexp.Compile(pat)` | Compile pattern, return error if invalid |
| `regexp.MustCompile(pat)` | Compile pattern, panic if invalid |
| `re.MatchString(s)` | `true` if any match found |
| `re.FindString(s)` | First match (empty string if none) |
| `re.FindStringIndex(s)` | `[start, end]` of first match |
| `re.FindAllString(s, n)` | All matches as `[]string` |
| `re.FindStringSubmatch(s)` | First match + capture groups |
| `re.FindAllStringSubmatch(s, n)` | All matches + capture groups |
| `re.ReplaceAllString(s, repl)` | Replace all matches |
| `re.ReplaceAllStringFunc(s, fn)` | Replace using a function |
| `re.Split(s, n)` | Split string by pattern |
| `re.SubexpNames()` | Names of all capture groups |
| `regexp.QuoteMeta(s)` | Escape a string to use literally in a pattern |

---

## 11 · Mistakes to Avoid

| ❌ Mistake | ✅ Fix |
|---|---|
| Compiling inside a loop/function | Declare `var re = regexp.MustCompile(...)` at package level |
| Not checking `FindStringSubmatch` result | Always check `if match == nil` before accessing `match[1]` |
| Using `FindString` when you need groups | Use `FindStringSubmatch` instead |
| Using user input directly in a pattern | Escape it first: `regexp.QuoteMeta(userInput)` |
| Forgetting case sensitivity | Add `(?i)` at the start of your pattern |

---

