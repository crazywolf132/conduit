package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazywolf132/conduit"
	"github.com/crazywolf132/conduit/client"
	"github.com/crazywolf132/conduit/server"
)

const socketPath = "/tmp/chat.sock"

// ChatMessage defines the JSON payload for chat messages.
type ChatMessage struct {
	Username string    `json:"username"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}

// runServer runs a simple chat server that broadcasts messages to all connected clients.
func runServer() error {
	config := conduit.DefaultServerConfig(socketPath)
	config.Logger = conduit.NewLogger(conduit.LogDebug, os.Stdout)

	s := server.NewServer(config)

	// Handle 'chat' messages by broadcasting them.
	s.Handle("chat", func(conn *server.Connection, msg *conduit.Message) error {
		var chatMsg ChatMessage
		if err := msg.UnmarshalPayload(&chatMsg); err != nil {
			return err
		}

		// Store the username in the connection's context for potential future use.
		conn.SetContext("username", chatMsg.Username)

		// Broadcast the chat message to all clients.
		return s.Broadcast("chat", chatMsg)
	})

	if err := s.Start(); err != nil {
		return err
	}

	// Wait for OS interrupt signals (e.g., Ctrl+C) to shut down gracefully.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return s.Stop()
}

// runClient runs a simple chat client that connects to the server, listens for chat messages,
// and sends user input as chat messages.
func runClient(username string) error {
	config := conduit.DefaultClientConfig(socketPath)
	config.Logger = conduit.NewLogger(conduit.LogInfo, os.Stdout)
	config.Reconnect = true
	config.ReconnectDelay = 5 * time.Second

	c := client.NewClient(config)

	// Handle incoming chat messages by printing them to stdout.
	c.Handle("chat", func(_ *client.Client, msg *conduit.Message) error {
		var chatMsg ChatMessage
		if err := msg.UnmarshalPayload(&chatMsg); err != nil {
			return err
		}

		fmt.Printf("[%s] %s: %s\n",
			chatMsg.Time.Format("15:04:05"),
			chatMsg.Username,
			chatMsg.Message,
		)
		return nil
	})

	if err := c.ConnectWithRetry(); err != nil {
		return err
	}
	defer c.Close()

	// Store the username in the client's context for potential use.
	c.SetContext("username", username)

	fmt.Println("Connected to chat server. Type your messages (Ctrl+C to quit):")

	// Handle Ctrl+C to gracefully shut down the client.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nDisconnecting...")
		c.Close()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := ChatMessage{
			Username: username,
			Message:  scanner.Text(),
			Time:     time.Now(),
		}

		if err := c.Send("chat", msg); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			if !c.IsConnected() {
				fmt.Println("Lost connection to server. Waiting for reconnection...")
			}
		}
	}

	return scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Start server: chat server")
		fmt.Println("  Start client: chat client <username>")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "server":
		err = runServer()
	case "client":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a username")
			os.Exit(1)
		}
		err = runClient(os.Args[2])
	default:
		fmt.Println("Invalid command. Use 'server' or 'client'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
