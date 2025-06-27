package parsers

import (
	"fmt"
	"strings"
)

// LsblkParser parses lsblk command output
type LsblkParser struct{}

// LsblkEntry represents a single lsblk output entry
type LsblkEntry struct {
	Name       string `json:"name"`
	MajMin     string `json:"maj_min"`
	Rm         string `json:"rm"`
	Size       string `json:"size"`
	Ro         string `json:"ro"`
	Type       string `json:"type"`
	Mountpoint string `json:"mountpoint"`
}

func (p *LsblkParser) Name() string {
	return "lsblk"
}

func (p *LsblkParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	if len(lines) == 0 {
		return []LsblkEntry{}, nil
	}

	// Skip header line
	dataLines := lines[1:]
	var entries []LsblkEntry

	for _, line := range dataLines {
		fields := splitFields(line)
		
		// lsblk output typically has 7 fields: NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINT
		if len(fields) < 6 {
			continue
		}

		entry := LsblkEntry{
			Name:   fields[0],
			MajMin: fields[1],
			Rm:     fields[2],
			Size:   fields[3],
			Ro:     fields[4],
			Type:   fields[5],
		}

		// Mountpoint is optional (can be empty)
		if len(fields) >= 7 {
			entry.Mountpoint = fields[6]
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
