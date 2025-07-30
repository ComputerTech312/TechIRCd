package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn       net.Conn
	nick       string
	user       string
	realname   string
	host       string
	server     *Server
	channels   map[string]*Channel
	modes      map[rune]bool
	away       string
	oper       bool
	operClass  string // Operator class name
	ssl        bool
	registered bool
	account    string // Services account name
	connectTime time.Time // When the client connected
	lastActivity time.Time // Last time client sent a message

	// Flood protection
	lastMessage  time.Time
	messageCount int

	// SASL authentication
	saslMech string
	saslData string

	// IRCv3 capabilities
	capabilities map[string]bool

	// Server Notice Masks (snomasks) for operators
	snomasks map[rune]bool

	// Ping timeout tracking
	lastPong       time.Time
	waitingForPong bool

	mu sync.RWMutex
}

func NewClient(conn net.Conn, server *Server) *Client {
	host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

	// Check if connection is SSL
	isSSL := false
	if _, ok := conn.(*tls.Conn); ok {
		isSSL = true
	}

	client := &Client{
		conn:           conn,
		host:           host,
		server:         server,
		channels:       make(map[string]*Channel),
		modes:          make(map[rune]bool),
		capabilities:   make(map[string]bool),
		snomasks:       make(map[rune]bool),
		ssl:            isSSL,
		connectTime:    time.Now(),
		lastActivity:   time.Now(),
		lastMessage:    time.Now(),
		lastPong:       time.Now(),
		waitingForPong: false,
	}

	// Set SSL user mode if connected via SSL
	if isSSL {
		client.SetMode('z', true)
	}

	return client
}

func (c *Client) SendMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add connection health check
	if c.conn == nil {
		return
	}

	// Set write deadline to prevent hanging
	c.conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	defer c.conn.SetWriteDeadline(time.Time{}) // Clear deadline

	_, err := fmt.Fprintf(c.conn, "%s\r\n", message)
	if err != nil {
		// Log the error but don't panic - connection will be cleaned up
		if c.server != nil {
			log.Printf("Error sending message to %s: %v", c.Nick(), err)
		}
	}
}

func (c *Client) SendFrom(source, message string) {
	c.SendMessage(fmt.Sprintf(":%s %s", source, message))
}

func (c *Client) SendNumeric(code int, message string) {
	if c.server == nil || c.server.config == nil {
		return
	}
	c.SendFrom(c.server.config.Server.Name, fmt.Sprintf("%03d %s %s", code, c.Nick(), message))
}

func (c *Client) Nick() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nick
}

func (c *Client) SetNick(nick string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nick = nick
}

func (c *Client) User() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.user
}

func (c *Client) SetUser(user string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.user = user
}

func (c *Client) Realname() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.realname
}

func (c *Client) SetRealname(realname string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.realname = realname
}

func (c *Client) Host() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.host
}

// HostForUser returns the appropriate hostname to show to a requesting user
// based on privacy settings and the requester's privileges
func (c *Client) HostForUser(requester *Client) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// If host hiding is disabled, always show real host
	if !c.server.config.Privacy.HideHostsFromUsers {
		return c.host
	}
	
	// If requester is an operator and bypass is enabled, show real host
	if requester != nil && requester.IsOper() && c.server.config.Privacy.OperBypassHostHide {
		return c.host
	}
	
	// If requester is viewing themselves, show real host
	if requester != nil && requester.Nick() == c.Nick() {
		return c.host
	}
	
	// Check if user has +x mode set (host masking)
	if c.HasMode('x') {
		// Return masked hostname
		return c.nick + "." + c.server.config.Privacy.MaskedHostSuffix
	}
	
	// Default behavior: show masked host when privacy is enabled
	return c.nick + "." + c.server.config.Privacy.MaskedHostSuffix
}

