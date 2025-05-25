package socket_sync

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxMessageSize  = 8192
	ReadTimeout     = 5 * time.Second
	WriteTimeout    = 5 * time.Second
	DefaultProtocol = "unix"
)

// SocketError represents errors specific to socket operations
type SocketError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *SocketError) Error() string {
	return fmt.Sprintf("socket %s error: %v", e.Op, e.Err)
}

func (e *SocketError) Unwrap() error {
	return e.Err
}

// EnsureSocketDirectory ensures the directory for the socket file exists
func EnsureSocketDirectory(socketPath string) error {
	dir := filepath.Dir(socketPath)

	// Check if directory exists
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return nil // Directory exists
	}

	if err != nil && !os.IsNotExist(err) {
		return &SocketError{Op: "check_dir", Err: err}
	}

	// Create directory with appropriate permissions
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &SocketError{Op: "create_dir", Err: err}
	}

	return nil
}

// CleanupSocket removes any stale socket file
func CleanupSocket(socketPath string) error {
	_, err := os.Stat(socketPath)
	if err == nil {
		// Socket file exists, try to connect to test if it's active
		conn, err := net.Dial(DefaultProtocol, socketPath)
		if err != nil {
			// Can't connect, so it's stale - remove it
			if err := os.Remove(socketPath); err != nil {
				return &SocketError{Op: "cleanup", Err: err}
			}
		} else {
			// Connection successful, socket is active
			conn.Close()
		}
	} else if !os.IsNotExist(err) {
		// Some error other than "not exists"
		return &SocketError{Op: "check_file", Err: err}
	}

	return nil
}

// WriteMessage writes a notification message to the connection
func WriteMessage(conn net.Conn, notification Notification) error {
	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return &SocketError{Op: "set_deadline", Err: err}
	}

	// Encode message to JSON
	data, err := json.Marshal(notification)
	if err != nil {
		return &SocketError{Op: "encode", Err: err}
	}

	// Add newline for message framing
	data = append(data, '\n')

	// Check message size
	if len(data) > MaxMessageSize {
		return &SocketError{
			Op:  "size_check",
			Err: errors.New("message exceeds maximum allowed size"),
		}
	}

	// Write message to connection
	_, err = conn.Write(data)
	if err != nil {
		return &SocketError{Op: "write", Err: err}
	}

	return nil
}

// ReadMessage reads a notification message from the connection
func ReadMessage(conn net.Conn) (Notification, error) {
	_ = conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	var notification Notification

	// Create buffered reader
	reader := bufio.NewReader(conn)

	// Read message up to newline
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return notification, &SocketError{Op: "read", Err: err}
	}

	if len(data) > MaxMessageSize {
		return notification, &SocketError{Op: "size_check", Err: errors.New("message too large")}
	}

	// Decode JSON message
	if err := json.Unmarshal(data, &notification); err != nil {
		return notification, &SocketError{Op: "decode", Err: err}
	}

	return notification, nil
}

// IsSocketClosed determines if an error indicates a closed connection
func IsSocketClosed(err error) bool {
	if err == nil {
		return false
	}

	// Unwrap socket error if needed
	var socketErr *SocketError
	if errors.As(err, &socketErr) {
		err = socketErr.Err
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return false
	}

	return errors.Is(err, net.ErrClosed) ||
		errors.Is(err, os.ErrClosed) ||
		errors.Is(err, io.EOF) ||
		strings.Contains(err.Error(), "use of closed network connection") ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "broken pipe")
}

func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	// Unwrap socket error if needed
	var socketErr *SocketError
	if errors.As(err, &socketErr) {
		err = socketErr.Err
	}

	// Check if it's a timeout error
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}

	return false
}

// ValidateSocketPath checks if the socket path is valid
func ValidateSocketPath(path string) error {
	// Check path length - Unix domain sockets have path length limitations
	// typically around 104-108 characters depending on the OS
	if len(path) > 100 {
		return &SocketError{
			Op:  "validate",
			Err: errors.New("socket path too long"),
		}
	}

	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return &SocketError{
			Op:  "validate",
			Err: errors.New("socket path must be absolute"),
		}
	}

	return nil
}
