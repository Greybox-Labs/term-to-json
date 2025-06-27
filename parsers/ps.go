package parsers

import (
	"strconv"
	"strings"
)

// PsParser parses ps command output
type PsParser struct{}

// PsEntry represents a single ps output entry
type PsEntry struct {
	PID     int     `json:"pid"`
	PPID    int     `json:"ppid,omitempty"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu_percent,omitempty"`
	Memory  float64 `json:"memory_percent,omitempty"`
	VSZ     int64   `json:"vsz,omitempty"`
	RSS     int64   `json:"rss,omitempty"`
	TTY     string  `json:"tty"`
	Stat    string  `json:"stat,omitempty"`
	Start   string  `json:"start,omitempty"`
	Time    string  `json:"time,omitempty"`
	Command string  `json:"command"`
}

func (p *PsParser) Name() string {
	return "ps"
}

func (p *PsParser) Parse(input string) (interface{}, error) {
	lines := splitLines(input)
	if len(lines) == 0 {
		return []PsEntry{}, nil
	}

	// Skip header line
	dataLines := lines[1:]
	var entries []PsEntry

	for _, line := range dataLines {
		fields := splitFields(line)
		
		if len(fields) < 4 {
			continue
		}

		entry := PsEntry{}

		// Basic ps output: PID TTY TIME CMD
		if len(fields) >= 4 {
			if pid, err := strconv.Atoi(fields[0]); err == nil {
				entry.PID = pid
			}
			entry.TTY = fields[1]
			entry.Time = fields[2]
			entry.Command = strings.Join(fields[3:], " ")
		}

		// Extended ps output (ps aux format): USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND
		if len(fields) >= 11 {
			entry.User = fields[0]
			if pid, err := strconv.Atoi(fields[1]); err == nil {
				entry.PID = pid
			}
			if cpu, err := strconv.ParseFloat(fields[2], 64); err == nil {
				entry.CPU = cpu
			}
			if mem, err := strconv.ParseFloat(fields[3], 64); err == nil {
				entry.Memory = mem
			}
			if vsz, err := strconv.ParseInt(fields[4], 10, 64); err == nil {
				entry.VSZ = vsz
			}
			if rss, err := strconv.ParseInt(fields[5], 10, 64); err == nil {
				entry.RSS = rss
			}
			entry.TTY = fields[6]
			entry.Stat = fields[7]
			entry.Start = fields[8]
			entry.Time = fields[9]
			entry.Command = strings.Join(fields[10:], " ")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
