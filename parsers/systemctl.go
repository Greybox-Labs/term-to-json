package parsers

import (
	"fmt"
	"strings"
)

// SystemctlParser parses systemctl command output
type SystemctlParser struct{}

// SystemctlEntry represents a systemctl service entry
type SystemctlEntry struct {
	Unit        string `json:"unit"`
	Load        string `json:"load"`
	Active      string `json:"active"`
	Sub         string `json:"sub"`
	Description string `json:"description"`
	// For systemctl status output
	Status      string `json:"status,omitempty"`
	Main        string `json:"main,omitempty"`
	Tasks       string `json:"tasks,omitempty"`
	Memory      string `json:"memory,omitempty"`
	CPU         string `json:"cpu,omitempty"`
	ProcessID   string `json:"process_id,omitempty"`
}

func (p *SystemctlParser) Name() string {
	return "systemctl"
}

func (p *SystemctlParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	if len(lines) == 0 {
		return []SystemctlEntry{}, nil
	}

	// Check if this is status output (contains ●)
	if strings.Contains(input, "●") || strings.Contains(input, "Active:") {
		return p.parseStatus(lines)
	}

	// Otherwise parse list-units output
	return p.parseListUnits(lines)
}

func (p *SystemctlParser) parseListUnits(lines []string) (interface{}, error) {
	var entries []SystemctlEntry
	
	// Skip header lines until we find one starting with UNIT
	startIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "UNIT") {
			startIdx = i + 1
			break
		}
	}

	if startIdx == -1 {
		startIdx = 0
	}

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		
		// Skip summary lines
		if strings.Contains(line, "loaded units listed") || 
		   strings.Contains(line, "units listed") ||
		   strings.HasPrefix(line, "LOAD") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		entry := SystemctlEntry{
			Unit:   fields[0],
			Load:   fields[1],
			Active: fields[2],
			Sub:    fields[3],
		}

		// Description is the rest of the fields
		if len(fields) > 4 {
			entry.Description = strings.Join(fields[4:], " ")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (p *SystemctlParser) parseStatus(lines []string) (interface{}, error) {
	entry := SystemctlEntry{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Parse service name from first line with ●
		if strings.Contains(line, "●") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				entry.Unit = parts[1]
				if len(parts) > 3 {
					entry.Description = strings.Join(parts[3:], " ")
				}
			}
		}
		
		// Parse key-value pairs
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				switch key {
				case "Loaded":
					entry.Load = value
				case "Active":
					entry.Active = value
					entry.Status = value
				case "Main PID":
					entry.Main = value
					entry.ProcessID = value
				case "Tasks":
					entry.Tasks = value
				case "Memory":
					entry.Memory = value
				case "CPU":
					entry.CPU = value
				}
			}
		}
	}

	return entry, nil
}
