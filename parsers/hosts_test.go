package parsers

import (
	"encoding/json"
	"testing"
)

func TestHostsParser(t *testing.T) {
	parser := &HostsParser{}

	testInput := `127.0.0.1	localhost localhost.localdomain
::1		localhost ip6-localhost ip6-loopback
192.168.1.100	server1.example.com server1	# web server
192.168.1.101	server2.example.com server2
# This is a comment line
10.0.0.1	internal.example.com`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]HostsEntry)
	if !ok {
		t.Fatalf("Expected []HostsEntry, got %T", result)
	}

	if len(entries) != 6 {
		t.Fatalf("Expected 6 entries, got %d", len(entries))
	}

	// Test first entry (localhost)
	if entries[0].IP != "127.0.0.1" {
		t.Errorf("Expected IP '127.0.0.1', got '%s'", entries[0].IP)
	}
	if len(entries[0].Hostnames) != 2 {
		t.Fatalf("Expected 2 hostnames, got %d", len(entries[0].Hostnames))
	}
	if entries[0].Hostnames[0] != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", entries[0].Hostnames[0])
	}
	if entries[0].Hostnames[1] != "localhost.localdomain" {
		t.Errorf("Expected hostname 'localhost.localdomain', got '%s'", entries[0].Hostnames[1])
	}

	// Test IPv6 entry
	if entries[1].IP != "::1" {
		t.Errorf("Expected IP '::1', got '%s'", entries[1].IP)
	}

	// Test entry with comment
	if entries[2].IP != "192.168.1.100" {
		t.Errorf("Expected IP '192.168.1.100', got '%s'", entries[2].IP)
	}
	if entries[2].Comment != "web server" {
		t.Errorf("Expected comment 'web server', got '%s'", entries[2].Comment)
	}
	if len(entries[2].Hostnames) != 2 {
		t.Fatalf("Expected 2 hostnames, got %d", len(entries[2].Hostnames))
	}

	// Test comment-only line
	if entries[4].IP != "" {
		t.Errorf("Expected empty IP for comment line, got '%s'", entries[4].IP)
	}
	if entries[4].Comment != "This is a comment line" {
		t.Errorf("Expected comment 'This is a comment line', got '%s'", entries[4].Comment)
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

func TestHostsParserIPOnly(t *testing.T) {
	parser := &HostsParser{}

	testInput := `192.168.1.1
10.0.0.1		# IP only with comment`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]HostsEntry)
	if !ok {
		t.Fatalf("Expected []HostsEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test IP-only entry
	if entries[0].IP != "192.168.1.1" {
		t.Errorf("Expected IP '192.168.1.1', got '%s'", entries[0].IP)
	}
	if len(entries[0].Hostnames) != 0 {
		t.Errorf("Expected 0 hostnames, got %d", len(entries[0].Hostnames))
	}

	// Test IP with comment
	if entries[1].Comment != "IP only with comment" {
		t.Errorf("Expected comment 'IP only with comment', got '%s'", entries[1].Comment)
	}
}

func TestHostsParserEmpty(t *testing.T) {
	parser := &HostsParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only empty lines and comments
	result, err := parser.Parse("\n# Comment\n\n")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]HostsEntry)
	if !ok {
		t.Fatalf("Expected []HostsEntry, got %T", result)
	}
	
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry (comment), got %d", len(entries))
	}
}
