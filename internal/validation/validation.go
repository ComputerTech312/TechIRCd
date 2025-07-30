package config

import (
	"fmt"
	"strings"
)

// ValidateConfig performs comprehensive validation of the server configuration
func (c *Config) Validate() error {
	// Validate server settings
	if c.Server.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	if c.Server.Network == "" {
		return fmt.Errorf("network name cannot be empty")
	}

	if c.Server.Listen.Host == "" {
		c.Server.Listen.Host = "localhost" // Default fallback
	}

	if c.Server.Listen.Port <= 0 || c.Server.Listen.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Listen.Port)
	}

	if c.Server.Listen.EnableSSL && (c.Server.Listen.SSLPort <= 0 || c.Server.Listen.SSLPort > 65535) {
		return fmt.Errorf("invalid SSL port number: %d", c.Server.Listen.SSLPort)
	}

	// Validate limits
	if c.Limits.MaxClients <= 0 {
		c.Limits.MaxClients = 1000 // Default
	}

	if c.Limits.MaxChannels <= 0 {
		c.Limits.MaxChannels = 100 // Default
	}

	if c.Limits.MaxNickLength <= 0 || c.Limits.MaxNickLength > 50 {
		c.Limits.MaxNickLength = 30 // Default
	}

	if c.Limits.PingTimeout <= 0 {
		c.Limits.PingTimeout = 300 // Default 5 minutes
	}

	if c.Limits.FloodLines <= 0 {
		c.Limits.FloodLines = 10 // Default
	}

	if c.Limits.FloodSeconds <= 0 {
		c.Limits.FloodSeconds = 60 // Default
	}

	// Validate channels
	for _, channelName := range c.Channels.AutoJoin {
		if !isChannelName(channelName) {
			return fmt.Errorf("invalid channel name in auto_join: %s", channelName)
		}
	}

	// Validate default modes
	validChannelModes := "mntisp"
	for _, mode := range c.Channels.DefaultModes {
		if mode != '+' && !strings.ContainsRune(validChannelModes, mode) {
			return fmt.Errorf("invalid default channel mode: %c", mode)
		}
	}

	// Validate operators
	for i, oper := range c.Opers {
		if oper.Name == "" {
			return fmt.Errorf("operator %d: name cannot be empty", i)
		}
		if oper.Password == "" {
			return fmt.Errorf("operator %s: password cannot be empty", oper.Name)
		}
		if oper.Host == "" {
			return fmt.Errorf("operator %s: host cannot be empty", oper.Name)
		}
	}

	return nil
}

// SanitizeConfig applies safe defaults and sanitizes configuration values
func (c *Config) SanitizeConfig() {
	// Ensure reasonable limits
	if c.Limits.MaxClients > 10000 {
		c.Limits.MaxClients = 10000
	}

	if c.Limits.MaxChannels > 1000 {
		c.Limits.MaxChannels = 1000
	}

	// Ensure reasonable string lengths
	if c.Limits.MaxTopicLength > 2048 {
		c.Limits.MaxTopicLength = 2048
	}

	if c.Limits.MaxKickLength > 2048 {
		c.Limits.MaxKickLength = 2048
	}

	// Ensure MOTD isn't too long
	if len(c.MOTD) > 50 {
		c.MOTD = c.MOTD[:50]
	}

	// Sanitize channel modes
	if c.Channels.DefaultModes == "" {
		c.Channels.DefaultModes = "+nt"
	}
}
