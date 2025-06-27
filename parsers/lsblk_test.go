package parsers

import (
	"encoding/json"
	"testing"
)

func TestLsblkParser(t *testing.T) {
	parser := &LsblkParser{}

	testInput := `NAME      MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
sda         8:0    0 931.5G  0 disk 
├─sda1      8:1    0   512M  0 part /boot/efi
├─sda2      8:2    0     1G  0 part /boot
└─sda3      8:3    0   930G  0 part /
sr0        11:0    1  1024M  0 rom  
nvme0n1   259:0    0   256G  0 disk 
└─nvme0n1p1 259:1  0   256G  0 part /home`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]LsblkEntry)
	if !ok {
		t.Fatalf("Expected []LsblkEntry, got %T", result)
	}

	if len(entries) != 7 {
		t.Fatalf("Expected 7 entries, got %d", len(entries))
	}

	// Test first entry (sda)
	if entries[0].Name != "sda" {
		t.Errorf("Expected name 'sda', got '%s'", entries[0].Name)
	}
	if entries[0].MajMin != "8:0" {
		t.Errorf("Expected maj:min '8:0', got '%s'", entries[0].MajMin)
	}
	if entries[0].Rm != "0" {
		t.Errorf("Expected rm '0', got '%s'", entries[0].Rm)
	}
	if entries[0].Size != "931.5G" {
		t.Errorf("Expected size '931.5G', got '%s'", entries[0].Size)
	}
	if entries[0].Ro != "0" {
		t.Errorf("Expected ro '0', got '%s'", entries[0].Ro)
	}
	if entries[0].Type != "disk" {
		t.Errorf("Expected type 'disk', got '%s'", entries[0].Type)
	}
	if entries[0].Mountpoint != "" {
		t.Errorf("Expected empty mountpoint, got '%s'", entries[0].Mountpoint)
	}

	// Test partition with mountpoint (sda1)
	if entries[1].Name != "├─sda1" {
		t.Errorf("Expected name '├─sda1', got '%s'", entries[1].Name)
	}
	if entries[1].Type != "part" {
		t.Errorf("Expected type 'part', got '%s'", entries[1].Type)
	}
	if entries[1].Mountpoint != "/boot/efi" {
		t.Errorf("Expected mountpoint '/boot/efi', got '%s'", entries[1].Mountpoint)
	}

	// Test ROM device
	romFound := false
	for _, entry := range entries {
		if entry.Type == "rom" {
			romFound = true
			if entry.Name != "sr0" {
				t.Errorf("Expected ROM device name 'sr0', got '%s'", entry.Name)
			}
			break
		}
	}
	if !romFound {
		t.Error("ROM device not found")
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

func TestLsblkParserSimple(t *testing.T) {
	parser := &LsblkParser{}

	testInput := `NAME MAJ:MIN RM SIZE RO TYPE
sda    8:0    0  100G  0 disk
sdb    8:16   0   50G  0 disk`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]LsblkEntry)
	if !ok {
		t.Fatalf("Expected []LsblkEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Both should have empty mountpoints
	for _, entry := range entries {
		if entry.Mountpoint != "" {
			t.Errorf("Expected empty mountpoint, got '%s'", entry.Mountpoint)
		}
	}
}

func TestLsblkParserEmpty(t *testing.T) {
	parser := &LsblkParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with only header
	result, err := parser.Parse("NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINT")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]LsblkEntry)
	if !ok {
		t.Fatalf("Expected []LsblkEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
