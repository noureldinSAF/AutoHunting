package nuclei

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cyinnove/logify"
	"github.com/google/uuid"
)

func NewNucleiScan(target []string, severities []Severity, tags []string, templateID string) *NucleiScan {
	return &NucleiScan{
		Target:     target,
		Severities: severities,
		Tags:       tags,
		TemplateID: templateID,
	}
}

func (s *NucleiScan) Run() *Output {

	s.tmpFile = fmt.Sprintf("/tmp/%s_tmp_subdomains.txt", uuid.New().String())

	s.getNucleiCommand()
	results, err := s.execAndGetResults()
	if err != nil {
		logify.Errorf("Error while executing nuclei command %s", err.Error())
		return nil
	}
	
	if results.Results == nil {
		results.Results = make([]*JSONCLIOutput, 0)
	}

	output := &Output{
		Results: make([]*SubdomainOutput, 0),
	}

	subdomainMap := make(map[string]*SubdomainOutput)

	for _, result := range results.Results {
		clientResult := &extractedResult{
			References: make([]string, 0),
		}

		clientResult.Severity = result.Info.Severity
		clientResult.Description = result.Info.Description
		clientResult.POC = result.MatchedAt
		clientResult.Subdomain = result.Host
		clientResult.Name = result.Info.Name
		clientResult.TemplateID = result.TemplateId
		clientResult.References = result.Info.Reference

		if clientResult.POC == "" {
			clientResult.POC = result.Host
		}

		if subdomainOutput, exists := subdomainMap[clientResult.Subdomain]; exists {
			subdomainOutput.ExtractedData = append(subdomainOutput.ExtractedData, clientResult)
		} else {
			subdomainOutput := &SubdomainOutput{
				Subdomain:     clientResult.Subdomain,
				ExtractedData: []*extractedResult{clientResult},
			}
			subdomainMap[clientResult.Subdomain] = subdomainOutput
			output.Results = append(output.Results, subdomainOutput)
		}
	}

	return output
}

func (o *Output) GetJSONResult() string {
	if len(o.Results) == 0 {
		log.Println("No results to write")
		return ""
	}
	jsonData, err := json.MarshalIndent(o.Results, "", "    ")
	if err != nil {
		log.Printf("Error while marshalling the results %s", err.Error())
		logify.Errorf("Error while marshalling the results %s", err.Error())
		return ""
	}
	return string(jsonData)
}