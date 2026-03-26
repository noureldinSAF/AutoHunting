package nuclei

//type NucleiResult struct {
//	Results map[string][]*JSONCLIOutput
//}

type Severity int8

const (
	Info Severity = iota
	Low
	Medium
	High
	Critical
)

type JSONCLIOutput struct {
	Template   string `json:"template"`
	TemplateId string `json:"template-id"`
	Info       struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Reference   []string `json:"reference"`
		Severity    string   `json:"severity"`
	} `json:"info"`
	Host      string `json:"host"`
	MatchedAt string `json:"matched-at"`
}

type NucleiScan struct {
	Target     []string
	Severities []Severity
	Tags       []string
	TemplateID string
	tmpFile    string
	nucleiCmd  string
}

type extractedResult struct {
	TemplateID  string   `json:"template_id"`
	Subdomain   string   `json:"subdomain"`
	Name        string   `json:"template_name"`
	Description string   `json:"template_description"`
	Severity    string   `json:"template_severity"`
	POC         string   `json:"template_poc"`
	References  []string `json:"template_references"`
}

type NucleiResult struct {
	Results []*JSONCLIOutput
}

type SubdomainOutput struct {
	Subdomain     string             `json:"subdomain"`
	ExtractedData []*extractedResult `json:"extracted_data"`
}

type Output struct {
	Results []*SubdomainOutput `json:"results"`
}