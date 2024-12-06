package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazywolf132/conduit"
	"github.com/crazywolf132/conduit/client"
	"github.com/crazywolf132/conduit/server"
)

const socketPath = "/tmp/pingpong.sock"

// runServer starts a simple server that responds to "ping" messages with "pong".
func runServer() error {
	cfg := conduit.DefaultServerConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogInfo, os.Stdout)

	s := server.NewServer(cfg)

	// Handle "ping" messages by replying with "pong".
	s.Handle("ping", func(conn *server.Connection, msg *conduit.Message) error {
		return conn.Send("pong", "pong")
	})

	if err := s.Start(); err != nil {
		return err
	}

	// Wait for OS interrupt signals to shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return s.Stop()
}

// runClient connects to the server, sends a "ping" message and prints any "pong" responses.
func runClient() error {
	cfg := conduit.DefaultClientConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogInfo, os.Stdout)

	c := client.NewClient(cfg)

	// Handle "pong" responses from server
	c.Handle("pong", func(_ *client.Client, msg *conduit.Message) error {
		var response string
		if err := msg.UnmarshalPayload(&response); err != nil {
			return err
		}
		fmt.Println("Received from server:", response)
		return nil
	})

	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	// Send "ping" to the server
	if err := c.Send("ping", "ping"); err != nil {
		return fmt.Errorf("failed to send ping: %w", err)
	}

	// Give some time to receive "pong"
	time.Sleep(1 * time.Second)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  pingpong server")
		fmt.Println("  pingpong client")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "server":
		err = runServer()
	case "client":
		err = runClient()
	default:
		fmt.Println("Invalid command. Use 'server' or 'client'.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
