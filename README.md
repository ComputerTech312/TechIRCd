# TechIRCd

A modern, high-performance IRC server written in Go with comprehensive RFC compliance and advanced features.

## Features

### 🚀 **Core IRC Protocol**
- Full RFC 2812 compliance (Internet Relay Chat: Client Protocol)
- User registration and authentication
- Channel management with comprehensive modes
- Private messaging and notices
- **Ultra-flexible WHOIS system** with granular privacy controls
- WHO and NAMES commands
- Ping/Pong keepalive mechanism

### 🔍 **Revolutionary WHOIS System**
- **Granular Privacy Controls**: Configure exactly what information is visible to everyone, operators, or only the user themselves
- **User Modes Visibility**: Control who can see user modes (+i, +w, +s, etc.)
- **SSL Status Display**: Show secure connection status
- **Idle Time & Signon Time**: Configurable time information
- **Real Host vs Masked Host**: Smart hostname display based on permissions
- **Channel Privacy**: Hide secret/private channels with fine-grained control
- **Operator Information**: Show operator class and privileges
- **Services Integration**: Account name display for services
- **Client Information**: Show IRC client details
- **Custom Fields**: Add your own WHOIS fields
- See [WHOIS Configuration Guide](docs/WHOIS_CONFIGURATION.md) for details

### 👑 **Advanced Channel Management**
- **Operator Hierarchy**: Owners (~), Operators (@), Half-ops (%), Voice (+)
- **Channel Modes**: 
  - `+m` (moderated) - Only voiced users can speak
  - `+n` (no external messages)
  - `+t` (topic protection)
  - `+i` (invite only)
  - `+s` (secret channel)
  - `+p` (private channel)
  - `+k` (channel key/password)
  - `+l` (user limit)
  - `+b` (ban list)
- **Extended Ban System**: Support for quiet mode (`~q:mask`) and other extended ban types

### 🔐 **IRC Operator Features**
- Comprehensive operator authentication system
- **Hierarchical Operator Classes** with rank-based permissions and inheritance
- **Completely Customizable Rank Names** (Gaming, Corporate, Fantasy themes, etc.)
- **Server Notice Masks (SNOmasks)**:
  - `+c` (connection notices)
  - `+k` (kill notices)
  - `+o` (oper notices)
  - `+x` (ban/quiet notices)
  - `+f` (flood notices)
  - `+n` (nick change notices)
  - `+s` (server notices)
  - `+d` (debug notices)
- **Unique Operator Commands**:
  - **`/GODMODE`** - ⚡ Toggle ultimate channel override powers
  - **`/STEALTH`** - 👤 Toggle invisibility to regular users
- **Advanced Operator Commands**: KILL, GLINE, REHASH, WALLOPS, OPERWALL
  - `+n` (nick change notices)
  - `+s` (server notices)
  - `+d` (debug notices)
- **Operator Commands**:
  - `KILL` - Disconnect users
  - `GLOBALNOTICE` - Send notices to all users
  - `OPERWALL` - Send messages to all operators
  - `WALLOPS` - Send wallops messages
  - `REHASH` - Reload configuration
  - `TRACE` - Network trace information

### 👤 **User Modes**
- `+i` (invisible) - Hide from WHO listings
- `+w` (wallops) - Receive wallops messages
- `+s` (server notices) - Receive server notices (opers only)
- `+o` (operator) - IRC operator status
- `+x` (host masking) - Hide real hostname
- `+B` (bot) - Mark as a bot
- `+z` (SSL) - Connected via SSL/TLS
- `+r` (registered) - Registered with services

### 🛡️ **Security & Stability**
- Advanced flood protection with operator exemption
- Connection timeout management
- Input validation and sanitization
- Panic recovery and error handling
- Resource monitoring and health checks
- Graceful shutdown capabilities

### ⚡ **Unique Features (Not Found in Other IRCds)**

#### 🌟 **God Mode (+G)**
- **Ultimate Channel Override**: Join any channel regardless of bans, limits, keys, or invite-only mode
- **Kick Immunity**: Cannot be kicked by any user
- **Mode Override**: Set any channel mode without operator privileges
- **Complete Bypass**: Ignore all channel restrictions and limitations

#### 👻 **Stealth Mode (+S)**
- **User Invisibility**: Completely hidden from regular users in WHO, NAMES, and WHOIS
- **Operator Visibility**: Other operators can still see stealth users
- **Covert Monitoring**: Watch channels without being detected
- **Security Operations**: Investigate issues invisibly

#### 🎨 **Customizable Rank Names**
- **Gaming Themes**: Cadet → Sergeant → Lieutenant → Captain → General
- **Corporate Themes**: Intern → Associate → Manager → Director → CEO
- **Fantasy Themes**: Apprentice → Guardian → Knight → Lord → King
- **Unlimited Creativity**: Create any rank system you can imagine

#### 🔍 **Revolutionary WHOIS**
- **Granular Privacy**: 15+ configurable information types
- **Three-Tier Permissions**: Everyone/Operators/Self visibility controls
- **Custom Fields**: Add your own WHOIS information
- **Complete Flexibility**: Control every aspect of user information display

### 📊 **Monitoring & Logging**
- Real-time health monitoring
- Memory usage tracking
- Goroutine count monitoring
- Performance metrics logging
- Configurable log levels and rotation
- Private messaging
- WHO/WHOIS commands
- Configurable server settings
- Extensible architecture
- Ping/Pong keep-alive mechanism
- Graceful shutdown handling

## Requirements

- Go 1.21 or higher

## Building

```bash
# Install dependencies
go mod tidy

# Build the server
go run tools/build.go -build

# Or use the traditional method
go build -o techircd
```

## Build Options

The `tools/build.go` script provides several build options:

```bash
# Build and run
go run tools/build.go -run

# Run tests
go run tools/build.go -test

# Format code
go run tools/build.go -fmt

# Cross-platform builds
go run tools/build.go -build-all

# Optimized release build
go run tools/build.go -release

# Clean build artifacts
go run tools/build.go -clean

# Show all options
go run tools/build.go -help
```

## Usage

```bash
./techircd
```

The server will start on localhost:6667 by default. You can configure the host and port in `config.go`.

## Configuration

Edit `config.go` to customize:
- Server host and port
- Server name and description
- Maximum connections
- Ping timeout
- Message of the Day (MOTD)

## Connecting

You can connect using any IRC client:
```
/server localhost 6667
```

Or use the included test client:
```bash
go run test_client.go
```

## Architecture

- `main.go` - Entry point and server initialization
- `server.go` - Main IRC server implementation with concurrent connection handling
- `client.go` - Client connection handling and state management
- `channel.go` - Channel management with thread-safe operations
- `commands.go` - IRC command implementations (NICK, USER, JOIN, PART, etc.)
- `config.go` - Server configuration structure
- `test_client.go` - Simple test client for development
