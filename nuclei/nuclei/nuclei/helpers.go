package nuclei

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cyinnove/logify"
)

func (s *NucleiScan) getNucleiCommand() {
	if err := WriteLinesToFile(s.tmpFile, s.Target); err != nil {
		logify.Errorf("Error while writing the subdomains to the tmp file %s", err.Error())
		return
	}

	baseCmd := fmt.Sprintf("nuclei -l \"%s\" -jsonl -silent -nh -ni", s.tmpFile)

	// Build severity flag, if any
	sevFlag := ""
	for _, severity := range s.Severities {
		switch severity {
		case Info:
			sevFlag += "info,"
		case Low:
			sevFlag += "low,"
		case Medium:
			sevFlag += "medium,"
		case High:
			sevFlag += "high,"
		case Critical:
			sevFlag += "critical,"
		}
	}
	sevFlag = strings.TrimSuffix(sevFlag, ",")

	// Build tags flag, if any
	tagsFlag := ""
	if len(s.Tags) > 0 {
		tagsFlag = strings.Join(s.Tags, ",")
	}

	// Prefer template ID if provided, otherwise tags, otherwise fallback to severity
	switch {
	case s.TemplateID != "":
		s.nucleiCmd = fmt.Sprintf("%s -id %s", baseCmd, s.TemplateID)
	case tagsFlag != "" && sevFlag != "":
		s.nucleiCmd = fmt.Sprintf("%s -tags %s -severity %s", baseCmd, tagsFlag, sevFlag)
	case tagsFlag != "":
		s.nucleiCmd = fmt.Sprintf("%s -tags %s", baseCmd, tagsFlag)
	case sevFlag != "":
		s.nucleiCmd = fmt.Sprintf("%s -severity %s", baseCmd, sevFlag)
	default:
		s.nucleiCmd = baseCmd
	}

	logify.Infof("%v", s.nucleiCmd)
}

func (s *NucleiScan) execAndGetResults() (NucleiResult, error) {
	output, err := s.execCmd()
	if err != nil {
		logify.Errorf("Error while executing nuclei command: %s", err.Error())
		if output != "" {
			logify.Errorf("Nuclei output: %s", output)
		}
		return NucleiResult{}, err
	}

	results, err := parseOutput(output)
	if err != nil {
		logify.Errorf("Error while parsing nuclei output %s", err.Error())
		return NucleiResult{}, err
	}
	return results, nil
}

func (s *NucleiScan) execCmd() (string, error) {
	output, err := exec.Command("bash", "-c", s.nucleiCmd).CombinedOutput()
	return string(output), err
}

func parseOutput(output string) (NucleiResult, error) {
	results := NucleiResult{
		Results: make([]*JSONCLIOutput, 0),
	}

	reader := strings.NewReader(output)
	decoder := json.NewDecoder(reader)
	for {
		var data JSONCLIOutput
		err := decoder.Decode(&data)
		if err != nil {
			if err.Error() != "EOF" {
				logify.Errorf("Error decoding JSON: %s", err)
			}
			break
		}
		results.Results = append(results.Results, &data)
	}
	return results, nil
}

func SetSeverities(severities []string) []Severity {
	severityList := []Severity{}
	for _, severity := range severities {
		switch severity {
		case "info":
			severityList = append(severityList, Info)
		case "low":
			severityList = append(severityList, Low)
		case "medium":
			severityList = append(severityList, Medium)
		case "high":
			severityList = append(severityList, High)
		case "critical":
			severityList = append(severityList, Critical)
		}
	}
	return severityList
}

func main() {
	file, err := os.Open("./test/urls.txt")
	if err != nil {
		logify.Errorf("Error while opening the file %s", err.Error())
		return
	}
	defer file.Close()

	var target []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			target = append(target, scanner.Text())
		}
	}

	// Example: use severity only
	nucleiScan := NewNucleiScan(target, SetSeverities([]string{"medium"}), nil, "")
	result := nucleiScan.Run()

	if result != nil {
		logify.Debugf("JSON Data : %s", result.GetJSONResult())
	} else {
		logify.Errorf("Failed to run NucleiScan")
	}
}
