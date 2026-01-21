package utils

import	(
	"os"
	"strings"
	"bufio"
)

func ReadInputFromFile(file string) ([]string, error ) {

	fileData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileData), "\n")

	return lines, nil
}

func ExtractDomainsFromString(input string) []string {
	return strings.Split(input, ",")
}

func WriteOutputToFile(file string, data []string) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)

	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}




