package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FindParser parses find command output (find -ls format)
type FindParser struct{}

// FindEntry represents a single find output entry
type FindEntry struct {
	Path         string    `json:"path"`
	Type         string    `json:"type"`
	Permissions  string    `json:"permissions"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modified_time"`
	Inode        int64     `json:"inode"`
	Links        int       `json:"links"`
	Owner        string    `json:"owner"`
	Group        string    `json:"group"`
}

func (p *FindParser) Name() string {
	return "find"
}

func (p *FindParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	var entries []FindEntry

	for _, line := range lines {
		// Handle different find output formats
		if strings.Contains(line, " ") {
			// Assume find -ls format or similar detailed output
			if entry := parseFindLsLine(line); entry != nil {
				entries = append(entries, *entry)
			}
		} else {
			// Simple path-only format
			entry := FindEntry{
				Path: strings.TrimSpace(line),
				Type: "unknown",
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// parseFindLsLine parses a line from find -ls output
func parseFindLsLine(line string) *FindEntry {
	fields := splitFields(line)
	
	// find -ls format: inode blocks permissions links owner group size date time path
	if len(fields) < 11 {
		return nil
	}

	entry := &FindEntry{}

	// Parse inode
	if inode, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
		entry.Inode = inode
	}

	// Parse permissions and determine type
	entry.Permissions = fields[2]
	if strings.HasPrefix(entry.Permissions, "d") {
		entry.Type = "directory"
	} else if strings.HasPrefix(entry.Permissions, "l") {
		entry.Type = "symlink"
	} else if strings.HasPrefix(entry.Permissions, "-") {
		entry.Type = "file"
	} else {
		entry.Type = "special"
	}

	// Parse links
	if links, err := strconv.Atoi(fields[3]); err == nil {
		entry.Links = links
	}

	// Parse owner and group
	entry.Owner = fields[4]
	entry.Group = fields[5]

	// Parse size
	if size, err := strconv.ParseInt(fields[6], 10, 64); err == nil {
		entry.Size = size
	}

	// Parse date/time (fields 7, 8, 9)
	if len(fields) >= 10 {
		dateStr := strings.Join(fields[7:10], " ")
		if parsedTime, err := parseDate(dateStr); err == nil {
			entry.ModifiedTime = parsedTime
		}
	}

	// Parse path (remaining fields)
	if len(fields) > 10 {
		entry.Path = strings.Join(fields[10:], " ")
	}

	return entry
}
