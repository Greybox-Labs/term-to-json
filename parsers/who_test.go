package parsers

import (
	"encoding/json"
	"testing"
)

func TestWhoParser(t *testing.T) {
	parser := &WhoParser{}

	testInput := `user     pts/0        2023-01-15 14:30 (192.168.1.100)
root     tty1         2023-01-15 09:00
testuser pts/1        Jan 15 13:45 (:0)
admin    pts/2        2023-01-15 16:20 (server.example.com)`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]WhoEntry)
	if !ok {
		t.Fatalf("Expected []WhoEntry, got %T", result)
	}

	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries, got %d", len(entries))
	}

	// Test first entry with host
	if entries[0].User != "user" {
		t.Errorf("Expected user 'user', got '%s'", entries[0].User)
	}
	if entries[0].TTY != "pts/0" {
		t.Errorf("Expected TTY 'pts/0', got '%s'", entries[0].TTY)
	}
	if entries[0].Host != "192.168.1.100" {
		t.Errorf("Expected host '192.168.1.100', got '%s'", entries[0].Host)
	}
	if entries[0].LoginTime.IsZero() {
		t.Error("Expected login time to be parsed")
	}

	// Test root entry without host
	if entries[1].User != "root" {
		t.Errorf("Expected user 'root', got '%s'", entries[1].User)
	}
	if entries[1].TTY != "tty1" {
		t.Errorf("Expected TTY 'tty1', got '%s'", entries[1].TTY)
	}
	if entries[1].Host != "" {
		t.Errorf("Expected empty host, got '%s'", entries[1].Host)
	}

	// Test entry with X11 session
	if entries[2].User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", entries[2].User)
	}
	if entries[2].Host != ":0" {
		t.Errorf("Expected host ':0', got '%s'", entries[2].Host)
	}

	// Test entry with hostname
	if entries[3].Host != "server.example.com" {
		t.Errorf("Expected host 'server.example.com', got '%s'", entries[3].Host)
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

func TestWhoParserSimple(t *testing.T) {
	parser := &WhoParser{}

	testInput := `user     console      Jan 15 10:00
admin    pts/0        Jan 15 11:30`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]WhoEntry)
	if !ok {
		t.Fatalf("Expected []WhoEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test console user
	if entries[0].User != "user" {
		t.Errorf("Expected user 'user', got '%s'", entries[0].User)
	}
	if entries[0].TTY != "console" {
		t.Errorf("Expected TTY 'console', got '%s'", entries[0].TTY)
	}
	if entries[0].Host != "" {
		t.Errorf("Expected empty host, got '%s'", entries[0].Host)
	}
}

func TestWhoParserWithComments(t *testing.T) {
	parser := &WhoParser{}

	testInput := `user     pts/0        2023-01-15 14:30 old
admin    pts/1        2023-01-15 15:00 (remote.host) active`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]WhoEntry)
	if !ok {
		t.Fatalf("Expected []WhoEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test entry with comment
	if entries[0].Comment != "old" {
		t.Errorf("Expected comment 'old', got '%s'", entries[0].Comment)
	}

	// Test entry with host and comment
	if entries[1].Host != "remote.host" {
		t.Errorf("Expected host 'remote.host', got '%s'", entries[1].Host)
	}
	if entries[1].Comment != "active" {
		t.Errorf("Expected comment 'active', got '%s'", entries[1].Comment)
	}
}

func TestWhoParserDateFormats(t *testing.T) {
	parser := &WhoParser{}

	testInput := `user1    pts/0        2023-01-15 14:30
user2    pts/1        Jan 15 14:30
user3    pts/2        Jan  1 09:00`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]WhoEntry)
	if !ok {
		t.Fatalf("Expected []WhoEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Check that users were parsed correctly
	for i, entry := range entries {
		if entry.User == "" {
			t.Errorf("Expected user to be populated for entry %d", i)
		}
	}
}

func TestWhoParserEmpty(t *testing.T) {
	parser := &WhoParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with insufficient fields
	result, err := parser.Parse("user pts/0")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]WhoEntry)
	if !ok {
		t.Fatalf("Expected []WhoEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for insufficient fields, got %d", len(entries))
	}
}
