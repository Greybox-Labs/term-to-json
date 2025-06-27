package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// IdParser parses id command output
type IdParser struct{}

// IdEntry represents id command output
type IdEntry struct {
	UID      int      `json:"uid"`
	User     string   `json:"user"`
	GID      int      `json:"gid"`
	Group    string   `json:"group"`
	Groups   []IdGroup `json:"groups"`
	Context  string   `json:"context,omitempty"`
}

// IdGroup represents a group entry
type IdGroup struct {
	GID  int    `json:"gid"`
	Name string `json:"name"`
}

func (p *IdParser) Name() string {
	return "id"
}

func (p *IdParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	entry := IdEntry{}

	// Example: uid=1000(user) gid=1000(user) groups=1000(user),4(adm),24(cdrom),27(sudo)
	
	// Parse UID
	uidRe := regexp.MustCompile(`uid=(\d+)\(([^)]+)\)`)
	if matches := uidRe.FindStringSubmatch(input); len(matches) > 2 {
		if uid, err := strconv.Atoi(matches[1]); err == nil {
			entry.UID = uid
		}
		entry.User = matches[2]
	}

	// Parse GID
	gidRe := regexp.MustCompile(`gid=(\d+)\(([^)]+)\)`)
	if matches := gidRe.FindStringSubmatch(input); len(matches) > 2 {
		if gid, err := strconv.Atoi(matches[1]); err == nil {
			entry.GID = gid
		}
		entry.Group = matches[2]
	}

	// Parse groups
	groupsRe := regexp.MustCompile(`groups=([^,\s]+(?:,[^,\s]+)*)`)
	if matches := groupsRe.FindStringSubmatch(input); len(matches) > 1 {
		groupStr := matches[1]
		groupEntries := strings.Split(groupStr, ",")
		
		for _, groupEntry := range groupEntries {
			groupRe := regexp.MustCompile(`(\d+)\(([^)]+)\)`)
			if groupMatches := groupRe.FindStringSubmatch(groupEntry); len(groupMatches) > 2 {
				group := IdGroup{}
				if gid, err := strconv.Atoi(groupMatches[1]); err == nil {
					group.GID = gid
				}
				group.Name = groupMatches[2]
				entry.Groups = append(entry.Groups, group)
			}
		}
	}

	// Parse context (SELinux)
	contextRe := regexp.MustCompile(`context=([^\s]+)`)
	if matches := contextRe.FindStringSubmatch(input); len(matches) > 1 {
		entry.Context = matches[1]
	}

	return entry, nil
}
