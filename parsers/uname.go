package parsers

import (
	"fmt"
	"strings"
)

// UnameParser parses uname command output
type UnameParser struct{}

// UnameEntry represents uname output
type UnameEntry struct {
	KernelName    string `json:"kernel_name"`
	NodeName      string `json:"node_name"`
	KernelRelease string `json:"kernel_release"`
	KernelVersion string `json:"kernel_version"`
	Machine       string `json:"machine"`
	Processor     string `json:"processor,omitempty"`
	OS            string `json:"operating_system,omitempty"`
}

func (p *UnameParser) Name() string {
	return "uname"
}

func (p *UnameParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	fields := splitFields(input)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields found")
	}

	entry := UnameEntry{}

	// uname -a output: Linux hostname 5.4.0-74-generic #83-Ubuntu SMP Sat May 8 02:35:39 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
	if len(fields) >= 1 {
		entry.KernelName = fields[0]
	}
	if len(fields) >= 2 {
		entry.NodeName = fields[1]
	}
	if len(fields) >= 3 {
		entry.KernelRelease = fields[2]
	}
	if len(fields) >= 6 {
		// Kernel version is typically multiple fields
		entry.KernelVersion = strings.Join(fields[3:len(fields)-3], " ")
	}
	if len(fields) >= 7 {
		entry.Machine = fields[len(fields)-3]
	}
	if len(fields) >= 8 {
		entry.Processor = fields[len(fields)-2]
	}
	if len(fields) >= 9 {
		entry.OS = fields[len(fields)-1]
	}

	return entry, nil
}
