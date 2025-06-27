package parsers

import (
	"encoding/json"
	"testing"
)

func TestVmstatParser(t *testing.T) {
	parser := &VmstatParser{}

	testInput := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0      0 4096000 204800 1024000    0    0     5    10  150  300 12  3 85  0  0
 0  0      0 4000000 205000 1025000    0    0     3     8  140  280 10  2 88  0  0`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]VmstatEntry)
	if !ok {
		t.Fatalf("Expected []VmstatEntry, got %T", result)
	}

	// Should parse actual data lines, skipping headers
	if len(entries) < 2 {
		t.Fatalf("Expected at least 2 entries, got %d", len(entries))
	}

	// Find a valid data entry (not header)
	var entry VmstatEntry
	found := false
	for _, e := range entries {
		if e.CPU.IdleTime > 0 { // Valid data entry
			entry = e
			found = true
			break
		}
	}
	if !found {
		t.Fatal("No valid vmstat entry found")
	}

	// Test processes (check that data was parsed)
	if entry.Memory.Free == 0 {
		t.Error("Expected free memory to be parsed")
	}
	if entry.Processes.Blocked != 0 {
		t.Errorf("Expected blocked processes 0, got %d", entry.Processes.Blocked)
	}

	// Test memory
	if entry.Memory.SwapUsed != 0 {
		t.Errorf("Expected swap used 0, got %d", entry.Memory.SwapUsed)
	}
	if entry.Memory.Free != 4096000 {
		t.Errorf("Expected free memory 4096000, got %d", entry.Memory.Free)
	}
	if entry.Memory.Buffers != 204800 {
		t.Errorf("Expected buffers 204800, got %d", entry.Memory.Buffers)
	}
	if entry.Memory.Cache != 1024000 {
		t.Errorf("Expected cache 1024000, got %d", entry.Memory.Cache)
	}

	// Test swap
	if entry.Swap.In != 0 {
		t.Errorf("Expected swap in 0, got %d", entry.Swap.In)
	}
	if entry.Swap.Out != 0 {
		t.Errorf("Expected swap out 0, got %d", entry.Swap.Out)
	}

	// Test I/O
	if entry.IO.BlocksIn != 5 {
		t.Errorf("Expected blocks in 5, got %d", entry.IO.BlocksIn)
	}
	if entry.IO.BlocksOut != 10 {
		t.Errorf("Expected blocks out 10, got %d", entry.IO.BlocksOut)
	}

	// Test system
	if entry.System.Interrupts != 150 {
		t.Errorf("Expected interrupts 150, got %d", entry.System.Interrupts)
	}
	if entry.System.ContextSwitches != 300 {
		t.Errorf("Expected context switches 300, got %d", entry.System.ContextSwitches)
	}

	// Test CPU
	if entry.CPU.UserTime != 12 {
		t.Errorf("Expected user time 12, got %d", entry.CPU.UserTime)
	}
	if entry.CPU.SystemTime != 3 {
		t.Errorf("Expected system time 3, got %d", entry.CPU.SystemTime)
	}
	if entry.CPU.IdleTime != 85 {
		t.Errorf("Expected idle time 85, got %d", entry.CPU.IdleTime)
	}
	if entry.CPU.WaitTime != 0 {
		t.Errorf("Expected wait time 0, got %d", entry.CPU.WaitTime)
	}
	if entry.CPU.StolenTime != 0 {
		t.Errorf("Expected stolen time 0, got %d", entry.CPU.StolenTime)
	}

	// Test that we got valid entries
	if len(entries) >= 2 {
		entry2 := entries[len(entries)-1] // Last entry
		if entry2.CPU.IdleTime == 0 {
			t.Error("Expected idle time to be parsed")
		}
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON output is empty")
	}
}

func TestVmstatParserWithStolenTime(t *testing.T) {
	parser := &VmstatParser{}

	testInput := `procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  1   1024 3000000 150000  800000    1    2     8    15  200  400 15  5 75  3  2`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]VmstatEntry)
	if !ok {
		t.Fatalf("Expected []VmstatEntry, got %T", result)
	}

	if len(entries) < 1 {
		t.Fatalf("Expected at least 1 entry, got %d", len(entries))
	}

	// Find a valid data entry
	var entry VmstatEntry
	found := false
	for _, e := range entries {
		if e.Memory.Free > 0 { // Valid data entry
			entry = e
			found = true
			break
		}
	}
	if !found {
		t.Fatal("No valid vmstat entry found")
	}

	// Test that data was parsed correctly
	if entry.Memory.Free == 0 {
		t.Error("Expected free memory to be parsed")
	}
	if entry.Processes.Blocked != 1 {
		t.Errorf("Expected blocked processes 1, got %d", entry.Processes.Blocked)
	}
	if entry.Memory.SwapUsed != 1024 {
		t.Errorf("Expected swap used 1024, got %d", entry.Memory.SwapUsed)
	}
	if entry.Swap.In != 1 {
		t.Errorf("Expected swap in 1, got %d", entry.Swap.In)
	}
	if entry.CPU.WaitTime != 3 {
		t.Errorf("Expected wait time 3, got %d", entry.CPU.WaitTime)
	}
	if entry.CPU.StolenTime != 2 {
		t.Errorf("Expected stolen time 2, got %d", entry.CPU.StolenTime)
	}
}

func TestVmstatParserEmpty(t *testing.T) {
	parser := &VmstatParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only headers
	result, err := parser.Parse("procs memory\n r b swpd free")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]VmstatEntry)
	if !ok {
		t.Fatalf("Expected []VmstatEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
