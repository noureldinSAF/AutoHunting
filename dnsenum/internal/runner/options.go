package runner

type Options struct {
	Targets             []string
	customResolvers     []string
	Records             []string
	TargetsFile         string
	CustomResolversFile string
	Strategy            string // fast/deep
	OutputFile          string
	OutputFormat         string // "json" (default) or "lines"
	Timeout             int
	Concurrency         int
	Silent              bool
}
