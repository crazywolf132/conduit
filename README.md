# Conduit 🚀

> "Because life's too short for complicated socket programming!" 

**Effortless, High-Performance Unix Socket Communication for Go**

Ever tried to get two processes to talk to each other and ended up feeling like you're herding cats? 🐱 We've been there! That's why we created Conduit – your friendly neighborhood Unix socket library that makes inter-process communication as easy as sending a text message. 

## 🤔 The Problem

Picture this: You've got multiple Go services running on the same machine, all trying to have a nice chat with each other. But setting up TCP servers feels like organizing a dinner party where everyone speaks a different language. And HTTP? That's like sending a letter through the post office when your friend lives next door! 

## 💡 The Solution

Enter Conduit! It's like having a magical pipe that connects your processes, but without the plumbing nightmares. Using Unix domain sockets (fancy term for "really fast local communication"), Conduit lets your programs chat with each other at lightning speed. ⚡

## ✨ Why Developers Love Conduit

- **🎯 Simple but Powerful:** Like a Swiss Army knife that you can actually figure out how to use
- **🔄 Bi-Directional Chat:** Programs can gossip back and forth, just like your favorite messaging app
- **🔒 Secure & Speedy:** Faster than TCP, safer than leaving your front door open
- **🛠️ Production-Ready:** Built like a tank, but drives like a Tesla
- **📝 Type-Safe:** Because nobody likes runtime surprises (except on their birthday)
- **😊 Developer Friendly:** Documentation that doesn't read like ancient hieroglyphics

## 🚀 Quick Start

### Installation
```bash
go get github.com/crazywolf132/conduit   # Your journey begins here!
```

### The "Hello World" of Socket Communication

#### 🎭 Server Side (The Host of the Party)
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

#### 👋 Client Side (The Party Guest)
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

## 🎨 Best Practices & Patterns

### 🎯 Message Types (Keep it Organized!)
```go
// Your message types, as organized as Marie Kondo's closet
const (
    MsgTypeHeartbeat = "heartbeat"  // The "you still there?" message
    MsgTypeCommand   = "command"    // The "please do this" message
    MsgTypeResponse  = "response"   // The "here's what happened" message
    MsgTypeEvent     = "event"      // The "guess what just happened!" message
)
```

### 🎭 Error Handling (Because Things Happen)
```go
// Client-side error handling (with style!)
if err := client.Send("command", cmd); err != nil {
    switch {
    case errors.Is(err, conduit.ErrConnectionClosed):
        // Connection went on vacation 🏖️
    case errors.Is(err, conduit.ErrTimeout):
        // Time moves slower than your grandmother's internet
    default:
        // Something unexpected happened 🤷
    }
}
```

## 🚦 Production Tips

### 📊 Monitoring (Keep Your Eyes on the Prize)
- Set up health checks (like a regular doctor's appointment for your app)
- Use logging (because print statements are so 2010)
- Watch those connections (they're like teenagers - need constant monitoring)

### 🔐 Security (Better Safe Than Sorry)
- Lock down those socket files (no party crashers allowed!)
- Validate your messages (trust no one, not even yourself)
- Keep your secrets secret (obvious, but you'd be surprised...)

## 🤔 FAQ (The "I'm Glad You Asked" Section)

### 🏃‍♂️ Performance

**Q: Is it fast?**
A: Faster than your coffee break! We're talking 2-3x quicker than TCP/IP for local chats. Sub-millisecond latency that would make The Flash jealous.

**Q: Can it handle big messages?**
A: Up to 10MB by default - that's like 5,000 pages of text! Need more? Just adjust the `MaxMessageSize`. But remember, bigger isn't always better!

### 🏢 Enterprise Stuff

**Q: Can I trust it in production?**
A: As reliable as your favorite coffee machine! Used by companies processing millions of messages daily without breaking a sweat.

**Q: What about support?**
A: We've got your back! While it's open-source (free as in both beer AND speech 🍺), we're happy to discuss enterprise support if you need the VIP treatment.

## 🤝 Join the Fun!

- 📚 [Docs](https://godoc.org/github.com/crazywolf132/conduit) (Actually readable!)
- 🐛 [Issues](https://github.com/crazywolf132/conduit/issues) (Found a bug? Let's squash it!)

## 📜 License

MIT License - Free as a bird! Use it, break it, fix it, share it... just don't blame us if your cat videos app goes viral! 

## 🤝 Contributing

Got ideas? We love ideas! Whether you're fixing typos or adding killer features, we're here for it. Join our merry band of socket enthusiasts and help make Conduit even more awesome! 

---

Made with ❤️ by developers who got tired of complicated socket programming.

Remember: Life's too short for bad APIs! 🌟
