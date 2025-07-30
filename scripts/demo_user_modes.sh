#!/bin/bash

# God Mode and Stealth Mode User Mode Demo
# This script demonstrates the new user mode implementation

echo "=== God Mode and Stealth Mode User Mode Demo ==="
echo ""
echo "God Mode and Stealth Mode are now implemented as proper IRC user modes:"
echo ""

echo "🔱 GOD MODE (+G):"
echo "  /mode mynick +G     # Enable God Mode - Ultimate channel override powers"
echo "  /mode mynick -G     # Disable God Mode"
echo ""
echo "  Capabilities:"
echo "  • Join banned/invite-only/password/full channels"
echo "  • Cannot be kicked by anyone"
echo "  • Override all channel restrictions"
echo ""

echo "👻 STEALTH MODE (+S):"
echo "  /mode mynick +S     # Enable Stealth Mode - Invisible to regular users"
echo "  /mode mynick -S     # Disable Stealth Mode"
echo ""
echo "  Effects:"
echo "  • Hidden from /who and /names for regular users"
echo "  • Still visible to other operators"
echo "  • Invisible channel presence"
echo ""

echo "⚡ COMBINED USAGE:"
echo "  /mode mynick +GS    # Enable both modes simultaneously"
echo "  /mode mynick -GS    # Disable both modes"
echo ""

echo "📋 REQUIREMENTS:"
echo "  • Must be an IRC operator (/oper)"
echo "  • Operator class must have 'god_mode' and/or 'stealth_mode' permissions"
echo "  • Can only set modes on yourself"
echo ""

echo "🛡️ PERMISSIONS:"
echo "  Add to your operator class in configs/opers.conf:"
echo '  "permissions": ["*", "god_mode", "stealth_mode"]'
echo ""

echo "🔍 CHECKING CURRENT MODES:"
echo "  /mode mynick        # Show your current user modes"
echo ""

echo "Example operator class with both permissions:"
echo '{'
echo '  "name": "admin",'
echo '  "rank": 4,'
echo '  "permissions": ["*", "god_mode", "stealth_mode"]'
echo '}'
echo ""

echo "This follows proper IRC conventions - user modes instead of custom commands!"
echo "Use /mode +G and /mode +S just like any other IRC user mode."
