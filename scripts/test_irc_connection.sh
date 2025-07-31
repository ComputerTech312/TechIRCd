#!/bin/bash

# Simple IRC client test to debug connection issues
# This script connects to the IRC server and sends basic registration commands

echo "=== TechIRCd Connection Test ==="
echo "Testing IRC registration sequence..."
echo ""

# Create a temporary file for the test
TEST_FILE="/tmp/irc_test.txt"

# Connect and send basic registration commands
{
    echo "NICK testuser"
    echo "USER testuser 0 * :Test User"
    sleep 2
    echo "QUIT :Test complete"
} | nc -q 3 127.0.0.1 6667 | tee "$TEST_FILE"

echo ""
echo "=== Raw IRC Server Response ==="
cat "$TEST_FILE"
echo ""

echo "=== Analysis ==="
if grep -q "001" "$TEST_FILE"; then
    echo "✅ RPL_WELCOME (001) found - Server sent welcome message"
else
    echo "❌ RPL_WELCOME (001) NOT found - Registration may have failed"
fi

if grep -q "002" "$TEST_FILE"; then
    echo "✅ RPL_YOURHOST (002) found"
else
    echo "❌ RPL_YOURHOST (002) NOT found"
fi

if grep -q "003" "$TEST_FILE"; then
    echo "✅ RPL_CREATED (003) found"
else
    echo "❌ RPL_CREATED (003) NOT found"
fi

if grep -q "004" "$TEST_FILE"; then
    echo "✅ RPL_MYINFO (004) found"
else
    echo "❌ RPL_MYINFO (004) NOT found"
fi

echo ""
echo "If the pyechat client connects but shows nothing, it might be:"
echo "1. Not sending NICK/USER commands automatically"
echo "2. Expecting different IRC numeric responses"
echo "3. Not parsing the IRC protocol properly"
echo "4. Waiting for specific server capabilities"

# Cleanup
rm -f "$TEST_FILE"
