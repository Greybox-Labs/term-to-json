package parsers

import (
	"encoding/json"
	"testing"
)

func TestIdParser(t *testing.T) {
	parser := &IdParser{}

	testInput := "uid=1000(user) gid=1000(user) groups=1000(user),4(adm),24(cdrom),27(sudo),46(plugdev),116(lpadmin)"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(IdEntry)
	if !ok {
		t.Fatalf("Expected IdEntry, got %T", result)
	}

	// Test UID
	if entry.UID != 1000 {
		t.Errorf("Expected UID 1000, got %d", entry.UID)
	}
	if entry.User != "user" {
		t.Errorf("Expected user 'user', got '%s'", entry.User)
	}

	// Test GID
	if entry.GID != 1000 {
		t.Errorf("Expected GID 1000, got %d", entry.GID)
	}
	if entry.Group != "user" {
		t.Errorf("Expected group 'user', got '%s'", entry.Group)
	}

	// Test groups
	if len(entry.Groups) != 6 {
		t.Fatalf("Expected 6 groups, got %d", len(entry.Groups))
	}

	// Test first group
	if entry.Groups[0].GID != 1000 {
		t.Errorf("Expected group GID 1000, got %d", entry.Groups[0].GID)
	}
	if entry.Groups[0].Name != "user" {
		t.Errorf("Expected group name 'user', got '%s'", entry.Groups[0].Name)
	}

	// Test sudo group
	sudoFound := false
	for _, group := range entry.Groups {
		if group.Name == "sudo" {
			sudoFound = true
			if group.GID != 27 {
				t.Errorf("Expected sudo group GID 27, got %d", group.GID)
			}
			break
		}
	}
	if !sudoFound {
		t.Error("sudo group not found")
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON output is empty")
	}
}

func TestIdParserWithContext(t *testing.T) {
	parser := &IdParser{}

	testInput := "uid=0(root) gid=0(root) groups=0(root) context=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(IdEntry)
	if !ok {
		t.Fatalf("Expected IdEntry, got %T", result)
	}

	// Test root user
	if entry.UID != 0 {
		t.Errorf("Expected UID 0, got %d", entry.UID)
	}
	if entry.User != "root" {
		t.Errorf("Expected user 'root', got '%s'", entry.User)
	}

	// Test context
	if entry.Context != "unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023" {
		t.Errorf("Expected context 'unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023', got '%s'", entry.Context)
	}
}

func TestIdParserSimple(t *testing.T) {
	parser := &IdParser{}

	testInput := "uid=500(testuser) gid=500(testuser)"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(IdEntry)
	if !ok {
		t.Fatalf("Expected IdEntry, got %T", result)
	}

	if entry.UID != 500 {
		t.Errorf("Expected UID 500, got %d", entry.UID)
	}
	if entry.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", entry.User)
	}
	if len(entry.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(entry.Groups))
	}
}

func TestIdParserEmpty(t *testing.T) {
	parser := &IdParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
