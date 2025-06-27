package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"term-to-json/parsers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <parser> [input]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Available parsers:\n")
		fmt.Fprintf(os.Stderr, "  System: uname, uptime, who, w, id, env\n")
		fmt.Fprintf(os.Stderr, "  Process: ps, free, vmstat\n")
		fmt.Fprintf(os.Stderr, "  Network: ping, netstat, arp, dig\n")
		fmt.Fprintf(os.Stderr, "  Files: ls, df, du, mount, lsblk, find, stat\n")
		fmt.Fprintf(os.Stderr, "  Services: systemctl\n")
		fmt.Fprintf(os.Stderr, "  Utilities: date, wc\n")
		fmt.Fprintf(os.Stderr, "  Config: hosts, passwd\n")
		os.Exit(1)
	}

	parserName := os.Args[1]
	var input string

	if len(os.Args) > 2 {
		input = os.Args[2]
	} else {
		// Read from stdin
		buf, err := os.ReadFile("/dev/stdin")
		if err != nil {
			log.Fatalf("Error reading stdin: %v", err)
		}
		input = string(buf)
	}

	result, err := parsers.Parse(parserName, input)
	if err != nil {
		log.Fatalf("Error parsing: %v", err)
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
