package parsers

import (
	"encoding/json"
	"testing"
)

func TestEnvParser(t *testing.T) {
	parser := &EnvParser{}

	testInput := `PATH=/usr/local/bin:/usr/bin:/bin
HOME=/home/user
USER=testuser
SHELL=/bin/bash
TERM=xterm-256color
PWD=/home/user/projects
LANG=en_US.UTF-8`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]EnvEntry)
	if !ok {
		t.Fatalf("Expected []EnvEntry, got %T", result)
	}

	if len(entries) != 7 {
		t.Fatalf("Expected 7 entries, got %d", len(entries))
	}

	// Test first entry (PATH)
	if entries[0].Name != "PATH" {
		t.Errorf("Expected name 'PATH', got '%s'", entries[0].Name)
	}
	if entries[0].Value != "/usr/local/bin:/usr/bin:/bin" {
		t.Errorf("Expected value '/usr/local/bin:/usr/bin:/bin', got '%s'", entries[0].Value)
	}

	// Test HOME entry
	homeFound := false
	for _, entry := range entries {
		if entry.Name == "HOME" {
			homeFound = true
			if entry.Value != "/home/user" {
				t.Errorf("Expected HOME value '/home/user', got '%s'", entry.Value)
			}
			break
		}
	}
	if !homeFound {
		t.Error("HOME environment variable not found")
	}

	// Test entry with equals in value
	testInputWithEquals := `MY_VAR=value=with=equals
ANOTHER_VAR=simple_value`
	
	result2, err := parser.Parse(testInputWithEquals)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries2, ok := result2.([]EnvEntry)
	if !ok {
		t.Fatalf("Expected []EnvEntry, got %T", result2)
	}

	if entries2[0].Name != "MY_VAR" {
		t.Errorf("Expected name 'MY_VAR', got '%s'", entries2[0].Name)
	}
	if entries2[0].Value != "value=with=equals" {
		t.Errorf("Expected value 'value=with=equals', got '%s'", entries2[0].Value)
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

func TestEnvParserEmpty(t *testing.T) {
	parser := &EnvParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only empty lines (this would cause an error)
	_, err = parser.Parse("\n\n")
	if err == nil {
		t.Error("Expected error for input with only empty lines")
	}
}
