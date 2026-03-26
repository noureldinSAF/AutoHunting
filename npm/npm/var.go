package main

import "regexp"

var secretPatterns = regexp.MustCompile(
    `(?i)` + // case insensitive
    
    // Optional prefixes (language-specific modifiers/keywords)
    `(?:` +
        `(?:const|let|var|export|public|private|protected|static|final|readonly|def)\s+|` + // JS/Java/C#/Python keywords
        `this\.|window\.|process\.env\.|` + // JS object properties
        `[@$]?` + // Ruby/PHP sigils
    `)?` +
    
    // Optional type declaration (Java/C#/C++/Go/TypeScript)
    `(?:[a-zA-Z_][a-zA-Z0-9_<>[\]]*\s+)?` +
    
    // Variable name (CAPTURE GROUP 1)
    // Handles: JS ($, _), PHP ($), Ruby (@, $), normal identifiers
    `([@$]?[a-zA-Z_][\w$]*)` +
    
    // Optional TypeScript type annotation
    `(?:\s*:\s*[a-zA-Z_][\w<>[\]|&\s]*)?` +
    
    // Assignment operator (=, :=, :, =>)
    `\s*(?::=|=>|=|:)\s*` +
    
    // Optional quote or brace before value
    `(?:["'\x60{])?` + // handles ", ', `, {
    
    // Value (CAPTURE GROUP 2) - non-greedy to stop at terminators
    `(.+?)` +
    
    // Optional closing quote or brace
    `(?:["'\x60}])?` +
    
    // Line terminators
    `(?:;|,|\n|$|//|#|/\*|}|]|\))`,
)

// FalsePositiveFilter - Regex pattern to exclude common false positives
var FalsePositiveFilter = regexp.MustCompile(
    `(?i)` +
    
    // 1. Function calls (any pattern with parentheses)
    `\w+\s*\([^)]*\)|` +
    
    // 2. Array access with brackets
    `\w+\s*\[[^\]]+\]|` +
    
    // 3. Mathematical operations
    `[+\-*/]\s*\w+|` +
    
    // 4. Comparisons and assignments to variables (not literals)
    `=\s*\w+\s*[;,]?\s*$|` +
    
    // 5. Object/property access chains
    `\w+\.\w+(?:\.\w+)*\s*[;,})\]]|` +
    
    // 6. String literals that are just variable names
    `=\s*["'][A-Z_]+["']\s*[;,]?\s*$|` +
    
    // 7. Imports and requires
    `require|import|export|__import|module\.exports|exports\.|` +
    
    // 8. Type checks and operators
    `typeof|instanceof|void\s+0|===|!==|\|\||&&|` +
    
    // 9. Documentation patterns
    `https?://|</?[\w]+|name=|description:|messageId:|` +
    
    // 10. Configuration/constant declarations
    `const\s+[A-Z_]+\s*=\s*["']?[A-Z_]+["']?\s*[;,]|` +
    
    // 11. Common non-secret keywords
    `null|undefined|default:|type:|return\s+|` +
    
    // 12. Empty or simple assignments
    `=\s*[\[\]{}]\s*[;,]?|` +
    
    // 13. Markdown/HTML content
    `^\s*`+"`"+`|^\s*a\s+name=|` +
    
    // 14. Template literals and interpolation
    `\$\{|\}|` +
    
    // 15. Ternary operators
    `\?\s*.*\s*:|` +
    
    // 16. New keyword (constructors)
    `new\s+\w+\(|` +
    
    // 17. Control flow
    `if\s*\(|for\s*\(|while\s*\(|switch\s*\(|case\s+|` +
    
    // 18. Spread and rest operators
    `\.\.\.|` +
    
    // 19. Regular expressions in code
    `new\s+RegExp|\/.*\/[gimuy]`,
)
