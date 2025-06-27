package parsers

import (
	"encoding/json"
	"testing"
)

func TestPasswdParser(t *testing.T) {
	parser := &PasswdParser{}

	testInput := `root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
sys:x:3:3:sys:/dev:/usr/sbin/nologin
user:x:1000:1000:User,,,:/home/user:/bin/bash
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]PasswdEntry)
	if !ok {
		t.Fatalf("Expected []PasswdEntry, got %T", result)
	}

	if len(entries) != 6 {
		t.Fatalf("Expected 6 entries, got %d", len(entries))
	}

	// Test root entry
	if entries[0].Username != "root" {
		t.Errorf("Expected username 'root', got '%s'", entries[0].Username)
	}
	if entries[0].Password != "x" {
		t.Errorf("Expected password 'x', got '%s'", entries[0].Password)
	}
	if entries[0].UID != 0 {
		t.Errorf("Expected UID 0, got %d", entries[0].UID)
	}
	if entries[0].GID != 0 {
		t.Errorf("Expected GID 0, got %d", entries[0].GID)
	}
	if entries[0].GECOS != "root" {
		t.Errorf("Expected GECOS 'root', got '%s'", entries[0].GECOS)
	}
	if entries[0].HomeDir != "/root" {
		t.Errorf("Expected home dir '/root', got '%s'", entries[0].HomeDir)
	}
	if entries[0].Shell != "/bin/bash" {
		t.Errorf("Expected shell '/bin/bash', got '%s'", entries[0].Shell)
	}

	// Test regular user entry
	userFound := false
	for _, entry := range entries {
		if entry.Username == "user" {
			userFound = true
			if entry.UID != 1000 {
				t.Errorf("Expected UID 1000, got %d", entry.UID)
			}
			if entry.GID != 1000 {
				t.Errorf("Expected GID 1000, got %d", entry.GID)
			}
			if entry.GECOS != "User,,," {
				t.Errorf("Expected GECOS 'User,,,', got '%s'", entry.GECOS)
			}
			if entry.HomeDir != "/home/user" {
				t.Errorf("Expected home dir '/home/user', got '%s'", entry.HomeDir)
			}
			break
		}
	}
	if !userFound {
		t.Error("user entry not found")
	}

	// Test nobody entry (high UID)
	if entries[5].Username != "nobody" {
		t.Errorf("Expected username 'nobody', got '%s'", entries[5].Username)
	}
	if entries[5].UID != 65534 {
		t.Errorf("Expected UID 65534, got %d", entries[5].UID)
	}
	if entries[5].Shell != "/usr/sbin/nologin" {
		t.Errorf("Expected shell '/usr/sbin/nologin', got '%s'", entries[5].Shell)
	}

	// Test original field
	if entries[0].Original != "root:x:0:0:root:/root:/bin/bash" {
		t.Errorf("Expected original to match input line")
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

func TestPasswdParserWithComments(t *testing.T) {
	parser := &PasswdParser{}

	testInput := `# This is a comment
root:x:0:0:root:/root:/bin/bash

# Another comment
user:x:1000:1000:Regular User:/home/user:/bin/bash`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]PasswdEntry)
	if !ok {
		t.Fatalf("Expected []PasswdEntry, got %T", result)
	}

	// Should skip comments and empty lines
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	if entries[0].Username != "root" {
		t.Errorf("Expected username 'root', got '%s'", entries[0].Username)
	}
	if entries[1].Username != "user" {
		t.Errorf("Expected username 'user', got '%s'", entries[1].Username)
	}
}

func TestPasswdParserEmpty(t *testing.T) {
	parser := &PasswdParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only comments
	result, err := parser.Parse("# Comment only\n# Another comment")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]PasswdEntry)
	if !ok {
		t.Fatalf("Expected []PasswdEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
