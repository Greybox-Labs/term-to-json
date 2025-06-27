package parsers

import (
	"fmt"
	"strings"
)

// Parser defines the interface for all command line parsers
type Parser interface {
	Parse(input string) (interface{}, error)
	Name() string
}

// Parse is the main entry point for parsing command line output
func Parse(parserName, input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	var parser Parser
	switch parserName {
	case "ls":
		parser = &LsParser{}
	case "ps":
		parser = &PsParser{}
	case "df":
		parser = &DfParser{}
	case "mount":
		parser = &MountParser{}
	case "lsblk":
		parser = &LsblkParser{}
	case "du":
		parser = &DuParser{}
	case "find":
		parser = &FindParser{}
	case "stat":
		parser = &StatParser{}
	case "uname":
		parser = &UnameParser{}
	case "uptime":
		parser = &UptimeParser{}
	case "who":
		parser = &WhoParser{}
	case "w":
		parser = &WParser{}
	case "id":
		parser = &IdParser{}
	case "ping":
		parser = &PingParser{}
	case "netstat":
		parser = &NetstatParser{}
	case "arp":
		parser = &ArpParser{}
	case "free":
		parser = &FreeParser{}
	case "vmstat":
		parser = &VmstatParser{}
	case "date":
		parser = &DateParser{}
	case "systemctl":
		parser = &SystemctlParser{}
	case "hosts":
		parser = &HostsParser{}
	case "passwd":
		parser = &PasswdParser{}
	case "env":
		parser = &EnvParser{}
	case "wc":
		parser = &WcParser{}
	case "dig":
		parser = &DigParser{}
	default:
		return nil, fmt.Errorf("unknown parser: %s", parserName)
	}

	return parser.Parse(input)
}

// splitLines splits input into lines and filters out empty lines
func splitLines(input string) []string {
	rawLines := strings.Split(input, "\n")
	var lines []string
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// splitFields splits a line into fields by whitespace
func splitFields(line string) []string {
	return strings.Fields(line)
}
