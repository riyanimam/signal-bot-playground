# Signal Bot Features Guide 🚀

This document provides a comprehensive guide to all features in the Signal Bot, especially those designed for managing large groups (50-700+ people) and Slack-like functionality.

## Table of Contents
- [Large Group Management](#large-group-management)
- [Slack-Like Features](#slack-like-features)
- [Basic Commands](#basic-commands)
- [Use Cases](#use-cases)
- [Best Practices](#best-practices)

---

## Large Group Management

### 📊 Polls & Voting

**Purpose**: Make group decisions democratically with structured voting.

**Commands**:
- `!poll <question> | <option1> | <option2> | ...` - Create a poll
- `!vote <poll_id> <option>` - Cast your vote
- `!pollresults <poll_id>` - View current results
- `!polls` - List all active polls

**Examples**:
```
!poll What time works best for the meeting? | 2 PM | 3 PM | 4 PM | 5 PM
!vote poll_1234567890 3 PM
!pollresults poll_1234567890
```

**Use Cases**:
- Schedule meetings that work for most people
- Choose venues for events
- Make purchasing decisions for group resources
- Gather opinions on important topics
- Vote on group rules or policies

**Features**:
- Visual bar charts showing vote distribution
- Real-time vote counting
- Percentage calculations
- Support for unlimited options
- One vote per person (can change vote)

---

### 📈 Group Statistics

**Purpose**: Understand group dynamics and member engagement.

**Commands**:
- `!stats` - View comprehensive group statistics
- `!mystats` - View your personal activity

**What You'll See**:
- Total message count
- Active members today
- Total member count
- Top 10 contributors with percentages
- Individual contribution metrics
- Last activity timestamps

**Use Cases**:
- Identify most engaged members
- Understand group activity patterns
- Recognize active contributors
- Track group growth over time
- Monitor participation rates

**Example Output**:
```
📊 Group Statistics

Total messages: 5,234
Active members today: 87
Total members: 156

Top contributors:
1. ****5678: 523 messages (10.0%)
2. ****9012: 445 messages (8.5%)
3. ****3456: 389 messages (7.4%)
...
```

---

### 📌 Announcements

**Purpose**: Pin important messages that need to stay visible.

**Commands**:
- `!announce <message>` - Create a pinned announcement
- `!announcements` - View all pinned announcements

**Use Cases**:
- Share important updates everyone needs to see
- Post group rules and guidelines
- Announce events and deadlines
- Share emergency information
- Keep contact information visible

**Features**:
- Up to 10 announcements per group
- Shows who created each announcement
- Displays how long ago it was posted
- Recent announcements appear first
- Persistent across all group members

**Example**:
```
!announce Next meeting: March 15th at 3 PM. Location: Conference Room A

📌 Announcement pinned! ID: announcement_1234567890
```

---

## Slack-Like Features

### 💬 @Mentions

**Purpose**: Get someone's attention in busy group chats.

**How to Use**:
Simply include `@+phonenumber` in your message:
```
Hey @+1234567890, can you review the document?
```

**Commands**:
- `!mentions` - View your recent mentions (last 10)
- `!clearmentions` - Clear your mention notifications

**Features**:
- Automatic mention detection
- Tracks up to 50 mentions per user
- Shows who mentioned you
- Displays when you were mentioned
- Shows the full message context

**Use Cases**:
- Direct questions to specific people
- Request reviews or feedback
- Tag team members in discussions
- Ensure important messages are seen
- Coordinate with specific members

---

### ⏰ Reminders

**Purpose**: Never forget important tasks or events.

**Commands**:
- `!remind <time> <message>` - Set a reminder
- `!reminders` - List your active reminders
- `!cancelreminder <id>` - Cancel a reminder

**Time Formats**:
- `30s` - 30 seconds
- `5m` - 5 minutes
- `2h` - 2 hours
- `1d` - 1 day

**Examples**:
```
!remind 30m Call the client
!remind 2h Submit the report
!remind 1d Team meeting preparation
```

**Use Cases**:
- Personal task reminders
- Meeting preparation
- Deadline tracking
- Follow-up reminders
- Event notifications

**Features**:
- Automatic delivery at scheduled time
- Personal reminders (private)
- Multiple active reminders
- Time-until-due display
- Easy cancellation

---

## Basic Commands

### 🔧 Utility Commands

- `!help` - Display all available commands with descriptions
- `!ping` - Check if the bot is responsive
- `!echo <text>` - Echo back your message (useful for testing)
- `!about` - Get information about the bot

---

## Use Cases

### For Large Community Groups (100-700 people)

**Challenges**:
- Important messages get lost quickly
- Hard to make group decisions
- Difficult to gauge member engagement
- No way to track who's active

**Solutions**:
1. Use **polls** for community decisions
2. Pin **announcements** for rules and important updates
3. Check **stats** to see engagement levels
4. Use **@mentions** to reach specific members

---

### For Project Teams (20-100 people)

**Challenges**:
- Coordinating schedules
- Tracking action items
- Ensuring everyone sees critical info
- Managing multiple sub-discussions

**Solutions**:
1. **Polls** for scheduling meetings
2. **Reminders** for deadlines and tasks
3. **Announcements** for project milestones
4. **@Mentions** for task assignments
5. **Stats** to track team participation

---

### For Event Planning Groups (50-200 people)

**Challenges**:
- Getting consensus on details
- Keeping everyone informed
- Managing RSVPs and preferences
- Coordinating volunteers

**Solutions**:
1. **Polls** for venue, timing, menu choices
2. **Announcements** for event details and updates
3. **Reminders** for event countdown
4. **Stats** to see participation levels
5. **@Mentions** to coordinate with volunteers

---

## Best Practices

### Polls
✅ **Do**:
- Keep questions clear and concise
- Provide 2-5 options for best results
- Close polls once decision is made
- Share results with the group

❌ **Don't**:
- Create too many simultaneous polls
- Use vague or confusing options
- Leave polls open indefinitely

### Announcements
✅ **Do**:
- Keep announcements brief and important
- Update outdated announcements
- Use for truly important info only
- Include deadlines and dates

❌ **Don't**:
- Pin casual messages
- Create too many announcements
- Use for conversations
- Leave outdated info pinned

### Mentions
✅ **Do**:
- Use for direct questions or requests
- Include context in your message
- Check your mentions regularly
- Clear old mentions periodically

❌ **Don't**:
- Over-mention people
- Use for general questions
- Expect instant responses
- Mention multiple people unnecessarily

### Stats
✅ **Do**:
- Check periodically to gauge health
- Recognize top contributors
- Use to identify inactive periods
- Share insights with moderators

❌ **Don't**:
- Compare members judgmentally
- Use stats to shame low contributors
- Check obsessively
- Share individual stats publicly

### Reminders
✅ **Do**:
- Set reminders for important tasks
- Use clear, actionable messages
- Cancel completed reminders
- Set appropriate lead times

❌ **Don't**:
- Set too many reminders
- Use for trivial things
- Set reminders too far in advance
- Forget to act on reminders

---

## Tips for Managing Large Groups

### 1. Set Clear Expectations
Use announcements to share group rules and guidelines when members join.

### 2. Make Decisions Transparently
Use polls for any group decision to ensure everyone has a voice.

### 3. Monitor Engagement
Check stats regularly to understand if the group is healthy and active.

### 4. Keep Important Info Visible
Pin critical information like meeting schedules, contact lists, or resources.

### 5. Facilitate Direct Communication
Encourage use of @mentions for questions directed at specific people.

### 6. Celebrate Participation
Use stats to recognize and thank active contributors.

### 7. Stay Organized
Use reminders for group milestones and regular activities.

---

## Feature Comparison: Signal vs Slack

| Feature | Signal Bot | Slack |
|---------|-----------|-------|
| Polls | ✅ Built-in | ✅ Via apps |
| Mentions | ✅ @+phone | ✅ @user |
| Pinned Messages | ✅ Announcements | ✅ Pins |
| Reminders | ✅ Personal | ✅ Personal + Channel |
| Statistics | ✅ Built-in | ✅ Via analytics |
| Threads | ❌ Not supported | ✅ Native |
| File Sharing | 🔶 Via Signal | ✅ Native |
| Search | 🔶 Basic via Signal | ✅ Advanced |
| Integrations | 🔶 Limited | ✅ Extensive |

**Key Advantages**:
- ✅ End-to-end encryption (Signal's core strength)
- ✅ No cost for large groups
- ✅ Privacy-focused
- ✅ No data collection
- ✅ Simple and focused feature set

---

## Getting Started

1. **Start with basics**: Try `!help` and `!ping` to get familiar
2. **Create a test poll**: Use `!poll` to see how voting works
3. **Pin your first announcement**: Share group rules or welcome message
4. **Try mentions**: Use @+phonenumber to see notifications
5. **Set a reminder**: Test `!remind 5m Test reminder`
6. **Check stats**: Run `!stats` to see your group metrics

---

## Support and Feedback

Found a bug? Have a feature request? 
Open an issue on GitHub: https://github.com/riyanimam/signal-bot-playground

---

Made with ❤️ for large Signal communities
