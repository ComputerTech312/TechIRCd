package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	config        *Config
	clients       map[string]*Client
	channels      map[string]*Channel
	listener      net.Listener
	sslListener   net.Listener
	mu            sync.RWMutex
	shutdown      chan bool
	healthMonitor *HealthMonitor
}

func NewServer(config *Config) *Server {
	server := &Server{
		config:   config,
		clients:  make(map[string]*Client),
		channels: make(map[string]*Channel),
		shutdown: make(chan bool),
	}
	server.healthMonitor = NewHealthMonitor(server)
	return server
}

func (s *Server) Start() error {
	// Start regular listener
	addr := fmt.Sprintf("%s:%d", s.config.Server.Listen.Host, s.config.Server.Listen.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}
	s.listener = listener

	log.Printf("IRC server listening on %s", addr)

	// Start health monitoring
	s.healthMonitor.Start()

	// Start SSL listener if enabled
	if s.config.Server.Listen.EnableSSL {
		go s.startSSLListener()
	}

	// Auto-create configured channels
	for _, channelName := range s.config.Channels.AutoJoin {
		channel := NewChannel(channelName)
		// Set default modes
		for _, mode := range s.config.Channels.DefaultModes {
			if mode != '+' {
				channel.SetMode(rune(mode), true)
			}
		}
		s.channels[strings.ToLower(channelName)] = channel
	}

	// Start ping routine
	go s.pingRoutine()

	// Accept connections
	for {
		select {
		case <-s.shutdown:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			client := NewClient(conn, s)
			s.AddClient(client)
			go client.Handle()
		}
	}
}

func (s *Server) startSSLListener() {
	// Load SSL certificates
	cert, err := tls.LoadX509KeyPair(s.config.Server.SSL.CertFile, s.config.Server.SSL.KeyFile)
	if err != nil {
		log.Printf("Failed to load SSL certificates: %v", err)
		return
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	addr := fmt.Sprintf("%s:%d", s.config.Server.Listen.Host, s.config.Server.Listen.SSLPort)
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		log.Printf("Failed to start SSL listener on %s: %v", addr, err)
		return
	}
	s.sslListener = listener

	log.Printf("IRC SSL server listening on %s", addr)

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			client := NewClient(conn, s)
			s.AddClient(client)
			go client.Handle()
		}
	}
}

func (s *Server) pingRoutine() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdown:
			return
		case <-ticker.C:
			s.mu.RLock()
			for _, client := range s.clients {
				if client.IsRegistered() {
					client.SendMessage(fmt.Sprintf("PING :%s", s.config.Server.Name))
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *Server) AddClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check client limit
	if len(s.clients) >= s.config.Limits.MaxClients {
		client.SendMessage("ERROR :Server full")
		client.conn.Close()
		return
	}

	s.clients[client.Host()] = client
}

func (s *Server) RemoveClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client.Host())
	
	// Send snomask notification for client disconnect
	if client.IsRegistered() {
		s.sendSnomask('c', fmt.Sprintf("Client disconnect: %s (%s@%s)", 
			client.Nick(), client.User(), client.Host()))
	}
}

// sendSnomask sends a server notice to operators watching a specific snomask
func (s *Server) sendSnomask(snomask rune, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		if client.IsOper() && client.HasSnomask(snomask) {
			client.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** %s", 
				s.config.Server.Name, client.Nick(), message))
		}
	}
}

// ReloadConfig reloads the server configuration
func (s *Server) ReloadConfig() error {
	config, err := LoadConfig("config.json")
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	
	s.mu.Lock()
	s.config = config
	s.mu.Unlock()
	
	return nil
}

func (s *Server) GetClient(nick string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		if strings.EqualFold(client.Nick(), nick) {
			return client
		}
	}
	return nil
}

func (s *Server) GetClientByHost(host string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[host]
}

func (s *Server) GetClients() map[string]*Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	clients := make(map[string]*Client)
	for key, client := range s.clients {
		clients[key] = client
	}
	return clients
}

func (s *Server) GetChannel(name string) *Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.channels[strings.ToLower(name)]
}

func (s *Server) GetOrCreateChannel(name string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	channelName := strings.ToLower(name)
	if channel, exists := s.channels[channelName]; exists {
		return channel
	}

	// Create new channel
	channel := NewChannel(name)
	// Set default modes
	for _, mode := range s.config.Channels.DefaultModes {
		if mode != '+' {
			channel.SetMode(rune(mode), true)
		}
	}
	s.channels[channelName] = channel
	return channel
}

func (s *Server) RemoveChannel(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.channels, strings.ToLower(name))
}

func (s *Server) GetChannels() map[string]*Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels := make(map[string]*Channel)
	for name, channel := range s.channels {
		channels[name] = channel
	}
	return channels
}

func (s *Server) CreateChannel(name string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.channels) >= s.config.Limits.MaxChannels {
		return nil
	}

	channel := NewChannel(name)
	s.channels[strings.ToLower(name)] = channel
	return channel
}

func (s *Server) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

func (s *Server) GetChannelCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.channels)
}

func (s *Server) IsNickInUse(nick string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		if strings.EqualFold(client.Nick(), nick) {
			return true
		}
	}
	return false
}

func (s *Server) HandleMessage(client *Client, message string) {
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return
	}

	command := strings.ToUpper(parts[0])

	// Log the command for debugging
	log.Printf("Client %s: %s", client.Host(), message)

	switch command {
	case "CAP":
		// TODO: implement CAP handling
	case "NICK":
		client.handleNick(parts)
	case "USER":
		client.handleUser(parts)
	case "PING":
		client.handlePing(parts)
	case "PONG":
		client.handlePong(parts)
	case "JOIN":
		client.handleJoin(parts)
	case "PART":
		client.handlePart(parts)
	case "PRIVMSG":
		client.handlePrivmsg(parts)
	case "NOTICE":
		client.handleNotice(parts)
	case "WHO":
		client.handleWho(parts)
	case "WHOIS":
		client.handleWhois(parts)
	case "NAMES":
		client.handleNames(parts)
	case "MODE":
		client.handleMode(parts)
	case "OPER":
		client.handleOper(parts)
	case "SNOMASK":
		client.handleSnomask(parts)
	case "GLOBALNOTICE":
		client.handleGlobalNotice(parts)
	case "OPERWALL":
		client.handleOperWall(parts)
	case "WALLOPS":
		client.handleWallops(parts)
	case "REHASH":
		client.handleRehash(parts)
	case "TRACE":
		client.handleTrace(parts)
	case "TOPIC":
		client.handleTopic(parts)
	case "KICK":
		client.handleKick(parts)
	case "INVITE":
		client.handleInvite(parts)
	case "AWAY":
		client.handleAway(parts)
	case "LIST":
		client.handleList(parts)
	case "KILL":
		client.handleKill(parts)
	case "QUIT":
		client.handleQuit(parts)
	default:
		client.SendNumeric(ERR_UNKNOWNCOMMAND, command+" :Unknown command")
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	log.Println("Initiating graceful shutdown...")

	// Stop health monitoring
	if s.healthMonitor != nil {
		s.healthMonitor.Stop()
	}

	// Signal shutdown
	close(s.shutdown)

	// Notify all clients
	s.mu.RLock()
	for _, client := range s.clients {
		client.SendMessage("ERROR :Server shutting down")
	}
	s.mu.RUnlock()

	// Close listeners
	if s.listener != nil {
		s.listener.Close()
	}
	if s.sslListener != nil {
		s.sslListener.Close()
	}

	log.Println("Server shutdown complete")
}
