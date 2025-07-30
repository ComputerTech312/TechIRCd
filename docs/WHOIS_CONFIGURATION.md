# WHOIS Configuration Guide

TechIRCd provides extremely flexible WHOIS configuration that allows administrators to control exactly what information is visible to different types of users.

## Overview

The WHOIS system in TechIRCd uses a three-tier permission system:
- **to_everyone**: Information visible to all users
- **to_opers**: Information visible only to IRC operators
- **to_self**: Information visible when users query themselves

## Configuration Options

### Basic Information Controls

#### User Modes (`show_user_modes`)
Controls who can see what user modes (+i, +w, +s, etc.) a user has set.

```json
"show_user_modes": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": true
}
```

#### SSL Status (`show_ssl_status`)
Shows whether a user is connected via SSL/TLS.

```json
"show_ssl_status": {
  "to_everyone": true,
  "to_opers": true,
  "to_self": true
}
```

#### Idle Time (`show_idle_time`)
Shows how long a user has been idle and their signon time.

```json
"show_idle_time": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": true
}
```

#### Signon Time (`show_signon_time`)
Shows when a user connected to the server (alternative to idle time).

```json
"show_signon_time": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": true
}
```

#### Real Host (`show_real_host`)
Shows the user's real IP address/hostname (bypasses host masking).

```json
"show_real_host": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": true
}
```

### Channel Information (`show_channels`)

Controls channel visibility with additional granular options:

```json
"show_channels": {
  "to_everyone": true,
  "to_opers": true,
  "to_self": true,
  "hide_secret_channels": true,
  "hide_private_channels": false,
  "show_membership_levels": true
}
```

- `hide_secret_channels`: Hide channels with mode +s from non-members
- `hide_private_channels`: Hide channels with mode +p from non-members  
- `show_membership_levels`: Show @/+/% prefixes for ops/voice/halfop

### Operator Information (`show_oper_class`)
Shows IRC operator class/type information.

```json
"show_oper_class": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": true
}
```

### Client Information (`show_client_info`)
Shows information about the IRC client software being used.

```json
"show_client_info": {
  "to_everyone": false,
  "to_opers": true,
  "to_self": false
}
```

### Account Name (`show_account_name`)
Shows services account name (for networks with services integration).

```json
"show_account_name": {
  "to_everyone": true,
  "to_opers": true,
  "to_self": true
}
```

## Advanced Features (Future)

These features are planned for future implementation:

- `show_activity_stats`: User activity analytics
- `show_github_integration`: GitHub profile integration
- `show_geolocation`: Approximate location information
- `show_performance_stats`: Connection performance metrics
- `show_device_info`: Device and OS information
- `show_social_graph`: Mutual channels and connections
- `show_security_score`: Account security rating

## Custom Fields

TechIRCd supports custom WHOIS fields for maximum flexibility:

```json
"custom_fields": [
  {
    "name": "website",
    "to_everyone": true,
    "to_opers": true,
    "to_self": true,
    "format": "Website: %s",
    "description": "User's personal website"
  }
]
```

## Example Configurations

### Maximum Privacy
Only show basic information to everyone, detailed info to opers:

```json
"whois_features": {
  "show_user_modes": {"to_everyone": false, "to_opers": true, "to_self": true},
  "show_ssl_status": {"to_everyone": false, "to_opers": true, "to_self": true},
  "show_idle_time": {"to_everyone": false, "to_opers": true, "to_self": true},
  "show_real_host": {"to_everyone": false, "to_opers": true, "to_self": true},
  "show_channels": {
    "to_everyone": true, "to_opers": true, "to_self": true,
    "hide_secret_channels": true, "hide_private_channels": true,
    "show_membership_levels": false
  }
}
```

### Maximum Transparency
Show most information to everyone:

```json
"whois_features": {
  "show_user_modes": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_ssl_status": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_idle_time": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_real_host": {"to_everyone": false, "to_opers": true, "to_self": true},
  "show_channels": {
    "to_everyone": true, "to_opers": true, "to_self": true,
    "hide_secret_channels": false, "hide_private_channels": false,
    "show_membership_levels": true
  }
}
```

### Development/Testing
Show everything to everyone for debugging:

```json
"whois_features": {
  "show_user_modes": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_ssl_status": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_idle_time": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_real_host": {"to_everyone": true, "to_opers": true, "to_self": true},
  "show_client_info": {"to_everyone": true, "to_opers": true, "to_self": true}
}
```

## Notes

- The WHOIS system respects IRC operator privileges and host masking settings
- Secret and private channel hiding works in conjunction with channel membership
- All settings are hot-reloadable (restart server after config changes)
- The system is designed to be extremely flexible while maintaining IRC protocol compliance

This configuration system makes TechIRCd one of the most configurable IRC servers available!
