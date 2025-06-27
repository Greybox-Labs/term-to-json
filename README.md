# term-to-json

A Go library that parses command line output from common Linux utilities into JSON format, inspired by the [jc](https://github.com/kellyjonbrazil/jc) tool.

## Features

- **25+ parsers** for common Linux utilities
- Clean, extensible parser interface
- Comprehensive test coverage
- Support for most common command output formats

### Supported Parsers

**System Information:**
- `uname` - System information
- `uptime` - System uptime and load
- `who` - Logged in users
- `w` - User activity
- `id` - User and group IDs
- `env` - Environment variables

**Process & System Monitoring:**
- `ps` - Process listing
- `free` - Memory usage
- `vmstat` - Virtual memory statistics

**Network:**
- `ping` - Network connectivity test
- `netstat` - Network connections
- `arp` - ARP table
- `dig` - DNS lookups

**Filesystem:**
- `ls` - File listings
- `df` - Disk usage
- `du` - Directory usage
- `mount` - Mounted filesystems
- `lsblk` - Block devices
- `find` - File search results
- `stat` - File statistics

**System Services:**
- `systemctl` - Systemd service status

**Utilities:**
- `date` - Date/time information
- `wc` - Word, line, character counts

**Configuration Files:**
- `hosts` - /etc/hosts entries
- `passwd` - /etc/passwd entries

## Installation

```bash
go get github.com/yourusername/term-to-json
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "log"
    "term-to-json/parsers"
)

func main() {
    lsOutput := `total 48
drwxr-xr-x  3 user group  4096 Jan 15 10:30 docs
-rw-r--r--  1 user group  1234 Jan 14 09:15 file.txt
lrwxrwxrwx  1 user group     8 Jan 13 14:20 link -> file.txt`

    result, err := parsers.Parse("ls", lsOutput)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", result)
}
```

### As a Command Line Tool

```bash
# Build the tool
go build -o term-to-json

# Parse ls output
ls -l | ./term-to-json ls

# Parse ps output  
ps aux | ./term-to-json ps

# Parse df output
df -h | ./term-to-json df

# Parse system info
uname -a | ./term-to-json uname

# Parse network connectivity
ping -c 3 google.com | ./term-to-json ping

# Parse memory usage
free -h | ./term-to-json free
```

## Examples

### ls Parser

Input:
```
drwxr-xr-x  3 user group  4096 Jan 15 10:30 docs
-rw-r--r--  1 user group  1234 Jan 14 09:15 file.txt
lrwxrwxrwx  1 user group     8 Jan 13 14:20 link -> file.txt
```

Output:
```json
[
  {
    "permissions": "drwxr-xr-x",
    "links": 3,
    "owner": "user",
    "group": "group",
    "size": 4096,
    "modified": "2024-01-15T10:30:00Z",
    "name": "docs",
    "is_directory": true,
    "is_symlink": false
  },
  {
    "permissions": "-rw-r--r--",
    "links": 1,
    "owner": "user", 
    "group": "group",
    "size": 1234,
    "modified": "2024-01-14T09:15:00Z",
    "name": "file.txt",
    "is_directory": false,
    "is_symlink": false
  },
  {
    "permissions": "lrwxrwxrwx",
    "links": 1,
    "owner": "user",
    "group": "group", 
    "size": 8,
    "modified": "2024-01-13T14:20:00Z",
    "name": "link",
    "is_directory": false,
    "is_symlink": true,
    "link_target": "file.txt"
  }
]
```

### ps Parser

Input:
```
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1  225316  9876 ?        Ss   Jan01   0:01 /sbin/init
user      1234  2.5  1.2  123456 12345 pts/0    R+   10:30   0:05 python script.py
```

Output:
```json
[
  {
    "pid": 1,
    "user": "root",
    "cpu_percent": 0.0,
    "memory_percent": 0.1,
    "vsz": 225316,
    "rss": 9876,
    "tty": "?",
    "stat": "Ss",
    "start": "Jan01",
    "time": "0:01",
    "command": "/sbin/init"
  },
  {
    "pid": 1234,
    "user": "user",
    "cpu_percent": 2.5,
    "memory_percent": 1.2,
    "vsz": 123456,
    "rss": 12345,
    "tty": "pts/0",
    "stat": "R+", 
    "start": "10:30",
    "time": "0:05",
    "command": "python script.py"
  }
]
```

### df Parser

Input:
```
Filesystem     1K-blocks    Used Available Use% Mounted on
/dev/sda1       20511312  123456  19365472   1% /
tmpfs            4096000       0   4096000   0% /tmp
```

Output:
```json
[
  {
    "filesystem": "/dev/sda1",
    "size": 20511312,
    "used": 123456,
    "available": 19365472,
    "use_percent": 1,
    "mount_point": "/",
    "used_bytes": 126419456, 
    "avail_bytes": 19838439424,
    "size_bytes": 20995583488
  },
  {
    "filesystem": "tmpfs",
    "size": 4096000,
    "used": 0,
    "available": 4096000,
    "use_percent": 0,
    "mount_point": "/tmp",
    "used_bytes": 0,
    "avail_bytes": 4194304000,
    "size_bytes": 4194304000
  }
]
```

## Adding New Parsers

To add a new parser:

1. Create a new file in the `parsers/` directory (e.g., `parsers/mycommand.go`)
2. Implement the `Parser` interface:

```go
type MyCommandParser struct{}

func (p *MyCommandParser) Name() string {
    return "mycommand"
}

func (p *MyCommandParser) Parse(input string) (interface{}, error) {
    // Parse logic here
    return result, nil
}
```

3. Add the parser to the switch statement in `parsers/parser.go`
4. Add comprehensive tests in `parsers/mycommand_test.go`

## Testing

```bash
go test ./parsers
```

## License

MIT License
