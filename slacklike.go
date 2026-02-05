package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Announcement represents a pinned/important message
type Announcement struct {
	ID        string
	Message   string
	Creator   string
	GroupID   string
	CreatedAt time.Time
}

// AnnouncementManager manages group announcements
type AnnouncementManager struct {
	announcements map[string][]*Announcement // groupID -> announcements
	mu            sync.RWMutex
	maxPerGroup   int
}

// NewAnnouncementManager creates a new announcement manager
func NewAnnouncementManager(maxPerGroup int) *AnnouncementManager {
	return &AnnouncementManager{
		announcements: make(map[string][]*Announcement),
		maxPerGroup:   maxPerGroup,
	}
}

// AddAnnouncement adds a new announcement
func (am *AnnouncementManager) AddAnnouncement(message, creator, groupID string) string {
	am.mu.Lock()
	defer am.mu.Unlock()

	announcementID := fmt.Sprintf("announcement_%d", time.Now().Unix())
	announcement := &Announcement{
		ID:        announcementID,
		Message:   message,
		Creator:   creator,
		GroupID:   groupID,
		CreatedAt: time.Now(),
	}

	announcements := am.announcements[groupID]
	announcements = append(announcements, announcement)

	// Keep only the most recent announcements
	if len(announcements) > am.maxPerGroup {
		announcements = announcements[len(announcements)-am.maxPerGroup:]
	}

	am.announcements[groupID] = announcements
	return announcementID
}

// GetAnnouncements returns all announcements for a group
func (am *AnnouncementManager) GetAnnouncements(groupID string) string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	announcements, exists := am.announcements[groupID]
	if !exists || len(announcements) == 0 {
		return "No announcements in this group."
	}

	result := "📌 Pinned Announcements:\n\n"
	for i := len(announcements) - 1; i >= 0; i-- {
		ann := announcements[i]
		timeAgo := time.Since(ann.CreatedAt)
		result += fmt.Sprintf("%d. %s\n   (by %s, %s ago)\n\n",
			len(announcements)-i,
			ann.Message,
			maskPhoneNumber(ann.Creator),
			formatDuration(timeAgo))
	}

	return result
}

// RemoveAnnouncement removes an announcement
func (am *AnnouncementManager) RemoveAnnouncement(announcementID, groupID, requester string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	announcements, exists := am.announcements[groupID]
	if !exists {
		return fmt.Errorf("no announcements in this group")
	}

	for i, ann := range announcements {
		if ann.ID == announcementID {
			// Only creator can remove their announcement
			if ann.Creator != requester {
				return fmt.Errorf("only the creator can remove this announcement")
			}

			// Remove announcement
			am.announcements[groupID] = append(announcements[:i], announcements[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("announcement not found")
}

// MentionTracker tracks mentions in messages
type MentionTracker struct {
	mentions map[string][]Mention // recipient -> mentions
	mu       sync.RWMutex
}

// Mention represents a mention in a message
type Mention struct {
	Sender    string
	Message   string
	GroupID   string
	Timestamp time.Time
}

// NewMentionTracker creates a new mention tracker
func NewMentionTracker() *MentionTracker {
	return &MentionTracker{
		mentions: make(map[string][]Mention),
	}
}

// TrackMention tracks a mention
func (mt *MentionTracker) TrackMention(sender, recipient, message, groupID string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mention := Mention{
		Sender:    sender,
		Message:   message,
		GroupID:   groupID,
		Timestamp: time.Now(),
	}

	mentions := mt.mentions[recipient]
	mentions = append(mentions, mention)

	// Keep only last 50 mentions per user
	if len(mentions) > 50 {
		mentions = mentions[len(mentions)-50:]
	}

	mt.mentions[recipient] = mentions
}

// GetMentions returns unread mentions for a user
func (mt *MentionTracker) GetMentions(recipient string, limit int) string {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	mentions, exists := mt.mentions[recipient]
	if !exists || len(mentions) == 0 {
		return "You have no mentions."
	}

	if limit <= 0 || limit > len(mentions) {
		limit = len(mentions)
	}

	result := fmt.Sprintf("💬 Your Mentions (%d):\n\n", len(mentions))
	
	// Show most recent mentions first
	for i := len(mentions) - 1; i >= len(mentions)-limit && i >= 0; i-- {
		mention := mentions[i]
		timeAgo := time.Since(mention.Timestamp)
		result += fmt.Sprintf("%d. From %s (%s ago)\n   %s\n\n",
			len(mentions)-i,
			maskPhoneNumber(mention.Sender),
			formatDuration(timeAgo),
			mention.Message)
	}

	return result
}

// ClearMentions clears mentions for a user
func (mt *MentionTracker) ClearMentions(recipient string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	delete(mt.mentions, recipient)
}

// ParseMentions extracts phone number mentions from a message
// Looks for patterns like @+1234567890
func ParseMentions(message string) []string {
	var mentions []string
	parts := strings.Fields(message)

	for _, part := range parts {
		if strings.HasPrefix(part, "@+") {
			phone := strings.TrimPrefix(part, "@")
			if isValidPhoneNumber(phone) {
				mentions = append(mentions, phone)
			}
		}
	}

	return mentions
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	} else {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
}