// canSeeWhoisInfo checks if the requester can see specific WHOIS information about the target
func (c *Client) canSeeWhoisInfo(target *Client, infoType string) bool {
	config := c.server.config.WhoisFeatures
	
	var toEveryone, toOpers, toSelf bool
	
	switch infoType {
	case "user_modes":
		toEveryone = config.ShowUserModes.ToEveryone
		toOpers = config.ShowUserModes.ToOpers
		toSelf = config.ShowUserModes.ToSelf
	case "ssl_status":
		toEveryone = config.ShowSSLStatus.ToEveryone
		toOpers = config.ShowSSLStatus.ToOpers
		toSelf = config.ShowSSLStatus.ToSelf
	case "idle_time":
		toEveryone = config.ShowIdleTime.ToEveryone
		toOpers = config.ShowIdleTime.ToOpers
		toSelf = config.ShowIdleTime.ToSelf
	case "signon_time":
		toEveryone = config.ShowSignonTime.ToEveryone
		toOpers = config.ShowSignonTime.ToOpers
		toSelf = config.ShowSignonTime.ToSelf
	case "real_host":
		toEveryone = config.ShowRealHost.ToEveryone
		toOpers = config.ShowRealHost.ToOpers
		toSelf = config.ShowRealHost.ToSelf
	case "oper_class":
		toEveryone = config.ShowOperClass.ToEveryone
		toOpers = config.ShowOperClass.ToOpers
		toSelf = config.ShowOperClass.ToSelf
	case "client_info":
		toEveryone = config.ShowClientInfo.ToEveryone
		toOpers = config.ShowClientInfo.ToOpers
		toSelf = config.ShowClientInfo.ToSelf
	case "account_name":
		toEveryone = config.ShowAccountName.ToEveryone
		toOpers = config.ShowAccountName.ToOpers
		toSelf = config.ShowAccountName.ToSelf
	default:
		return false
	}
	
	// Check if target is viewing themselves
	if c.Nick() == target.Nick() && toSelf {
		return true
	}
	
	// Check if requester is an operator
	if c.IsOper() && toOpers {
		return true
	}
	
	// Check if everyone can see this info
	if toEveryone {
		return true
	}
	
	return false
}

// canSeeChannels checks if the requester can see channel information
func (c *Client) canSeeChannels(target *Client) bool {
	config := c.server.config.WhoisFeatures.ShowChannels
	
	// Check if target is viewing themselves
	if c.Nick() == target.Nick() && config.ToSelf {
		return true
	}
	
	// Check if requester is an operator
	if c.IsOper() && config.ToOpers {
		return true
	}
	
	// Check if everyone can see this info
	if config.ToEveryone {
		return true
	}
	
	return false
}

// UpdateActivity updates the last activity time for the client
func (c *Client) UpdateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

// SetAccount sets the account name for services integration
func (c *Client) SetAccount(account string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.account = account
}

// Account returns the account name
func (c *Client) Account() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.account
}

// ConnectTime returns when the client connected
func (c *Client) ConnectTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connectTime
}

// LastActivity returns the last activity time
func (c *Client) LastActivity() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastActivity
}

// SetOperClass sets the operator class for this client
func (c *Client) SetOperClass(class string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.operClass = class
}

// OperClass returns the operator class
func (c *Client) OperClass() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.operClass
}

// HasOperPermission checks if the client has a specific operator permission
func (c *Client) HasOperPermission(permission string) bool {
	if !c.IsOper() {
		return false
	}
	
	// Load oper config and check permissions
	operConfig, err := LoadOperConfig(c.server.config.OperConfig.ConfigFile)
	if err != nil || !c.server.config.OperConfig.Enable {
		// Fallback to legacy system
		return true // Basic oper permissions
	}
	
	return operConfig.HasPermission(c.Nick(), permission)
}

// GetOperRank returns the operator rank (higher number = higher authority)
func (c *Client) GetOperRank() int {
	if !c.IsOper() {
		return 0
	}
	
	operConfig, err := LoadOperConfig(c.server.config.OperConfig.ConfigFile)
	if err != nil || !c.server.config.OperConfig.Enable {
		return 1 // Basic rank for legacy
	}
	
	return operConfig.GetOperRank(c.Nick())
}

// CanOperateOn checks if this operator can perform actions on another operator
func (c *Client) CanOperateOn(target *Client) bool {
	if !c.IsOper() {
		return false
	}
	
	if !target.IsOper() {
		return true // Opers can operate on regular users
	}
	
	operConfig, err := LoadOperConfig(c.server.config.OperConfig.ConfigFile)
	if err != nil || !c.server.config.OperConfig.Enable {
		return true // Legacy behavior
	}
	
	return operConfig.CanOperateOn(c.Nick(), target.Nick())
}

