package parsers

import (
	"fmt"
	"strings"
)

// MountParser parses mount command output
type MountParser struct{}

// MountEntry represents a single mount output entry
type MountEntry struct {
	Device      string `json:"device"`
	MountPoint  string `json:"mount_point"`
	FilesystemType string `json:"filesystem_type"`
	Options     string `json:"options"`
}

func (p *MountParser) Name() string {
	return "mount"
}

func (p *MountParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	
	lines := splitLines(input)
	var entries []MountEntry

	for _, line := range lines {
		// Mount output format: device on mountpoint type filesystem (options)
		// Example: /dev/sda1 on / type ext4 (rw,relatime)
		
		// Find " on " to split device from rest
		onIndex := strings.Index(line, " on ")
		if onIndex == -1 {
			continue
		}
		
		device := line[:onIndex]
		remainder := line[onIndex+4:] // Skip " on "
		
		// Find " type " to split mountpoint from filesystem and options
		typeIndex := strings.Index(remainder, " type ")
		if typeIndex == -1 {
			continue
		}
		
		mountPoint := remainder[:typeIndex]
		fsAndOptions := remainder[typeIndex+6:] // Skip " type "
		
		// Find space and opening paren to split filesystem from options
		openParenIndex := strings.Index(fsAndOptions, " (")
		if openParenIndex == -1 {
			// No options provided
			entry := MountEntry{
				Device:         device,
				MountPoint:     mountPoint,
				FilesystemType: fsAndOptions,
				Options:        "",
			}
			entries = append(entries, entry)
			continue
		}
		
		filesystem := fsAndOptions[:openParenIndex]
		optionsWithParen := fsAndOptions[openParenIndex+2:] // Skip " ("
		
		// Remove closing parenthesis
		options := strings.TrimSuffix(optionsWithParen, ")")
		
		entry := MountEntry{
			Device:         device,
			MountPoint:     mountPoint,
			FilesystemType: filesystem,
			Options:        options,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
