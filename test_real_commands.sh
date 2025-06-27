#!/bin/bash

# Test script using real system commands
# This demonstrates the library working with actual command output

echo "🔍 Testing with real system commands..."
echo

# Test with real commands that should be available on most systems

# Test uname
echo "🖥️  Real uname output:"
if command -v uname >/dev/null 2>&1; then
    uname -a | ./term-to-json uname | jq -r '.kernel_name + " " + .node_name + " " + .kernel_release'
    echo
fi

# Test who
echo "👥 Real who output:"
if command -v who >/dev/null 2>&1; then
    who | ./term-to-json who | jq -r '.[] | .user + " on " + .tty'
    echo
fi

# Test df
echo "💾 Real df output:"
if command -v df >/dev/null 2>&1; then
    df | ./term-to-json df | jq -r '.[] | .filesystem + ": " + (.use_percent | tostring) + "% used"'
    echo
fi

# Test ps (just a few processes)
echo "🔍 Real ps output (top 3):"
if command -v ps >/dev/null 2>&1; then
    ps aux | head -4 | ./term-to-json ps | jq -r '.[] | .user + " " + (.pid | tostring) + " " + .command' | head -3
    echo
fi

# Test env (just a few variables)
echo "🌍 Real env output (first 3):"
if command -v env >/dev/null 2>&1; then
    env | head -3 | ./term-to-json env | jq -r '.[] | .name + "=" + .value'
    echo
fi

# Test date
echo "📅 Real date output:"
if command -v date >/dev/null 2>&1; then
    date | ./term-to-json date | jq -r '.iso + " (" + .weekday + ")"'
    echo
fi

# Test mount (first 3 mounts)
echo "🗂️  Real mount output (first 3):"
if command -v mount >/dev/null 2>&1; then
    mount | head -3 | ./term-to-json mount | jq -r '.[] | .device + " -> " + .mount_point + " (" + .filesystem_type + ")"'
    echo
fi

echo "✅ Real command tests completed!"
echo
echo "💡 Tips:"
echo "  • Pipe any command output directly: 'command | ./term-to-json parser'"
echo "  • Use jq to filter JSON output: 'command | ./term-to-json parser | jq .field'"
echo "  • Combine with other tools: 'ls -la | ./term-to-json ls | jq \".[] | select(.is_directory)\"'"
