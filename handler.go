package main

import (
	"fmt"
	"log"
	"strings"
	"time"
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
	config              *Config
	pollManager         *PollManager
	reminderManager     *ReminderManager
	statsManager        *StatsManager
	announcementManager *AnnouncementManager
	mentionTracker      *MentionTracker
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(config *Config, sendFunc func(string, string, string) error) *MessageHandler {
	return &MessageHandler{
		config:              config,
		pollManager:         NewPollManager(),
		reminderManager:     NewReminderManager(sendFunc),
		statsManager:        NewStatsManager(),
		announcementManager: NewAnnouncementManager(10),
		mentionTracker:      NewMentionTracker(),
	}
}

// HandleMessage processes an incoming message and returns a response
func (h *MessageHandler) HandleMessage(msg *Message) (string, error) {
	// Record message in stats if it's a group message
	if msg.GroupID != "" {
		h.statsManager.RecordMessage(msg.GroupID, msg.Sender)
	}

	// Check for mentions
	mentions := ParseMentions(msg.Text)
	for _, mentioned := range mentions {
		h.mentionTracker.TrackMention(msg.Sender, mentioned, msg.Text, msg.GroupID)
	}

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
	
	// Poll commands
	case "poll":
		return h.handlePoll(args, msg), nil
	case "vote":
		return h.handleVote(args, msg), nil
	case "pollresults":
		return h.handlePollResults(args), nil
	case "polls":
		return h.handleListPolls(msg), nil
	
	// Reminder commands
	case "remind", "reminder":
		return h.handleReminder(args, msg), nil
	case "reminders":
		return h.handleListReminders(msg), nil
	case "cancelreminder":
		return h.handleCancelReminder(args, msg), nil
	
	// Stats commands
	case "stats":
		return h.handleStats(msg), nil
	case "mystats":
		return h.handleMyStats(msg), nil
	
	// Announcement commands
	case "announce", "pin":
		return h.handleAnnounce(args, msg), nil
	case "announcements", "pins":
		return h.handleListAnnouncements(msg), nil
	
	// Mention commands
	case "mentions":
		return h.handleMentions(msg), nil
	case "clearmentions":
		return h.handleClearMentions(msg), nil
	
	default:
		return h.handleUnknownCommand(command), nil
	}
}

