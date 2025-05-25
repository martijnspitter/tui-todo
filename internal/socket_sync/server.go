package socket_sync

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// Server handles the socket server functionality for the primary app instance
type Server struct {
	socketPath      string
	listener        net.Listener
	clients         map[string]net.Conn
	clientsMutex    sync.RWMutex
	shutdown        chan struct{}
	shutdownWg      sync.WaitGroup
	notifListener   NotificationListener
	broadcastLock   sync.Mutex
	heartbeatTicker *time.Ticker
}

// NewServer creates a new socket server
func NewServer(socketPath string, listener NotificationListener) (*Server, error) {
	// Validate socket path
	if err := ValidateSocketPath(socketPath); err != nil {
		return nil, err
	}

	// Clean up any existing socket file
	if err := CleanupSocket(socketPath); err != nil {
		return nil, err
	}

	// Ensure directory exists
	if err := EnsureSocketDirectory(socketPath); err != nil {
		return nil, err
	}

	// Create the server
	server := &Server{
		socketPath:    socketPath,
		clients:       make(map[string]net.Conn),
		shutdown:      make(chan struct{}),
		notifListener: listener,
	}

	// Create the socket listener
	var err error
	server.listener, err = net.Listen(DefaultProtocol, socketPath)
	if err != nil {
		return nil, &SocketError{Op: "listen", Err: err}
	}

	return server, nil
}

// Start begins accepting client connections
func (s *Server) Start() error {
	// Set permission on socket file to ensure all users can connect
	if err := os.Chmod(s.socketPath, 0666); err != nil {
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	// Start the accept loop
	s.shutdownWg.Add(1)
	go s.acceptLoop()

	s.heartbeatTicker = time.NewTicker(10 * time.Second)
	s.shutdownWg.Add(1)
	go s.heartbeatLoop()

	return nil
}

func (s *Server) heartbeatLoop() {
	defer s.shutdownWg.Done()

	for {
		select {
		case <-s.shutdown:
			if s.heartbeatTicker != nil {
				s.heartbeatTicker.Stop()
			}
			return
		case <-s.heartbeatTicker.C:
			s.sendHeartbeat()
		}
	}
}

func (s *Server) sendHeartbeat() {
	// Check if we have any clients first
	s.clientsMutex.RLock()
	clientCount := len(s.clients)
	s.clientsMutex.RUnlock()

	if clientCount == 0 {
		return // No clients to send heartbeats to
	}

	heartbeat := Notification{
		Type:      Heartbeat,
		Timestamp: time.Now(),
		ID:        0,
	}

	// Send heartbeat to all clients
	if err := s.broadcastToOthers(heartbeat, ""); err != nil {
		log.Error("Failed to send heartbeat", "error", err)
	}
}

// acceptLoop accepts new client connections
func (s *Server) acceptLoop() {
	defer s.shutdownWg.Done()

	// Create a ticker for logging connection status periodically
	statusTicker := time.NewTicker(5 * time.Minute)
	defer statusTicker.Stop()

	for {
		select {
		case <-s.shutdown:
			log.Warn("Shutting down accept loop")
			return
		case <-statusTicker.C:
			s.logStatus()
		default:
			// Set accept timeout so we can check for shutdown
			err := s.listener.(*net.UnixListener).SetDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				log.Error("Failed to set listener deadline", "error", err)
				continue
			}

			// Accept a new connection
			conn, err := s.listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue // accept deadline reached
				}
				if IsSocketClosed(err) {
					// Socket was closed, probably during shutdown
					return
				}
				log.Error("Failed to accept connection", "error", err)
				continue
			}

			// Handle the new client connection
			clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
			log.Info("New client connected", "id", clientID)
			s.addClient(clientID, conn)

			// Start a goroutine to handle this client
			s.shutdownWg.Add(1)
			go s.handleClient(clientID, conn)
		}
	}
}

// addClient adds a client to the clients map
func (s *Server) addClient(id string, conn net.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	s.clients[id] = conn
}

// removeClient removes a client from the clients map
func (s *Server) removeClient(id string) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	delete(s.clients, id)
}

// handleClient processes messages from a connected client
func (s *Server) handleClient(clientID string, conn net.Conn) {
	defer func() {
		conn.Close()
		s.removeClient(clientID)
		s.shutdownWg.Done()
	}()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			// Set read deadline
			if err := conn.SetReadDeadline(time.Now().Add(120 * time.Second)); err != nil {
				log.Error("Failed to set read deadline", "id", clientID, "error", err)
				return
			}

			// Read a message from the client
			notification, err := ReadMessage(conn)
			if err != nil {
				// Handle timeouts gracefully
				if IsTimeout(err) {
					continue
				}

				if IsSocketClosed(err) {
					return
				}

				log.Error("Error reading from client", "id", clientID, "error", err)
				return
			}

			if notification.Type == Heartbeat {
				continue
			}

			// Notify our listener first
			s.notifListener.OnNotification(notification)

			// Broadcast to all other clients
			if err := s.broadcastToOthers(notification, clientID); err != nil {
				// Remove the failing client so that subsequent broadcasts don’t keep erroring
				log.Error("Broadcast to others failed", "sender", clientID, "error", err)
			}
		}
	}
}

// Broadcast sends a notification to all connected clients
func (s *Server) Broadcast(notification Notification) error {
	// First notify our own listener
	s.notifListener.OnNotification(notification)

	return s.broadcastToOthers(notification, "")
}

// broadcastToOthers sends a notification to all clients except the sender
func (s *Server) broadcastToOthers(notification Notification, senderID string) error {
	s.broadcastLock.Lock()
	defer s.broadcastLock.Unlock()
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	if len(s.clients) == 0 {
		return nil // No clients to send to
	}

	var lastErr error
	for id, conn := range s.clients {
		if id == senderID {
			continue // Skip the sender
		}

		if err := WriteMessage(conn, notification); err != nil {
			// Remove dead client to prevent noisy logs & leaks
			go func(c net.Conn, clientID string) {
				c.Close()
				s.removeClient(clientID)
			}(conn, id)
		}
	}

	return lastErr
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	// Signal all goroutines to shut down
	close(s.shutdown)

	// Close the listener to stop accepting connections
	if s.listener != nil {
		s.listener.Close()
	}

	if s.heartbeatTicker != nil {
		s.heartbeatTicker.Stop()
	}

	// Close all client connections
	func() {
		s.clientsMutex.Lock()
		defer s.clientsMutex.Unlock()

		for id, conn := range s.clients {
			log.Warn("Closing client connection", "id", id)
			conn.Close()
		}
		s.clients = make(map[string]net.Conn)
	}()

	// Wait for all goroutines to exit (with timeout)
	done := make(chan struct{})
	go func() {
		s.shutdownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Clean shutdown completed
	case <-time.After(5 * time.Second):
		log.Warn("Server shutdown timed out")
	}

	// Clean up socket file
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// logStatus outputs the current number of connected clients
func (s *Server) logStatus() {
	s.clientsMutex.RLock()
	count := len(s.clients)
	s.clientsMutex.RUnlock()
	log.Debug("Sync server status", "connected_clients", count)
}
