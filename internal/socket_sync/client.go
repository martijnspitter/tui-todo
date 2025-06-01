package socket_sync

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

const (
	// MaxReconnectAttempts is the maximum number of reconnection attempts
	MaxReconnectAttempts = 5

	// InitialReconnectDelay is the initial delay between reconnection attempts
	InitialReconnectDelay = 500 * time.Millisecond

	// MaxReconnectDelay is the maximum delay between reconnection attempts
	MaxReconnectDelay = 30 * time.Second
)

// Client handles the socket client functionality for secondary app instances
type Client struct {
	socketPath     string
	conn           net.Conn
	connMutex      sync.RWMutex
	notifListener  NotificationListener
	shutdown       chan struct{}
	shutdownWg     sync.WaitGroup
	reconnecting   bool
	reconnectMutex sync.Mutex
	writeMutex     sync.Mutex
}

// NewClient creates a new socket client
func NewClient(socketPath string, listener NotificationListener) (*Client, error) {
	// Validate socket path
	if err := ValidateSocketPath(socketPath); err != nil {
		return nil, err
	}

	client := &Client{
		socketPath:    socketPath,
		notifListener: listener,
		shutdown:      make(chan struct{}),
	}

	// Initial connection attempt
	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a connection to the server
func (c *Client) connect() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	// Clean up any existing connection
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	// Connect to the socket
	conn, err := net.DialTimeout(DefaultProtocol, c.socketPath, 2*time.Second)
	if err != nil {
		return &SocketError{Op: "connect", Err: err}
	}

	c.conn = conn
	log.Info("Connected to sync server", "socket", c.socketPath)

	return nil
}

// Start begins the message receiving loop
func (c *Client) Start() error {
	c.shutdownWg.Add(1)
	go c.receiveLoop()

	return nil
}

// receiveLoop continuously reads messages from the server
func (c *Client) receiveLoop() {
	defer c.shutdownWg.Done()

	for {
		select {
		case <-c.shutdown:
			log.Info("Shutting down receive loop")
			return
		default:
			// Get current connection
			c.connMutex.RLock()
			conn := c.conn
			c.connMutex.RUnlock()

			if conn == nil {
				// No connection, wait a bit before trying again
				time.Sleep(1 * time.Second)
				continue
			}

			err := conn.SetReadDeadline(time.Now().Add(90 * time.Second))
			if err != nil {
				log.Error("Failed to set read deadline", "error", err)
			}

			// Read a message
			notification, err := ReadMessage(conn)
			if err != nil {
				if IsTimeout(err) {
					continue
				}
				if IsSocketClosed(err) || c.isShuttingDown() {
					// If we're shutting down, just exit
					if c.isShuttingDown() {
						return
					}

					// Otherwise try to reconnect
					log.Warn("Lost connection to server, attempting to reconnect")
					go c.reconnect()

					// Wait for reconnection or shutdown
					time.Sleep(1 * time.Second)
					continue
				}

				log.Error("Error reading from server", "error", err)
				continue
			}

			if notification.Type == Heartbeat {
				continue
			}

			// Process the notification
			log.Info("Received notification from server", "type", notification.Type, "todoID", notification.ID)
			c.notifListener.OnNotification(notification)
		}
	}
}

// reconnect attempts to reestablish the connection to the server
func (c *Client) reconnect() {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	// Avoid multiple simultaneous reconnection attempts
	if c.reconnecting {
		return
	}
	c.reconnecting = true
	defer func() { c.reconnecting = false }()

	// Use exponential backoff for reconnection attempts
	delay := InitialReconnectDelay

	for attempt := 1; attempt <= MaxReconnectAttempts; attempt++ {
		// Check if we're shutting down
		if c.isShuttingDown() {
			return
		}

		// Try to connect
		err := c.connect()
		if err == nil {
			log.Info("Successfully reconnected to server")
			return
		}

		log.Error("Reconnection failed", "attempt", attempt, "error", err)

		// Wait before next attempt, with exponential backoff
		select {
		case <-c.shutdown:
			return
		case <-time.After(delay):
			// Double the delay for next attempt, but cap it
			delay *= 2
			if delay > MaxReconnectDelay {
				delay = MaxReconnectDelay
			}
		}
	}

	log.Error("Failed to reconnect after maximum attempts", "attempts", MaxReconnectAttempts)
}

// SendNotification sends a notification to the server
func (c *Client) SendNotification(notification Notification) error {
	c.connMutex.RLock()
	conn := c.conn
	c.connMutex.RUnlock()

	if conn == nil {
		return errors.New("not connected to server")
	}

	c.writeMutex.Lock()
	err := WriteMessage(conn, notification)
	c.writeMutex.Unlock()

	if err != nil {
		// If the connection failed, try to reconnect
		if IsSocketClosed(err) && !c.isShuttingDown() {
			log.Warn("Connection lost while sending, attempting to reconnect")
			go c.reconnect()
		}
		return err
	}

	return nil
}

// Stop gracefully shuts down the client
func (c *Client) Stop() error {
	// Signal all goroutines to shut down
	close(c.shutdown)

	// Close the connection
	c.connMutex.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connMutex.Unlock()

	// Wait for all goroutines to exit (with timeout)
	done := make(chan struct{})
	go func() {
		c.shutdownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Clean shutdown completed
		log.Info("Sync client stopped successfully")
	case <-time.After(5 * time.Second):
		log.Warn("Client shutdown timed out")
	}

	return nil
}

// isShuttingDown checks if the client is in the process of shutting down
func (c *Client) isShuttingDown() bool {
	select {
	case <-c.shutdown:
		return true
	default:
		return false
	}
}
