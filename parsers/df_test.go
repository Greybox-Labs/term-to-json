package parsers

import (
	"encoding/json"
	"testing"
)

func TestDfParser(t *testing.T) {
	parser := &DfParser{}

	testInput := `Filesystem     1K-blocks    Used Available Use% Mounted on
/dev/sda1       20511312  123456  19365472   1% /
tmpfs            4096000       0   4096000   0% /tmp
/dev/sdb1      104857600 52428800  47185920  53% /home`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]DfEntry)
	if !ok {
		t.Fatalf("Expected []DfEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test first entry
	if entries[0].Filesystem != "/dev/sda1" {
		t.Errorf("Expected filesystem '/dev/sda1', got '%s'", entries[0].Filesystem)
	}
	if entries[0].Size != 20511312 {
		t.Errorf("Expected size 20511312, got %d", entries[0].Size)
	}
	if entries[0].SizeBytes != 20511312*1024 {
		t.Errorf("Expected size bytes %d, got %d", 20511312*1024, entries[0].SizeBytes)
	}
	if entries[0].UsePercent != 1 {
		t.Errorf("Expected use percent 1, got %d", entries[0].UsePercent)
	}
	if entries[0].MountPoint != "/" {
		t.Errorf("Expected mount point '/', got '%s'", entries[0].MountPoint)
	}

	// Test second entry (tmpfs with 0% usage)
	if entries[1].UsePercent != 0 {
		t.Errorf("Expected use percent 0, got %d", entries[1].UsePercent)
	}
	if entries[1].Used != 0 {
		t.Errorf("Expected used 0, got %d", entries[1].Used)
	}

	// Test third entry with higher usage
	if entries[2].UsePercent != 53 {
		t.Errorf("Expected use percent 53, got %d", entries[2].UsePercent)
	}
	if entries[2].MountPoint != "/home" {
		t.Errorf("Expected mount point '/home', got '%s'", entries[2].MountPoint)
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

func TestDfParserEmpty(t *testing.T) {
	parser := &DfParser{}
	
	result, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	result, err = parser.Parse("Filesystem     1K-blocks    Used Available Use% Mounted on")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]DfEntry)
	if !ok {
		t.Fatalf("Expected []DfEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
