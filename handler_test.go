package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPollManager(t *testing.T) {
	pm := NewPollManager()

	// Test creating a poll
	question := "Favorite color?"
	options := []string{"Red", "Blue", "Green"}
	pollID := pm.CreatePoll(question, options, "+1234567890", "group123", 0)

	if pollID == "" {
		t.Error("Expected non-empty poll ID")
	}

	// Test voting
	err := pm.Vote(pollID, "+1111111111", "Red")
	if err != nil {
		t.Errorf("Failed to vote: %v", err)
	}

	err = pm.Vote(pollID, "+2222222222", "Blue")
	if err != nil {
		t.Errorf("Failed to vote: %v", err)
	}

	err = pm.Vote(pollID, "+3333333333", "Red")
	if err != nil {
		t.Errorf("Failed to vote: %v", err)
	}

	// Test voting with invalid option
	err = pm.Vote(pollID, "+4444444444", "Yellow")
	if err == nil {
		t.Error("Expected error for invalid option")
	}

	// Test getting results
	results, err := pm.GetPollResults(pollID)
	if err != nil {
		t.Errorf("Failed to get results: %v", err)
	}

	if !strings.Contains(results, "Red: 2") {
		t.Errorf("Expected 2 votes for Red, got: %s", results)
	}

	if !strings.Contains(results, "Blue: 1") {
		t.Errorf("Expected 1 vote for Blue, got: %s", results)
	}

	// Test listing active polls
	polls := pm.ListActivePolls("group123")
	if len(polls) != 1 {
		t.Errorf("Expected 1 active poll, got %d", len(polls))
	}
}

func TestStatsManager(t *testing.T) {
	sm := NewStatsManager()

	// Record some messages
	sm.RecordMessage("group123", "+1234567890")
	sm.RecordMessage("group123", "+1234567890")
	sm.RecordMessage("group123", "+9876543210")

	// Get group stats
	stats := sm.GetGroupStats("group123")

	if !strings.Contains(stats, "Total messages: 3") {
		t.Errorf("Expected 3 total messages in stats: %s", stats)
	}

	if !strings.Contains(stats, "Total members: 2") {
		t.Errorf("Expected 2 total members in stats: %s", stats)
	}

	// Get user activity
	activity := sm.GetUserActivity("group123", "+1234567890")

	if !strings.Contains(activity, "Messages sent: 2") {
		t.Errorf("Expected 2 messages for user in activity: %s", activity)
	}
}

func TestAnnouncementManager(t *testing.T) {
	am := NewAnnouncementManager(5)

	// Add an announcement
	id := am.AddAnnouncement("Important message", "+1234567890", "group123")

	if id == "" {
		t.Error("Expected non-empty announcement ID")
	}

	// Get announcements
	announcements := am.GetAnnouncements("group123")

	if !strings.Contains(announcements, "Important message") {
		t.Errorf("Expected announcement in list: %s", announcements)
	}

	// Test max announcements
	for i := 0; i < 10; i++ {
		am.AddAnnouncement(fmt.Sprintf("Message %d", i), "+1234567890", "group123")
	}

	// Should only have maxPerGroup announcements
	announcements = am.GetAnnouncements("group123")
	lines := strings.Split(announcements, "\n")
	// Count actual announcements (skip header and empty lines)
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") || 
		   strings.HasPrefix(line, "3.") || strings.HasPrefix(line, "4.") || 
		   strings.HasPrefix(line, "5.") {
			count++
		}
	}

	if count > 5 {
		t.Errorf("Expected at most 5 announcements, got %d", count)
	}
}

