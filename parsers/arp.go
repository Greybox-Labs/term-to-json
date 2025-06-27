package parsers

import (
	"fmt"
	"strings"
)

// ArpParser parses arp command output
type ArpParser struct{}

// ArpEntry represents a single ARP table entry
type ArpEntry struct {
	Address    string `json:"address"`
	HWType     string `json:"hw_type,omitempty"`
	HWAddress  string `json:"hw_address"`
	Flags      string `json:"flags,omitempty"`
	Mask       string `json:"mask,omitempty"`
	Interface  string `json:"interface"`
	Incomplete bool   `json:"incomplete,omitempty"`
}

func (p *ArpParser) Name() string {
	return "arp"
}

func (p *ArpParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	var entries []ArpEntry

	for _, line := range lines {
		// Skip header lines
		if strings.HasPrefix(line, "Address") ||
		   strings.Contains(line, "HWtype") ||
		   strings.Contains(line, "HWaddress") {
			continue
		}

		fields := splitFields(line)
		if len(fields) < 3 {
			continue
		}

		entry := ArpEntry{}
		entry.Address = fields[0]

		// Check if this is an incomplete entry
		if strings.Contains(line, "<incomplete>") {
			entry.Incomplete = true
			entry.HWAddress = "<incomplete>"
			// Find interface (usually last field)
			if len(fields) > 1 {
				entry.Interface = fields[len(fields)-1]
			}
		} else {
			// Standard format: Address HWtype HWaddress Flags Mask Iface
			if len(fields) >= 2 {
				entry.HWType = fields[1]
			}
			if len(fields) >= 3 {
				entry.HWAddress = fields[2]
			}
			if len(fields) >= 4 {
				entry.Flags = fields[3]
			}
			if len(fields) >= 5 {
				entry.Mask = fields[4]
			}
			if len(fields) >= 6 {
				entry.Interface = fields[5]
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
