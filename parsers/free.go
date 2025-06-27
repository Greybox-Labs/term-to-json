package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// FreeParser parses free command output
type FreeParser struct{}

// FreeEntry represents free command output
type FreeEntry struct {
	Type      string `json:"type"`
	Total     int64  `json:"total"`
	Used      int64  `json:"used"`
	Free      int64  `json:"free"`
	Shared    int64  `json:"shared,omitempty"`
	Buffers   int64  `json:"buffers,omitempty"`
	Cache     int64  `json:"cache,omitempty"`
	Available int64  `json:"available,omitempty"`
}

// FreeOutput represents the complete free command output
type FreeOutput struct {
	Memory []FreeEntry `json:"memory"`
	Swap   *FreeEntry  `json:"swap,omitempty"`
}

func (p *FreeParser) Name() string {
	return "free"
}

func (p *FreeParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	if len(lines) < 2 {
		return nil, fmt.Errorf("insufficient lines in free output")
	}

	output := FreeOutput{}
	var memoryEntries []FreeEntry

	for i, line := range lines {
		// Skip header line
		if i == 0 && (strings.Contains(line, "total") || strings.Contains(line, "used")) {
			continue
		}

		fields := splitFields(line)
		if len(fields) < 4 {
			continue
		}

		entry := FreeEntry{}
		
		// Determine type (Mem:, -/+ buffers/cache:, Swap:)
		entryType := strings.TrimSuffix(fields[0], ":")
		entry.Type = entryType

		// Parse numeric fields
		if len(fields) >= 2 {
			if total, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
				entry.Total = total
			}
		}
		if len(fields) >= 3 {
			if used, err := strconv.ParseInt(fields[2], 10, 64); err == nil {
				entry.Used = used
			}
		}
		if len(fields) >= 4 {
			if free, err := strconv.ParseInt(fields[3], 10, 64); err == nil {
				entry.Free = free
			}
		}

		// Handle different free output formats
		if entryType == "Mem" {
			if len(fields) >= 5 {
				if shared, err := strconv.ParseInt(fields[4], 10, 64); err == nil {
					entry.Shared = shared
				}
			}
			if len(fields) >= 6 {
				if buffers, err := strconv.ParseInt(fields[5], 10, 64); err == nil {
					entry.Buffers = buffers
				}
			}
			if len(fields) >= 7 {
				if cache, err := strconv.ParseInt(fields[6], 10, 64); err == nil {
					entry.Cache = cache
				}
			}
			if len(fields) >= 8 {
				if available, err := strconv.ParseInt(fields[7], 10, 64); err == nil {
					entry.Available = available
				}
			}
		}

		if entryType == "Swap" {
			output.Swap = &entry
		} else {
			memoryEntries = append(memoryEntries, entry)
		}
	}

	output.Memory = memoryEntries
	return output, nil
}
