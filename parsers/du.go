package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// DuParser parses du command output
type DuParser struct{}

// DuEntry represents a single du output entry
type DuEntry struct {
	Size      int64  `json:"size"`
	SizeBytes int64  `json:"size_bytes"`
	Path      string `json:"path"`
}

func (p *DuParser) Name() string {
	return "du"
}

func (p *DuParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	var entries []DuEntry

	for _, line := range lines {
		fields := splitFields(line)
		
		// du output format: size path
		if len(fields) < 2 {
			continue
		}

		entry := DuEntry{}

		// Parse size (typically in KB by default)
		if size, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
			entry.Size = size
			// Convert to bytes (assuming KB input by default)
			entry.SizeBytes = size * 1024
		}

		// Path is everything after the first field
		entry.Path = strings.Join(fields[1:], " ")

		entries = append(entries, entry)
	}

	return entries, nil
}
