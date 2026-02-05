package main

import (
	"fmt"
	"sync"
	"time"
)

// Reminder represents a scheduled reminder
type Reminder struct {
	ID        string
	Message   string
	Recipient string
	GroupID   string
	DueTime   time.Time
	Recurring bool
	Interval  time.Duration
}

// ReminderManager manages reminders
type ReminderManager struct {
	reminders map[string]*Reminder
	mu        sync.RWMutex
	sendFunc  func(groupID, recipient, message string) error
}

// NewReminderManager creates a new reminder manager
func NewReminderManager(sendFunc func(groupID, recipient, message string) error) *ReminderManager {
	rm := &ReminderManager{
		reminders: make(map[string]*Reminder),
		sendFunc:  sendFunc,
	}
	go rm.checkReminders()
	return rm
}

// AddReminder adds a new reminder
func (rm *ReminderManager) AddReminder(message, recipient, groupID string, duration time.Duration, recurring bool) string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	reminderID := fmt.Sprintf("reminder_%d", time.Now().Unix())
	reminder := &Reminder{
		ID:        reminderID,
		Message:   message,
		Recipient: recipient,
		GroupID:   groupID,
		DueTime:   time.Now().Add(duration),
		Recurring: recurring,
		Interval:  duration,
	}

	rm.reminders[reminderID] = reminder
	return reminderID
}

// checkReminders continuously checks for due reminders
func (rm *ReminderManager) checkReminders() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rm.mu.Lock()
		now := time.Now()

		for id, reminder := range rm.reminders {
			if now.After(reminder.DueTime) {
				// Send reminder
				message := fmt.Sprintf("⏰ Reminder: %s", reminder.Message)
				if rm.sendFunc != nil {
					if err := rm.sendFunc(reminder.GroupID, reminder.Recipient, message); err != nil {
						// Log error but continue processing other reminders
						fmt.Printf("Failed to send reminder %s: %v\n", id, err)
					}
				}

				// Handle recurring reminders
				if reminder.Recurring {
					reminder.DueTime = now.Add(reminder.Interval)
				} else {
					delete(rm.reminders, id)
				}
			}
		}
		rm.mu.Unlock()
	}
}

// ListReminders returns all active reminders for a user
func (rm *ReminderManager) ListReminders(recipient string) []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var reminders []string
	for _, reminder := range rm.reminders {
		if reminder.Recipient == recipient {
			timeUntil := time.Until(reminder.DueTime)
			reminders = append(reminders, fmt.Sprintf("%s - %s (in %s)", 
				reminder.ID, reminder.Message, timeUntil.Round(time.Minute)))
		}
	}

	return reminders
}

// CancelReminder cancels a reminder
func (rm *ReminderManager) CancelReminder(reminderID, recipient string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	reminder, exists := rm.reminders[reminderID]
	if !exists {
		return fmt.Errorf("reminder not found")
	}

	if reminder.Recipient != recipient {
		return fmt.Errorf("not authorized to cancel this reminder")
	}

	delete(rm.reminders, reminderID)
	return nil
}
