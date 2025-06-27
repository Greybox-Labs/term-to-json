package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StatParser parses stat command output
type StatParser struct{}

// StatEntry represents a single stat output entry
type StatEntry struct {
	File        string    `json:"file"`
	Size        int64     `json:"size"`
	Blocks      int64     `json:"blocks"`
	IOBlock     int64     `json:"io_block"`
	Device      string    `json:"device"`
	Inode       int64     `json:"inode"`
	Links       int       `json:"links"`
	Permissions string    `json:"permissions"`
	UID         int       `json:"uid"`
	GID         int       `json:"gid"`
	AccessTime  time.Time `json:"access_time"`
	ModifyTime  time.Time `json:"modify_time"`
	ChangeTime  time.Time `json:"change_time"`
}

func (p *StatParser) Name() string {
	return "stat"
}

func (p *StatParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	var entries []StatEntry
	var currentEntry *StatEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a new file entry (starts with "File:")
		if strings.HasPrefix(line, "File:") {
			// Save previous entry if exists
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
			}
			// Start new entry
			currentEntry = &StatEntry{}
			// Extract filename from "File: 'filename'" or "File: filename"
			fileStr := strings.TrimPrefix(line, "File:")
			fileStr = strings.TrimSpace(fileStr)
			fileStr = strings.Trim(fileStr, "'\"")
			currentEntry.File = fileStr
		} else if currentEntry != nil {
			// Parse other stat fields
			parseStatField(line, currentEntry)
		}
	}

	// Add the last entry
	if currentEntry != nil {
		entries = append(entries, *currentEntry)
	}

	return entries, nil
}

// parseStatField parses individual stat output fields
func parseStatField(line string, entry *StatEntry) {
	// Handle different stat output formats
	if strings.Contains(line, "Size:") {
		// Size: 1024 Blocks: 8 IO Block: 4096 regular file
		fields := splitFields(line)
		for i, field := range fields {
			switch field {
			case "Size:":
				if i+1 < len(fields) {
					if size, err := strconv.ParseInt(fields[i+1], 10, 64); err == nil {
						entry.Size = size
					}
				}
			case "Blocks:":
				if i+1 < len(fields) {
					if blocks, err := strconv.ParseInt(fields[i+1], 10, 64); err == nil {
						entry.Blocks = blocks
					}
				}
			case "Block:":
				if i+1 < len(fields) {
					if ioBlock, err := strconv.ParseInt(fields[i+1], 10, 64); err == nil {
						entry.IOBlock = ioBlock
					}
				}
			}
		}
	} else if strings.Contains(line, "Device:") {
		// Device: 801h/2049d Inode: 123456 Links: 1
		fields := splitFields(line)
		for i, field := range fields {
			switch field {
			case "Device:":
				if i+1 < len(fields) {
					entry.Device = fields[i+1]
				}
			case "Inode:":
				if i+1 < len(fields) {
					if inode, err := strconv.ParseInt(fields[i+1], 10, 64); err == nil {
						entry.Inode = inode
					}
				}
			case "Links:":
				if i+1 < len(fields) {
					if links, err := strconv.Atoi(fields[i+1]); err == nil {
						entry.Links = links
					}
				}
			}
		}
	} else if strings.Contains(line, "Access: (") {
		// Access: (0644/-rw-r--r--) Uid: (1000/user) Gid: (1000/group)
		// Extract permissions
		if start := strings.Index(line, "("); start != -1 {
			if end := strings.Index(line[start:], ")"); end != -1 {
				permsStr := line[start+1 : start+end]
				if slash := strings.Index(permsStr, "/"); slash != -1 {
					entry.Permissions = permsStr[slash+1:]
				}
			}
		}
		
		// Extract UID and GID
		if uidStart := strings.Index(line, "Uid: ("); uidStart != -1 {
			uidEnd := strings.Index(line[uidStart:], ")")
			if uidEnd != -1 {
				uidStr := line[uidStart+6 : uidStart+uidEnd]
				if slash := strings.Index(uidStr, "/"); slash != -1 {
					if uid, err := strconv.Atoi(uidStr[:slash]); err == nil {
						entry.UID = uid
					}
				}
			}
		}
		
		if gidStart := strings.Index(line, "Gid: ("); gidStart != -1 {
			gidEnd := strings.Index(line[gidStart:], ")")
			if gidEnd != -1 {
				gidStr := line[gidStart+6 : gidStart+gidEnd]
				if slash := strings.Index(gidStr, "/"); slash != -1 {
					if gid, err := strconv.Atoi(gidStr[:slash]); err == nil {
						entry.GID = gid
					}
				}
			}
		}
	} else if strings.HasPrefix(line, "Access: ") && !strings.Contains(line, "(") {
		// Access: 2023-01-01 12:00:00.000000000 +0000
		timeStr := strings.TrimPrefix(line, "Access: ")
		if t, err := parseStatTime(timeStr); err == nil {
			entry.AccessTime = t
		}
	} else if strings.HasPrefix(line, "Modify: ") {
		// Modify: 2023-01-01 12:00:00.000000000 +0000
		timeStr := strings.TrimPrefix(line, "Modify: ")
		if t, err := parseStatTime(timeStr); err == nil {
			entry.ModifyTime = t
		}
	} else if strings.HasPrefix(line, "Change: ") {
		// Change: 2023-01-01 12:00:00.000000000 +0000
		timeStr := strings.TrimPrefix(line, "Change: ")
		if t, err := parseStatTime(timeStr); err == nil {
			entry.ChangeTime = t
		}
	}
}

// parseStatTime parses timestamp from stat output
func parseStatTime(timeStr string) (time.Time, error) {
	// Common stat time formats
	formats := []string{
		"2006-01-02 15:04:05.000000000 -0700",
		"2006-01-02 15:04:05.000000000 +0000",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 +0000",
		"2006-01-02 15:04:05.000000000",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
