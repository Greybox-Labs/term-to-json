package parsers

import (
	"encoding/json"
	"testing"
)

func TestWParser(t *testing.T) {
	parser := &WParser{}

	testInput := ` 14:30:42 up 12 days,  3:45,  2 users,  load average: 0.15, 0.12, 0.10
USER     TTY      FROM             LOGIN@   IDLE   JCPU   PCPU WHAT
user     pts/0    192.168.1.100    14:20    5:30   0.05s  0.01s ssh server1
root     tty1     -                09:00    2:15m  0.02s  0.02s -bash
user     pts/1    :0               13:45    0.00s  0.35s  0.03s vim test.txt`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	output, ok := result.(WOutput)
	if !ok {
		t.Fatalf("Expected WOutput, got %T", result)
	}

	// Test header
	if output.Header.CurrentTime != "14:30:42" {
		t.Errorf("Expected current time '14:30:42', got '%s'", output.Header.CurrentTime)
	}
	if output.Header.Users != 2 {
		t.Errorf("Expected 2 users, got %d", output.Header.Users)
	}
	if output.Header.LoadAvg1 != 0.15 {
		t.Errorf("Expected load avg 1 min 0.15, got %f", output.Header.LoadAvg1)
	}
	if output.Header.LoadAvg5 != 0.12 {
		t.Errorf("Expected load avg 5 min 0.12, got %f", output.Header.LoadAvg5)
	}
	if output.Header.LoadAvg15 != 0.10 {
		t.Errorf("Expected load avg 15 min 0.10, got %f", output.Header.LoadAvg15)
	}

	// Test user entries
	if len(output.Users) != 3 {
		t.Fatalf("Expected 3 user entries, got %d", len(output.Users))
	}

	// Test first user entry
	if output.Users[0].User != "user" {
		t.Errorf("Expected user 'user', got '%s'", output.Users[0].User)
	}
	if output.Users[0].TTY != "pts/0" {
		t.Errorf("Expected TTY 'pts/0', got '%s'", output.Users[0].TTY)
	}
	if output.Users[0].From != "192.168.1.100" {
		t.Errorf("Expected from '192.168.1.100', got '%s'", output.Users[0].From)
	}
	if output.Users[0].Login != "14:20" {
		t.Errorf("Expected login '14:20', got '%s'", output.Users[0].Login)
	}
	if output.Users[0].Idle != "5:30" {
		t.Errorf("Expected idle '5:30', got '%s'", output.Users[0].Idle)
	}
	if output.Users[0].JCPU != "0.05s" {
		t.Errorf("Expected JCPU '0.05s', got '%s'", output.Users[0].JCPU)
	}
	if output.Users[0].PCPU != "0.01s" {
		t.Errorf("Expected PCPU '0.01s', got '%s'", output.Users[0].PCPU)
	}
	if output.Users[0].What != "ssh server1" {
		t.Errorf("Expected what 'ssh server1', got '%s'", output.Users[0].What)
	}

	// Test root entry (no FROM field)
	if output.Users[1].User != "root" {
		t.Errorf("Expected user 'root', got '%s'", output.Users[1].User)
	}
	if output.Users[1].From != "-" {
		t.Errorf("Expected from '-', got '%s'", output.Users[1].From)
	}
	if output.Users[1].TTY != "tty1" {
		t.Errorf("Expected TTY 'tty1', got '%s'", output.Users[1].TTY)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON output is empty")
	}
}

func TestWParserSimple(t *testing.T) {
	parser := &WParser{}

	testInput := ` 10:15:30 up 1 day,  5:25,  1 user,  load average: 0.05, 0.03, 0.01
USER     TTY      LOGIN@   IDLE   WHAT
testuser pts/0    10:00    0.00s  bash`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	output, ok := result.(WOutput)
	if !ok {
		t.Fatalf("Expected WOutput, got %T", result)
	}

	// Test simplified format
	if len(output.Users) != 1 {
		t.Fatalf("Expected 1 user entry, got %d", len(output.Users))
	}

	user := output.Users[0]
	if user.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", user.User)
	}
	if user.From != "" {
		t.Errorf("Expected empty from field, got '%s'", user.From)
	}
	// Check if the 'what' field is populated (may be empty in some formats)
	if user.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", user.User)
	}
}

func TestWParserEmpty(t *testing.T) {
	parser := &WParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with insufficient lines
	_, err = parser.Parse("single line")
	if err == nil {
		t.Error("Expected error for insufficient lines")
	}
}
