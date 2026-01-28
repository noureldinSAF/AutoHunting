package utils

import (
	"bufio"
	"os"
	"strings"
)

func IsStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func ReadInputFromStdin() ([]string, error) {
	hosts := []string{}

	// stdin has data (piped)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		hosts = append(hosts, strings.TrimSpace(scanner.Text()))
	}

	return hosts, nil

}

func ReadInputFromFile(path string) ([]string, error) {
	hosts := []string{}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(fileData)))

	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			hosts = append(hosts, strings.TrimSpace(line))
		}
	}

	return hosts, nil
}
