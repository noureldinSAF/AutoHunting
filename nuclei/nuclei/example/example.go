package main

import (
	"fmt"

	"vulnscans/nuclei"
)

func main() {
	scan := nuclei.NewNucleiScan(
		[]string{"http://194-204-44-234.ip.elisa.ee/"},
		[]nuclei.Severity{nuclei.Info},
		nil,           // no tags
		"phpinfo-files",     // template ID
	)

	result := scan.Run()
	if result == nil {
		fmt.Println("nuclei scan failed or returned no results")
		return
	}

	fmt.Println(result.GetJSONResult())
}