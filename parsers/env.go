package parsers

import (
	"fmt"
	"strings"
)

// EnvParser parses env command output
type EnvParser struct{}

// EnvEntry represents a single environment variable
type EnvEntry struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Original string `json:"original"`
}

func (p *EnvParser) Name() string {
	return "env"
}

func (p *EnvParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")
	var entries []EnvEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}

		entry := EnvEntry{
			Original: line,
		}

		// Environment variables are in format NAME=VALUE
		eqIdx := strings.Index(line, "=")
		if eqIdx == -1 {
			// No equals sign, treat whole line as name with empty value
			entry.Name = line
			entry.Value = ""
		} else {
			entry.Name = line[:eqIdx]
			entry.Value = line[eqIdx+1:]
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
