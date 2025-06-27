package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// PasswdParser parses /etc/passwd file format
type PasswdParser struct{}

// PasswdEntry represents a single passwd file entry
type PasswdEntry struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	GECOS    string `json:"gecos"`
	HomeDir  string `json:"home_dir"`
	Shell    string `json:"shell"`
	Original string `json:"original"`
}

func (p *PasswdParser) Name() string {
	return "passwd"
}

func (p *PasswdParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")
	var entries []PasswdEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry := PasswdEntry{
			Original: line,
		}

		// passwd format: username:password:UID:GID:GECOS:directory:shell
		fields := strings.Split(line, ":")
		if len(fields) != 7 {
			continue
		}

		entry.Username = fields[0]
		entry.Password = fields[1]
		
		if uid, err := strconv.Atoi(fields[2]); err == nil {
			entry.UID = uid
		}
		
		if gid, err := strconv.Atoi(fields[3]); err == nil {
			entry.GID = gid
		}
		
		entry.GECOS = fields[4]
		entry.HomeDir = fields[5]
		entry.Shell = fields[6]

		entries = append(entries, entry)
	}

	return entries, nil
}
