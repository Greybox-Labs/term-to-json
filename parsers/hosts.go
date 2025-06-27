package parsers

import (
	"fmt"
	"strings"
)

// HostsParser parses /etc/hosts file format
type HostsParser struct{}

// HostsEntry represents a single hosts file entry
type HostsEntry struct {
	IP        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
	Comment   string   `json:"comment,omitempty"`
	Original  string   `json:"original"`
}

func (p *HostsParser) Name() string {
	return "hosts"
}

func (p *HostsParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")
	var entries []HostsEntry

	for _, line := range lines {
		original := line
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}

		entry := HostsEntry{
			Original: original,
		}

		// Handle comment-only lines
		if strings.HasPrefix(line, "#") {
			entry.Comment = strings.TrimSpace(line[1:])
			entries = append(entries, entry)
			continue
		}

		// Check for inline comments
		commentIdx := strings.Index(line, "#")
		if commentIdx != -1 {
			entry.Comment = strings.TrimSpace(line[commentIdx+1:])
			line = strings.TrimSpace(line[:commentIdx])
		}

		// Parse IP and hostnames
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			entry.IP = fields[0]
			entry.Hostnames = fields[1:]
		} else if len(fields) == 1 {
			// Only IP, no hostnames
			entry.IP = fields[0]
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