// GetOperSymbol returns the symbol for this operator class
func (c *Client) GetOperSymbol() string {
	if !c.IsOper() {
		return ""
	}
	
	operConfig, err := LoadOperConfig(c.server.config.OperConfig.ConfigFile)
	if err != nil || !c.server.config.OperConfig.Enable {
		return "*" // Default symbol
	}
	
	class := operConfig.GetOperClass(c.OperClass())
	if class == nil {
		return "*"
	}
	
	return class.Symbol
}

func (c *Client) IsRegistered() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.registered
}

func (c *Client) SetRegistered(registered bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.registered = registered
}

func (c *Client) IsOper() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.oper
}

func (c *Client) SetOper(oper bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.oper = oper
}

func (c *Client) IsSSL() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ssl
}

func (c *Client) Away() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.away
}

func (c *Client) SetAway(away string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.away = away
}

func (c *Client) HasMode(mode rune) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.modes[mode]
}

func (c *Client) SetMode(mode rune, set bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if set {
		c.modes[mode] = true
	} else {
		delete(c.modes, mode)
	}
}

func (c *Client) GetModes() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var modes []rune
	for mode := range c.modes {
		modes = append(modes, mode)
	}

	if len(modes) == 0 {
		return ""
	}

	return "+" + string(modes)
}

func (c *Client) HasSnomask(snomask rune) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snomasks[snomask]
}

func (c *Client) SetSnomask(snomask rune, set bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if set {
		c.snomasks[snomask] = true
	} else {
		delete(c.snomasks, snomask)
	}
}

func (c *Client) GetSnomasks() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var snomasks []rune
	for snomask := range c.snomasks {
		snomasks = append(snomasks, snomask)
	}

	if len(snomasks) == 0 {
		return ""
	}

	return "+" + string(snomasks)
}

// HasGodMode returns true if the client has god mode enabled (user mode +G)
func (c *Client) HasGodMode() bool {
	return c.HasMode('G') && c.HasOperPermission("god_mode")
}

// HasStealthMode returns true if the client has stealth mode enabled (user mode +S)
func (c *Client) HasStealthMode() bool {
	return c.HasMode('S') && c.HasOperPermission("stealth_mode")
}

// IsVisibleTo returns true if this client should be visible to the target client
func (c *Client) IsVisibleTo(target *Client) bool {
	// If client doesn't have stealth mode, always visible
	if !c.HasStealthMode() {
		return true
	}
	
	// If target is an operator, stealth users are visible
	if target.IsOper() {
		return true
	}
	
	// If target is the stealth user themselves, always visible
	if target == c {
		return true
	}
	
	// Otherwise, stealth users are invisible to normal users
	return false
}

// CanBypassChannelRestrictions returns true if the client can bypass channel restrictions
func (c *Client) CanBypassChannelRestrictions() bool {
	return c.HasGodMode()
}

func (c *Client) AddChannel(channel *Channel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.channels[strings.ToLower(channel.name)] = channel
}

func (c *Client) RemoveChannel(channelName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.channels, strings.ToLower(channelName))
}

func (c *Client) IsInChannel(channelName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.channels[strings.ToLower(channelName)]
	return exists
}

func (c *Client) GetChannels() map[string]*Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channels := make(map[string]*Channel)
	for name, channel := range c.channels {
		channels[name] = channel
	}
	return channels
}

func (c *Client) Prefix() string {
	return fmt.Sprintf("%s!%s@%s", c.Nick(), c.User(), c.Host())
}

// sendSnomask sends a server notice to all operators with the specified snomask
func (c *Client) sendSnomask(snomask rune, message string) {
	if c.server == nil {
		return
	}
	
	serverName := "localhost"
	if c.server.config != nil {
		serverName = c.server.config.Server.Name
	}
	
	// Send to all operators with this snomask
	c.server.mu.RLock()
	clients := make([]*Client, 0, len(c.server.clients))
	for _, client := range c.server.clients {
		clients = append(clients, client)
	}
	c.server.mu.RUnlock()
	
	for _, client := range clients {
		if client.IsOper() && client.HasSnomask(snomask) {
			client.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** %s", serverName, client.Nick(), message))
		}
	}
}

// getServerConfig safely returns the server config, or nil if not available
func (c *Client) getServerConfig() *Config {
	if c.server == nil {
		return nil
	}
	return c.server.config
}

// getRegistrationTimeout safely gets the registration timeout duration
func (c *Client) getRegistrationTimeout() time.Duration {
	config := c.getServerConfig()
	if config == nil {
		return 60 * time.Second // default 60 seconds
	}
	return config.RegistrationTimeoutDuration()
}