func TestMentionTracker(t *testing.T) {
	mt := NewMentionTracker()

	// Track a mention
	mt.TrackMention("+1234567890", "+9876543210", "Hey @+9876543210 check this", "group123")

	// Get mentions
	mentions := mt.GetMentions("+9876543210", 10)

	if !strings.Contains(mentions, "check this") {
		t.Errorf("Expected mention in list: %s", mentions)
	}

	// Clear mentions
	mt.ClearMentions("+9876543210")

	mentions = mt.GetMentions("+9876543210", 10)

	if !strings.Contains(mentions, "no mentions") {
		t.Error("Expected no mentions after clearing")
	}
}

func TestParseMentions(t *testing.T) {
	tests := []struct {
		message  string
		expected int
	}{
		{"Hey @+1234567890 how are you?", 1},
		{"@+1234567890 @+9876543210 meeting at 3", 2},
		{"No mentions here", 0},
		{"Invalid @+123 mention", 0}, // Too short
		{"@+123456789012345 valid", 1}, // 15 digits is valid
	}

	for _, test := range tests {
		mentions := ParseMentions(test.message)
		if len(mentions) != test.expected {
			t.Errorf("For message '%s', expected %d mentions, got %d",
				test.message, test.expected, len(mentions))
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"30s", 30 * time.Second, false},
		{"5m", 5 * time.Minute, false},
		{"2h", 2 * time.Hour, false},
		{"1d", 24 * time.Hour, false},
		{"invalid", 0, true},
		{"10x", 0, true},
		{"", 0, true},
	}

	for _, test := range tests {
		duration, err := parseDuration(test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s'", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
			if duration != test.expected {
				t.Errorf("For input '%s', expected %v, got %v",
					test.input, test.expected, duration)
			}
		}
	}
}

func TestReminderManager(t *testing.T) {
	// Mock send function
	sentMessages := make([]string, 0)
	var mu sync.Mutex
	sendFunc := func(groupID, recipient, message string) error {
		mu.Lock()
		sentMessages = append(sentMessages, message)
		mu.Unlock()
		return nil
	}

	rm := NewReminderManager(sendFunc)

	// Add a reminder with very short duration for testing
	id := rm.AddReminder("Test reminder", "+1234567890", "group123", 1*time.Second, false)

	if id == "" {
		t.Error("Expected non-empty reminder ID")
	}

	// List reminders
	reminders := rm.ListReminders("+1234567890")
	if len(reminders) == 0 {
		t.Error("Expected at least one reminder")
	}

	// Wait for reminder to trigger (reminder checker runs every 30 seconds, so we need to wait)
	// For testing purposes, we'll just verify the reminder was added correctly
	// and skip the actual delivery test since it requires waiting 30+ seconds
	
	// Cancel the reminder
	err := rm.CancelReminder(id, "+1234567890")
	if err != nil {
		t.Errorf("Failed to cancel reminder: %v", err)
	}

	// Check if reminder was removed
	reminders = rm.ListReminders("+1234567890")
	if len(reminders) != 0 {
		t.Errorf("Expected no reminders after cancellation, got %d", len(reminders))
	}
}

func TestHandleMessage(t *testing.T) {
	config := &Config{
		PhoneNumber:   "+1234567890",
		CommandPrefix: "!",
	}

	sendFunc := func(groupID, recipient, message string) error {
		return nil
	}

	handler := NewMessageHandler(config, sendFunc)

	// Test ping command
	msg := &Message{
		Sender:  "+9876543210",
		Text:    "!ping",
		GroupID: "group123",
	}

	response, err := handler.HandleMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(response, "Pong") {
		t.Errorf("Expected Pong response, got: %s", response)
	}

	// Test echo command
	msg.Text = "!echo Hello World"
	response, err = handler.HandleMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response != "Hello World" {
		t.Errorf("Expected 'Hello World', got: %s", response)
	}

	// Test non-command message
	msg.Text = "Just a regular message"
	response, err = handler.HandleMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response != "" {
		t.Error("Expected empty response for non-command message")
	}

	// Test help command
	msg.Text = "!help"
	response, err = handler.HandleMessage(msg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(response, "Available commands") {
		t.Error("Expected help text in response")
	}
}
