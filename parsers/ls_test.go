package parsers

import (
	"encoding/json"
	"testing"
)

func TestLsParser(t *testing.T) {
	parser := &LsParser{}

	testInput := `total 48
drwxr-xr-x  3 user group  4096 Jan 15 10:30 docs
-rw-r--r--  1 user group  1234 Jan 14 09:15 file.txt
lrwxrwxrwx  1 user group     8 Jan 13 14:20 link -> file.txt
-rwxr-xr-x  1 user group  5678 Jan 12 16:45 script.sh`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]LsEntry)
	if !ok {
		t.Fatalf("Expected []LsEntry, got %T", result)
	}

	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries, got %d", len(entries))
	}

	// Test directory entry
	if !entries[0].IsDirectory {
		t.Error("Expected first entry to be a directory")
	}
	if entries[0].Name != "docs" {
		t.Errorf("Expected name 'docs', got '%s'", entries[0].Name)
	}
	if entries[0].Size != 4096 {
		t.Errorf("Expected size 4096, got %d", entries[0].Size)
	}

	// Test file entry
	if entries[1].IsDirectory {
		t.Error("Expected second entry to be a file")
	}
	if entries[1].Name != "file.txt" {
		t.Errorf("Expected name 'file.txt', got '%s'", entries[1].Name)
	}

	// Test symlink entry
	if !entries[2].IsSymlink {
		t.Error("Expected third entry to be a symlink")
	}
	if entries[2].LinkTarget != "file.txt" {
		t.Errorf("Expected link target 'file.txt', got '%s'", entries[2].LinkTarget)
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

func TestLsParserEmpty(t *testing.T) {
	parser := &LsParser{}
	
	result, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	result, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	result, err = parser.Parse("total 0")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]LsEntry)
	if !ok {
		t.Fatalf("Expected []LsEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
