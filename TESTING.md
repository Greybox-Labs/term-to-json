# Testing Guide for term-to-json

## Quick Start Testing

### 1. Run All Tests
```bash
# Run Go unit tests
go test ./parsers -v

# Run demonstration script with sample data
./test_parsers.sh

# Test with real system commands
./test_real_commands.sh
```

### 2. Test Individual Parsers

#### Basic Usage
```bash
# Pipe command output directly
ls -la | ./term-to-json ls

# Or provide sample data
echo "sample_output" | ./term-to-json <parser_name>
```

#### Quick Test Commands

**System Information:**
```bash
# Test uname
uname -a | ./term-to-json uname

# Test uptime
uptime | ./term-to-json uptime

# Test who
who | ./term-to-json who

# Test current user ID
id | ./term-to-json id

# Test environment variables (first 5)
env | head -5 | ./term-to-json env
```

**Process & System Monitoring:**
```bash
# Test process list
ps aux | head -10 | ./term-to-json ps

# Test memory usage
free -h | ./term-to-json free

# Test virtual memory stats
vmstat | ./term-to-json vmstat
```

**Network:**
```bash
# Test ping (3 packets)
ping -c 3 google.com | ./term-to-json ping

# Test network connections
netstat -tuln | ./term-to-json netstat

# Test ARP table
arp -a | ./term-to-json arp
```

**Filesystem:**
```bash
# Test file listing
ls -la | ./term-to-json ls

# Test disk usage
df -h | ./term-to-json df

# Test directory sizes
du -h /tmp/* | head -5 | ./term-to-json du

# Test mounted filesystems
mount | ./term-to-json mount

# Test block devices (Linux)
lsblk | ./term-to-json lsblk

# Test file statistics
stat /etc/passwd | ./term-to-json stat
```

**Utilities:**
```bash
# Test date
date | ./term-to-json date

# Test word count
echo "hello world" | wc | ./term-to-json wc
```

**Configuration Files:**
```bash
# Test hosts file
cat /etc/hosts | ./term-to-json hosts

# Test passwd file (first 3 lines)
head -3 /etc/passwd | ./term-to-json passwd
```

### 3. Advanced Testing with jq

Filter and format JSON output:

```bash
# Show only directories from ls
ls -la | ./term-to-json ls | jq '.[] | select(.is_directory) | .name'

# Show memory usage percentage
free | ./term-to-json free | jq '.memory[0] | (.used / .total * 100)'

# Show high CPU processes
ps aux | ./term-to-json ps | jq '.[] | select(.cpu_percent > 1) | {user, pid, cpu_percent, command}'

# Show filesystem usage over 50%
df | ./term-to-json df | jq '.[] | select(.use_percent > 50) | {filesystem, mount_point, use_percent}'

# Show ping response times
ping -c 5 google.com | ./term-to-json ping | jq '.packets[] | .time_ms'
```

### 4. Performance Testing

Test with larger datasets:

```bash
# Large file listing
find /usr -type f | head -1000 | ./term-to-json find

# Many processes
ps aux | ./term-to-json ps | jq 'length'

# Large environment
env | ./term-to-json env | jq 'length'
```

### 5. Error Testing

Test error handling:

```bash
# Empty input
echo "" | ./term-to-json ls

# Invalid parser
echo "test" | ./term-to-json invalid_parser

# Malformed input
echo "not-a-valid-ls-output" | ./term-to-json ls
```

### 6. Integration Testing

Combine with other tools:

```bash
# Monitor and parse in real-time
watch -n 1 'ps aux | head -5 | ./term-to-json ps | jq ".[].cpu_percent"'

# Store parsed data
df | ./term-to-json df > disk_usage.json

# Process multiple files
for file in /etc/passwd /etc/group; do
    echo "=== $file ==="
    cat "$file" | head -3 | ./term-to-json passwd
done
```

## Available Parsers

Run `./term-to-json` without arguments to see all available parsers:

- **System:** uname, uptime, who, w, id, env
- **Process:** ps, free, vmstat  
- **Network:** ping, netstat, arp, dig
- **Files:** ls, df, du, mount, lsblk, find, stat
- **Services:** systemctl
- **Utilities:** date, wc
- **Config:** hosts, passwd

## Writing Tests

### Unit Test Pattern

```go
func TestMyParser(t *testing.T) {
    parser := &MyParser{}
    
    testInput := "sample command output"
    
    result, err := parser.Parse(testInput)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    entry, ok := result.(MyEntry)
    if !ok {
        t.Fatalf("Expected MyEntry, got %T", result)
    }
    
    // Test specific fields
    if entry.Field != "expected_value" {
        t.Errorf("Expected 'expected_value', got '%s'", entry.Field)
    }
}
```

### Integration Test Pattern

```bash
# Test real command output
actual_output=$(command_here)
parsed_output=$(echo "$actual_output" | ./term-to-json parser_name)
echo "$parsed_output" | jq . > /dev/null # Validate JSON
```

## Troubleshooting

### Common Issues

1. **Parser not found:** Check available parsers with `./term-to-json`
2. **Invalid JSON:** Check input format matches expected command output
3. **Missing fields:** Some parsers handle different output formats - check the parser source
4. **Performance:** For very large outputs, consider streaming or filtering first

### Debug Mode

Add debug output to parsers:
```go
fmt.Fprintf(os.Stderr, "Debug: parsing line: %s\n", line)
```
