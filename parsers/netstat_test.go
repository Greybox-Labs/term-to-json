package parsers

import (
	"encoding/json"
	"testing"
)

func TestNetstatParser(t *testing.T) {
	parser := &NetstatParser{}

	testInput := `Active Internet connections (w/o servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 192.168.1.100:22        192.168.1.1:54321       ESTABLISHED 1234/ssh
tcp        0      0 192.168.1.100:80        0.0.0.0:*               LISTEN      5678/apache2
udp        0      0 0.0.0.0:53              0.0.0.0:*                           9012/systemd-resolve`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]NetstatEntry)
	if !ok {
		t.Fatalf("Expected []NetstatEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test TCP ESTABLISHED connection
	if entries[0].Protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got '%s'", entries[0].Protocol)
	}
	if entries[0].RecvQ != 0 {
		t.Errorf("Expected RecvQ 0, got %d", entries[0].RecvQ)
	}
	if entries[0].SendQ != 0 {
		t.Errorf("Expected SendQ 0, got %d", entries[0].SendQ)
	}
	if entries[0].LocalAddress != "192.168.1.100:22" {
		t.Errorf("Expected local address '192.168.1.100:22', got '%s'", entries[0].LocalAddress)
	}
	if entries[0].ForeignAddress != "192.168.1.1:54321" {
		t.Errorf("Expected foreign address '192.168.1.1:54321', got '%s'", entries[0].ForeignAddress)
	}
	if entries[0].State != "ESTABLISHED" {
		t.Errorf("Expected state 'ESTABLISHED', got '%s'", entries[0].State)
	}
	if entries[0].PID != 1234 {
		t.Errorf("Expected PID 1234, got %d", entries[0].PID)
	}
	if entries[0].Program != "ssh" {
		t.Errorf("Expected program 'ssh', got '%s'", entries[0].Program)
	}

	// Test TCP LISTEN connection
	if entries[1].State != "LISTEN" {
		t.Errorf("Expected state 'LISTEN', got '%s'", entries[1].State)
	}
	if entries[1].PID != 5678 {
		t.Errorf("Expected PID 5678, got %d", entries[1].PID)
	}
	if entries[1].Program != "apache2" {
		t.Errorf("Expected program 'apache2', got '%s'", entries[1].Program)
	}

	// Test UDP connection (no state)
	if entries[2].Protocol != "udp" {
		t.Errorf("Expected protocol 'udp', got '%s'", entries[2].Protocol)
	}
	if entries[2].State != "" {
		t.Errorf("Expected empty state for UDP, got '%s'", entries[2].State)
	}
	if entries[2].Program != "systemd-resolve" {
		t.Errorf("Expected program 'systemd-resolve', got '%s'", entries[2].Program)
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

func TestNetstatParserSimple(t *testing.T) {
	parser := &NetstatParser{}

	// Test simpler netstat output without RecvQ/SendQ
	testInput := `Proto Local Address           Foreign Address         State
tcp   127.0.0.1:3306          0.0.0.0:*               LISTEN
tcp   192.168.1.100:443       192.168.1.50:12345      ESTABLISHED`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]NetstatEntry)
	if !ok {
		t.Fatalf("Expected []NetstatEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test that RecvQ and SendQ are 0 when not provided
	if entries[0].RecvQ != 0 {
		t.Errorf("Expected RecvQ 0, got %d", entries[0].RecvQ)
	}
	if entries[0].SendQ != 0 {
		t.Errorf("Expected SendQ 0, got %d", entries[0].SendQ)
	}

	if entries[0].LocalAddress != "127.0.0.1:3306" {
		t.Errorf("Expected local address '127.0.0.1:3306', got '%s'", entries[0].LocalAddress)
	}
}

func TestNetstatParserEmpty(t *testing.T) {
	parser := &NetstatParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only headers
	result, err := parser.Parse("Proto Local Address Foreign Address State")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]NetstatEntry)
	if !ok {
		t.Fatalf("Expected []NetstatEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
