// Package netutil provides utilities for network operations.
package netutil

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

// TimeoutConn is a wrapper around a net.Conn that applies a timeout
// to each Read and Write operation.
type TimeoutConn struct {
	conn    net.Conn
	timeout time.Duration
}

// NewTimeoutConn creates a new TimeoutConn.
// A timeout of 0 means no deadline will be set.
func NewTimeoutConn(conn net.Conn, timeout time.Duration) *TimeoutConn {
	return &TimeoutConn{
		conn:    conn,
		timeout: timeout,
	}
}

// Read reads data from the connection, setting a read deadline first.
// It reads up to len(p) bytes and returns the number of bytes read.
// It may return fewer bytes than requested without an error.
func (tc *TimeoutConn) Read(p []byte) (n int, err error) {
	if tc.timeout > 0 {
		if err := tc.conn.SetReadDeadline(time.Now().Add(tc.timeout)); err != nil {
			return 0, fmt.Errorf("failed to set read deadline: %w", err)
		}
	}
	return tc.conn.Read(p)
}

// ReadFull reads exactly len(p) bytes from the connection, setting a read deadline first.
// It is a convenience wrapper around io.ReadFull.
func (tc *TimeoutConn) ReadFull(p []byte) (n int, err error) {
	if tc.timeout > 0 {
		if err := tc.conn.SetReadDeadline(time.Now().Add(tc.timeout)); err != nil {
			return 0, fmt.Errorf("failed to set read deadline: %w", err)
		}
	}
	return io.ReadFull(tc.conn, p)
}

// Write writes data to the connection, setting a write deadline first.
func (tc *TimeoutConn) Write(p []byte) (n int, err error) {
	if tc.timeout > 0 {
		if err := tc.conn.SetWriteDeadline(time.Now().Add(tc.timeout)); err != nil {
			return 0, fmt.Errorf("failed to set write deadline: %w", err)
		}
	}
	return tc.conn.Write(p)
}

// Close closes the underlying connection.
func (tc *TimeoutConn) Close() error {
	return tc.conn.Close()
}

// LocalAddr returns the local network address.
func (tc *TimeoutConn) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (tc *TimeoutConn) RemoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}

// --- Connection Helpers ---

// DialTCP connects to a TCP address. The returned net.Conn has no timeouts by default.
// Wrap it with NewTimeoutConn to enforce timeouts on I/O operations.
func DialTCP(ipAddress string, port int, dialTimeout time.Duration) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", ipAddress, port)
	dialer := net.Dialer{Timeout: dialTimeout}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", address, err)
	}
	return conn, nil
}

// DialUnix connects to a Unix domain socket. The returned net.Conn has no timeouts by default.
// Wrap it with NewTimeoutConn to enforce timeouts on I/O operations.
func DialUnix(path string, dialTimeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "unix", path)
	if err != nil {
		return nil, fmt.Errorf("unable to dial unix socket at %s: %w", path, err)
	}
	return conn, nil
}
