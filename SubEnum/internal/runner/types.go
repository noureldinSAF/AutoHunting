package runner

type Options struct {
	Query            string
	queries          []string
	InputFile        string
	OutputFile       string
	Concurrency      int
	Timeout          int
	MaxMutationsSize int
	Silent           bool
	Verbose          bool
	ActiveEnabled    bool
	Enrich           bool
}
