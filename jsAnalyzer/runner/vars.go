package runner

import "regexp"

// RE2 (Go regexp) does NOT support lookbehind/lookahead.
// So patterns must avoid (?<=...), (?<!...), etc.

var (
	// SubdomainRegex matches domains/subdomains appearing in text.
	SubdomainRegex = regexp.MustCompile(
		`(?i)\b(?:` +
			`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+` +
			`[a-z]{2,63}` +
			`)\b`,
	)

	// CloudBucketRegex matches common bucket/object storage URL styles (S3/GCS/Azure).
	CloudBucketRegex = regexp.MustCompile(
		`(?i)\b(?:` +
			// AWS S3
			`(?:s3://([a-z0-9][a-z0-9.-]{1,61}[a-z0-9])(?:/[^\s"'<>]*)?)` +
			`|(?:https?://([a-z0-9][a-z0-9.-]{1,61}[a-z0-9])\.(?:s3(?:[.-][a-z0-9-]+)?|s3-website(?:[.-][a-z0-9-]+)?)\.amazonaws\.com(?:/[^\s"'<>]*)?)` +
			`|(?:https?://(?:s3(?:[.-][a-z0-9-]+)?|s3-website(?:[.-][a-z0-9-]+)?)\.amazonaws\.com/([a-z0-9][a-z0-9.-]{1,61}[a-z0-9])(?:/[^\s"'<>]*)?)` +
			// GCS
			`|(?:gs://([a-z0-9][a-z0-9._-]{1,221}[a-z0-9])(?:/[^\s"'<>]*)?)` +
			`|(?:https?://(?:storage\.googleapis\.com|storage\.cloud\.google\.com)/([a-z0-9][a-z0-9._-]{1,221}[a-z0-9])(?:/[^\s"'<>]*)?)` +
			// Azure Blob Storage
			`|(?:https?://([a-z0-9]{3,24})\.blob\.core\.windows\.net/([a-z0-9](?:[a-z0-9-]{1,61}[a-z0-9])?)(?:/[^\s"'<>]*)?)` +
			`)\b`,
	)

	// EndpointRegex matches:
	// - http(s)/ws(s) URLs
	// - scheme-relative //example.com/...
	// - common absolute API paths like "/api/...", "/v1/...", "/graphql", etc.
	//
	// IMPORTANT: No lookbehind. For path endpoints we use a "left boundary"
	// group: (^|[^A-Za-z0-9_]) then the path. When extracting, take submatch[2].
	EndpointRegex = regexp.MustCompile(
		`(?i)(?:` +
			// Full URLs
			`\b(?:https?|wss?)://` +
			`[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?` +
			`(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*` +
			`(?::\d{2,5})?` +
			`(?:/[^\s"'<>)]*)?` +
			`|` +
			// Scheme-relative URLs
			`\b//` +
			`[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?` +
			`(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*` +
			`(?::\d{2,5})?` +
			`(?:/[^\s"'<>)]*)?` +
			`|` +
			// Common API-like absolute paths (capture group 2)
			`(^|[^a-z0-9_])` +
			`(` +
			`/(?:api|apis|v\d+|graphql|rest|rpc|oauth2?|auth|login|logout|token|sessions?)` +
			`(?:/[a-z0-9._~!$&'()*+,;=:@%-]+)*` +
			`)` +
			`)`,
	)

	// ParameterRegex matches query parameters and code/config style key-value.
	ParameterRegex = regexp.MustCompile(
		`(?i)(?:` +
			`[?&]([a-z_][a-z0-9_.-]{0,127})=([^&#\s"'<>]*)` +
			`|` +
			`\b([a-z_][a-z0-9_.-]{0,127})\s*(?:=|:)\s*(?:` +
			`"([^"\r\n]{0,500})"` +
			`|'([^'\r\n]{0,500})'` +
			`|([^\s,;}\]\r\n]{1,500})` +
			`)` +
			`)`,
	)

	// NodeModulesRegex matches node_modules paths (POSIX + Windows).
	NodeModulesRegex = regexp.MustCompile(
		`(?i)(?:^|[\s"'(=,:])(?:` +
			`(?:[a-z]:\\(?:[^\\\r\n]+\\)*)?node_modules\\[^ \r\n"'()<>]+` +
			`|` +
			`(?:\./|\../|/)?node_modules/[^ \r\n"'()<>]+` +
			`)(?:$|[\s"')\],;<>])`,
	)

	// PackageNameRegex matches npm package names (scoped or unscoped).
	PackageNameRegex = regexp.MustCompile(
		`(?i)\b(?:@([a-z0-9][a-z0-9-._]{0,50})/)?([a-z0-9][a-z0-9-._]{0,213})\b`,
	)
)
