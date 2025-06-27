package parsers

import (
	"encoding/json"
	"testing"
)

func TestSystemctlParser(t *testing.T) {
	parser := &SystemctlParser{}

	testInput := `UNIT                               LOAD   ACTIVE SUB     DESCRIPTION
accounts-daemon.service            loaded active running Accounts Service
acpid.service                      loaded active running ACPI event daemon
apache2.service                    loaded active running The Apache HTTP Server
cron.service                       loaded active running Regular background program processing daemon
ssh.service                        loaded active running OpenBSD Secure Shell server

LOAD   = Reflects whether the unit definition was properly loaded.
ACTIVE = The high-level unit activation state, i.e. generalization of SUB.
SUB    = The low-level unit activation state, values depend on unit type.

5 loaded units listed. Pass --all to see loaded but inactive units, too.`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entries, ok := result.([]SystemctlEntry)
	if !ok {
		t.Fatalf("Expected []SystemctlEntry, got %T", result)
	}

	// Should skip summary lines and load explanation
	if len(entries) < 5 {
		t.Fatalf("Expected at least 5 entries, got %d", len(entries))
	}

	// Test first entry
	if entries[0].Unit != "accounts-daemon.service" {
		t.Errorf("Expected unit 'accounts-daemon.service', got '%s'", entries[0].Unit)
	}
	if entries[0].Load != "loaded" {
		t.Errorf("Expected load 'loaded', got '%s'", entries[0].Load)
	}
	if entries[0].Active != "active" {
		t.Errorf("Expected active 'active', got '%s'", entries[0].Active)
	}
	if entries[0].Sub != "running" {
		t.Errorf("Expected sub 'running', got '%s'", entries[0].Sub)
	}
	if entries[0].Description != "Accounts Service" {
		t.Errorf("Expected description 'Accounts Service', got '%s'", entries[0].Description)
	}

	// Test Apache entry
	apacheFound := false
	for _, entry := range entries {
		if entry.Unit == "apache2.service" {
			apacheFound = true
			if entry.Description != "The Apache HTTP Server" {
				t.Errorf("Expected description 'The Apache HTTP Server', got '%s'", entry.Description)
			}
			break
		}
	}
	if !apacheFound {
		t.Error("apache2.service not found")
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

func TestSystemctlParserStatus(t *testing.T) {
	parser := &SystemctlParser{}

	testInput := `● apache2.service - The Apache HTTP Server
   Loaded: loaded (/lib/systemd/system/apache2.service; enabled; vendor preset: enabled)
   Active: active (running) since Mon 2023-01-15 14:30:25 UTC; 2h 15min ago
     Docs: https://httpd.apache.org/docs/2.4/
   Main PID: 12345 (apache2)
    Tasks: 55 (limit: 4915)
   Memory: 123.4M
      CPU: 2.345s
   CGroup: /system.slice/apache2.service`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(SystemctlEntry)
	if !ok {
		t.Fatalf("Expected SystemctlEntry, got %T", result)
	}

	// Test service info
	if entry.Unit != "apache2.service" {
		t.Errorf("Expected unit 'apache2.service', got '%s'", entry.Unit)
	}
	if entry.Description != "The Apache HTTP Server" {
		t.Errorf("Expected description 'The Apache HTTP Server', got '%s'", entry.Description)
	}

	// Test status fields
	if entry.Load == "" {
		t.Error("Expected load to be populated")
	}
	if entry.Active == "" {
		t.Error("Expected active to be populated")
	}
	if entry.Main != "12345 (apache2)" {
		t.Errorf("Expected main '12345 (apache2)', got '%s'", entry.Main)
	}
	if entry.ProcessID != "12345 (apache2)" {
		t.Errorf("Expected process ID '12345 (apache2)', got '%s'", entry.ProcessID)
	}
	if entry.Tasks != "55 (limit: 4915)" {
		t.Errorf("Expected tasks '55 (limit: 4915)', got '%s'", entry.Tasks)
	}
	if entry.Memory != "123.4M" {
		t.Errorf("Expected memory '123.4M', got '%s'", entry.Memory)
	}
	if entry.CPU != "2.345s" {
		t.Errorf("Expected CPU '2.345s', got '%s'", entry.CPU)
	}
}

func TestSystemctlParserEmpty(t *testing.T) {
	parser := &SystemctlParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
	
	// Test with header only
	result, err := parser.Parse("UNIT LOAD ACTIVE SUB DESCRIPTION")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	entries, ok := result.([]SystemctlEntry)
	if !ok {
		t.Fatalf("Expected []SystemctlEntry, got %T", result)
	}
	
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
