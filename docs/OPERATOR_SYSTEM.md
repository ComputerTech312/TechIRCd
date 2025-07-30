# Operator Configuration Guide

TechIRCd features a sophisticated hierarchical operator system with separate configuration files and granular permission control.

## Overview

The operator system consists of:
- **Operator Classes**: Define ranks, permissions, and inheritance
- **Individual Operators**: Specific users with assigned classes
- **Permission System**: Granular control over what operators can do
- **Rank Hierarchy**: Higher ranks can operate on lower ranks

## Configuration Files

### Main Config (`config.json`)
```json
"oper_config": {
  "config_file": "configs/opers.conf",
  "enable": true
}
```

### Operator Config (`configs/opers.conf`)
Contains the complete operator hierarchy and individual operator definitions.

## Operator Classes

### Customizable Rank Names

TechIRCd allows you to completely customize the names of operator ranks. Instead of using the standard Helper/Moderator/Operator names, you can create themed naming schemes:

```json
"rank_names": {
  "rank_1": "Cadet",       // Instead of "Helper" 
  "rank_2": "Sergeant",    // Instead of "Moderator"
  "rank_3": "Lieutenant",  // Instead of "Operator"
  "rank_4": "Captain",     // Instead of "Administrator"
  "rank_5": "General",     // Instead of "Owner"
  "custom_ranks": {
    "Field Marshal": 6,    // Custom ranks beyond 5
    "Supreme Commander": 10
  }
}
```

### Popular Naming Themes

#### Gaming/Military Theme
- Rank 1: `Cadet`, `Recruit`, `Private`
- Rank 2: `Sergeant`, `Corporal`, `Specialist` 
- Rank 3: `Lieutenant`, `Captain`, `Major`
- Rank 4: `Colonel`, `General`, `Commander`
- Rank 5: `Admiral`, `Marshal`, `Supreme Commander`

#### Corporate/Business Theme
- Rank 1: `Intern`, `Assistant`, `Associate`
- Rank 2: `Specialist`, `Senior Associate`, `Team Lead`
- Rank 3: `Manager`, `Senior Manager`, `Department Head`
- Rank 4: `Director`, `VP`, `Executive`
- Rank 5: `CEO`, `President`, `Chairman`

#### Fantasy/Medieval Theme
- Rank 1: `Apprentice`, `Squire`, `Page`
- Rank 2: `Guardian`, `Warrior`, `Mage`
- Rank 3: `Knight`, `Paladin`, `Wizard`
- Rank 4: `Lord`, `Baron`, `Archmage`
- Rank 5: `King`, `Emperor`, `Divine Ruler`

#### Sci-Fi Theme
- Rank 1: `Ensign`, `Technician`, `Cadet`
- Rank 2: `Engineer`, `Pilot`, `Specialist`
- Rank 3: `Commander`, `Captain`, `Leader`
- Rank 4: `Admiral`, `Fleet Commander`, `Director`
- Rank 5: `Supreme Admiral`, `Galactic Emperor`, `AI Overlord`

### Built-in Classes (Default Names)

#### Helper (Rank 1)
- **Symbol**: `%`
- **Color**: Green
- **Permissions**: Basic moderation
  - `kick` - Kick users from channels
  - `topic` - Change channel topics
  - `mode_channel` - Change channel modes

#### Moderator (Rank 2)  
- **Symbol**: `@`
- **Color**: Blue
- **Inherits**: Helper permissions
- **Additional Permissions**:
  - `ban` / `unban` - Manage channel bans
  - `mute` - Mute users
  - `mode_user` - Change user modes
  - `who_override` - See hidden users in WHO

#### Operator (Rank 3)
- **Symbol**: `*`
- **Color**: Red
- **Inherits**: Moderator permissions
- **Additional Permissions**:
  - `kill` - Kill user connections
  - `gline` - Global bans
  - `rehash` - Reload configuration
  - `connect` / `squit` - Server linking
  - `wallops` / `operwall` - Send operator messages

#### Administrator (Rank 4)
- **Symbol**: `&`
- **Color**: Purple
- **Permissions**: `*` (All permissions)

#### Owner (Rank 5)
- **Symbol**: `~`
- **Color**: Gold
- **Special Permissions**:
  - `*` - All permissions
  - `override_rank` - Can operate on same/higher ranks
  - `shutdown` / `restart` - Server control

### Custom Classes

Create your own operator classes:

```json
{
  "name": "security",
  "rank": 3,
  "description": "Security Officer - Network protection",
  "permissions": [
    "kill",
    "gline",
    "scan_network",
    "access_logs"
  ],
  "inherits": "moderator",
  "color": "orange",
  "symbol": "!"
}
```

## Individual Operators

### Basic Operator Definition

```json
{
  "name": "alice",
  "password": "secure_password_here",
  "host": "*@trusted.example.com",
  "class": "moderator",
  "flags": ["extra_channels"],
  "max_clients": 500,
  "contact": "alice@example.com"
}
```

