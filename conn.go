// Package connxray provides net.Listener and net.Conn introspection by wrapping
// around these objects and provide callbacks which can be invoked for any
// method of the unerlying object. net.Conn callbacks can be passed to a
// net.Listener which will be passed onwards to all net.Conn objects that are
// created as a result of net.Listener#Accept calls.
//
// Also of note is the fact that connxray.Listener#Accept creates an instance
// of connxray.Conn whose callbacks can be set dynamically at any stage of the
// connection lifetime. This means that only certain net.Conn objects can be
// explicitly monitored (eg. sampling) or that monitoring behavior can be
// adjusted on the fly.
package connxray

import (
	"net"
	"time"
)

// Conn wraps a net.Conn and presents the same interface while allowing
// callback functions to be injected that will be called as the underlying
// net.Conn calls are invoked - with both their arguments and return values
// as callback arguments. All callbacks are void functions.
type Conn struct {
	// Underlying net.Conn.
	NetConn net.Conn

	// ReadCallback is called every time the underlying net.Conn's 'Read'
	// method is called, with its argument ([]byte) and return values (int
	// and error) as arguments.
	ReadCallback func(net.Conn, []byte, int, error)

	// WriteCallback is called every time the underlying net.Conn's 'Write'
	// method is called, with its argument ([]byte) and return values (int
	// and error) as arguments.
	WriteCallback func(net.Conn, []byte, int, error)

	// CloseCallback is called every time the underlying net.Conn's 'Close'
	// method is called, with its its return value (error) as an argument.
	CloseCallback func(net.Conn, error)

	// LocalAddrCallback is called every time the underlying net.Conn's
	// 'LocalAddr' method is called, with its its return value (net.Addr) as
	// an argument.
	LocalAddrCallback func(net.Conn, net.Addr)

	// RemoteAddrCallback is called every time the underlying net.Conn's
	// 'RemoteAddr' method is called, with its its return value (net.Addr)
	// as an argument.
	RemoteAddrCallback func(net.Conn, net.Addr)

	// RemoteAddrCallback is called every time underlying net.Conn's
	// 'SetDeadline' method is called, with its argument (time.Time) and
	// return value (error) as arguments.
	SetDeadlineCallback func(net.Conn, time.Time, error)

	// ReadDeadlineCallback is called every time underlying net.Conn's
	// 'SetReadDeadline' method is called, with its argument (time.Time) and
	// return value (error) as arguments.
	SetReadDeadlineCallback func(net.Conn, time.Time, error)

	// SetWriteDeadlineCallback is called every time underlying net.Conn's
	// 'SetWriteDeadline' method is called, with its argument (time.Time)
	// and return value (error) as arguments.
	SetWriteDeadlineCallback func(net.Conn, time.Time, error)
}

// Read reads from the underlying net.Conn and runs a ReadCallback if one was
// specified.
func (c *Conn) Read(b []byte) (int, error) {
	n, err := c.NetConn.Read(b)
	if c.ReadCallback != nil {
		defer c.ReadCallback(c.NetConn, b, n, err)
	}
	return n, err
}

// Write writes to the underlying net.Conn and runs a WriteCallback if one was
// specified.
func (c *Conn) Write(b []byte) (int, error) {
	n, err := c.NetConn.Write(b)
	if c.WriteCallback != nil {
		defer c.WriteCallback(c.NetConn, b, n, err)
	}
	return n, err
}

// Close closes the underlying net.Conn and runs a CloseCallback if one was
// specified.
func (c *Conn) Close() error {
	err := c.NetConn.Close()
	if c.CloseCallback != nil {
		defer c.CloseCallback(c.NetConn, err)
	}
	return err
}

// LocalAddr gets the local address from the underlying net.Conn and runs a
// LocalAddrCallback if one was specified.
func (c *Conn) LocalAddr() net.Addr {
	addr := c.NetConn.LocalAddr()
	if c.LocalAddrCallback != nil {
		defer c.LocalAddrCallback(c.NetConn, addr)
	}
	return addr
}

// RemoteAddr gets the remote address from the underlying net.Conn and runs a
// RemoteAddrCallback if one was specified.
func (c *Conn) RemoteAddr() net.Addr {
	addr := c.NetConn.RemoteAddr()
	if c.RemoteAddrCallback != nil {
		defer c.RemoteAddrCallback(c.NetConn, addr)
	}
	return addr
}

// SetDeadline sets a deadline on the underlying net.Conn and runs a
// SetDeadlineCallback if one was specified.
func (c *Conn) SetDeadline(t time.Time) error {
	err := c.NetConn.SetDeadline(t)
	if c.SetDeadlineCallback != nil {
		defer c.SetDeadlineCallback(c.NetConn, t, err)
	}
	return err
}

// SetReadDeadline sets a read deadline on the underlying net.Conn and runs a
// SetReadDeadlineCallback if one was specified.
func (c *Conn) SetReadDeadline(t time.Time) error {
	err := c.NetConn.SetReadDeadline(t)
	if c.SetReadDeadlineCallback != nil {
		defer c.SetReadDeadlineCallback(c.NetConn, t, err)
	}
	return err
}

// SetWriteDeadline sets a write deadline on the underlying net.Conn and runs a
// SetWriteDeadlineCallback if one was specified.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	err := c.NetConn.SetWriteDeadline(t)
	if c.SetWriteDeadlineCallback != nil {
		defer c.SetWriteDeadlineCallback(c.NetConn, t, err)
	}
	return err
}
