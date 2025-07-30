#!/bin/bash

# TechIRCd Rank Customization Demo
# This script demonstrates how easy it is to change operator rank names

echo "=== TechIRCd Rank Customization Demo ==="
echo ""

echo "1. Default Configuration (configs/opers.conf):"
echo "   Rank 1: Helper"
echo "   Rank 2: Moderator"
echo "   Rank 3: Operator"
echo "   Rank 4: Administrator"
echo "   Rank 5: Owner"
echo ""

echo "2. Gaming Theme (configs/examples/opers-gaming-theme.conf):"
echo "   Rank 1: Cadet"
echo "   Rank 2: Sergeant"
echo "   Rank 3: Lieutenant"
echo "   Rank 4: Captain"
echo "   Rank 5: General"
echo "   Rank 6: Field Marshal (custom)"
echo "   Rank 10: Supreme Commander (custom)"
echo ""

echo "3. Corporate Theme (configs/examples/opers-corporate-theme.conf):"
echo "   Rank 1: Intern"
echo "   Rank 2: Associate"
echo "   Rank 3: Manager"
echo "   Rank 4: Director"
echo "   Rank 5: CEO"
echo "   Rank 6: Board Member (custom)"
echo "   Rank 8: Chairman (custom)"
echo "   Rank 10: Founder (custom)"
echo ""

echo "4. Fantasy Theme (configs/examples/opers-fantasy-theme.conf):"
echo "   Rank 1: Apprentice"
echo "   Rank 2: Guardian"
echo "   Rank 3: Knight"
echo "   Rank 4: Lord"
echo "   Rank 5: King"
echo "   Rank 6: Archmage (custom)"
echo "   Rank 8: High King (custom)"
echo "   Rank 10: Divine Emperor (custom)"
echo ""

echo "To use a theme:"
echo "1. Copy one of the example configs to configs/opers.conf"
echo "2. Update your main config.json to point to it:"
echo '   "oper_config": {"config_file": "configs/opers.conf", "enable": true}'
echo "3. Restart TechIRCd"
echo ""

echo "Custom ranks can go up to any number (6, 7, 8, 10, 15, 100, etc.)"
echo "The rank name will appear in WHOIS responses like:"
echo "[313] player is an IRC operator (sergeant - Squad leader) [Sergeant]"
echo ""

echo "Make your IRC network unique with custom operator hierarchies!"