### Advanced Features

```json
{
  "name": "bob",
  "password": "another_secure_password",
  "host": "admin@192.168.1.*",
  "class": "admin",
  "flags": ["debug_access", "special_channels"],
  "max_clients": 1000,
  "expires": "2025-12-31",
  "contact": "bob@company.com",
  "last_seen": "2025-07-30T10:30:00Z"
}
```

## Permission System

### Standard Permissions

#### Channel Management
- `kick` - Kick users from channels
- `ban` / `unban` - Manage channel bans
- `topic` - Change channel topics
- `mode_channel` - Change channel modes
- `mode_user` - Change user modes

#### User Management
- `kill` - Disconnect users
- `gline` - Global network bans
- `mute` - Silence users
- `who_override` - See hidden information

#### Server Management
- `rehash` - Reload configuration
- `connect` / `squit` - Server linking
- `wallops` / `operwall` - Operator communications
- `shutdown` / `restart` - Server control

#### Special Permissions
- `*` - All permissions (wildcard)
- `override_rank` - Ignore rank restrictions
- `debug_access` - Debug commands
- `log_access` - View server logs

### Custom Permissions

Add your own permissions for custom commands:

```json
"permissions": [
  "custom_report",
  "special_database_access",
  "network_monitoring"
]
```

## Security Features

### Settings Configuration

```json
"settings": {
  "require_ssl": true,
  "max_failed_attempts": 3,
  "lockout_duration_minutes": 30,
  "log_oper_actions": true,
  "notify_on_oper_up": true,
  "auto_expire_inactive_days": 365,
  "require_two_factor": false
}
```

### Security Options

- **SSL Requirement**: Force SSL for operator authentication
- **Failed Attempt Tracking**: Lock accounts after failed attempts
- **Action Logging**: Log all operator actions
- **Auto-Expiration**: Remove inactive operators
- **Two-Factor Authentication**: (Future feature)

## Rank System

### How Ranks Work

- **Higher numbers = Higher authority**
- **Same rank cannot operate on each other**
- **Override permission bypasses rank restrictions**

### Example Hierarchy

```
Owner (5)     -> Can do anything to anyone
Admin (4)     -> Can operate on Operator (3) and below
Operator (3)  -> Can operate on Moderator (2) and below
Moderator (2) -> Can operate on Helper (1) and below
Helper (1)    -> Can only operate on regular users
```

## WHOIS Integration

The operator system integrates with WHOIS to show:

- Operator status with custom symbols
- Class names and descriptions
- **Custom rank names** based on your configuration
- Permission levels

Example WHOIS output with custom rank names:
```
[313] alice is an IRC operator (moderator - Channel and user management) [Sergeant]
[313] bob is an IRC operator (captain - Senior officer with full authority) [Captain]
```

With fantasy theme:
```
[313] merlin is an IRC operator (lord - Ruler of vast territories) [Lord]
```

## Configuration Examples

### Quick Setup - Gaming Theme

```json
{
  "rank_names": {
    "rank_1": "Cadet",
    "rank_2": "Sergeant", 
    "rank_3": "Lieutenant",
    "rank_4": "Captain",
    "rank_5": "General"
  },
  "classes": [
    {
      "name": "cadet",
      "rank": 1,
      "description": "New recruit with basic training",
      "permissions": ["kick", "topic"],
      "symbol": "+",
      "color": "green"
    }
  ]
}
```

See `/configs/examples/` for complete themed configurations:
- `opers-gaming-theme.conf` - Military/Gaming ranks
- `opers-corporate-theme.conf` - Business hierarchy
- `opers-fantasy-theme.conf` - Medieval/Fantasy theme

## Migration from Legacy

TechIRCd automatically falls back to the legacy operator system if:
- `oper_config.enable` is `false`
- The opers.conf file cannot be loaded
- No matching operator is found in the new system

Legacy operators get basic permissions and rank 1.

## Commands for Operators

### Checking Permissions
```
/WHOIS operator_name  # See operator class and rank
/OPERLIST             # List all operators (future)
/OPERHELP             # Show operator commands (future)
```

### Management Commands
```
/OPER name password   # Become an operator
/OPERWALL message     # Send message to all operators
/REHASH               # Reload configuration
```

## Best Practices

### Security
1. Use strong passwords for all operators
2. Restrict host masks to trusted IPs
3. Enable SSL requirement for sensitive ranks
4. Regularly audit operator lists
5. Set expiration dates for temporary operators

### Organization
1. Create classes that match your network structure
2. Use inheritance to avoid permission duplication
3. Document custom permissions clearly
4. Use descriptive class names and descriptions
5. Assign appropriate symbols and colors

### Maintenance
1. Review operator activity regularly
2. Remove inactive operators
3. Update permissions as needed
4. Monitor operator actions through logs
5. Keep opers.conf in version control

This operator system makes TechIRCd extremely flexible for networks of any size, from small communities to large networks with complex hierarchies!
