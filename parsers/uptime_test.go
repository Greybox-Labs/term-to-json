package parsers

import (
	"encoding/json"
	"testing"
)

func TestUptimeParser(t *testing.T) {
	parser := &UptimeParser{}

	testInput := " 14:30:42 up 12 days,  3:45,  2 users,  load average: 0.15, 0.12, 0.10"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(UptimeEntry)
	if !ok {
		t.Fatalf("Expected UptimeEntry, got %T", result)
	}

	// Test current time
	if entry.CurrentTime != "14:30:42" {
		t.Errorf("Expected current time '14:30:42', got '%s'", entry.CurrentTime)
	}

	// Test uptime
	if entry.Uptime != "12 days,  3:45" {
		t.Errorf("Expected uptime '12 days,  3:45', got '%s'", entry.Uptime)
	}

	// Test uptime in seconds (12 days + 3 hours + 45 minutes)
	expectedSeconds := 12*24*3600 + 3*3600 + 45*60
	if entry.UptimeSeconds != expectedSeconds {
		t.Errorf("Expected uptime seconds %d, got %d", expectedSeconds, entry.UptimeSeconds)
	}

	// Test users
	if entry.Users != 2 {
		t.Errorf("Expected users 2, got %d", entry.Users)
	}

	// Test load averages
	if entry.LoadAvg1 != 0.15 {
		t.Errorf("Expected load avg 1 min 0.15, got %f", entry.LoadAvg1)
	}
	if entry.LoadAvg5 != 0.12 {
		t.Errorf("Expected load avg 5 min 0.12, got %f", entry.LoadAvg5)
	}
	if entry.LoadAvg15 != 0.10 {
		t.Errorf("Expected load avg 15 min 0.10, got %f", entry.LoadAvg15)
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

func TestUptimeParserShortUptime(t *testing.T) {
	parser := &UptimeParser{}

	testInput := " 09:15:30 up  1:23,  1 user,  load average: 0.25, 0.20, 0.18"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(UptimeEntry)
	if !ok {
		t.Fatalf("Expected UptimeEntry, got %T", result)
	}

	// Test short uptime (1 hour 23 minutes)
	expectedSeconds := 1*3600 + 23*60
	if entry.UptimeSeconds != expectedSeconds {
		t.Errorf("Expected uptime seconds %d, got %d", expectedSeconds, entry.UptimeSeconds)
	}

	// Test single user
	if entry.Users != 1 {
		t.Errorf("Expected users 1, got %d", entry.Users)
	}
}

func TestUptimeParserMinutesOnly(t *testing.T) {
	parser := &UptimeParser{}

	testInput := " 10:05:15 up 45 min,  3 users,  load average: 1.05, 0.95, 0.85"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(UptimeEntry)
	if !ok {
		t.Fatalf("Expected UptimeEntry, got %T", result)
	}

	// Test minutes-only uptime
	expectedSeconds := 45 * 60
	if entry.UptimeSeconds != expectedSeconds {
		t.Errorf("Expected uptime seconds %d, got %d", expectedSeconds, entry.UptimeSeconds)
	}

	// Test higher load averages
	if entry.LoadAvg1 != 1.05 {
		t.Errorf("Expected load avg 1 min 1.05, got %f", entry.LoadAvg1)
	}
}

func TestUptimeParserEmpty(t *testing.T) {
	parser := &UptimeParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
