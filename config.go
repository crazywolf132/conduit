package conduit

import "time"

// ServerConfig holds configuration options for the server.
//
// Fields:
//   - SocketPath: Filesystem path to the Unix domain socket.
//   - SocketPermissions: Filesystem permissions for the socket file.
//   - Logger: A Logger interface for outputting server logs. Defaults to a basic logger if not set.
//   - ReadTimeout: Maximum duration for reading a single message from a client.
//   - WriteTimeout: Maximum duration for writing a single message to a client.
//   - MaxMessageSize: Maximum allowed size of a single message in bytes.
type ServerConfig struct {
	SocketPath        string
	SocketPermissions uint32
	Logger            Logger
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxMessageSize    int64
}

// DefaultServerConfig returns a ServerConfig with standard default values.
//
// Example usage:
//
//	cfg := conduit.DefaultServerConfig("/tmp/app.sock")
//	s := server.NewServer(cfg)
func DefaultServerConfig(socketPath string) *ServerConfig {
	return &ServerConfig{
		SocketPath:        socketPath,
		SocketPermissions: 0666,
		Logger:            NewLogger(LogInfo, nil),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		MaxMessageSize:    32 * 1024 * 1024, // 32MB default
	}
}

// ClientConfig holds configuration options for the client.
//
// Fields:
//   - SocketPath: Filesystem path to the Unix domain socket the client connects to.
//   - Logger: A Logger interface for outputting client logs. Defaults to a basic logger if not set.
//   - ReadTimeout: Maximum duration for reading a single message from the server.
//   - WriteTimeout: Maximum duration for writing a single message to the server.
//   - MaxMessageSize: Maximum allowed size of a single message in bytes.
//   - Reconnect: If true, the client will attempt to reconnect on connection loss.
//   - ReconnectDelay: Delay between reconnection attempts if Reconnect is true.
type ClientConfig struct {
	SocketPath     string
	Logger         Logger
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxMessageSize int64
	Reconnect      bool
	ReconnectDelay time.Duration
}

// DefaultClientConfig returns a ClientConfig with standard default values.
//
// Example usage:
//
//	cfg := conduit.DefaultClientConfig("/tmp/app.sock")
//	c := client.NewClient(cfg)
//	if err := c.Connect(); err != nil { ... }
func DefaultClientConfig(socketPath string) *ClientConfig {
	return &ClientConfig{
		SocketPath:     socketPath,
		Logger:         NewLogger(LogInfo, nil),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxMessageSize: 32 * 1024 * 1024, // 32MB default
		Reconnect:      true,
		ReconnectDelay: 5 * time.Second,
	}
}