// handleHelp returns help text
func (h *MessageHandler) handleHelp() string {
	return fmt.Sprintf(`Available commands:

📝 Basic Commands:
%shelp - Show this help message
%sping - Check if bot is alive
%secho <text> - Echo back your message
%sabout - Information about this bot

📊 Polls & Voting (Large Groups):
%spoll <question> | <option1> | <option2> ... - Create a poll
%svote <poll_id> <option> - Vote in a poll
%spollresults <poll_id> - View poll results
%spolls - List active polls

⏰ Reminders:
%sremind <time> <message> - Set a reminder (e.g., 1h, 30m, 1d)
%sreminders - List your reminders
%scancelreminder <id> - Cancel a reminder

📊 Group Statistics:
%sstats - View group statistics
%smystats - View your activity stats

📌 Announcements:
%sannounce <message> - Pin an announcement
%sannouncements - View pinned announcements

💬 Mentions:
%smentions - View your mentions (use @+1234567890 to mention)
%sclearmentions - Clear your mentions

Use these features to manage large groups effectively!`,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
		h.config.CommandPrefix,
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

// Poll command handlers
func (h *MessageHandler) handlePoll(args []string, msg *Message) string {
	if msg.GroupID == "" {
		return "Polls can only be created in group chats."
	}

	if len(args) == 0 {
		return "Usage: !poll <question> | <option1> | <option2> | ...\nExample: !poll Where to eat? | Pizza | Sushi | Burgers"
	}

	// Parse poll data from args
	fullText := strings.Join(args, " ")
	parts := strings.Split(fullText, "|")

	if len(parts) < 3 {
		return "A poll needs a question and at least 2 options. Use | to separate them."
	}

	question := strings.TrimSpace(parts[0])
	options := make([]string, 0, len(parts)-1)

	for i := 1; i < len(parts); i++ {
		option := strings.TrimSpace(parts[i])
		if option != "" {
			options = append(options, option)
		}
	}

	if len(options) < 2 {
		return "Please provide at least 2 valid options."
	}

	pollID := h.pollManager.CreatePoll(question, options, msg.Sender, msg.GroupID, 0)

	result := fmt.Sprintf("📊 Poll created! ID: %s\n\n%s\n\nOptions:\n", pollID, question)
	for i, opt := range options {
		result += fmt.Sprintf("%d. %s\n", i+1, opt)
	}
	result += fmt.Sprintf("\nVote with: %svote %s <option>", h.config.CommandPrefix, pollID)

	return result
}

func (h *MessageHandler) handleVote(args []string, msg *Message) string {
	if len(args) < 2 {
		return "Usage: !vote <poll_id> <option>"
	}

	pollID := args[0]
	option := strings.Join(args[1:], " ")

	err := h.pollManager.Vote(pollID, msg.Sender, option)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return "✅ Vote recorded!"
}

func (h *MessageHandler) handlePollResults(args []string) string {
	if len(args) == 0 {
		return "Usage: !pollresults <poll_id>"
	}

	pollID := args[0]
	results, err := h.pollManager.GetPollResults(pollID)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return results
}

func (h *MessageHandler) handleListPolls(msg *Message) string {
	if msg.GroupID == "" {
		return "Polls are only available in group chats."
	}

	polls := h.pollManager.ListActivePolls(msg.GroupID)
	if len(polls) == 0 {
		return "No active polls in this group."
	}

	result := "📊 Active Polls:\n\n"
	for _, poll := range polls {
		result += fmt.Sprintf("• %s\n", poll)
	}

	return result
}

// Reminder command handlers
func (h *MessageHandler) handleReminder(args []string, msg *Message) string {
	if len(args) < 2 {
		return "Usage: !remind <time> <message>\nExamples: !remind 30m Meeting | !remind 2h Call John | !remind 1d Submit report"
	}

	timeStr := args[0]
	message := strings.Join(args[1:], " ")

	duration, err := parseDuration(timeStr)
	if err != nil {
		return fmt.Sprintf("Invalid time format. Use: 30s, 5m, 2h, 1d\nError: %v", err)
	}

	reminderID := h.reminderManager.AddReminder(message, msg.Sender, msg.GroupID, duration, false)
	return fmt.Sprintf("⏰ Reminder set! ID: %s\nI'll remind you in %s", reminderID, timeStr)
}

func (h *MessageHandler) handleListReminders(msg *Message) string {
	reminders := h.reminderManager.ListReminders(msg.Sender)
	if len(reminders) == 0 {
		return "You have no active reminders."
	}

	result := "⏰ Your Reminders:\n\n"
	for _, reminder := range reminders {
		result += fmt.Sprintf("• %s\n", reminder)
	}

	return result
}

func (h *MessageHandler) handleCancelReminder(args []string, msg *Message) string {
	if len(args) == 0 {
		return "Usage: !cancelreminder <reminder_id>"
	}

	reminderID := args[0]
	err := h.reminderManager.CancelReminder(reminderID, msg.Sender)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return "✅ Reminder cancelled."
}

// Stats command handlers
func (h *MessageHandler) handleStats(msg *Message) string {
	if msg.GroupID == "" {
		return "Statistics are only available in group chats."
	}

	return h.statsManager.GetGroupStats(msg.GroupID)
}

func (h *MessageHandler) handleMyStats(msg *Message) string {
	if msg.GroupID == "" {
		return "Statistics are only available in group chats."
	}

	return h.statsManager.GetUserActivity(msg.GroupID, msg.Sender)
}

// Announcement command handlers
func (h *MessageHandler) handleAnnounce(args []string, msg *Message) string {
	if msg.GroupID == "" {
		return "Announcements can only be made in group chats."
	}

	if len(args) == 0 {
		return "Usage: !announce <message>"
	}

	message := strings.Join(args, " ")
	announcementID := h.announcementManager.AddAnnouncement(message, msg.Sender, msg.GroupID)

	return fmt.Sprintf("📌 Announcement pinned! ID: %s", announcementID)
}

func (h *MessageHandler) handleListAnnouncements(msg *Message) string {
	if msg.GroupID == "" {
		return "Announcements are only available in group chats."
	}

	return h.announcementManager.GetAnnouncements(msg.GroupID)
}

// Mention command handlers
func (h *MessageHandler) handleMentions(msg *Message) string {
	return h.mentionTracker.GetMentions(msg.Sender, 10)
}

func (h *MessageHandler) handleClearMentions(msg *Message) string {
	h.mentionTracker.ClearMentions(msg.Sender)
	return "✅ Mentions cleared."
}

// parseDuration parses duration strings like "30m", "2h", "1d"
func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit: %s (use s, m, h, or d)", unit)
	}
}
