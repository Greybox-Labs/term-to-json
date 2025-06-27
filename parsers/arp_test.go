package parsers

import (
	"encoding/json"
	"testing"
)

func TestArpParser(t *testing.T) {
	parser := &ArpParser{}

	testInput := `Address                  HWtype  HWaddress           Flags Mask            Iface
192.168.1.1              ether   aa:bb:cc:dd:ee:ff   C     *               eth0
192.168.1.100            ether   11:22:33:44:55:66   C     *               eth0
192.168.1.200                    <incomplete>                              eth0`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]ArpEntry)
	if !ok {
		t.Fatalf("Expected []ArpEntry, got %T", result)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Test first entry
	if entries[0].Address != "192.168.1.1" {
		t.Errorf("Expected address '192.168.1.1', got '%s'", entries[0].Address)
	}
	if entries[0].HWType != "ether" {
		t.Errorf("Expected hwtype 'ether', got '%s'", entries[0].HWType)
	}
	if entries[0].HWAddress != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("Expected hwaddress 'aa:bb:cc:dd:ee:ff', got '%s'", entries[0].HWAddress)
	}
	if entries[0].Interface != "eth0" {
		t.Errorf("Expected interface 'eth0', got '%s'", entries[0].Interface)
	}
	if entries[0].Incomplete {
		t.Error("Expected first entry to not be incomplete")
	}

	// Test incomplete entry
	if !entries[2].Incomplete {
		t.Error("Expected third entry to be incomplete")
	}
	if entries[2].HWAddress != "<incomplete>" {
		t.Errorf("Expected hwaddress '<incomplete>', got '%s'", entries[2].HWAddress)
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

func TestArpParserEmpty(t *testing.T) {
	parser := &ArpParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with header only
	result, err := parser.Parse("Address                  HWtype  HWaddress           Flags Mask            Iface")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]ArpEntry)
	if !ok {
		t.Fatalf("Expected []ArpEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
