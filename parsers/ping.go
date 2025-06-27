package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PingParser parses ping command output
type PingParser struct{}

// PingEntry represents ping output
type PingEntry struct {
	Type           string       `json:"type"`
	Pattern        string       `json:"pattern,omitempty"`
	Timestamp      string       `json:"timestamp,omitempty"`
	ResponseTime   float64      `json:"response_time_ms,omitempty"`
	Bytes          int          `json:"bytes,omitempty"`
	Response       string       `json:"response,omitempty"`
	Destination    string       `json:"destination"`
	DestinationIP  string       `json:"destination_ip,omitempty"`
	Packets        []PingPacket `json:"packets,omitempty"`
	Statistics     *PingStats   `json:"statistics,omitempty"`
}

// PingPacket represents individual ping packet
type PingPacket struct {
	Bytes        int     `json:"bytes"`
	Destination  string  `json:"destination"`
	DestinationIP string `json:"destination_ip"`
	ICMPSeq      int     `json:"icmp_seq"`
	TTL          int     `json:"ttl"`
	Time         float64 `json:"time_ms"`
	Duplicate    bool    `json:"duplicate,omitempty"`
}

// PingStats represents ping statistics
type PingStats struct {
	PacketsTransmitted int     `json:"packets_transmitted"`
	PacketsReceived    int     `json:"packets_received"`
	PacketLoss         float64 `json:"packet_loss_percent"`
	Time               int     `json:"time_ms"`
	RTTMin             float64 `json:"rtt_min_ms,omitempty"`
	RTTAvg             float64 `json:"rtt_avg_ms,omitempty"`
	RTTMax             float64 `json:"rtt_max_ms,omitempty"`
	RTTMdev            float64 `json:"rtt_mdev_ms,omitempty"`
}

func (p *PingParser) Name() string {
	return "ping"
}

func (p *PingParser) Parse(input string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	lines := splitLines(input)
	entry := PingEntry{Type: "ping"}
	
	var packets []PingPacket
	var destination, destinationIP string

	for _, line := range lines {
		// Parse PING header
		if strings.HasPrefix(line, "PING ") {
			pingHeaderRe := regexp.MustCompile(`PING ([^\s]+) \(([^)]+)\)`)
			if matches := pingHeaderRe.FindStringSubmatch(line); len(matches) > 2 {
				destination = matches[1]
				destinationIP = matches[2]
				entry.Destination = destination
				entry.DestinationIP = destinationIP
			}
			continue
		}

		// Parse ping responses
		if strings.Contains(line, "bytes from") {
			packet := parsePingResponse(line)
			if packet.Destination == "" {
				packet.Destination = destination
				packet.DestinationIP = destinationIP
			}
			packets = append(packets, packet)
			continue
		}

		// Parse statistics
		if strings.Contains(line, "packets transmitted") {
			stats := parsePingStats(line)
			entry.Statistics = &stats
			continue
		}

		// Parse RTT statistics
		if strings.HasPrefix(line, "round-trip") || strings.HasPrefix(line, "rtt") {
			if entry.Statistics != nil {
				parseRTTStats(line, entry.Statistics)
			}
			continue
		}
	}

	entry.Packets = packets
	return entry, nil
}

func parsePingResponse(line string) PingPacket {
	packet := PingPacket{}

	// Example: "64 bytes from google.com (142.250.191.14): icmp_seq=1 ttl=55 time=12.3 ms"
	bytesRe := regexp.MustCompile(`(\d+) bytes from`)
	if matches := bytesRe.FindStringSubmatch(line); len(matches) > 1 {
		if bytes, err := strconv.Atoi(matches[1]); err == nil {
			packet.Bytes = bytes
		}
	}

	hostRe := regexp.MustCompile(`from ([^\s]+) \(([^)]+)\)`)
	if matches := hostRe.FindStringSubmatch(line); len(matches) > 2 {
		packet.Destination = matches[1]
		packet.DestinationIP = matches[2]
	}

	seqRe := regexp.MustCompile(`icmp_seq=(\d+)`)
	if matches := seqRe.FindStringSubmatch(line); len(matches) > 1 {
		if seq, err := strconv.Atoi(matches[1]); err == nil {
			packet.ICMPSeq = seq
		}
	}

	ttlRe := regexp.MustCompile(`ttl=(\d+)`)
	if matches := ttlRe.FindStringSubmatch(line); len(matches) > 1 {
		if ttl, err := strconv.Atoi(matches[1]); err == nil {
			packet.TTL = ttl
		}
	}

	timeRe := regexp.MustCompile(`time=([0-9.]+)`)
	if matches := timeRe.FindStringSubmatch(line); len(matches) > 1 {
		if time, err := strconv.ParseFloat(matches[1], 64); err == nil {
			packet.Time = time
		}
	}

	if strings.Contains(line, "DUP!") {
		packet.Duplicate = true
	}

	return packet
}

func parsePingStats(line string) PingStats {
	stats := PingStats{}

	// Example: "5 packets transmitted, 5 received, 0% packet loss, time 4005ms"
	statsRe := regexp.MustCompile(`(\d+) packets transmitted, (\d+) received, ([0-9.]+)% packet loss, time (\d+)ms`)
	if matches := statsRe.FindStringSubmatch(line); len(matches) > 4 {
		if transmitted, err := strconv.Atoi(matches[1]); err == nil {
			stats.PacketsTransmitted = transmitted
		}
		if received, err := strconv.Atoi(matches[2]); err == nil {
			stats.PacketsReceived = received
		}
		if loss, err := strconv.ParseFloat(matches[3], 64); err == nil {
			stats.PacketLoss = loss
		}
		if time, err := strconv.Atoi(matches[4]); err == nil {
			stats.Time = time
		}
	}

	return stats
}

func parseRTTStats(line string, stats *PingStats) {
	// Example: "rtt min/avg/max/mdev = 12.123/15.456/18.789/2.345 ms"
	rttRe := regexp.MustCompile(`([0-9.]+)/([0-9.]+)/([0-9.]+)/([0-9.]+)`)
	if matches := rttRe.FindStringSubmatch(line); len(matches) > 4 {
		if min, err := strconv.ParseFloat(matches[1], 64); err == nil {
			stats.RTTMin = min
		}
		if avg, err := strconv.ParseFloat(matches[2], 64); err == nil {
			stats.RTTAvg = avg
		}
		if max, err := strconv.ParseFloat(matches[3], 64); err == nil {
			stats.RTTMax = max
		}
		if mdev, err := strconv.ParseFloat(matches[4], 64); err == nil {
			stats.RTTMdev = mdev
		}
	}
}
