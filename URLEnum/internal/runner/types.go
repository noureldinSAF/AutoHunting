package runner

type Options struct {
	Domain        string
	Timeout       int
	Concurrency   int
	Input         string
	Output        string
	ActiveEnabled bool
	queries []string
	IncludeSubdomains bool
}