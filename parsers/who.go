package parsers

import (
	"fmt"
	"strings"
	"time"
)

// WhoParser parses who command output
type WhoParser struct{}

// WhoEntry represents a single who output entry
type WhoEntry struct {
	User     string    `json:"user"`
	TTY      string    `json:"tty"`
	LoginTime time.Time `json:"login_time"`
	Host     string    `json:"host,omitempty"`
	Comment  string    `json:"comment,omitempty"`
}

func (p *WhoParser) Name() string {
	return "who"
}

func (p *WhoParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	var entries []WhoEntry

	for _, line := range lines {
		fields := splitFields(line)
		if len(fields) < 3 {
			continue
		}

		entry := WhoEntry{}
		entry.User = fields[0]
		entry.TTY = fields[1]

		// Parse date and time (format varies)
		// Example: "2023-01-15 14:30" or "Jan 15 14:30"
		dateTimeStr := ""
		if len(fields) >= 4 {
			dateTimeStr = fields[2] + " " + fields[3]
		}

		if parsedTime, err := parseWhoDateTime(dateTimeStr); err == nil {
			entry.LoginTime = parsedTime
		}

		// Extract host if present (usually in parentheses)
		restFields := fields[4:]
		for _, field := range restFields {
			if strings.HasPrefix(field, "(") && strings.HasSuffix(field, ")") {
				entry.Host = strings.Trim(field, "()")
			} else {
				if entry.Comment == "" {
					entry.Comment = field
				} else {
					entry.Comment += " " + field
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func parseWhoDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"Jan 2 15:04",
		"Jan  2 15:04",
		"2006-01-02 15:04:05",
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

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
