package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// WcParser parses wc command output
type WcParser struct{}

// WcEntry represents wc command output
type WcEntry struct {
	Lines      int    `json:"lines"`
	Words      int    `json:"words"`
	Characters int    `json:"characters"`
	Bytes      int    `json:"bytes"`
	Filename   string `json:"filename,omitempty"`
	Original   string `json:"original"`
}

func (p *WcParser) Name() string {
	return "wc"
}

func (p *WcParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	var entries []WcEntry

	for _, line := range lines {
		entry := WcEntry{
			Original: line,
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		// Standard wc output formats:
		// wc -l: "123 filename" or just "123"
		// wc -w: "456 filename" or just "456"
		// wc -c: "789 filename" or just "789"
		// wc (default): "123 456 789 filename" or "123 456 789"

		switch len(fields) {
		case 1:
			// Single number, could be lines, words, or characters
			if val, err := strconv.Atoi(fields[0]); err == nil {
				entry.Lines = val
				entry.Words = val
				entry.Characters = val
			}
		case 2:
			// Number and filename
			if val, err := strconv.Atoi(fields[0]); err == nil {
				entry.Lines = val
				entry.Words = val
				entry.Characters = val
			}
			entry.Filename = fields[1]
		case 3:
			// lines words chars (no filename)
			if lines, err := strconv.Atoi(fields[0]); err == nil {
				entry.Lines = lines
			}
			if words, err := strconv.Atoi(fields[1]); err == nil {
				entry.Words = words
			}
			if chars, err := strconv.Atoi(fields[2]); err == nil {
				entry.Characters = chars
				entry.Bytes = chars
			}
		case 4:
			// lines words chars filename
			if lines, err := strconv.Atoi(fields[0]); err == nil {
				entry.Lines = lines
			}
			if words, err := strconv.Atoi(fields[1]); err == nil {
				entry.Words = words
			}
			if chars, err := strconv.Atoi(fields[2]); err == nil {
				entry.Characters = chars
				entry.Bytes = chars
			}
			entry.Filename = fields[3]
		}

		entries = append(entries, entry)
	}

	// If only one entry and it looks like a total line, return just that entry
	if len(entries) == 1 {
		return entries[0], nil
	}

	return entries, nil
}
