# Signal Bot Quick Start Guide

This is a quick reference for getting your Signal bot up and running.

## Prerequisites Checklist

- [ ] Java 17+ installed
- [ ] signal-cli installed
- [ ] Go 1.21+ installed (or Docker)
- [ ] Phone number for the bot (can use Google Voice, Twilio, etc.)

## Setup Steps (5 minutes)

### 1. Initial Setup
```bash
# Clone and enter directory
git clone https://github.com/riyanimam/signal-bot-playground.git
cd signal-bot-playground

# Run setup script
./setup.sh
```

### 2. Register with Signal
```bash
# Register your number
signal-cli -a +1234567890 register

# Verify (you'll get a code via SMS)
signal-cli -a +1234567890 verify 123-456
```

### 3. Run the Bot

**Option A: Direct Run**
```bash
make build
make run
```

**Option B: Docker**
```bash
make docker-build
make docker-up
```

## Test Your Bot

Send a message to your bot's number:
```
!ping
```

You should receive:
```
🏓 Pong!
```

## Quick Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `!help` | Show help message | `!help` |
| `!ping` | Test if bot is alive | `!ping` |
| `!echo <text>` | Echo back message | `!echo Hello World` |
| `!about` | Bot information | `!about` |

## Common Issues

**"signal-cli not found"**
```bash
# Install on macOS
brew install signal-cli

# On Linux, see README.md
```

**"SIGNAL_PHONE_NUMBER is required"**
```bash
# Edit .env file and add your number
nano .env
# Set: SIGNAL_PHONE_NUMBER=+1234567890
```

**Bot not responding**
1. Check bot is running: `ps aux | grep signal-bot`
2. Check logs for errors
3. Verify phone number in .env matches registered number
4. Try restarting the bot

## Next Steps

1. **Add Custom Commands**: Edit `handler.go` to add your own commands
2. **Deploy to Production**: See README.md for deployment guides
3. **Monitor Logs**: Use `make docker-logs` or check terminal output

## Getting Help

- Full documentation: See [README.md](README.md)
- Issues: Open an issue on GitHub
- Signal-CLI docs: https://github.com/AsamK/signal-cli

---

Happy botting! 🤖
