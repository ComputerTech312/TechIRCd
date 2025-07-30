# God Mode and Stealth Mode

TechIRCd supports two advanced operator features via user modes:

## God Mode (+G)

God Mode gives operators ultimate channel override capabilities.

### Usage
```
/mode nickname +G    # Enable God Mode
/mode nickname -G    # Disable God Mode
```

### Capabilities
- **Channel Access**: Can join any channel regardless of bans, limits, keys, or invite-only status
- **Kick Immunity**: Cannot be kicked from channels by anyone, including other operators
- **Ban Immunity**: Can join and remain in channels even when banned
- **Invite Override**: Can join invite-only channels without being invited
- **Limit Override**: Can join channels that have reached their user limit
- **Key Override**: Can join channels with passwords/keys without providing them

### Requirements
- Must be an IRC operator (`/oper`)
- Must have `god_mode` permission in operator class configuration
- Only the user themselves can set/unset their God Mode

### Security Notes
- God Mode actions are logged to operator snomasks (`+o`)
- Use responsibly - this bypasses all normal channel protections
- Intended for emergency situations and network administration

## Stealth Mode (+S)

Stealth Mode makes operators invisible to regular users while remaining visible to other operators.

### Usage
```
/mode nickname +S    # Enable Stealth Mode
/mode nickname -S    # Disable Stealth Mode
```

### Effects
- **WHO Command**: Stealth users don't appear in `/who` results for regular users
- **NAMES Command**: Stealth users don't appear in channel user lists for regular users
- **Channel Lists**: Regular users can't see stealth operators in channels
- **Operator Visibility**: Other operators can always see stealth users

### Requirements
- Must be an IRC operator (`/oper`)
- Must have `stealth_mode` permission in operator class configuration
- Only the user themselves can set/unset their Stealth Mode

### Use Cases
- Undercover moderation and monitoring
- Reduced operator visibility during investigations
- Network administration without user awareness

## Configuration

Add permissions to operator classes in `configs/opers.conf`:

```json
{
  "classes": [
    {
      "name": "admin",
      "rank": 4,
      "description": "Administrator with special powers",
      "permissions": [
        "*",
        "god_mode",
        "stealth_mode"
      ]
    }
  ]
}
```

## Examples

### Basic Usage
```
# As an operator with god_mode permission:
/mode mynick +G
# *** GOD MODE enabled - You have ultimate power!

# Join a banned/invite-only channel:
/join #private-channel
# Successfully joins despite restrictions

# Disable God Mode:
/mode mynick -G
# *** GOD MODE disabled

# Enable Stealth Mode:
/mode mynick +S
# *** STEALTH MODE enabled - You are now invisible to users

# Regular users won't see you in:
/who #channel
/names #channel

# Other operators will still see you
```

### Combined Usage
```
# Enable both modes simultaneously:
/mode mynick +GS
# *** GOD MODE enabled - You have ultimate power!
# *** STEALTH MODE enabled - You are now invisible to users

# Now you can:
# - Join any channel (God Mode)
# - Remain invisible to regular users (Stealth Mode)
# - Be visible to other operators
```

### Permission Checking
```
# Attempting without proper permissions:
/mode mynick +G
# :server 481 mynick :Permission Denied - You need god_mode permission

# Must be oper first:
/mode mynick +S
# :server 481 mynick :Permission Denied- You're not an IRC operator
```

## Technical Implementation

### God Mode
- Stored as user mode `+G`
- Checked via `HasGodMode()` method
- Bypasses channel restrictions in:
  - `JOIN` command (bans, limits, keys, invite-only)
  - `KICK` command (cannot be kicked)
  - Channel access validation

### Stealth Mode
- Stored as user mode `+S`
- Checked via `HasStealthMode()` method
- Filtered in:
  - `WHO` command responses
  - `NAMES` command responses
  - Channel member visibility

### Mode Persistence
- User modes are stored per-client session
- Lost on disconnect/reconnect
- Must be re-enabled after each connection

## Security Considerations

1. **Audit Trail**: All God Mode and Stealth Mode activations are logged
2. **Permission Based**: Requires explicit operator class permissions
3. **Self-Only**: Users can only set modes on themselves
4. **Operator Level**: Requires existing operator privileges
5. **Reversible**: Can be disabled at any time

## Troubleshooting

### Mode Not Setting
- Verify you are opered (`/oper`)
- Check operator class has required permissions
- Ensure using correct syntax (`/mode nickname +G`)

### Not Working as Expected
- God Mode only affects channel restrictions, not other commands
- Stealth Mode only hides from regular users, not operators
- Modes are case-sensitive (`+G` not `+g`)

### Permission Denied
- Contact network administrator to add permissions to your operator class
- Verify operator class configuration in `configs/opers.conf`
