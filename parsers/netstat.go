package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// NetstatParser parses netstat command output
type NetstatParser struct{}

// NetstatEntry represents a single netstat output entry
type NetstatEntry struct {
	Protocol     string `json:"protocol"`
	LocalAddress string `json:"local_address"`
	ForeignAddress string `json:"foreign_address"`
	State        string `json:"state,omitempty"`
	PID          int    `json:"pid,omitempty"`
	Program      string `json:"program,omitempty"`
	RecvQ        int    `json:"recv_q,omitempty"`
	SendQ        int    `json:"send_q,omitempty"`
}

func (p *NetstatParser) Name() string {
	return "netstat"
}

func (p *NetstatParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	var entries []NetstatEntry

	for _, line := range lines {
		// Skip header lines and empty lines
		if strings.HasPrefix(line, "Active") || 
		   strings.HasPrefix(line, "Proto") ||
		   strings.Contains(line, "Local Address") ||
		   strings.TrimSpace(line) == "" {
			continue
		}

		fields := splitFields(line)
		if len(fields) < 4 {
			continue
		}

		entry := NetstatEntry{}
		entry.Protocol = fields[0]

		fieldIndex := 1

		// Check if RecvQ and SendQ are present
		if len(fields) >= 6 && isNumeric(fields[1]) {
			if recvQ, err := strconv.Atoi(fields[1]); err == nil {
				entry.RecvQ = recvQ
			}
			if sendQ, err := strconv.Atoi(fields[2]); err == nil {
				entry.SendQ = sendQ
			}
			fieldIndex = 3
		}

		// Local and Foreign addresses
		if fieldIndex < len(fields) {
			entry.LocalAddress = fields[fieldIndex]
		}
		if fieldIndex+1 < len(fields) {
			entry.ForeignAddress = fields[fieldIndex+1]
		}

		// State (for TCP connections)
		if fieldIndex+2 < len(fields) && entry.Protocol == "tcp" {
			possibleState := fields[fieldIndex+2]
			if isNetstatState(possibleState) {
				entry.State = possibleState
				fieldIndex++
			}
		}

		// PID/Program (if present)
		if fieldIndex+2 < len(fields) {
			pidProgram := fields[fieldIndex+2]
			if strings.Contains(pidProgram, "/") {
				parts := strings.SplitN(pidProgram, "/", 2)
				if len(parts) == 2 {
					if pid, err := strconv.Atoi(parts[0]); err == nil {
						entry.PID = pid
					}
					entry.Program = parts[1]
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isNetstatState(s string) bool {
	states := []string{
		"ESTABLISHED", "SYN_SENT", "SYN_RECV", "FIN_WAIT1", "FIN_WAIT2",
		"TIME_WAIT", "CLOSE", "CLOSE_WAIT", "LAST_ACK", "LISTEN",
		"CLOSING", "UNKNOWN",
	}

	for _, state := range states {
		if s == state {
			return true
		}
	}
	return false
}
