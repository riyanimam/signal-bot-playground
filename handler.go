package main

import (
	"fmt"
	"log"
	"strings"
)

// Message represents a Signal message
type Message struct {
	Sender    string
	Text      string
	Timestamp int64
	GroupID   string
}

// MessageHandler processes incoming messages
type MessageHandler struct {
	config *Config
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(config *Config) *MessageHandler {
	return &MessageHandler{
		config: config,
	}
}

// HandleMessage processes an incoming message and returns a response
func (h *MessageHandler) HandleMessage(msg *Message) (string, error) {
	// Check if message starts with command prefix
	if !strings.HasPrefix(msg.Text, h.config.CommandPrefix) {
		return "", nil // Not a command, ignore
	}

	// Remove prefix and parse command
	commandText := strings.TrimPrefix(msg.Text, h.config.CommandPrefix)
	parts := strings.Fields(commandText)

	if len(parts) == 0 {
		return "", nil
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Route to appropriate handler
	switch command {
	case "help":
		return h.handleHelp(), nil
	case "ping":
		return h.handlePing(), nil
	case "echo":
		return h.handleEcho(args), nil
	case "about":
		return h.handleAbout(), nil
	default:
		return h.handleUnknownCommand(command), nil
	}
}

// handleHelp returns help text
func (h *MessageHandler) handleHelp() string {
	return fmt.Sprintf(`Available commands:
%shelp - Show this help message
%sping - Check if bot is alive
%secho <text> - Echo back your message
%sabout - Information about this bot`,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
	)
}

// handlePing returns a pong response
func (h *MessageHandler) handlePing() string {
	return "🏓 Pong!"
}

// handleEcho echoes back the user's message
func (h *MessageHandler) handleEcho(args []string) string {
	if len(args) == 0 {
		return "Please provide a message to echo!"
	}
	return strings.Join(args, " ")
}

// handleAbout returns information about the bot
func (h *MessageHandler) handleAbout() string {
	return `Signal Bot v1.0
A simple Signal bot written in Go
Repository: github.com/riyanimam/signal-bot-playground`
}

// handleUnknownCommand handles unknown commands
func (h *MessageHandler) handleUnknownCommand(command string) string {
	return fmt.Sprintf("Unknown command: %s\nType %shelp for available commands",
		command, h.config.CommandPrefix)
}

// LogMessage logs an incoming message with phone number masking for privacy
func (h *MessageHandler) LogMessage(msg *Message) {
	maskedSender := maskPhoneNumber(msg.Sender)
	if msg.GroupID != "" {
		log.Printf("[Group: %s] %s: %s", msg.GroupID, maskedSender, msg.Text)
	} else {
		log.Printf("[Direct] %s: %s", maskedSender, msg.Text)
	}
}
