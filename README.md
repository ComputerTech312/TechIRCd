# TechIRCd

A modern, high-performance IRC server written in Go with comprehensive RFC compliance and advanced features.

## Features

### 🚀 **Core IRC Protocol**
- Full RFC 2812 compliance (Internet Relay Chat: Client Protocol)
- User registration and authentication
- Channel management with comprehensive modes
- Private messaging and notices
- WHOIS, WHO, and NAMES commands
- Ping/Pong keepalive mechanism

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
- **Server Notice Masks (SNOmasks)**:
  - `+c` (connection notices)
  - `+k` (kill notices)
  - `+o` (oper notices)
  - `+x` (ban/quiet notices)
  - `+f` (flood notices)
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
go mod tidy
go build -o techircd
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
