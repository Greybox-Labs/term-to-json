package parsers

import (
	"encoding/json"
	"testing"
)

func TestMountParser(t *testing.T) {
	parser := &MountParser{}

	testInput := `/dev/sda1 on / type ext4 (rw,relatime,errors=remount-ro)
/dev/sda2 on /boot type ext2 (rw,relatime)
tmpfs on /tmp type tmpfs (rw,nosuid,nodev,noexec,relatime,size=4096000k)
/dev/nvme0n1p1 on /home type xfs (rw,relatime,attr2,inode64,logbufs=8,logbsize=32k,noquota)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]MountEntry)
	if !ok {
		t.Fatalf("Expected []MountEntry, got %T", result)
	}

	if len(entries) != 5 {
		t.Fatalf("Expected 5 entries, got %d", len(entries))
	}

	// Test first entry (root filesystem)
	if entries[0].Device != "/dev/sda1" {
		t.Errorf("Expected device '/dev/sda1', got '%s'", entries[0].Device)
	}
	if entries[0].MountPoint != "/" {
		t.Errorf("Expected mount point '/', got '%s'", entries[0].MountPoint)
	}
	if entries[0].FilesystemType != "ext4" {
		t.Errorf("Expected filesystem type 'ext4', got '%s'", entries[0].FilesystemType)
	}
	if entries[0].Options != "rw,relatime,errors=remount-ro" {
		t.Errorf("Expected options 'rw,relatime,errors=remount-ro', got '%s'", entries[0].Options)
	}

	// Test tmpfs entry
	tmpfsFound := false
	for _, entry := range entries {
		if entry.Device == "tmpfs" {
			tmpfsFound = true
			if entry.MountPoint != "/tmp" {
				t.Errorf("Expected mount point '/tmp', got '%s'", entry.MountPoint)
			}
			if entry.FilesystemType != "tmpfs" {
				t.Errorf("Expected filesystem type 'tmpfs', got '%s'", entry.FilesystemType)
			}
			if entry.Options != "rw,nosuid,nodev,noexec,relatime,size=4096000k" {
				t.Errorf("Expected options 'rw,nosuid,nodev,noexec,relatime,size=4096000k', got '%s'", entry.Options)
			}
			break
		}
	}
	if !tmpfsFound {
		t.Error("tmpfs entry not found")
	}

	// Test proc entry
	if entries[4].Device != "proc" {
		t.Errorf("Expected device 'proc', got '%s'", entries[4].Device)
	}
	if entries[4].MountPoint != "/proc" {
		t.Errorf("Expected mount point '/proc', got '%s'", entries[4].MountPoint)
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

func TestMountParserNoOptions(t *testing.T) {
	parser := &MountParser{}

	testInput := `/dev/sda1 on / type ext4
/dev/sdb1 on /home type xfs`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]MountEntry)
	if !ok {
		t.Fatalf("Expected []MountEntry, got %T", result)
	}

	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}

	// Test entries without options
	for _, entry := range entries {
		if entry.Options != "" {
			t.Errorf("Expected empty options, got '%s'", entry.Options)
		}
	}

	if entries[0].Device != "/dev/sda1" {
		t.Errorf("Expected device '/dev/sda1', got '%s'", entries[0].Device)
	}
	if entries[0].FilesystemType != "ext4" {
		t.Errorf("Expected filesystem type 'ext4', got '%s'", entries[0].FilesystemType)
	}
}

func TestMountParserEmpty(t *testing.T) {
	parser := &MountParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with invalid format
	result, err := parser.Parse("invalid mount line")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]MountEntry)
	if !ok {
		t.Fatalf("Expected []MountEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
