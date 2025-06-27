package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DateParser parses date command output
type DateParser struct{}

// DateEntry represents parsed date command output
type DateEntry struct {
	Timestamp   int64  `json:"timestamp"`
	ISO         string `json:"iso"`
	RFC3339     string `json:"rfc3339"`
	Unix        int64  `json:"unix"`
	Weekday     string `json:"weekday"`
	Month       string `json:"month"`
	Day         int    `json:"day"`
	Time        string `json:"time"`
	Timezone    string `json:"timezone"`
	Year        int    `json:"year"`
	Original    string `json:"original"`
}

func (p *DateParser) Name() string {
	return "date"
}

func (p *DateParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	entry := DateEntry{
		Original: input,
	}

	// Try to parse common date formats
	if unixTime, err := strconv.ParseInt(input, 10, 64); err == nil {
		// Unix timestamp
		t := time.Unix(unixTime, 0)
		entry.Timestamp = unixTime
		entry.Unix = unixTime
		entry.ISO = t.Format(time.RFC3339)
		entry.RFC3339 = t.Format(time.RFC3339)
		entry.fillFromTime(t)
		return entry, nil
	}

	// Try parsing ISO format
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		entry.Timestamp = t.Unix()
		entry.Unix = t.Unix()
		entry.ISO = t.Format(time.RFC3339)
		entry.RFC3339 = t.Format(time.RFC3339)
		entry.fillFromTime(t)
		return entry, nil
	}

	// Parse standard date output format: "Wed Jan 15 14:30:25 PST 2025"
	dateRegex := regexp.MustCompile(`^(\w+)\s+(\w+)\s+(\d+)\s+(\d+:\d+:\d+)\s+(\w+)\s+(\d+)$`)
	matches := dateRegex.FindStringSubmatch(input)
	if len(matches) == 7 {
		entry.Weekday = matches[1]
		entry.Month = matches[2]
		if day, err := strconv.Atoi(matches[3]); err == nil {
			entry.Day = day
		}
		entry.Time = matches[4]
		entry.Timezone = matches[5]
		if year, err := strconv.Atoi(matches[6]); err == nil {
			entry.Year = year
		}

		// Try to parse into time.Time for timestamp
		layout := "Mon Jan 2 15:04:05 MST 2006"
		if t, err := time.Parse(layout, input); err == nil {
			entry.Timestamp = t.Unix()
			entry.Unix = t.Unix()
			entry.ISO = t.Format(time.RFC3339)
			entry.RFC3339 = t.Format(time.RFC3339)
		}
		return entry, nil
	}

	// Generic parsing attempt
	if t, err := time.Parse(time.RFC1123, input); err == nil {
		entry.Timestamp = t.Unix()
		entry.Unix = t.Unix()
		entry.ISO = t.Format(time.RFC3339)
		entry.RFC3339 = t.Format(time.RFC3339)
		entry.fillFromTime(t)
		return entry, nil
	}

	return entry, nil
}

func (d *DateEntry) fillFromTime(t time.Time) {
	d.Weekday = t.Weekday().String()
	d.Month = t.Month().String()
	d.Day = t.Day()
	d.Time = t.Format("15:04:05")
	d.Timezone = t.Format("MST")
	d.Year = t.Year()
}
