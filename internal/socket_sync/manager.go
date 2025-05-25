package socket_sync

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
	osoperations "github.com/martijnspitter/tui-todo/internal/os-operations"
)

// Manager coordinates the synchronization between multiple application instances
type Manager struct {
	socketPath    string
	notifListener NotificationListener
	server        *Server
	client        *Client
	isPrimary     bool
	started       atomic.Bool
	startMutex    sync.Mutex
	stopped       atomic.Bool
	notifBuffer   []Notification // Buffer notifications during initialization
	bufferMutex   sync.Mutex
	lastPollTime  time.Time
	pollMutex     sync.Mutex
}

// NewManager creates a new synchronization manager
func NewManager(version string, listener NotificationListener) (*Manager, error) {
	socketPath := osoperations.GetFilePath("todo.sock", version)

	if !filepath.IsAbs(socketPath) {
		absPath, err := filepath.Abs(socketPath)
		if err != nil {
			return nil, fmt.Errorf("failed to convert socket path to absolute: %w", err)
		}
		socketPath = absPath
	}

	manager := &Manager{
		socketPath:    socketPath,
		notifListener: listener,
		lastPollTime:  time.Now(),
	}

	return manager, nil
}

// Start initializes the sync system with leader election
func (m *Manager) Start() error {
	m.startMutex.Lock()
	defer m.startMutex.Unlock()

	if m.started.Load() {
		return nil // Already started
	}

	// Attempt to become the primary instance by starting a server
	server, err := NewServer(m.socketPath, m)
	if err == nil {
		// Successfully created server, we are the primary
		m.server = server
		m.isPrimary = true

		if err := m.server.Start(); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		log.Info("This instance is the primary (server)")
	} else {
		// Failed to create server, we are a secondary instance
		log.Debug("Failed to create server", "error", err)
		log.Info("This instance is a secondary (client)")

		// Wait a brief moment to ensure the server is ready
		time.Sleep(100 * time.Millisecond)

		// Create and start the client
		client, err := NewClient(m.socketPath, m)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		m.client = client
		m.isPrimary = false

		if err := m.client.Start(); err != nil {
			return fmt.Errorf("failed to start client: %w", err)
		}
	}

	m.started.Store(true)

	// Process any notifications that were buffered during initialization
	m.processBufferedNotifications()

	return nil
}

// Stop gracefully shuts down the sync system
func (m *Manager) Stop() error {
	m.startMutex.Lock()
	defer m.startMutex.Unlock()

	if !m.started.Load() || m.stopped.Load() {
		return nil // Already stopped or never started
	}

	var err error
	if m.isPrimary && m.server != nil {
		err = m.server.Stop()
	} else if m.client != nil {
		err = m.client.Stop()
	}

	m.stopped.Store(true)
	return err
}

// NotifyChange broadcasts a change notification to other instances
func (m *Manager) NotifyChange(notificationType NotificationType, id int64) error {
	if m.stopped.Load() {
		return nil // Already shut down
	}

	// If not yet started, buffer so we can replay after Start().
	if !m.started.Load() {
		m.bufferMutex.Lock()
		m.notifBuffer = append(m.notifBuffer, Notification{
			Type:      notificationType,
			ID:        id,
			Timestamp: time.Now(),
		})
		m.bufferMutex.Unlock()
		return nil
	}

	notification := Notification{
		Type:      notificationType,
		ID:        id,
		Timestamp: time.Now(),
	}

	if m.isPrimary && m.server != nil {
		return m.server.Broadcast(notification)
	} else if m.client != nil {
		return m.client.SendNotification(notification)
	}

	// If we're not fully initialized yet, buffer the notification
	m.bufferMutex.Lock()
	defer m.bufferMutex.Unlock()
	m.notifBuffer = append(m.notifBuffer, notification)

	return nil
}

// processBufferedNotifications sends any notifications that were queued before initialization
func (m *Manager) processBufferedNotifications() {
	m.bufferMutex.Lock()
	buffer := m.notifBuffer
	m.notifBuffer = nil
	m.bufferMutex.Unlock()

	for _, notification := range buffer {
		if m.isPrimary && m.server != nil {
			if err := m.server.Broadcast(notification); err != nil {
				log.Error("Failed to replay buffered notification", "error", err)
			}
		} else if m.client != nil {
			if err := m.client.SendNotification(notification); err != nil {
				log.Error("Failed to replay buffered notification", "error", err)
			}
		}
	}
}

// IsPrimary returns whether this instance is the primary (server)
func (m *Manager) IsPrimary() bool {
	return m.isPrimary
}

// UpdatePollingTime updates the timestamp of the last polling operation
func (m *Manager) UpdatePollingTime() {
	m.pollMutex.Lock()
	m.lastPollTime = time.Now()
	m.pollMutex.Unlock()
}

// TimeSinceLastPoll returns the duration since the last polling operation
func (m *Manager) TimeSinceLastPoll() time.Duration {
	m.pollMutex.Lock()
	defer m.pollMutex.Unlock()
	return time.Since(m.lastPollTime)
}

// GetSocketPath returns the path to the Unix socket used by this manager
func (m *Manager) GetSocketPath() string {
	return m.socketPath
}

// OnNotification implements the NotificationListener interface
// This will be called when a notification is received from another instance
func (m *Manager) OnNotification(notification Notification) {
	// Ignore heartbeat messages - they're just for connection maintenance
	if notification.Type == Heartbeat {
		return
	}

	// Update the polling time to avoid unnecessary polling right after a notification
	m.UpdatePollingTime()

	// Forward the notification to the registered listener
	if m.notifListener != nil {
		m.notifListener.OnNotification(notification)
	}
}

// GetProcessID returns the current process ID for debugging
func (m *Manager) GetProcessID() int {
	return os.Getpid()
}
