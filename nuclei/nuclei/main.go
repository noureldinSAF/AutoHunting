package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type NucleiOutput struct {
	TemplateUrl string `json:"template-url"`
	Info        struct {
		Name        string `json:"name"`
		Author []string `json:"author"`
		Tags   []string `json:"tags"`
		Description string `json:"description"`
		Severity    string `json:"severity"`
		Remediation string `json:"remediation"`
	} `json:"info"`
	Url string `json:"matched-at"`
	POC string `json:"curl-command"`
}

func runNuclei(url string, templateId string, tag string) []*NucleiOutput {

	fmtCommand := []string{"nuclei", "-u", url, "-jsonl", "-silent"}

	if templateId == "" {
		fmtCommand = append(fmtCommand, "-tags", tag)
	} else {
		fmtCommand = append(fmtCommand, "-id", templateId)
	}


	cmd := exec.Command(fmtCommand[0], fmtCommand[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parseNucleiOutput(output)
}

func parseNucleiOutput(output []byte) []*NucleiOutput {
    var nucleiOutputs []*NucleiOutput

    scanner := bufio.NewScanner(bytes.NewReader(output))
    // Set buffer to 1MB max line size
    buf := make([]byte, 5*1024*1024)
    scanner.Buffer(buf, 5*1024*1024)

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue
        }
        var nucleiOutput NucleiOutput
        err := json.Unmarshal([]byte(line), &nucleiOutput)
        if err != nil {
            fmt.Printf("failed to parse line: %v\n", err)
            continue
        }
        nucleiOutputs = append(nucleiOutputs, &nucleiOutput)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("scanner error: %v\n", err)
        return nil
    }

    return nucleiOutputs
}

func main() {
	nucleiOutputs := runNuclei("http://194-204-44-234.ip.elisa.ee/ ", "", "phpinfo")
	for _, nucleiOutput := range nucleiOutputs {
		fmt.Println(nucleiOutput.Url)
		fmt.Println(nucleiOutput.POC)
		fmt.Println(nucleiOutput.Info.Description)
	}
}