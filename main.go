package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Printf("Failed to load config from %s, using defaults: %v", *configFile, err)
		// Create default config if file doesn't exist
		config = DefaultConfig()
		if err := SaveConfig(config, *configFile); err != nil {
			log.Printf("Failed to save default config: %v", err)
		}
	}

	// Validate and sanitize configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	config.SanitizeConfig()
	log.Println("Configuration validated successfully")

	// Create and start the server
	server := NewServer(config)

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		server.Shutdown()
		os.Exit(0)
	}()

	// Start the server
	log.Printf("Starting TechIRCd %s on %s:%d", config.Server.Version, config.Server.Listen.Host, config.Server.Listen.Port)
	if config.Server.Listen.EnableSSL {
		log.Printf("SSL enabled with cert: %s, key: %s", config.Server.SSL.CertFile, config.Server.SSL.KeyFile)
	}
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
