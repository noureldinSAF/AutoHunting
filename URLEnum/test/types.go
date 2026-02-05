package test

type Options struct {
	Domain             string
	Timeout            int
	Input              string
	Output             string
	ActiveEnabled      bool
	IncludeSubdomains  bool

	// NEW: split concurrency
	PassiveConcurrency int
	ActiveConcurrency  int

	queries []string
}