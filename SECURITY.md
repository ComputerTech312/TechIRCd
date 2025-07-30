# Security Policy

## Supported Versions

We currently support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of TechIRCd seriously. If you believe you have found a security vulnerability, please report it to us as described below.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: security@techircd.org (or the maintainer's email)

Please include the following information (as much as you can provide) to help us better understand the nature and scope of the possible issue:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

## Preferred Languages

We prefer all communications to be in English.

## Response Time

We will respond to your report within 48 hours and provide regular updates at least every 72 hours.

## Security Measures

TechIRCd implements several security measures:

- Input validation and sanitization
- Flood protection and rate limiting
- Connection timeout management
- Panic recovery mechanisms
- Memory usage monitoring
- Secure configuration validation

## Responsible Disclosure

We follow responsible disclosure practices:

1. **Report**: Submit vulnerability report privately
2. **Acknowledge**: We acknowledge receipt within 48 hours
3. **Investigate**: We investigate and develop a fix
4. **Coordinate**: We coordinate disclosure timeline with reporter
5. **Release**: We release security update and public disclosure
