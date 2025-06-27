package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// DigParser parses dig DNS command output
type DigParser struct{}

// DigEntry represents dig command output
type DigEntry struct {
	Query     DigQuery      `json:"query"`
	Answer    []DigAnswer   `json:"answer"`
	Authority []DigAnswer   `json:"authority,omitempty"`
	Additional []DigAnswer  `json:"additional,omitempty"`
	Stats     DigStats      `json:"stats"`
	Original  string        `json:"original"`
}

// DigQuery represents the query section
type DigQuery struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

// DigAnswer represents a DNS answer record
type DigAnswer struct {
	Name  string `json:"name"`
	TTL   int    `json:"ttl"`
	Class string `json:"class"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DigStats represents query statistics
type DigStats struct {
	QueryTime   int    `json:"query_time_ms"`
	Server      string `json:"server"`
	When        string `json:"when"`
	MessageSize int    `json:"message_size"`
	Flags       string `json:"flags"`
	Status      string `json:"status"`
}

func (p *DigParser) Name() string {
	return "dig"
}

func (p *DigParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	entry := DigEntry{
		Original: input,
		Answer:   []DigAnswer{},
		Authority: []DigAnswer{},
		Additional: []DigAnswer{},
	}

	lines := strings.Split(input, "\n")
	section := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			// Parse query from comment
			if strings.Contains(line, ";; QUESTION SECTION:") {
				section = "question"
			} else if strings.Contains(line, ";; ANSWER SECTION:") {
				section = "answer"
			} else if strings.Contains(line, ";; AUTHORITY SECTION:") {
				section = "authority"
			} else if strings.Contains(line, ";; ADDITIONAL SECTION:") {
				section = "additional"
			} else if strings.HasPrefix(line, ";;") && strings.Contains(line, "IN") {
				// Parse query line like ";; example.com. IN A"
				parts := strings.Fields(line)
				if len(parts) >= 4 {
					entry.Query.Name = strings.TrimSuffix(parts[1], ".")
					entry.Query.Class = parts[2]
					entry.Query.Type = parts[3]
				}
			}
			continue
		}

		// Parse answer records
		if section == "answer" || section == "authority" || section == "additional" {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				answer := DigAnswer{
					Name:  fields[0],
					Class: fields[2],
					Type:  fields[3],
					Value: strings.Join(fields[4:], " "),
				}
				
				if ttl, err := strconv.Atoi(fields[1]); err == nil {
					answer.TTL = ttl
				}

				switch section {
				case "answer":
					entry.Answer = append(entry.Answer, answer)
				case "authority":
					entry.Authority = append(entry.Authority, answer)
				case "additional":
					entry.Additional = append(entry.Additional, answer)
				}
			}
		}

		// Parse statistics
		if strings.HasPrefix(line, ";; Query time:") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				if queryTime, err := strconv.Atoi(parts[3]); err == nil {
					entry.Stats.QueryTime = queryTime
				}
			}
		} else if strings.HasPrefix(line, ";; SERVER:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				entry.Stats.Server = parts[2]
			}
		} else if strings.HasPrefix(line, ";; WHEN:") {
			entry.Stats.When = strings.Join(strings.Fields(line)[2:], " ")
		} else if strings.HasPrefix(line, ";; MSG SIZE") {
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				if msgSize, err := strconv.Atoi(parts[4]); err == nil {
					entry.Stats.MessageSize = msgSize
				}
			}
		} else if strings.Contains(line, "status:") && strings.Contains(line, "flags:") {
			// Parse status line
			if statusIdx := strings.Index(line, "status:"); statusIdx != -1 {
				statusPart := line[statusIdx+7:]
				if commaIdx := strings.Index(statusPart, ","); commaIdx != -1 {
					entry.Stats.Status = strings.TrimSpace(statusPart[:commaIdx])
				} else {
					entry.Stats.Status = strings.TrimSpace(statusPart)
				}
			}
			if flagsIdx := strings.Index(line, "flags:"); flagsIdx != -1 {
				flagsPart := line[flagsIdx+6:]
				if semicolonIdx := strings.Index(flagsPart, ";"); semicolonIdx != -1 {
					entry.Stats.Flags = strings.TrimSpace(flagsPart[:semicolonIdx])
				} else {
					entry.Stats.Flags = strings.TrimSpace(flagsPart)
				}
			}
		}
	}

	return entry, nil
}
