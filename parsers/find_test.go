package parsers

import (
	"encoding/json"
	"testing"
)

func TestFindParser(t *testing.T) {
	parser := &FindParser{}

	testInput := `   123456       8 -rw-r--r--   1 user group     1024 Jan 15 14:30 ./file.txt
   789012      16 drwxr-xr-x   2 user group     4096 Jan 14 10:20 ./directory
   345678       4 lrwxrwxrwx   1 user group        8 Jan 13 16:45 ./link`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]FindEntry)
	if !ok {
		t.Fatalf("Expected []FindEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test file entry
	if entries[0].Inode != 123456 {
		t.Errorf("Expected inode 123456, got %d", entries[0].Inode)
	}
	if entries[0].Type != "file" {
		t.Errorf("Expected type 'file', got '%s'", entries[0].Type)
	}
	if entries[0].Permissions != "-rw-r--r--" {
		t.Errorf("Expected permissions '-rw-r--r--', got '%s'", entries[0].Permissions)
	}
	if entries[0].Links != 1 {
		t.Errorf("Expected links 1, got %d", entries[0].Links)
	}
	if entries[0].Owner != "user" {
		t.Errorf("Expected owner 'user', got '%s'", entries[0].Owner)
	}
	if entries[0].Group != "group" {
		t.Errorf("Expected group 'group', got '%s'", entries[0].Group)
	}
	if entries[0].Size != 1024 {
		t.Errorf("Expected size 1024, got %d", entries[0].Size)
	}
	if entries[0].Path != "./file.txt" {
		t.Errorf("Expected path './file.txt', got '%s'", entries[0].Path)
	}

	// Test directory entry
	if entries[1].Type != "directory" {
		t.Errorf("Expected type 'directory', got '%s'", entries[1].Type)
	}
	if entries[1].Permissions != "drwxr-xr-x" {
		t.Errorf("Expected permissions 'drwxr-xr-x', got '%s'", entries[1].Permissions)
	}

	// Test symlink entry
	if entries[2].Type != "symlink" {
		t.Errorf("Expected type 'symlink', got '%s'", entries[2].Type)
	}
	if entries[2].Permissions != "lrwxrwxrwx" {
		t.Errorf("Expected permissions 'lrwxrwxrwx', got '%s'", entries[2].Permissions)
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

func TestFindParserSimple(t *testing.T) {
	parser := &FindParser{}

	// Test simple path-only format
	testInput := `./file1.txt
./dir/file2.txt
./another_dir/`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]FindEntry)
	if !ok {
		t.Fatalf("Expected []FindEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	if entries[0].Path != "./file1.txt" {
		t.Errorf("Expected path './file1.txt', got '%s'", entries[0].Path)
	}
	if entries[0].Type != "unknown" {
		t.Errorf("Expected type 'unknown', got '%s'", entries[0].Type)
	}
}

func TestFindParserEmpty(t *testing.T) {
	parser := &FindParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
