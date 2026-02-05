package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Poll represents a poll/vote in a group
type Poll struct {
	ID        string
	Question  string
	Options   []string
	Votes     map[string]string // voter phone -> option
	Creator   string
	GroupID   string
	CreatedAt time.Time
	Duration  time.Duration
}

// PollManager manages active polls
type PollManager struct {
	polls map[string]*Poll // pollID -> Poll
	mu    sync.RWMutex
}

// NewPollManager creates a new poll manager
func NewPollManager() *PollManager {
	return &PollManager{
		polls: make(map[string]*Poll),
	}
}

// CreatePoll creates a new poll
func (pm *PollManager) CreatePoll(question string, options []string, creator, groupID string, duration time.Duration) string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pollID := fmt.Sprintf("poll_%d", time.Now().Unix())
	poll := &Poll{
		ID:        pollID,
		Question:  question,
		Options:   options,
		Votes:     make(map[string]string),
		Creator:   creator,
		GroupID:   groupID,
		CreatedAt: time.Now(),
		Duration:  duration,
	}

	pm.polls[pollID] = poll

	// Auto-close poll after duration
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			pm.ClosePoll(pollID)
		}()
	}

	return pollID
}

// Vote records a vote in a poll
func (pm *PollManager) Vote(pollID, voter, option string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	poll, exists := pm.polls[pollID]
	if !exists {
		return fmt.Errorf("poll not found")
	}

	// Check if poll is still active
	if poll.Duration > 0 && time.Since(poll.CreatedAt) > poll.Duration {
		return fmt.Errorf("poll has closed")
	}

	// Validate option
	validOption := false
	for _, opt := range poll.Options {
		if strings.EqualFold(opt, option) {
			validOption = true
			option = opt // Use the original casing
			break
		}
	}

	if !validOption {
		return fmt.Errorf("invalid option")
	}

	poll.Votes[voter] = option
	return nil
}

// GetPollResults returns the results of a poll
func (pm *PollManager) GetPollResults(pollID string) (string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	poll, exists := pm.polls[pollID]
	if !exists {
		return "", fmt.Errorf("poll not found")
	}

	// Count votes
	counts := make(map[string]int)
	for _, option := range poll.Options {
		counts[option] = 0
	}

	for _, vote := range poll.Votes {
		counts[vote]++
	}

	// Format results
	result := fmt.Sprintf("📊 Poll Results: %s\n\n", poll.Question)
	totalVotes := len(poll.Votes)

	for _, option := range poll.Options {
		count := counts[option]
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(count) / float64(totalVotes) * 100
		}

		bar := strings.Repeat("█", count)
		result += fmt.Sprintf("%s: %d votes (%.1f%%)\n%s\n\n", option, count, percentage, bar)
	}

	result += fmt.Sprintf("Total votes: %d", totalVotes)

	if poll.Duration > 0 {
		timeLeft := poll.Duration - time.Since(poll.CreatedAt)
		if timeLeft > 0 {
			result += fmt.Sprintf("\nTime left: %s", timeLeft.Round(time.Minute))
		} else {
			result += "\n🔒 Poll closed"
		}
	}

	return result, nil
}

// ListActivePolls returns all active polls for a group
func (pm *PollManager) ListActivePolls(groupID string) []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var polls []string
	for _, poll := range pm.polls {
		if poll.GroupID == groupID {
			if poll.Duration == 0 || time.Since(poll.CreatedAt) <= poll.Duration {
				polls = append(polls, fmt.Sprintf("%s: %s", poll.ID, poll.Question))
			}
		}
	}

	return polls
}

// ClosePoll manually closes a poll
func (pm *PollManager) ClosePoll(pollID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	poll, exists := pm.polls[pollID]
	if !exists {
		return fmt.Errorf("poll not found")
	}

	// Set duration to expired
	poll.Duration = 0

	return nil
}
