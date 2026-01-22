# Signal Bot Architecture

## Overview

This document describes the architecture and flow of the Signal Bot.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Signal Messenger                         │
│                  (Users send messages)                       │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ Signal Protocol
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                      signal-cli                              │
│              (Java-based Signal client)                      │
│  - Handles Signal protocol communication                     │
│  - Receives messages in JSON format                          │
│  - Sends messages on behalf of the bot                       │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ JSON output (stdout)
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                  Signal Bot (Go)                             │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ main.go - Main application loop                      │   │
│  │  - Spawns signal-cli process                         │   │
│  │  - Reads JSON messages from stdout                   │   │
│  │  - Handles graceful shutdown                         │   │
│  └──────────────┬───────────────────────────────────────┘   │
│                 │                                            │
│                 ↓                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ handler.go - Message processing                      │   │
│  │  - Parses incoming messages                          │   │
│  │  - Routes commands to handlers                       │   │
│  │  - Generates responses                               │   │
│  └──────────────┬───────────────────────────────────────┘   │
│                 │                                            │
│                 ↓                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ config.go - Configuration management                 │   │
│  │  - Loads environment variables                       │   │
│  │  - Validates configuration                           │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Message Flow

### Receiving Messages

```
1. User sends message to bot's Signal number
   ↓
2. Signal servers receive and forward message
   ↓
3. signal-cli receives message via Signal protocol
   ↓
4. signal-cli outputs message as JSON to stdout
   ↓
5. Go bot reads JSON from signal-cli's stdout
   ↓
6. Go bot parses JSON into Message struct
   ↓
7. MessageHandler checks for command prefix
   ↓
8. If command: route to appropriate handler
   ↓
9. Handler generates response string
```

### Sending Messages

```
1. Handler returns response string
   ↓
2. main.go receives response
   ↓
3. sendMessage() function called with response
   ↓
4. Spawns signal-cli process with 'send' command
   ↓
5. signal-cli sends message via Signal protocol
   ↓
6. Signal servers deliver message to user
   ↓
7. User receives bot's response
```

## Component Details

### main.go
**Responsibilities:**
- Application initialization
- Spawning and managing signal-cli process
- Reading messages from signal-cli stdout
- Coordinating message handling
- Graceful shutdown handling

**Key Functions:**
- `main()` - Entry point, orchestrates the bot
- `sendMessage()` - Sends messages via signal-cli

### handler.go
**Responsibilities:**
- Message parsing and validation
- Command routing
- Response generation
- Logging

**Key Types:**
- `Message` - Represents a Signal message
- `MessageHandler` - Processes messages

**Key Functions:**
- `HandleMessage()` - Main message processing function
- `handleHelp()`, `handlePing()`, etc. - Individual command handlers

### config.go
**Responsibilities:**
- Loading configuration from environment
- Configuration validation
- Providing defaults

**Key Types:**
- `Config` - Configuration structure

**Key Functions:**
- `LoadConfig()` - Loads and validates configuration
- `getEnv()` - Helper for environment variables

## Data Structures

### Message
```go
type Message struct {
    Sender    string  // Phone number of sender
    Text      string  // Message content
    Timestamp int64   // Unix timestamp
    GroupID   string  // Group ID (if group message)
}
```

### Config
```go
type Config struct {
    PhoneNumber   string  // Bot's Signal number
    DataDir       string  // signal-cli data directory
    CommandPrefix string  // Command prefix (e.g., "!")
    LogLevel      string  // Logging level
}
```

### SignalMessage (from signal-cli JSON)
```go
type SignalMessage struct {
    Envelope struct {
        Source       string
        SourceNumber string
        Timestamp    int64
        DataMessage  *struct {
            Message string
            GroupInfo *struct {
                GroupID string
            }
        }
    }
}
```

## Extension Points

### Adding New Commands

1. Add a new case in `HandleMessage()` switch statement
2. Implement handler function (e.g., `handleMyCommand()`)
3. Update help text in `handleHelp()`

Example:
```go
// In handler.go, HandleMessage():
case "weather":
    return h.handleWeather(args), nil

// Add new handler:
func (h *MessageHandler) handleWeather(args []string) string {
    // Your weather logic here
    return "Weather info..."
}
```

### Adding Middleware

You can add middleware for:
- Rate limiting
- Authentication
- Logging
- Analytics

Add middleware checks before command routing in `HandleMessage()`.

### Adding State Management

For stateful conversations:
1. Create a session manager
2. Store session data (in-memory map or database)
3. Check session state in message handler
4. Update state based on user responses

## Security Considerations

1. **Input Validation**: All user input is treated as untrusted
2. **Command Injection**: No user input is passed to shell commands
3. **Rate Limiting**: Small delay between messages (100ms)
4. **Error Handling**: Errors are logged but don't expose internals

## Performance Characteristics

- **Message Processing**: ~10-50ms per message
- **Startup Time**: ~1-2 seconds
- **Memory Usage**: ~20-40MB (Go process)
- **CPU Usage**: Minimal when idle, spike during message processing

## Dependencies

### Runtime Dependencies
- signal-cli (Java application)
- Java Runtime Environment 17+

### Go Dependencies
- github.com/joho/godotenv (environment variable loading)

## Deployment Options

1. **Standalone**: Run directly on a server
2. **Docker**: Containerized deployment
3. **systemd**: System service on Linux
4. **Cloud**: AWS, DigitalOcean, etc.

See README.md for detailed deployment instructions.

## Monitoring and Logging

The bot logs:
- Startup messages
- Configuration loaded
- Incoming messages (sender, text)
- Outgoing responses
- Errors and warnings

Use `LOG_LEVEL` environment variable to control verbosity:
- `debug`: Detailed information
- `info`: General information (default)
- `warn`: Warnings only
- `error`: Errors only

## Future Enhancements

Potential improvements:
- Database integration for persistent state
- Web dashboard for monitoring
- Plugin system for extensions
- Multi-language support
- Rich media handling (images, files)
- Advanced message formatting
- Group management commands
- Admin authentication

---

This architecture provides a solid foundation for a Signal bot while remaining simple and maintainable.
