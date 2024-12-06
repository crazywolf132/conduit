package test

import (
	"os"
	"testing"
	"time"

	"github.com/crazywolf132/conduit"
	"github.com/crazywolf132/conduit/client"
	"github.com/crazywolf132/conduit/server"
)

// TestClientServerInteraction tests a simple interaction between client and server.
func TestClientServerInteraction(t *testing.T) {
	socketPath := "/tmp/conduit_test.sock"
	defer os.RemoveAll(socketPath)

	serverCfg := conduit.DefaultServerConfig(socketPath)
	serverCfg.Logger = conduit.NewLogger(conduit.LogError, nil) // silent for test
	srv := server.NewServer(serverCfg)

	// server handler echoes back messages of type "echo"
	srv.Handle("echo", func(conn *server.Connection, msg *conduit.Message) error {
		var payload string
		if err := msg.UnmarshalPayload(&payload); err != nil {
			return err
		}
		return conn.Send("echo_response", payload+"_response")
	})

	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	clientCfg := conduit.DefaultClientConfig(socketPath)
	clientCfg.Logger = conduit.NewLogger(conduit.LogError, nil) // silent for test
	c := client.NewClient(clientCfg)

	received := make(chan string, 1)
	c.Handle("echo_response", func(_ *client.Client, msg *conduit.Message) error {
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

	// Send a message
	if err := c.Send("echo", "hello"); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	select {
	case resp := <-received:
		if resp != "hello_response" {
			t.Errorf("Expected 'hello_response', got '%s'", resp)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for response")
	}
}

// TestClientNotConnected ensures sending before connecting returns an error.
func TestClientNotConnected(t *testing.T) {
	cfg := conduit.DefaultClientConfig("/tmp/does_not_exist.sock")
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	c := client.NewClient(cfg)

	if err := c.Send("test", "payload"); err == nil {
		t.Error("Expected error when sending on not connected client, got none")
	}
}

// TestClientReconnect tests that the client reconnects automatically if Reconnect is enabled.
func TestClientReconnect(t *testing.T) {
	socketPath := "/tmp/conduit_reconnect_test.sock"
	defer os.RemoveAll(socketPath)

	serverCfg := conduit.DefaultServerConfig(socketPath)
	serverCfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	srv := server.NewServer(serverCfg)

	srv.Handle("test", func(conn *server.Connection, msg *conduit.Message) error {
		return conn.Send("test_response", "ok")
	})

	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	clientCfg := conduit.DefaultClientConfig(socketPath)
	clientCfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	clientCfg.Reconnect = true
	clientCfg.ReconnectDelay = 100 * time.Millisecond
	c := client.NewClient(clientCfg)

	received := make(chan string, 1)
	c.Handle("test_response", func(_ *client.Client, msg *conduit.Message) error {
		var resp string
		if err := msg.UnmarshalPayload(&resp); err != nil {
			return err
		}
		received <- resp
		return nil
	})

	// Connect the client
	if err := c.Connect(); err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}

	// Close the server to force a client disconnect
	if err := srv.Stop(); err != nil {
		t.Fatalf("Failed to stop server: %v", err)
	}

	// Restart the server to simulate server coming back online
	srv = server.NewServer(serverCfg)
	srv.Handle("test", func(conn *server.Connection, msg *conduit.Message) error {
		return conn.Send("test_response", "ok_reconnected")
	})
	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server again: %v", err)
	}
	defer srv.Stop()

	// The client should reconnect automatically
	time.Sleep(500 * time.Millisecond) // wait for reconnection attempts

	// Send a message again
	if err := c.Send("test", "ping"); err != nil {
		t.Fatalf("Failed to send message after reconnection: %v", err)
	}

	select {
	case resp := <-received:
		if resp != "ok_reconnected" {
			t.Errorf("Expected 'ok_reconnected', got '%s'", resp)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for response after reconnection")
	}

	c.Close()
}

// TestClientContext tests that setting context on the client works as expected.
func TestClientContext(t *testing.T) {
	cfg := conduit.DefaultClientConfig("/tmp/does_not_exist.sock")
	cfg.Logger = conduit.NewLogger(conduit.LogError, nil)
	c := client.NewClient(cfg)

	c.SetContext("user_id", 42)
	val, ok := c.GetContext("user_id")
	if !ok {
		t.Error("Expected to get a value from context, got none")
	}
	if val != 42 {
		t.Errorf("Expected context value 42, got %v", val)
	}
}
