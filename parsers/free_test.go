package parsers

import (
	"encoding/json"
	"testing"
)

func TestFreeParser(t *testing.T) {
	parser := &FreeParser{}

	testInput := `              total        used        free      shared     buffers       cache   available
Mem:        8048604     2048152     4096000      102400      204800     1904452     5600000
Swap:       2097148           0     2097148`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	output, ok := result.(FreeOutput)
	if !ok {
		t.Fatalf("Expected FreeOutput, got %T", result)
	}

	// Test memory entries
	if len(output.Memory) != 1 {
		t.Fatalf("Expected 1 memory entry, got %d", len(output.Memory))
	}

	memEntry := output.Memory[0]
	if memEntry.Type != "Mem" {
		t.Errorf("Expected type 'Mem', got '%s'", memEntry.Type)
	}
	if memEntry.Total != 8048604 {
		t.Errorf("Expected total 8048604, got %d", memEntry.Total)
	}
	if memEntry.Used != 2048152 {
		t.Errorf("Expected used 2048152, got %d", memEntry.Used)
	}
	if memEntry.Free != 4096000 {
		t.Errorf("Expected free 4096000, got %d", memEntry.Free)
	}
	if memEntry.Shared != 102400 {
		t.Errorf("Expected shared 102400, got %d", memEntry.Shared)
	}
	if memEntry.Buffers != 204800 {
		t.Errorf("Expected buffers 204800, got %d", memEntry.Buffers)
	}
	if memEntry.Cache != 1904452 {
		t.Errorf("Expected cache 1904452, got %d", memEntry.Cache)
	}
	if memEntry.Available != 5600000 {
		t.Errorf("Expected available 5600000, got %d", memEntry.Available)
	}

	// Test swap entry
	if output.Swap == nil {
		t.Fatal("Expected swap entry, got nil")
	}
	if output.Swap.Type != "Swap" {
		t.Errorf("Expected swap type 'Swap', got '%s'", output.Swap.Type)
	}
	if output.Swap.Total != 2097148 {
		t.Errorf("Expected swap total 2097148, got %d", output.Swap.Total)
	}
	if output.Swap.Used != 0 {
		t.Errorf("Expected swap used 0, got %d", output.Swap.Used)
	}
	if output.Swap.Free != 2097148 {
		t.Errorf("Expected swap free 2097148, got %d", output.Swap.Free)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON output is empty")
	}
}

func TestFreeParserSimple(t *testing.T) {
	parser := &FreeParser{}

	// Test simpler free output format
	testInput := `             total       used       free
Mem:       8048604    2048152    4096000`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	output, ok := result.(FreeOutput)
	if !ok {
		t.Fatalf("Expected FreeOutput, got %T", result)
	}

	if len(output.Memory) != 1 {
		t.Fatalf("Expected 1 memory entry, got %d", len(output.Memory))
	}

	memEntry := output.Memory[0]
	if memEntry.Total != 8048604 {
		t.Errorf("Expected total 8048604, got %d", memEntry.Total)
	}
}

func TestFreeParserEmpty(t *testing.T) {
	parser := &FreeParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only header
	_, err = parser.Parse("total used free")
	if err == nil {
		t.Error("Expected error for insufficient lines")
	}
}
