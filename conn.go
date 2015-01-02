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
//
// Last but not least connxray.Conn implements net.PacketConn so can be used
// with any code that expects one (eg. golang.org/x/net/ipv[46]). If the
// underlying connection object does not implement net.PacketConn a releant
// error (ErrNotPacketConn) will be returned by ReadFrom and WriteTo methodes
// as well as passed to their respective callbacks, if any were specified.
package connxray

import (
	"fmt"
	"net"
	"time"
)

var (
	// ErrNotPacketConn signifies that the underlying net.Conn is does not
	// implement the net.PacketConn interface.
	ErrNotPacketConn = fmt.Errorf("this net.Conn is not a net.PacketConn")
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

	// ReadFromCallback is called every time the underlying net.PacketConn's
	// 'ReadFrom' method is called, with its argument ([]byte) and return
	// values (int, net.Addr and error) as arguments.
	ReadFromCallback func(net.Conn, []byte, int, net.Addr, error)

	// WriteCallback is called every time the underlying net.Conn's 'Write'
	// method is called, with its argument ([]byte) and return values (int
	// and error) as arguments.
	WriteCallback func(net.Conn, []byte, int, error)

	// WriteToCallback is called every time the underlying net.PacketConn's
	// 'WriteTo' method is called, with its arguments ([]byte and net.Addr)
	// and return values (int and error) as arguments.
	WriteToCallback func(net.Conn, []byte, net.Addr, int, error)

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

// ReadFrom reads from the underlying net.PacketConn and runs a ReadFromCallback
// if one was specified.
func (c *Conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	pconn, implements := c.NetConn.(net.PacketConn)
	if !implements {
		err = ErrNotPacketConn
	} else {
		n, addr, err = pconn.ReadFrom(b)
	}
	if c.ReadFromCallback != nil {
		defer c.ReadFromCallback(c.NetConn, b, n, addr, err)
	}
	return n, addr, err
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

// WriteTo writes to the underlying net.PacketConn and runs a WriteToCallback
// if one was specified.
func (c *Conn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	pconn, implements := c.NetConn.(net.PacketConn)
	if !implements {
		err = ErrNotPacketConn
	} else {
		n, err = pconn.WriteTo(b, addr)
	}
	if c.WriteToCallback != nil {
		defer c.WriteToCallback(c.NetConn, b, addr, n, err)
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
