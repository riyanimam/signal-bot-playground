package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// SignalMessage represents a message from signal-cli JSON output
type SignalMessage struct {
	Envelope struct {
		Source       string `json:"source"`
		SourceNumber string `json:"sourceNumber"`
		SourceUUID   string `json:"sourceUuid"`
		Timestamp    int64  `json:"timestamp"`
		DataMessage  *struct {
			Timestamp int64  `json:"timestamp"`
			Message   string `json:"message"`
			GroupInfo *struct {
				GroupID string `json:"groupId"`
			} `json:"groupInfo"`
		} `json:"dataMessage"`
		SyncMessage *struct {
			SentMessage *struct {
				Timestamp   int64  `json:"timestamp"`
				Message     string `json:"message"`
				Destination string `json:"destination"`
			} `json:"sentMessage"`
		} `json:"syncMessage"`
	} `json:"envelope"`
}

func main() {
	log.Println("Starting Signal Bot...")

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Bot configured for number: %s", maskPhoneNumber(config.PhoneNumber))
	log.Printf("Command prefix: %s", config.CommandPrefix)

	// Create message handler
	handler := NewMessageHandler(config)

	// Create signal-cli command to receive messages in JSON mode
	cmd := exec.Command("signal-cli",
		"-a", config.PhoneNumber,
		"--config", config.DataDir,
		"receive",
		"--json",
	)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start signal-cli: %v", err)
	}

	log.Println("Bot is now listening for messages...")
	log.Println("Press Ctrl+C to stop")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nShutting down bot gracefully...")
		if cmd.Process != nil {
			// Try graceful shutdown first
			cmd.Process.Signal(os.Interrupt)
			// Wait a bit for graceful shutdown
			time.Sleep(2 * time.Second)
			// Force kill if still running
			cmd.Process.Kill()
		}
		os.Exit(0)
	}()

	// Read messages
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse JSON message
		var signalMsg SignalMessage
		if err := json.Unmarshal([]byte(line), &signalMsg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// Process data messages (incoming messages)
		if signalMsg.Envelope.DataMessage != nil {
			msg := &Message{
				Sender:    signalMsg.Envelope.SourceNumber,
				Text:      signalMsg.Envelope.DataMessage.Message,
				Timestamp: signalMsg.Envelope.DataMessage.Timestamp,
			}

			if signalMsg.Envelope.DataMessage.GroupInfo != nil {
				msg.GroupID = signalMsg.Envelope.DataMessage.GroupInfo.GroupID
			}

			// Log the message
			handler.LogMessage(msg)

			// Handle the message
			response, err := handler.HandleMessage(msg)
			if err != nil {
				log.Printf("Error handling message: %v", err)
				continue
			}

			// Send response if there is one
			if response != "" {
				if err := sendMessage(config, msg, response); err != nil {
					log.Printf("Failed to send response: %v", err)
				} else {
					log.Printf("Sent response: %s", response)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from signal-cli: %v", err)
	}

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("signal-cli exited with error: %v", err)
	}
}

// sendMessage sends a message via signal-cli
func sendMessage(config *Config, originalMsg *Message, text string) error {
	var cmd *exec.Cmd

	if originalMsg.GroupID != "" {
		// Validate group ID (should be base64-like string)
		if !isValidIdentifier(originalMsg.GroupID) {
			return fmt.Errorf("invalid group ID format")
		}
		// Send to group
		cmd = exec.Command("signal-cli",
			"-a", config.PhoneNumber,
			"--config", config.DataDir,
			"send",
			"-g", originalMsg.GroupID,
			"-m", text,
		)
	} else {
		// Validate sender phone number
		if !isValidPhoneNumber(originalMsg.Sender) {
			return fmt.Errorf("invalid sender phone number format")
		}
		// Send direct message
		cmd = exec.Command("signal-cli",
			"-a", config.PhoneNumber,
			"--config", config.DataDir,
			"send",
			originalMsg.Sender,
			"-m", text,
		)
	}

	// Add a small delay to avoid rate limiting
	time.Sleep(100 * time.Millisecond)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("signal-cli send failed: %v, output: %s", err, output)
	}

	return nil
}

// isValidPhoneNumber validates that a phone number is in expected format
func isValidPhoneNumber(phone string) bool {
	if len(phone) < 8 || len(phone) > 20 {
		return false
	}
	// Must start with + and contain only digits
	if phone[0] != '+' {
		return false
	}
	for _, c := range phone[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// isValidIdentifier validates identifiers (group IDs, UUIDs, etc.)
func isValidIdentifier(id string) bool {
	if len(id) == 0 || len(id) > 100 {
		return false
	}
	// Allow alphanumeric, dash, underscore, equals (base64)
	for _, c := range id {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '=' || c == '+' || c == '/') {
			return false
		}
	}
	return true
}

// maskPhoneNumber masks a phone number for logging (shows only last 4 digits)
func maskPhoneNumber(phone string) string {
	if len(phone) <= 4 {
		return "****"
	}
	return "****" + phone[len(phone)-4:]
}
