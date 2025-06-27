package parsers

import (
	"testing"
)

func TestPingParser(t *testing.T) {
	parser := &PingParser{}

	testInput := `PING google.com (142.250.191.14) 56(84) bytes of data.
64 bytes from google.com (142.250.191.14): icmp_seq=1 ttl=55 time=12.3 ms
64 bytes from google.com (142.250.191.14): icmp_seq=2 ttl=55 time=15.1 ms
--- google.com ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 12.3/13.7/15.1/1.4 ms`

	result, err := parser.Parse(testInput)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	entry, ok := result.(PingEntry)
	if !ok {
		t.Fatalf("Expected PingEntry, got %T", result)
	}

	if entry.Destination != "google.com" {
		t.Errorf("Expected destination 'google.com', got '%s'", entry.Destination)
	}

	if entry.DestinationIP != "142.250.191.14" {
		t.Errorf("Expected destination IP '142.250.191.14', got '%s'", entry.DestinationIP)
	}

	if len(entry.Packets) != 2 {
		t.Fatalf("Expected 2 packets, got %d", len(entry.Packets))
	}

	// Check first packet
	packet1 := entry.Packets[0]
	if packet1.ICMPSeq != 1 {
		t.Errorf("Expected icmp_seq 1, got %d", packet1.ICMPSeq)
	}
	if packet1.Time != 12.3 {
		t.Errorf("Expected time 12.3, got %f", packet1.Time)
	}
	if packet1.TTL != 55 {
		t.Errorf("Expected TTL 55, got %d", packet1.TTL)
	}

	// Check statistics
	if entry.Statistics == nil {
		t.Fatal("Expected statistics, got nil")
	}
	if entry.Statistics.PacketsTransmitted != 2 {
		t.Errorf("Expected 2 packets transmitted, got %d", entry.Statistics.PacketsTransmitted)
	}
	if entry.Statistics.PacketsReceived != 2 {
		t.Errorf("Expected 2 packets received, got %d", entry.Statistics.PacketsReceived)
	}
	if entry.Statistics.PacketLoss != 0 {
		t.Errorf("Expected 0%% packet loss, got %f", entry.Statistics.PacketLoss)
	}
	if entry.Statistics.RTTMin != 12.3 {
		t.Errorf("Expected RTT min 12.3, got %f", entry.Statistics.RTTMin)
	}
}
