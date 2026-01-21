package runner 

type Options struct {
	Query      string
	queries     []string
	InputFile   string
	OutputFile  string
	Concurrency int
	Timeout     int
	Silent      bool
	Verbose     bool
}
