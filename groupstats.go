package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// GroupStats tracks group activity statistics
type GroupStats struct {
	GroupID         string
	MessageCounts   map[string]int // sender -> count
	LastActivity    map[string]time.Time
	TotalMessages   int
	ActiveToday     int
	CreatedAt       time.Time
}

// StatsManager manages group statistics
type StatsManager struct {
	stats map[string]*GroupStats // groupID -> stats
	mu    sync.RWMutex
}

// NewStatsManager creates a new stats manager
func NewStatsManager() *StatsManager {
	return &StatsManager{
		stats: make(map[string]*GroupStats),
	}
}

// RecordMessage records a message in the stats
func (sm *StatsManager) RecordMessage(groupID, sender string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if groupID == "" {
		return // Skip DMs
	}

	stats, exists := sm.stats[groupID]
	if !exists {
		stats = &GroupStats{
			GroupID:       groupID,
			MessageCounts: make(map[string]int),
			LastActivity:  make(map[string]time.Time),
			CreatedAt:     time.Now(),
		}
		sm.stats[groupID] = stats
	}

	stats.MessageCounts[sender]++
	stats.LastActivity[sender] = time.Now()
	stats.TotalMessages++

	// Check if message is from today
	if time.Since(stats.LastActivity[sender]) < 24*time.Hour {
		// Recalculate active today
		stats.ActiveToday = 0
		today := time.Now().Truncate(24 * time.Hour)
		for _, lastSeen := range stats.LastActivity {
			if lastSeen.After(today) {
				stats.ActiveToday++
			}
		}
	}
}

// GetGroupStats returns statistics for a group
func (sm *StatsManager) GetGroupStats(groupID string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats, exists := sm.stats[groupID]
	if !exists {
		return "No statistics available for this group yet."
	}

	// Sort users by message count
	type userCount struct {
		sender string
		count  int
	}
	
	var userCounts []userCount
	for sender, count := range stats.MessageCounts {
		userCounts = append(userCounts, userCount{sender, count})
	}

	sort.Slice(userCounts, func(i, j int) bool {
		return userCounts[i].count > userCounts[j].count
	})

	result := fmt.Sprintf("📊 Group Statistics\n\n")
	result += fmt.Sprintf("Total messages: %d\n", stats.TotalMessages)
	result += fmt.Sprintf("Active members today: %d\n", stats.ActiveToday)
	result += fmt.Sprintf("Total members: %d\n\n", len(stats.MessageCounts))
	
	result += "Top contributors:\n"
	for i, uc := range userCounts {
		if i >= 10 {
			break
		}
		maskedSender := maskPhoneNumber(uc.sender)
		percentage := float64(uc.count) / float64(stats.TotalMessages) * 100
		result += fmt.Sprintf("%d. %s: %d messages (%.1f%%)\n", i+1, maskedSender, uc.count, percentage)
	}

	return result
}

// GetUserActivity returns activity info for a specific user
func (sm *StatsManager) GetUserActivity(groupID, sender string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats, exists := sm.stats[groupID]
	if !exists {
		return "No statistics available for this group yet."
	}

	count := stats.MessageCounts[sender]
	lastSeen := stats.LastActivity[sender]
	
	result := fmt.Sprintf("📊 Your Activity\n\n")
	result += fmt.Sprintf("Messages sent: %d\n", count)
	
	if !lastSeen.IsZero() {
		result += fmt.Sprintf("Last active: %s ago\n", time.Since(lastSeen).Round(time.Minute))
	}
	
	if stats.TotalMessages > 0 {
		percentage := float64(count) / float64(stats.TotalMessages) * 100
		result += fmt.Sprintf("Contribution: %.1f%% of group messages", percentage)
	}

	return result
}
