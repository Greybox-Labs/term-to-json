package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// WParser parses w command output
type WParser struct{}

// WEntry represents a single w output entry
type WEntry struct {
	User    string  `json:"user"`
	TTY     string  `json:"tty"`
	From    string  `json:"from,omitempty"`
	Login   string  `json:"login"`
	Idle    string  `json:"idle"`
	JCPU    string  `json:"jcpu,omitempty"`
	PCPU    string  `json:"pcpu,omitempty"`
	What    string  `json:"what"`
}

// WHeader represents the header info from w command
type WHeader struct {
	CurrentTime   string  `json:"current_time"`
	Uptime        string  `json:"uptime"`
	Users         int     `json:"users"`
	LoadAvg1      float64 `json:"load_avg_1"`
	LoadAvg5      float64 `json:"load_avg_5"`
	LoadAvg15     float64 `json:"load_avg_15"`
}

// WOutput represents the complete w command output
type WOutput struct {
	Header WHeader `json:"header"`
	Users  []WEntry `json:"users"`
}

func (p *WParser) Name() string {
	return "w"
}

func (p *WParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	if len(lines) < 2 {
		return nil, fmt.Errorf("insufficient lines in w output")
	}

	output := WOutput{}

	// Parse header (first line)
	headerLine := lines[0]
	output.Header = parseWHeader(headerLine)

	// Skip the column headers line (second line)
	dataLines := lines[2:]
	var entries []WEntry

	for _, line := range dataLines {
		fields := splitFields(line)
		if len(fields) < 4 {
			continue
		}

		entry := WEntry{}
		entry.User = fields[0]
		entry.TTY = fields[1]

		fieldIndex := 2
		// Check if FROM field is present (not always there)
		if len(fields) > 6 && !strings.Contains(fields[2], ":") {
			entry.From = fields[2]
			fieldIndex = 3
		}

		if fieldIndex < len(fields) {
			entry.Login = fields[fieldIndex]
		}
		if fieldIndex+1 < len(fields) {
			entry.Idle = fields[fieldIndex+1]
		}
		if fieldIndex+2 < len(fields) {
			entry.JCPU = fields[fieldIndex+2]
		}
		if fieldIndex+3 < len(fields) {
			entry.PCPU = fields[fieldIndex+3]
		}
		if fieldIndex+4 < len(fields) {
			entry.What = strings.Join(fields[fieldIndex+4:], " ")
		}

		entries = append(entries, entry)
	}

	output.Users = entries
	return output, nil
}

func parseWHeader(headerLine string) WHeader {
	header := WHeader{}

	// Example: " 14:30:42 up 12 days,  3:45,  2 users,  load average: 0.15, 0.12, 0.10"
	
	// Extract current time
	parts := strings.Fields(headerLine)
	if len(parts) > 0 {
		header.CurrentTime = parts[0]
	}

	// Parse users count
	for i, part := range parts {
		if strings.Contains(part, "user") && i > 0 {
			if users, err := strconv.Atoi(parts[i-1]); err == nil {
				header.Users = users
			}
			break
		}
	}

	// Parse uptime (simplified)
	if upIdx := findInSlice(parts, "up"); upIdx >= 0 && upIdx+1 < len(parts) {
		// Find the uptime portion
		uptimeParts := []string{}
		for i := upIdx + 1; i < len(parts); i++ {
			if strings.Contains(parts[i], "user") {
				break
			}
			uptimeParts = append(uptimeParts, parts[i])
		}
		header.Uptime = strings.Join(uptimeParts, " ")
		header.Uptime = strings.TrimSuffix(header.Uptime, ",")
	}

	// Parse load averages
	if avgIdx := findInSlice(parts, "average:"); avgIdx >= 0 {
		if avgIdx+1 < len(parts) {
			if load1, err := strconv.ParseFloat(strings.TrimSuffix(parts[avgIdx+1], ","), 64); err == nil {
				header.LoadAvg1 = load1
			}
		}
		if avgIdx+2 < len(parts) {
			if load5, err := strconv.ParseFloat(strings.TrimSuffix(parts[avgIdx+2], ","), 64); err == nil {
				header.LoadAvg5 = load5
			}
		}
		if avgIdx+3 < len(parts) {
			if load15, err := strconv.ParseFloat(parts[avgIdx+3], 64); err == nil {
				header.LoadAvg15 = load15
			}
		}
	}

	return header
}

func findInSlice(slice []string, target string) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}
