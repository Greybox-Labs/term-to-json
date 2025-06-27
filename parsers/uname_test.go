package parsers

import (
	"testing"
)

func TestUnameParser(t *testing.T) {
	parser := &UnameParser{}

	testInput := "Linux hostname 5.4.0-74-generic #83-Ubuntu SMP Sat May 8 02:35:39 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux"

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(UnameEntry)
	if !ok {
		t.Fatalf("Expected UnameEntry, got %T", result)
	}

	if entry.KernelName != "Linux" {
		t.Errorf("Expected kernel name 'Linux', got '%s'", entry.KernelName)
	}
	if entry.NodeName != "hostname" {
		t.Errorf("Expected node name 'hostname', got '%s'", entry.NodeName)
	}
	if entry.KernelRelease != "5.4.0-74-generic" {
		t.Errorf("Expected kernel release '5.4.0-74-generic', got '%s'", entry.KernelRelease)
	}
	if entry.Machine != "x86_64" {
		t.Errorf("Expected machine 'x86_64', got '%s'", entry.Machine)
	}
	if entry.OS != "GNU/Linux" {
		t.Errorf("Expected OS 'GNU/Linux', got '%s'", entry.OS)
	}
}

func TestUnameParserEmpty(t *testing.T) {
	parser := &UnameParser{}
	
	_, err := parser.Parse("")
	if err == nil {
		t.Error("Expected error for empty input")
	}
	
	_, err = parser.Parse("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only input")
	}
}
