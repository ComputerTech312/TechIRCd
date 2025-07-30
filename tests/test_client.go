package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// TestClient is a simple IRC client for testing
func testClient() {
	conn, err := net.Dial("tcp", "localhost:6667")
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	defer conn.Close()

	// Start reading responses
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("< %s\n", scanner.Text())
		}
	}()

	// Send registration
	fmt.Fprintf(conn, "NICK testuser\r\n")
	fmt.Fprintf(conn, "USER testuser 0 * :Test User\r\n")

	// Wait for registration
	time.Sleep(2 * time.Second)

	// Join a channel
	fmt.Fprintf(conn, "JOIN #test\r\n")

	// Send a message
	time.Sleep(1 * time.Second)
	fmt.Fprintf(conn, "PRIVMSG #test :Hello, world!\r\n")

	// Interactive mode
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter IRC commands (or 'quit' to exit):")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "quit" {
			fmt.Fprintf(conn, "QUIT :Goodbye\r\n")
			break
		}

		if line != "" {
			fmt.Fprintf(conn, "%s\r\n", line)
		}
	}
}
