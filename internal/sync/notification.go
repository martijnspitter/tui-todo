package sync

import "time"

type NotificationType string

const (
	TodoCreated         NotificationType = "TODO_CREATED"
	TodoUpdated         NotificationType = "TODO_UPDATED"
	TodoDeleted         NotificationType = "TODO_DELETED"
	TodoStatusChanged   NotificationType = "TODO_STATUS_CHANGED"
	TodoArchived        NotificationType = "TODO_ARCHIVED"
	TodoUnarchived      NotificationType = "TODO_UNARCHIVED"
	TodoTagAdded        NotificationType = "TODO_TAG_ADDED"
	TodoTagRemoved      NotificationType = "TODO_TAG_REMOVED"
	TodoDueDateSet      NotificationType = "TODO_DUE_DATE_SET"
	TodoDueDateCleared  NotificationType = "TODO_DUE_DATE_CLEARED"
	TodoPriorityChanged NotificationType = "TODO_PRIORITY_CHANGED"

	Heartbeat NotificationType = "HEARTBEAT"
)

// Notification represents a change to be broadcast to other instances
type Notification struct {
	Type      NotificationType `json:"type"`
	ID        int64            `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
}

// NotificationListener receives notifications from other app instances
type NotificationListener interface {
	OnNotification(notification Notification)
}
