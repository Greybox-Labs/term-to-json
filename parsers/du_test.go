package parsers

import (
	"encoding/json"
	"testing"
)

func TestDuParser(t *testing.T) {
	parser := &DuParser{}

	testInput := `1024	./docs
512	./src/main
2048	./src/tests
256	./config
4096	.`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]DuEntry)
	if !ok {
		t.Fatalf("Expected []DuEntry, got %T", result)
	}

	if len(entries) != 5 {
		t.Fatalf("Expected 5 entries, got %d", len(entries))
	}

	// Test first entry
	if entries[0].Size != 1024 {
		t.Errorf("Expected size 1024, got %d", entries[0].Size)
	}
	if entries[0].SizeBytes != 1024*1024 {
		t.Errorf("Expected size bytes %d, got %d", 1024*1024, entries[0].SizeBytes)
	}
	if entries[0].Path != "./docs" {
		t.Errorf("Expected path './docs', got '%s'", entries[0].Path)
	}

	// Test path with spaces
	testInputWithSpaces := `512	./My Documents/file.txt
256	./Program Files/app`
	
	result2, err := parser.Parse(testInputWithSpaces)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries2, ok := result2.([]DuEntry)
	if !ok {
		t.Fatalf("Expected []DuEntry, got %T", result2)
	}

	if entries2[0].Path != "./My Documents/file.txt" {
		t.Errorf("Expected path './My Documents/file.txt', got '%s'", entries2[0].Path)
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

func TestDuParserEmpty(t *testing.T) {
	parser := &DuParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with single field (insufficient data)
	result, err := parser.Parse("single_field")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]DuEntry)
	if !ok {
		t.Fatalf("Expected []DuEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
