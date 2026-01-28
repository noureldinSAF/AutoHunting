package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadLines(filename string) ([]string, error) {
	lines := []string{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func ParsePorts(ps string) []int {
	intPorts := []int{}

	ports := strings.Split(ps, ",")
	for _, port := range ports {
		if port == "" {
			continue
		}

		port, err := strconv.Atoi(port)
		if err != nil {
			continue
		}
		intPorts = append(intPorts, port)
	}
	return intPorts
}
