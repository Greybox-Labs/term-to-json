#!/bin/bash

# Quick test script for term-to-json parsers
# This script demonstrates each parser with sample data

echo "🚀 Testing term-to-json parsers..."
echo

# Build the tool
echo "📦 Building term-to-json..."
go build -o term-to-json
echo

# Test uname parser
echo "🖥️  Testing uname parser:"
echo "Linux hostname 5.4.0-74-generic #83-Ubuntu SMP Sat May 8 02:35:39 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux" | ./term-to-json uname | jq .
echo

# Test uptime parser
echo "⏰ Testing uptime parser:"
echo " 14:30:42 up 12 days,  3:45,  2 users,  load average: 0.15, 0.12, 0.10" | ./term-to-json uptime | jq .
echo

# Test who parser
echo "👥 Testing who parser:"
echo "user     pts/0        2023-01-15 14:30 (192.168.1.10)
root     tty1         2023-01-15 10:00" | ./term-to-json who | jq .
echo

# Test ps parser
echo "🔍 Testing ps parser:"
echo "USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1  225316  9876 ?        Ss   Jan01   0:01 /sbin/init
user      1234  2.5  1.2  123456 12345 pts/0    R+   10:30   0:05 python script.py" | ./term-to-json ps | jq .
echo

# Test df parser
echo "💾 Testing df parser:"
echo "Filesystem     1K-blocks    Used Available Use% Mounted on
/dev/sda1       20511312  123456  19365472   1% /
tmpfs            4096000       0   4096000   0% /tmp" | ./term-to-json df | jq .
echo

# Test ls parser
echo "📁 Testing ls parser:"
echo "drwxr-xr-x  3 user group  4096 Jan 15 10:30 docs
-rw-r--r--  1 user group  1234 Jan 14 09:15 file.txt
lrwxrwxrwx  1 user group     8 Jan 13 14:20 link -> file.txt" | ./term-to-json ls | jq .
echo

# Test free parser
echo "🧠 Testing free parser:"
echo "              total        used        free      shared  buff/cache   available
Mem:        8147720     2084264     4044900      123456     2018556     5719320
Swap:       2097148           0     2097148" | ./term-to-json free | jq .
echo

# Test ping parser
echo "🏓 Testing ping parser:"
echo "PING google.com (142.250.191.14) 56(84) bytes of data.
64 bytes from google.com (142.250.191.14): icmp_seq=1 ttl=55 time=12.3 ms
64 bytes from google.com (142.250.191.14): icmp_seq=2 ttl=55 time=15.1 ms
--- google.com ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 12.3/13.7/15.1/1.4 ms" | ./term-to-json ping | jq .
echo

# Test netstat parser
echo "🌐 Testing netstat parser:"
echo "Proto Recv-Q Send-Q Local Address           Foreign Address         State       
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      
tcp        0      0 127.0.0.1:3306          0.0.0.0:*               LISTEN" | ./term-to-json netstat | jq .
echo

# Test mount parser
echo "🗂️  Testing mount parser:"
echo "/dev/sda1 on / type ext4 (rw,relatime,errors=remount-ro)
tmpfs on /tmp type tmpfs (rw,nosuid,nodev)" | ./term-to-json mount | jq .
echo

# Test date parser
echo "📅 Testing date parser:"
echo "Mon Jan 15 14:30:45 UTC 2024" | ./term-to-json date | jq .
echo

# Test env parser
echo "🌍 Testing env parser:"
echo "PATH=/usr/local/bin:/usr/bin:/bin
HOME=/home/user
USER=user" | ./term-to-json env | jq .
echo

# Test hosts parser
echo "🏠 Testing hosts parser:"
echo "127.0.0.1	localhost
192.168.1.10	server1 server1.local
# This is a comment
::1	ip6-localhost ip6-loopback" | ./term-to-json hosts | jq .
echo

# Test passwd parser
echo "👤 Testing passwd parser:"
echo "root:x:0:0:root:/root:/bin/bash
user:x:1000:1000:User:/home/user:/bin/bash
daemon:x:2:2:daemon:/sbin:/usr/sbin/nologin" | ./term-to-json passwd | jq .
echo

echo "✅ All tests completed!"
echo
echo "To run individual tests:"
echo "  echo 'sample_output' | ./term-to-json <parser_name>"
echo
echo "To run Go unit tests:"
echo "  go test ./parsers -v"
echo
echo "Available parsers:"
echo "  System: uname, uptime, who, w, id, env"
echo "  Process: ps, free, vmstat"
echo "  Network: ping, netstat, arp, dig"
echo "  Files: ls, df, du, mount, lsblk, find, stat"
echo "  Services: systemctl"
echo "  Utilities: date, wc"
echo "  Config: hosts, passwd"
