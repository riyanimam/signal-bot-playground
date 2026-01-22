# Signal Bot Playground 🤖

A simple Signal messenger bot written in Go. This bot can respond to commands in both direct messages and group chats.

## Features

- ✅ Responds to commands in direct messages and group chats
- ✅ Easy to extend with custom commands
- ✅ Uses signal-cli for Signal protocol communication
- ✅ Docker support for easy deployment
- ✅ Environment-based configuration

## Prerequisites

Before you can run this bot, you need:

1. **Java Runtime Environment (JRE) 17+** - Required for signal-cli
2. **signal-cli** - Command-line interface for Signal
3. **Go 1.21+** - To build and run the bot
4. **A Signal account** - You'll need a phone number to register the bot

## Quick Start

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/riyanimam/signal-bot-playground.git
   cd signal-bot-playground
   ```

2. **Set up your environment**
   ```bash
   cp .env.example .env
   # Edit .env and add your Signal phone number
   ```

3. **Register your Signal account**
   ```bash
   docker-compose run --rm signal-bot signal-cli -a YOUR_PHONE_NUMBER register
   # You'll receive a verification code via SMS
   docker-compose run --rm signal-bot signal-cli -a YOUR_PHONE_NUMBER verify CODE
   ```

4. **Start the bot**
   ```bash
   docker-compose up -d
   ```

### Option 2: Manual Installation

#### Step 1: Install signal-cli

**On Linux:**
```bash
# Download and extract signal-cli
wget https://github.com/AsamK/signal-cli/releases/download/v0.13.1/signal-cli-0.13.1-Linux.tar.gz
tar xf signal-cli-0.13.1-Linux.tar.gz -C /opt
sudo ln -sf /opt/signal-cli-0.13.1/bin/signal-cli /usr/local/bin/
```

**On macOS:**
```bash
brew install signal-cli
```

**Verify installation:**
```bash
signal-cli --version
```

#### Step 2: Register your Signal number

```bash
# Register the phone number (you'll receive an SMS with a verification code)
signal-cli -a +1234567890 register

# Verify with the code you received
signal-cli -a +1234567890 verify 123-456
```

Replace `+1234567890` with your actual phone number in international format.

#### Step 3: Set up the bot

```bash
# Clone the repository
git clone https://github.com/riyanimam/signal-bot-playground.git
cd signal-bot-playground

# Copy environment template
cp .env.example .env

# Edit .env and set your phone number
# SIGNAL_PHONE_NUMBER=+1234567890
```

#### Step 4: Build and run

```bash
# Install dependencies
go mod download

# Build the bot
go build -o signal-bot

# Run the bot
./signal-bot
```

## What You Need to Provide

To get your Signal bot running, you need to:

1. **Phone Number**: A dedicated phone number for your bot (can be a virtual number from services like Twilio, Google Voice, etc.)
   - Must be in international format: `+1234567890`
   - This number will be registered with Signal and used exclusively for the bot

2. **Access to receive SMS**: During registration, you'll need to receive an SMS verification code

3. **Server/Computer**: A machine to run the bot
   - Can be your local machine for testing
   - A cloud server (AWS EC2, DigitalOcean, etc.) for production
   - Must be able to run continuously if you want 24/7 availability

## Usage

Once the bot is running, you can message it on Signal. The bot responds to these commands:

- `!help` - Show available commands
- `!ping` - Check if bot is alive
- `!echo <text>` - Echo back your message
- `!about` - Get information about the bot

**Example:**
```
You: !ping
Bot: 🏓 Pong!

You: !echo Hello World
Bot: Hello World
```

## Configuration

The bot is configured via environment variables. See `.env.example` for all options:

| Variable | Description | Default |
|----------|-------------|---------|
| `SIGNAL_PHONE_NUMBER` | Your bot's Signal phone number | *Required* |
| `SIGNAL_DATA_DIR` | Directory for Signal data | `./signal-data` |
| `BOT_COMMAND_PREFIX` | Prefix for bot commands | `!` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

## Project Structure

```
signal-bot-playground/
├── main.go           # Main application entry point
├── config.go         # Configuration management
├── handler.go        # Message handling and command logic
├── go.mod           # Go module definition
├── .env.example     # Environment variable template
├── .gitignore       # Git ignore rules
├── Dockerfile       # Docker image definition
├── docker-compose.yml # Docker Compose configuration
└── README.md        # This file
```

## Adding Custom Commands

To add your own commands, edit `handler.go`:

```go
// In the HandleMessage function, add a new case:
case "mycommand":
    return h.handleMyCommand(args), nil

// Then implement the handler:
func (h *MessageHandler) handleMyCommand(args []string) string {
    return "My custom response!"
}
```

## Troubleshooting

### signal-cli not found
Make sure signal-cli is installed and in your PATH:
```bash
which signal-cli
signal-cli --version
```

### Registration issues
- Ensure your phone number is in international format (+1234567890)
- Check that you can receive SMS on that number
- If using a VOIP number, some may not work with Signal

### Bot not responding
- Check that the bot is running: `docker-compose ps` or `ps aux | grep signal-bot`
- Check logs: `docker-compose logs -f` or look at terminal output
- Verify your phone number in `.env` matches the registered number
- Make sure you're using the correct command prefix (default is `!`)

### Connection issues
- Signal-cli needs internet connectivity
- Check firewall settings
- Ensure signal-cli data directory is writable

## Production Deployment

### Using systemd (Linux)

Create `/etc/systemd/system/signal-bot.service`:

```ini
[Unit]
Description=Signal Bot
After=network.target

[Service]
Type=simple
User=signalbot
WorkingDirectory=/opt/signal-bot-playground
ExecStart=/opt/signal-bot-playground/signal-bot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable signal-bot
sudo systemctl start signal-bot
sudo systemctl status signal-bot
```

### Using Docker in production

```bash
# Build and start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the bot
docker-compose down
```

## Security Considerations

1. **Keep your phone number private** - Anyone can message your bot
2. **Rate limiting** - The bot includes basic rate limiting, but consider adding more robust protections
3. **Input validation** - Always validate user input in command handlers
4. **Environment variables** - Never commit `.env` file to version control
5. **Updates** - Keep signal-cli and dependencies up to date

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Resources

- [signal-cli GitHub](https://github.com/AsamK/signal-cli)
- [Signal Protocol Documentation](https://signal.org/docs/)
- [Go Documentation](https://golang.org/doc/)

## Support

If you encounter any issues or have questions:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review signal-cli documentation
3. Open an issue on GitHub

---

Made with ❤️ using Go and Signal