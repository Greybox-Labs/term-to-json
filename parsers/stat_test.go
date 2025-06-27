package parsers

import (
	"encoding/json"
	"testing"
)

func TestStatParser(t *testing.T) {
	parser := &StatParser{}

	testInput := `  File: 'test.txt'
  Size: 1024      	Blocks: 8          IO Block: 4096   regular file
Device: 801h/2049d	Inode: 123456      Links: 1
Access: (0644/-rw-r--r--)  Uid: (1000/   user)   Gid: (1000/  group)
Access: 2023-01-15 14:30:25.000000000 +0000
Modify: 2023-01-15 14:25:10.000000000 +0000
Change: 2023-01-15 14:25:10.000000000 +0000`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]StatEntry)
	if !ok {
		t.Fatalf("Expected []StatEntry, got %T", result)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]

	// Test file info
	if entry.File != "test.txt" {
		t.Errorf("Expected file 'test.txt', got '%s'", entry.File)
	}
	if entry.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", entry.Size)
	}
	if entry.Blocks != 8 {
		t.Errorf("Expected blocks 8, got %d", entry.Blocks)
	}
	if entry.IOBlock != 4096 {
		t.Errorf("Expected IO block 4096, got %d", entry.IOBlock)
	}

	// Test device and inode
	if entry.Device != "801h/2049d" {
		t.Errorf("Expected device '801h/2049d', got '%s'", entry.Device)
	}
	if entry.Inode != 123456 {
		t.Errorf("Expected inode 123456, got %d", entry.Inode)
	}
	if entry.Links != 1 {
		t.Errorf("Expected links 1, got %d", entry.Links)
	}

	// Test permissions and ownership
	if entry.Permissions != "-rw-r--r--" {
		t.Errorf("Expected permissions '-rw-r--r--', got '%s'", entry.Permissions)
	}
	if entry.UID != 1000 {
		t.Errorf("Expected UID 1000, got %d", entry.UID)
	}
	if entry.GID != 1000 {
		t.Errorf("Expected GID 1000, got %d", entry.GID)
	}

	// Test timestamps
	if entry.AccessTime.IsZero() {
		t.Error("Expected access time to be parsed")
	}
	if entry.ModifyTime.IsZero() {
		t.Error("Expected modify time to be parsed")
	}
	if entry.ChangeTime.IsZero() {
		t.Error("Expected change time to be parsed")
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

func TestStatParserMultipleFiles(t *testing.T) {
	parser := &StatParser{}

	testInput := `  File: 'file1.txt'
  Size: 512       	Blocks: 4          IO Block: 4096   regular file
Device: 801h/2049d	Inode: 111111      Links: 1
Access: (0644/-rw-r--r--)  Uid: (1000/   user)   Gid: (1000/  group)

  File: 'dir1'
  Size: 4096      	Blocks: 8          IO Block: 4096   directory
Device: 801h/2049d	Inode: 222222      Links: 2
Access: (0755/drwxr-xr-x)  Uid: (1000/   user)   Gid: (1000/  group)`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]StatEntry)
	if !ok {
		t.Fatalf("Expected []StatEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test first file
	if entries[0].File != "file1.txt" {
		t.Errorf("Expected file 'file1.txt', got '%s'", entries[0].File)
	}
	if entries[0].Size != 512 {
		t.Errorf("Expected size 512, got %d", entries[0].Size)
	}
	if entries[0].Inode != 111111 {
		t.Errorf("Expected inode 111111, got %d", entries[0].Inode)
	}

	// Test directory
	if entries[1].File != "dir1" {
		t.Errorf("Expected file 'dir1', got '%s'", entries[1].File)
	}
	if entries[1].Size != 4096 {
		t.Errorf("Expected size 4096, got %d", entries[1].Size)
	}
	if entries[1].Links != 2 {
		t.Errorf("Expected links 2, got %d", entries[1].Links)
	}
	if entries[1].Permissions != "drwxr-xr-x" {
		t.Errorf("Expected permissions 'drwxr-xr-x', got '%s'", entries[1].Permissions)
	}
}

func TestStatParserEmpty(t *testing.T) {
	parser := &StatParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with minimal valid input
	result, err := parser.Parse("File: 'test'")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]StatEntry)
	if !ok {
		t.Fatalf("Expected []StatEntry, got %T", result)
	}
	
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	
	if entries[0].File != "test" {
		t.Errorf("Expected file 'test', got '%s'", entries[0].File)
	}
}
