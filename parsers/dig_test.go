package parsers

import (
	"encoding/json"
	"testing"
)

func TestDigParser(t *testing.T) {
	parser := &DigParser{}

	testInput := `;; QUESTION SECTION:
;; example.com.			IN	A

;; ANSWER SECTION:
example.com.		86400	IN	A	93.184.216.34
example.com.		86400	IN	A	93.184.216.35

;; AUTHORITY SECTION:
example.com.		172800	IN	NS	a.iana-servers.net.
example.com.		172800	IN	NS	b.iana-servers.net.

;; Query time: 45 msec
;; SERVER: 8.8.8.8#53(8.8.8.8)
;; WHEN: Mon Jan 15 14:30:25 PST 2025
;; MSG SIZE  rcvd: 123`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(DigEntry)
	if !ok {
		t.Fatalf("Expected DigEntry, got %T", result)
	}

	// Test query section
	if entry.Query.Name != "example.com" {
		t.Errorf("Expected query name 'example.com', got '%s'", entry.Query.Name)
	}
	if entry.Query.Class != "IN" {
		t.Errorf("Expected query class 'IN', got '%s'", entry.Query.Class)
	}
	if entry.Query.Type != "A" {
		t.Errorf("Expected query type 'A', got '%s'", entry.Query.Type)
	}

	// Test answer section
	if len(entry.Answer) != 2 {
		t.Fatalf("Expected 2 answer records, got %d", len(entry.Answer))
	}
	if entry.Answer[0].Name != "example.com." {
		t.Errorf("Expected answer name 'example.com.', got '%s'", entry.Answer[0].Name)
	}
	if entry.Answer[0].TTL != 86400 {
		t.Errorf("Expected TTL 86400, got %d", entry.Answer[0].TTL)
	}
	if entry.Answer[0].Value != "93.184.216.34" {
		t.Errorf("Expected value '93.184.216.34', got '%s'", entry.Answer[0].Value)
	}

	// Test authority section
	if len(entry.Authority) != 2 {
		t.Fatalf("Expected 2 authority records, got %d", len(entry.Authority))
	}
	if entry.Authority[0].Type != "NS" {
		t.Errorf("Expected authority type 'NS', got '%s'", entry.Authority[0].Type)
	}

	// Test that some stats are parsed
	if entry.Original == "" {
		t.Error("Expected original field to be populated")
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

func TestDigParserEmpty(t *testing.T) {
	parser := &DigParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with minimal valid input
	result, err := parser.Parse(";; QUESTION SECTION:\n;example.com. IN A")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entry, ok := result.(DigEntry)
	if !ok {
		t.Fatalf("Expected DigEntry, got %T", result)
	}
	
	if len(entry.Answer) != 0 {
		t.Errorf("Expected 0 answer records, got %d", len(entry.Answer))
	}
}
