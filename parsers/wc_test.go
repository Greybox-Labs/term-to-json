package parsers

import (
	"encoding/json"
	"testing"
)

func TestWcParser(t *testing.T) {
	parser := &WcParser{}

	testInput := `  123   456   789 test.txt`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(WcEntry)
	if !ok {
		t.Fatalf("Expected WcEntry, got %T", result)
	}

	// Test all fields for standard wc output
	if entry.Lines != 123 {
		t.Errorf("Expected lines 123, got %d", entry.Lines)
	}
	if entry.Words != 456 {
		t.Errorf("Expected words 456, got %d", entry.Words)
	}
	if entry.Characters != 789 {
		t.Errorf("Expected characters 789, got %d", entry.Characters)
	}
	if entry.Bytes != 789 {
		t.Errorf("Expected bytes 789, got %d", entry.Bytes)
	}
	if entry.Filename != "test.txt" {
		t.Errorf("Expected filename 'test.txt', got '%s'", entry.Filename)
	}
	if entry.Original == "" {
		t.Error("Expected original to be populated")
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

func TestWcParserMultipleFiles(t *testing.T) {
	parser := &WcParser{}

	testInput := `  10   20   30 file1.txt
  15   25   35 file2.txt
  25   45   65 total`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]WcEntry)
	if !ok {
		t.Fatalf("Expected []WcEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test first file
	if entries[0].Lines != 10 {
		t.Errorf("Expected lines 10, got %d", entries[0].Lines)
	}
	if entries[0].Filename != "file1.txt" {
		t.Errorf("Expected filename 'file1.txt', got '%s'", entries[0].Filename)
	}

	// Test second file
	if entries[1].Words != 25 {
		t.Errorf("Expected words 25, got %d", entries[1].Words)
	}
	if entries[1].Filename != "file2.txt" {
		t.Errorf("Expected filename 'file2.txt', got '%s'", entries[1].Filename)
	}

	// Test total line
	if entries[2].Lines != 25 {
		t.Errorf("Expected total lines 25, got %d", entries[2].Lines)
	}
	if entries[2].Filename != "total" {
		t.Errorf("Expected filename 'total', got '%s'", entries[2].Filename)
	}
}

func TestWcParserLinesOnly(t *testing.T) {
	parser := &WcParser{}

	testInput := `42 myfile.txt`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(WcEntry)
	if !ok {
		t.Fatalf("Expected WcEntry, got %T", result)
	}

	// Test single number format (could be lines, words, or characters)
	if entry.Lines != 42 {
		t.Errorf("Expected lines 42, got %d", entry.Lines)
	}
	if entry.Words != 42 {
		t.Errorf("Expected words 42, got %d", entry.Words)
	}
	if entry.Characters != 42 {
		t.Errorf("Expected characters 42, got %d", entry.Characters)
	}
	if entry.Filename != "myfile.txt" {
		t.Errorf("Expected filename 'myfile.txt', got '%s'", entry.Filename)
	}
}

func TestWcParserNoFilename(t *testing.T) {
	parser := &WcParser{}

	testInput := `   5   10   25`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(WcEntry)
	if !ok {
		t.Fatalf("Expected WcEntry, got %T", result)
	}

	// Test without filename
	if entry.Lines != 5 {
		t.Errorf("Expected lines 5, got %d", entry.Lines)
	}
	if entry.Words != 10 {
		t.Errorf("Expected words 10, got %d", entry.Words)
	}
	if entry.Characters != 25 {
		t.Errorf("Expected characters 25, got %d", entry.Characters)
	}
	if entry.Filename != "" {
		t.Errorf("Expected empty filename, got '%s'", entry.Filename)
	}
}

func TestWcParserSingleNumber(t *testing.T) {
	parser := &WcParser{}

	testInput := `100`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(WcEntry)
	if !ok {
		t.Fatalf("Expected WcEntry, got %T", result)
	}

	// Test single number (ambiguous - could be any count)
	if entry.Lines != 100 {
		t.Errorf("Expected lines 100, got %d", entry.Lines)
	}
	if entry.Words != 100 {
		t.Errorf("Expected words 100, got %d", entry.Words)
	}
	if entry.Characters != 100 {
		t.Errorf("Expected characters 100, got %d", entry.Characters)
	}
}

func TestWcParserEmpty(t *testing.T) {
	parser := &WcParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