func (c *Client) CheckFlood() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add nil checks
	if c.server == nil || c.server.config == nil {
		return false
	}

	// IRC operators are exempt from flood protection
	if c.oper {
		return false
	}

	// Be very lenient with flood protection for unregistered clients
	// during the initial connection phase (first 60 seconds)
	if !c.registered {
		// Allow up to 100 commands per minute for unregistered clients
		now := time.Now()
		if now.Sub(c.lastMessage) > 60*time.Second {
			c.messageCount = 0
		}
		c.messageCount++
		c.lastMessage = now
		return c.messageCount > 100
	}

	// For registered clients, use the configured limits but make them more reasonable
	now := time.Now()
	floodWindow := time.Duration(c.server.config.Limits.FloodSeconds) * time.Second
	if now.Sub(c.lastMessage) > floodWindow {
		c.messageCount = 0
	}

	c.messageCount++
	c.lastMessage = now

	// Use higher limits than configured for better user experience
	maxLines := c.server.config.Limits.FloodLines * 3 // Triple the configured limit
	return c.messageCount > maxLines
}

func (c *Client) HasCapability(cap string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.capabilities[cap]
}

func (c *Client) SetCapability(cap string, enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if enabled {
		c.capabilities[cap] = true
	} else {
		delete(c.capabilities, cap)
	}
}

func (c *Client) Handle() {
	defer func() {
		// Panic recovery
		if r := recover(); r != nil {
			log.Printf("Panic in client handler for %s: %v", c.Nick(), r)
		}

		// Cleanup
		if c.conn != nil {
			c.conn.Close()
		}
		if c.server != nil {
			c.server.RemoveClient(c)
		}

		// Part all channels
		for _, channel := range c.GetChannels() {
			channel.RemoveClient(c)
			if len(channel.GetClients()) == 0 && c.server != nil {
				c.server.RemoveChannel(channel.name)
			}
		}
	}()

	// Add nil checks
	if c.server == nil || c.server.config == nil {
		c.conn.Close()
		return
	}

	scanner := bufio.NewScanner(c.conn)

	// Set maximum line length to prevent memory exhaustion
	const maxLineLength = 4096
	scanner.Buffer(make([]byte, maxLineLength), maxLineLength)

	// Set read deadline for scanner
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	// Set registration timeout
	registrationTimer := time.NewTimer(c.getRegistrationTimeout())
	defer registrationTimer.Stop()
	registrationActive := true

	// Set up ping timeout mechanism
	pingInterval := 30 * time.Second // Send ping every 30 seconds
	pingTimeout := c.server.config.PingTimeoutDuration()

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	// Initialize ping state
	c.mu.Lock()
	c.lastPong = time.Now()
	c.waitingForPong = false
	c.mu.Unlock()

	for {
		select {
		case <-registrationTimer.C:
			if registrationActive && !c.IsRegistered() {
				c.SendMessage("ERROR :Registration timeout")
				return
			}
		case <-pingTicker.C:
			if c.IsRegistered() {
				c.mu.RLock()
				waitingForPong := c.waitingForPong
				lastPong := c.lastPong
				c.mu.RUnlock()

				if waitingForPong && time.Since(lastPong) > pingTimeout {
					c.SendMessage("ERROR :Ping timeout")
					return
				}
				// Send ping
				c.SendMessage(fmt.Sprintf("PING :%s", c.server.config.Server.Name))

				c.mu.Lock()
				c.waitingForPong = true
				c.mu.Unlock()
			}
		default:
			// Reset read deadline for each message
			c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

			if !scanner.Scan() {
				// Check for scanner error
				if err := scanner.Err(); err != nil {
					log.Printf("Scanner error for client %s: %v", c.Nick(), err)
				}
				return
			}

			// Enhanced flood checking
			if c.CheckFlood() {
				c.SendMessage("ERROR :Excess Flood")
				return
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// Additional input validation
			if len(line) > 4096 {
				c.SendMessage("ERROR :Line too long")
				continue
			}

			// Handle the message through the server command router
			if c.server != nil {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Panic handling message from %s: %v", c.Nick(), r)
							c.SendMessage("ERROR :Internal server error")
						}
					}()
					c.server.HandleMessage(c, line)
				}()
			}

			// Stop registration timer once registered
			if c.IsRegistered() && registrationActive {
				registrationTimer.Stop()
				registrationActive = false
			}
		}
	}
}
