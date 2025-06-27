package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// UptimeParser parses uptime command output
type UptimeParser struct{}

// UptimeEntry represents uptime output
type UptimeEntry struct {
	CurrentTime   string  `json:"current_time"`
	Uptime        string  `json:"uptime"`
	UptimeSeconds int     `json:"uptime_seconds"`
	Users         int     `json:"users"`
	LoadAvg1      float64 `json:"load_avg_1"`
	LoadAvg5      float64 `json:"load_avg_5"`
	LoadAvg15     float64 `json:"load_avg_15"`
}

func (p *UptimeParser) Name() string {
	return "uptime"
}

func (p *UptimeParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	entry := UptimeEntry{}

	// Example: " 14:30:42 up 12 days,  3:45,  2 users,  load average: 0.15, 0.12, 0.10"
	
	// Extract current time
	timeRe := regexp.MustCompile(`^\s*(\d{1,2}:\d{2}:\d{2})`)
	if matches := timeRe.FindStringSubmatch(input); len(matches) > 1 {
		entry.CurrentTime = matches[1]
	}

	// Extract uptime
	uptimeRe := regexp.MustCompile(`up\s+(.+?),\s+\d+\s+users?`)
	if matches := uptimeRe.FindStringSubmatch(input); len(matches) > 1 {
		entry.Uptime = strings.TrimSpace(matches[1])
		entry.UptimeSeconds = parseUptimeToSeconds(entry.Uptime)
	}

	// Extract users
	usersRe := regexp.MustCompile(`(\d+)\s+users?`)
	if matches := usersRe.FindStringSubmatch(input); len(matches) > 1 {
		if users, err := strconv.Atoi(matches[1]); err == nil {
			entry.Users = users
		}
	}

	// Extract load averages
	loadRe := regexp.MustCompile(`load average:\s*([0-9.]+),\s*([0-9.]+),\s*([0-9.]+)`)
	if matches := loadRe.FindStringSubmatch(input); len(matches) > 3 {
		if load1, err := strconv.ParseFloat(matches[1], 64); err == nil {
			entry.LoadAvg1 = load1
		}
		if load5, err := strconv.ParseFloat(matches[2], 64); err == nil {
			entry.LoadAvg5 = load5
		}
		if load15, err := strconv.ParseFloat(matches[3], 64); err == nil {
			entry.LoadAvg15 = load15
		}
	}

	return entry, nil
}

func parseUptimeToSeconds(uptime string) int {
	seconds := 0
	
	// Parse days
	dayRe := regexp.MustCompile(`(\d+)\s+days?`)
	if matches := dayRe.FindStringSubmatch(uptime); len(matches) > 1 {
		if days, err := strconv.Atoi(matches[1]); err == nil {
			seconds += days * 24 * 3600
		}
	}

	// Parse hours and minutes
	timeRe := regexp.MustCompile(`(\d+):(\d+)`)
	if matches := timeRe.FindStringSubmatch(uptime); len(matches) > 2 {
		if hours, err := strconv.Atoi(matches[1]); err == nil {
			seconds += hours * 3600
		}
		if minutes, err := strconv.Atoi(matches[2]); err == nil {
			seconds += minutes * 60
		}
	}

	// Parse minutes only (e.g., "5 min")
	minRe := regexp.MustCompile(`(\d+)\s+min`)
	if matches := minRe.FindStringSubmatch(uptime); len(matches) > 1 {
		if minutes, err := strconv.Atoi(matches[1]); err == nil {
			seconds += minutes * 60
		}
	}

	return seconds
}
