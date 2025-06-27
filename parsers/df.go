package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// DfParser parses df command output
type DfParser struct{}

// DfEntry represents a single df output entry
type DfEntry struct {
	Filesystem string  `json:"filesystem"`
	Size       int64   `json:"size"`
	Used       int64   `json:"used"`
	Available  int64   `json:"available"`
	UsePercent int     `json:"use_percent"`
	MountPoint string  `json:"mount_point"`
	UsedBytes  int64   `json:"used_bytes"`
	AvailBytes int64   `json:"avail_bytes"`
	SizeBytes  int64   `json:"size_bytes"`
}

func (p *DfParser) Name() string {
	return "df"
}

func (p *DfParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	if len(lines) == 0 {
		return []DfEntry{}, nil
	}

	// Skip header line
	dataLines := lines[1:]
	var entries []DfEntry

	for _, line := range dataLines {
		fields := splitFields(line)
		
		if len(fields) < 6 {
			continue
		}

		entry := DfEntry{}
		entry.Filesystem = fields[0]

		// Parse size (in 1K blocks by default)
		if size, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
			entry.Size = size
			entry.SizeBytes = size * 1024
		}

		// Parse used
		if used, err := strconv.ParseInt(fields[2], 10, 64); err == nil {
			entry.Used = used
			entry.UsedBytes = used * 1024
		}

		// Parse available
		if avail, err := strconv.ParseInt(fields[3], 10, 64); err == nil {
			entry.Available = avail
			entry.AvailBytes = avail * 1024
		}

		// Parse use percentage
		usePercentStr := strings.TrimSuffix(fields[4], "%")
		if usePercent, err := strconv.Atoi(usePercentStr); err == nil {
			entry.UsePercent = usePercent
		}

		// Mount point
		entry.MountPoint = fields[5]

		entries = append(entries, entry)
	}

	return entries, nil
}
