# Changelog

All notable changes to TechIRCd will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-07-30

### Added
- Initial release of TechIRCd
- Full IRC protocol implementation (RFC 2812 compliant)
- Advanced channel management with operator hierarchy
- Extended ban system with quiet mode support (~q:mask)
- Comprehensive operator features and SNOmasks
- User modes system with SSL detection
- Enhanced stability with panic recovery
- Real-time health monitoring and metrics
- Configuration validation and sanitization
- Graceful shutdown capabilities
- Flood protection with operator exemption

### Features
- Channel modes: +m +n +t +i +s +p +k +l +b
- User modes: +i +w +s +o +x +B +z +r
- Operator commands: KILL, GLOBALNOTICE, OPERWALL, WALLOPS, REHASH, TRACE
- SNOmasks: +c +k +o +x +f +n +s +d
- Extended ban types with quiet mode support
- Health monitoring with memory and goroutine tracking
- Comprehensive error handling and recovery systems

[1.0.0]: https://github.com/ComputerTech312/TechIRCd/releases/tag/v1.0.0
