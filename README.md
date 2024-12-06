# Conduit

**Effortless, High-Performance Unix Socket Communication for Go**
In a world of networked applications, sometimes you just want a faster, simpler way for local processes to talk to each other—without the overhead or complexity of TCP or HTTP. Enter Unix domain sockets: a powerful, low-latency mechanism built right into your operating system. But let's face it: working directly with raw sockets can feel like turning a Rubik’s Cube blindfolded.

**Conduit** is here to make your life easier. It transforms the often fiddly world of Unix sockets into a clean, developer-friendly, and production-grade experience. Need to build a local microservice architecture? A high-performance IPC channel for your modules? A robust internal messaging system? Conduit has your back.

## Why Conduit?
- **Simplicity Meets Power:** Straightforward APIs let you send and receive JSON-encoded messages without wrestling with byte streams.
- **Bi-Directional Communication:** Clients and servers talk both ways, naturally. Whether sending commands downstream or broadcasting updates upstream, Conduit keeps it easy.
- **Secure and Fast:** Unix domain sockets provide a secure, high-speed channel right on your filesystem—no random ports, no internet exposure.
- **Production-Ready:** With automatic reconnection, configurable timeouts, robust error handling, and built-in logging, Conduit is ready for real-world deployments.
- **Type-Safe Handlers:** Register handlers by message type. Decode your payloads directly into Go structs. No more manual parsing errors or messy switch statements.
- **Fun and Friendly:** Developer experience matters. Conduit’s code is well-documented, comprehensively tested, and easy to learn, so you can focus on building great software.

## Installation

```bash
go get github.com/crazywolf132/conduit
```

## Quick Peek

### The Problem

You’ve got multiple Go services on the same machine. They need to chatter—send commands, events, logs—back and forth, but setting up TCP servers or REST endpoints feels heavy and complicated. **You just need a simple, speedy, and reliable channel.**

### The Solution

Conduit leverages Unix domain sockets to give you low-latency, file-based endpoints. With just a few lines of code, you’ve got a robust messaging layer that "just works," letting you focus on building features, not plumbing.

## Tiny Example

### Server (Easy as Pie)

```go
package main

import (
    "fmt"
    "github.com/crazywolf132/conduit"
    "github.com/crazywolf132/conduit/server"
)

func main() {
    cfg := conduit.DefaultServerConfig("/tmp/app.sock")
    s := server.NewServer(cfg)

    // Define a handler for the "greeting" message type
    s.Handle("greeting", func(conn *server.Connection, msg *conduit.Message) error {
        var greeting string
        if err := msg.UnmarshalPayload(&greeting); err != nil {
            return err
        }

        // Echo it back with some pizzazz
        return conn.Send("greeting", fmt.Sprintf("Server says hi back: %s", greeting))
    })

    // Start the server
    if err := s.Start(); err != nil {
        panic(err)
    }
    defer s.Stop()

    // Block until interrupted
    select {}
}
```

### Client (No Muss, No Fuss)

```go
package main

import (
    "fmt"
    "github.com/crazywolf132/conduit"
    "github.com/crazywolf132/conduit/client"
)

func main() {
    cfg := conduit.DefaultClientConfig("/tmp/app.sock")
    c := client.NewClient(cfg)

    // Print responses from the server
    c.Handle("greeting", func(_ *client.Client, msg *conduit.Message) error {
        var response string
        if err := msg.UnmarshalPayload(&response); err != nil {
            return err
        }
        fmt.Println("Received:", response)
        return nil
    })

    // Connect
    if err := c.Connect(); err != nil {
        panic(err)
    }
    defer c.Close()

    // Say hello
    if err := c.Send("greeting", "Hello, Conduit!"); err != nil {
        panic(err)
    }
}

```

## Real-World Example: Chat App

Check out `examples/chat` for a full-featured demo—imagine multiple local clients chatting through a single server socket. Launch the server and fire up two clients (say, Alice and Bob). Everyone sees everyone else’s messages instantly. No complicated network setup, no extra frameworks—just good old-fashioned sockets, made convenient.

```bash
# Start the server
go run examples/chat/main.go server

# In another terminal, start a client named Alice
go run examples/chat/main.go client Alice

# In a third terminal, start a client named Bob
go run examples/chat/main.go client Bob

```
Now Alice and Bob can chat freely. Easy, right?

## Key Configuration & Extensions

- **Timeouts & Reconnects:** Need to handle flaky upstream services? Conduit’s client can automatically reconnect and retry, with configurable backoff.
- **Custom Loggers:** Bring your own logger or integrate with your favorite logging framework to keep an eye on the system.
- **Message Size Limits:** For safety and performance, set maximum allowed message sizes and rest easy knowing you won’t be flooded with giant payloads.

## Future-Proof & Tested

Conduit follows semver for predictable versioning. It’s backed by extensive unit tests and integration tests. We invite contributions—improve handlers, add new examples, help us refine the code. Your input helps Conduit stay robust and developer-friendly.


## License
MIT License

Conduit is free and open-source. Use it, modify it, break it, rebuild it—just respect the license, and share your improvements if you’d like.

## Contributing

All contributions are welcome. Submit a pull request, open an issue, or suggest new features. Together, we can make Conduit the go-to tool for elegant Unix socket communication in Go.

