package test

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/crazywolf132/conduit"
	"github.com/crazywolf132/conduit/client"
	"github.com/crazywolf132/conduit/server"
)

// TestServerStartStop tests that the server starts and stops cleanly.
func TestServerStartStop(t *testing.T) {
	socketPath := "/tmp/conduit_server_test.sock"
	defer os.RemoveAll(socketPath)

	cfg := conduit.DefaultServerConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	s := server.NewServer(cfg)

	if err := s.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Check that the socket file was created
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		t.Fatalf("Expected socket file %s, got none", socketPath)
	}

	// Stop the server
	if err := s.Stop(); err != nil {
		t.Fatalf("Failed to stop server: %v", err)
	}

	// Socket file should be removed
	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		t.Fatalf("Expected socket file to be removed, but it still exists")
	}
}

// TestServerBroadcast tests that the server can broadcast messages to all clients.
func TestServerBroadcast(t *testing.T) {
	socketPath := "/tmp/conduit_broadcast_test.sock"
	defer os.RemoveAll(socketPath)

	serverCfg := conduit.DefaultServerConfig(socketPath)
	serverCfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	srv := server.NewServer(serverCfg)

	// No special handler needed, we will broadcast manually.
	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	clientCfg := conduit.DefaultClientConfig(socketPath)
	clientCfg.Logger = conduit.NewLogger(conduit.LogError, nil)

	c1 := client.NewClient(clientCfg)
	if err := c1.Connect(); err != nil {
		t.Fatalf("Client1 failed to connect: %v", err)
	}
	defer c1.Close()

	c2 := client.NewClient(clientCfg)
	if err := c2.Connect(); err != nil {
		t.Fatalf("Client2 failed to connect: %v", err)
	}
	defer c2.Close()

	received1 := make(chan string, 1)
	c1.Handle("announcement", func(_ *client.Client, msg *conduit.Message) error {
		var m string
		if err := msg.UnmarshalPayload(&m); err != nil {
			return err
		}
		received1 <- m
		return nil
	})

	received2 := make(chan string, 1)
	c2.Handle("announcement", func(_ *client.Client, msg *conduit.Message) error {
		var m string
		if err := msg.UnmarshalPayload(&m); err != nil {
			return err
		}
		received2 <- m
		return nil
	})

	// Give clients some time to set handlers
	time.Sleep(100 * time.Millisecond)

	// Broadcast message
	if err := srv.Broadcast("announcement", "Hello, everyone!"); err != nil {
		t.Fatalf("Failed to broadcast: %v", err)
	}

	select {
	case msg := <-received1:
		if msg != "Hello, everyone!" {
			t.Errorf("Client1 expected 'Hello, everyone!', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for client1 to receive broadcast")
	}

	select {
	case msg := <-received2:
		if msg != "Hello, everyone!" {
			t.Errorf("Client2 expected 'Hello, everyone!', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for client2 to receive broadcast")
	}
}

// TestServerConnection tests if a client can connect to the server.
func TestServerConnection(t *testing.T) {
	socketPath := "/tmp/conduit_connection_test.sock"
	defer os.RemoveAll(socketPath)

	cfg := conduit.DefaultServerConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	s := server.NewServer(cfg)

	if err := s.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer s.Stop()

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close() // close the connection immediately

	// If we reached here, the server accepted the connection.
}

// TestServerUnsupportedMessageType tests that the server logs a warning when it receives a message type without a handler.
func TestServerUnsupportedMessageType(t *testing.T) {
	socketPath := "/tmp/conduit_unsupported_test.sock"
	defer os.RemoveAll(socketPath)

	cfg := conduit.DefaultServerConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil) // silent for tests
	s := server.NewServer(cfg)

	if err := s.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer s.Stop()

	// Connect to the server manually
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send a message type with no handler
	msg, err := conduit.NewMessage("unknown_type", "test")
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	if err := json.NewEncoder(conn).Encode(msg); err != nil {
		t.Fatalf("Failed to send message to server: %v", err)
	}

	// There's no direct assertion here since we rely on logging.
	// In a more advanced setup, you could provide a custom logger to check for expected warnings.
}

// TestServerContext tests that setting and getting connection context on the server works.
func TestServerContext(t *testing.T) {
	socketPath := "/tmp/conduit_context_test.sock"
	defer os.RemoveAll(socketPath)

	cfg := conduit.DefaultServerConfig(socketPath)
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	s := server.NewServer(cfg)

	s.Handle("set_context", func(conn *server.Connection, msg *conduit.Message) error {
		conn.SetContext("foo", "bar")
		return conn.Send("ack", "context_set")
	})

	s.Handle("get_context", func(conn *server.Connection, msg *conduit.Message) error {
		val, ok := conn.GetContext("foo")
		if !ok {
			return conn.Send("ack", "no_value")
		}
		return conn.Send("ack", val)
	})

	if err := s.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer s.Stop()

	clientCfg := conduit.DefaultClientConfig(socketPath)
	clientCfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	c := client.NewClient(clientCfg)

	received := make(chan string, 2)
	c.Handle("ack", func(_ *client.Client, msg *conduit.Message) error {
		var resp string
		if err := msg.UnmarshalPayload(&resp); err != nil {
			return err
		}
		received <- resp
		return nil
	})

	if err := c.Connect(); err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	defer c.Close()

	if err := c.Send("set_context", nil); err != nil {
		t.Fatalf("Failed to send set_context message: %v", err)
	}

	select {
	case resp := <-received:
		if resp != "context_set" {
			t.Errorf("Expected 'context_set', got '%s'", resp)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for set_context response")
	}

	if err := c.Send("get_context", nil); err != nil {
		t.Fatalf("Failed to send get_context message: %v", err)
	}

	select {
	case resp := <-received:
		if resp != "bar" {
			t.Errorf("Expected 'bar', got '%s'", resp)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for get_context response")
	}
}
