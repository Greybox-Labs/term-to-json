package parsers

import (
	"encoding/json"
	"testing"
)

func TestPsParser(t *testing.T) {
	parser := &PsParser{}

	testInput := `USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1  225316  9876 ?        Ss   Jan01   0:01 /sbin/init
user      1234  2.5  1.2  123456 12345 pts/0    R+   10:30   0:05 python script.py
user      5678  0.1  0.5   67890  5432 pts/1    S    09:15   0:00 bash`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]PsEntry)
	if !ok {
		t.Fatalf("Expected []PsEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test first entry
	if entries[0].User != "root" {
		t.Errorf("Expected user 'root', got '%s'", entries[0].User)
	}
	if entries[0].PID != 1 {
		t.Errorf("Expected PID 1, got %d", entries[0].PID)
	}
	if entries[0].Command != "/sbin/init" {
		t.Errorf("Expected command '/sbin/init', got '%s'", entries[0].Command)
	}

	// Test second entry with CPU/Memory
	if entries[1].CPU != 2.5 {
		t.Errorf("Expected CPU 2.5, got %f", entries[1].CPU)
	}
	if entries[1].Memory != 1.2 {
		t.Errorf("Expected Memory 1.2, got %f", entries[1].Memory)
	}
	if entries[1].VSZ != 123456 {
		t.Errorf("Expected VSZ 123456, got %d", entries[1].VSZ)
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

func TestPsParserSimple(t *testing.T) {
	parser := &PsParser{}

	// Test simple ps output format
	testInput := `  PID TTY          TIME CMD
 1234 pts/0    00:00:05 python
 5678 pts/1    00:00:00 bash`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]PsEntry)
	if !ok {
		t.Fatalf("Expected []PsEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	if entries[0].PID != 1234 {
		t.Errorf("Expected PID 1234, got %d", entries[0].PID)
	}
	if entries[0].TTY != "pts/0" {
		t.Errorf("Expected TTY 'pts/0', got '%s'", entries[0].TTY)
	}
}
