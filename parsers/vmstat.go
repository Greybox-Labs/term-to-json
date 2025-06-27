package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// VmstatParser parses vmstat command output
type VmstatParser struct{}

// VmstatEntry represents a single vmstat output entry
type VmstatEntry struct {
	Processes    VmstatProcesses `json:"processes"`
	Memory       VmstatMemory    `json:"memory"`
	Swap         VmstatSwap      `json:"swap"`
	IO           VmstatIO        `json:"io"`
	System       VmstatSystem    `json:"system"`
	CPU          VmstatCPU       `json:"cpu"`
}

type VmstatProcesses struct {
	Runnable int `json:"runnable"`
	Blocked  int `json:"blocked"`
}

type VmstatMemory struct {
	SwapUsed   int64 `json:"swap_used"`
	Free       int64 `json:"free"`
	Buffers    int64 `json:"buffers"`
	Cache      int64 `json:"cache"`
}

type VmstatSwap struct {
	In  int `json:"in"`
	Out int `json:"out"`
}

type VmstatIO struct {
	BlocksIn  int `json:"blocks_in"`
	BlocksOut int `json:"blocks_out"`
}

type VmstatSystem struct {
	Interrupts      int `json:"interrupts"`
	ContextSwitches int `json:"context_switches"`
}

type VmstatCPU struct {
	UserTime   int `json:"user_time"`
	SystemTime int `json:"system_time"`
	IdleTime   int `json:"idle_time"`
	WaitTime   int `json:"wait_time"`
	StolenTime int `json:"stolen_time,omitempty"`
}

func (p *VmstatParser) Name() string {
	return "vmstat"
}

func (p *VmstatParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	var entries []VmstatEntry

	headerSeen := false
	for _, line := range lines {
		// Skip header lines
		if strings.Contains(line, "procs") || strings.Contains(line, "memory") ||
		   strings.Contains(line, "r") && strings.Contains(line, "b") && !headerSeen {
			headerSeen = true
			continue
		}

		fields := splitFields(line)
		if len(fields) < 16 {
			continue
		}

		entry := VmstatEntry{}

		// Parse processes
		if r, err := strconv.Atoi(fields[0]); err == nil {
			entry.Processes.Runnable = r
		}
		if b, err := strconv.Atoi(fields[1]); err == nil {
			entry.Processes.Blocked = b
		}

		// Parse memory
		if swpd, err := strconv.ParseInt(fields[2], 10, 64); err == nil {
			entry.Memory.SwapUsed = swpd
		}
		if free, err := strconv.ParseInt(fields[3], 10, 64); err == nil {
			entry.Memory.Free = free
		}
		if buff, err := strconv.ParseInt(fields[4], 10, 64); err == nil {
			entry.Memory.Buffers = buff
		}
		if cache, err := strconv.ParseInt(fields[5], 10, 64); err == nil {
			entry.Memory.Cache = cache
		}

		// Parse swap
		if si, err := strconv.Atoi(fields[6]); err == nil {
			entry.Swap.In = si
		}
		if so, err := strconv.Atoi(fields[7]); err == nil {
			entry.Swap.Out = so
		}

		// Parse I/O
		if bi, err := strconv.Atoi(fields[8]); err == nil {
			entry.IO.BlocksIn = bi
		}
		if bo, err := strconv.Atoi(fields[9]); err == nil {
			entry.IO.BlocksOut = bo
		}

		// Parse system
		if in, err := strconv.Atoi(fields[10]); err == nil {
			entry.System.Interrupts = in
		}
		if cs, err := strconv.Atoi(fields[11]); err == nil {
			entry.System.ContextSwitches = cs
		}

		// Parse CPU
		if us, err := strconv.Atoi(fields[12]); err == nil {
			entry.CPU.UserTime = us
		}
		if sy, err := strconv.Atoi(fields[13]); err == nil {
			entry.CPU.SystemTime = sy
		}
		if id, err := strconv.Atoi(fields[14]); err == nil {
			entry.CPU.IdleTime = id
		}
		if wa, err := strconv.Atoi(fields[15]); err == nil {
			entry.CPU.WaitTime = wa
		}
		if len(fields) > 16 {
			if st, err := strconv.Atoi(fields[16]); err == nil {
				entry.CPU.StolenTime = st
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
