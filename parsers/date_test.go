package parsers

import (
	"testing"
)

func TestDateParser(t *testing.T) {
	parser := &DateParser{}

	testInput := "Wed Jan 15 14:30:25 PST 2025"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(DateEntry)
	if !ok {
		t.Fatalf("Expected DateEntry, got %T", result)
	}

	if entry.Weekday != "Wed" {
		t.Errorf("Expected weekday 'Wed', got '%s'", entry.Weekday)
	}
	if entry.Month != "Jan" {
		t.Errorf("Expected month 'Jan', got '%s'", entry.Month)
	}
	if entry.Day != 15 {
		t.Errorf("Expected day 15, got %d", entry.Day)
	}
	if entry.Time != "14:30:25" {
		t.Errorf("Expected time '14:30:25', got '%s'", entry.Time)
	}
	if entry.Timezone != "PST" {
		t.Errorf("Expected timezone 'PST', got '%s'", entry.Timezone)
	}
	if entry.Year != 2025 {
		t.Errorf("Expected year 2025, got %d", entry.Year)
	}
	if entry.Original != testInput {
		t.Errorf("Expected original '%s', got '%s'", testInput, entry.Original)
	}

	// Check that timestamp is populated
	if entry.Timestamp == 0 {
		t.Error("Expected timestamp to be populated")
	}
	if entry.Unix == 0 {
		t.Error("Expected unix timestamp to be populated")
	}
	if entry.ISO == "" {
		t.Error("Expected ISO format to be populated")
	}
}

func TestDateParserUnixTimestamp(t *testing.T) {
	parser := &DateParser{}

	testInput := "1705321825"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(DateEntry)
	if !ok {
		t.Fatalf("Expected DateEntry, got %T", result)
	}

	if entry.Unix != 1705321825 {
		t.Errorf("Expected unix timestamp 1705321825, got %d", entry.Unix)
	}
	if entry.Timestamp != 1705321825 {
		t.Errorf("Expected timestamp 1705321825, got %d", entry.Timestamp)
	}
	if entry.ISO == "" {
		t.Error("Expected ISO format to be populated")
	}
}

func TestDateParserEmpty(t *testing.T) {
	parser := &DateParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
