package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LsParser parses ls command output
type LsParser struct{}

// LsEntry represents a single ls output entry
type LsEntry struct {
	Permissions string    `json:"permissions"`
	Links       int       `json:"links"`
	Owner       string    `json:"owner"`
	Group       string    `json:"group"`
	Size        int64     `json:"size"`
	Modified    time.Time `json:"modified"`
	Name        string    `json:"name"`
	IsDirectory bool      `json:"is_directory"`
	IsSymlink   bool      `json:"is_symlink"`
	LinkTarget  string    `json:"link_target,omitempty"`
}

func (p *LsParser) Name() string {
	return "ls"
}

func (p *LsParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	var entries []LsEntry

	for _, line := range lines {
		fields := splitFields(line)
		
		// Skip lines that don't look like ls -l output
		if len(fields) < 9 {
			continue
		}

		entry := LsEntry{}

		// Parse permissions
		entry.Permissions = fields[0]
		entry.IsDirectory = strings.HasPrefix(entry.Permissions, "d")
		entry.IsSymlink = strings.HasPrefix(entry.Permissions, "l")

		// Parse links
		if links, err := strconv.Atoi(fields[1]); err == nil {
			entry.Links = links
		}

		// Parse owner and group
		entry.Owner = fields[2]
		entry.Group = fields[3]

		// Parse size
		if size, err := strconv.ParseInt(fields[4], 10, 64); err == nil {
			entry.Size = size
		}

		// Parse date/time (assumes format: Mon DD HH:MM or Mon DD YYYY)
		dateStr := strings.Join(fields[5:8], " ")
		if parsedTime, err := parseDate(dateStr); err == nil {
			entry.Modified = parsedTime
		}

		// Parse filename and link target
		nameFields := fields[8:]
		if entry.IsSymlink && len(nameFields) >= 3 {
			// Handle symlink: name -> target
			arrowIndex := -1
			for i, field := range nameFields {
				if field == "->" {
					arrowIndex = i
					break
				}
			}
			if arrowIndex > 0 {
				entry.Name = strings.Join(nameFields[:arrowIndex], " ")
				entry.LinkTarget = strings.Join(nameFields[arrowIndex+1:], " ")
			} else {
				entry.Name = strings.Join(nameFields, " ")
			}
		} else {
			entry.Name = strings.Join(nameFields, " ")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// parseDate attempts to parse various date formats from ls output
func parseDate(dateStr string) (time.Time, error) {
	// Common ls date formats
	formats := []string{
		"Jan 2 15:04",
		"Jan 2 2006",
		"Jan  2 15:04",
		"Jan  2 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// If year is not specified, assume current year
			if t.Year() == 0 {
				now := time.Now()
				t = t.AddDate(now.Year(), 0, 0)
			}
			return t, nil
		}
	}

	return time.Time{}, nil
}
