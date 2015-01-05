// Package connxray provides net.Listener and net.Conn introspection by wrapping
// around these objects and provide hooks which can be invoked for any
// method of the unerlying object. net.Conn hooks can be passed to a
// net.Listener which will be passed onwards to all net.Conn objects that are
// created as a result of net.Listener#Accept calls.
//
// Also of note is the fact that connxray.Listener#Accept creates an instance
// of connxray.Conn whose hooks can be set dynamically at any stage of the
// connection lifetime. This means that only certain net.Conn objects can be
// explicitly monitored (eg. sampling) or that monitoring behavior can be
// adjusted on the fly.
//
// Last but not least connxray.Conn implements net.PacketConn so can be used
// with any code that expects one (eg. golang.org/x/net/ipv[46]). If the
// underlying connection object does not implement net.PacketConn a releant
// error (ErrNotPacketConn) will be returned by ReadFrom and WriteTo methodes
// as well as passed to their respective hooks, if any were specified.
package connxray

import (
	"errors"
	"net"
	"time"
)

var (
	// ErrNotPacketConn signifies that the underlying net.Conn is does not
	// implement the net.PacketConn interface.
	ErrNotPacketConn = errors.New("this net.Conn is not a net.PacketConn")
)

// Conn wraps a net.Conn and presents the same interface while allowing
// hook functions to be injected that will be called before and/or after
// the underlying net.Conn calls are invoked. Please see the package top-level
// documentation for more information about hooks.
type Conn struct {
	// Underlying net.Conn.
	Base net.Conn

	// BeforeRead is a 'before' hook for the Read method.
	BeforeRead func(*Conn, []byte) error

	// AfterRead is an 'after' hook for the Read method.
	AfterRead func(*Conn, []byte, int, error)

	// BeforeReadFrom is a 'before' hook for the ReadFrom method.
	BeforeReadFrom func(*Conn, []byte) error

	// AfterReadFrom is an 'after' hook for the ReadFrom method.
	AfterReadFrom func(*Conn, []byte, int, net.Addr, error)

	// BeforeWrite is a 'before' hook for the Write method.
	BeforeWrite func(*Conn, []byte) error

	// AfterWrite is an 'after' hook for the Write method.
	AfterWrite func(*Conn, []byte, int, error)

	// BeforeWriteTo is a 'before' hook for the WriteTo method.
	BeforeWriteTo func(*Conn, []byte, net.Addr) error

	// AfterWriteTo is an 'after' hook for the WriteTo method.
	AfterWriteTo func(*Conn, []byte, net.Addr, int, error)

	// BeforeClose is a 'before' hook for the Close method.
	BeforeClose func(*Conn) error

	// AfterClose is an 'after' hook for the Close method.
	AfterClose func(*Conn, error)

	// AfterLocalAddr is an 'after' hook for the LocalAddr method.
	AfterLocalAddr func(*Conn, net.Addr)

	// AfterRemoteAddr is an 'after' hook for the RemoteAddr method.
	AfterRemoteAddr func(*Conn, net.Addr)

	// BeforeSetDeadline is a 'before' hook for the SetDeadline method.
	BeforeSetDeadline func(*Conn, time.Time) error

	// AfterSetDeadline is an 'after' hook for the SetDeadline method.
	AfterSetDeadline func(*Conn, time.Time, error)

	// BeforeSetReadDeadline is a 'before' hook for the SetReadDeadline
	// method.
	BeforeSetReadDeadline func(*Conn, time.Time) error

	// AfterSetReadDeadline is an 'after' hook for the SetReadDeadline
	// method.
	AfterSetReadDeadline func(*Conn, time.Time, error)

	// BeforeSetWriteDeadline is a 'before' hook for the SetWriteDeadline
	// method.
	BeforeSetWriteDeadline func(*Conn, time.Time) error

	// AfterSetWriteDeadline is an 'after' hook for the SetWriteDeadline
	// method.
	AfterSetWriteDeadline func(*Conn, time.Time, error)
}

// Read reads from the underlying net.Conn and invokes relevant hooks
// ('before' and 'after') that were set up.
func (c *Conn) Read(b []byte) (int, error) {
	if c.BeforeRead != nil {
		if err := c.BeforeRead(c, b); err != nil {
			return 0, err
		}
	}
	n, err := c.Base.Read(b)
	if c.AfterRead != nil {
		defer c.AfterRead(c, b, n, err)
	}
	return n, err
}

// ReadFrom reads from the underlying net.PacketConn and invokes relevant hooks
// ('before' and 'after') that were set up.
func (c *Conn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	pconn, implements := c.Base.(net.PacketConn)
	if !implements {
		err = ErrNotPacketConn
		return
	}
	if c.BeforeReadFrom != nil {
		err = c.BeforeReadFrom(c, b)
	}
	if err != nil {
		return
	}
	n, addr, err = pconn.ReadFrom(b)
	if c.AfterReadFrom != nil {
		defer c.AfterReadFrom(c, b, n, addr, err)
	}
	return n, addr, err
}

// Write writes to the underlying net.Conn and invokes relevant hooks ('before'
// and 'after') that were set up.
func (c *Conn) Write(b []byte) (int, error) {
	if c.BeforeWrite != nil {
		if err := c.BeforeWrite(c, b); err != nil {
			return 0, err
		}
	}
	n, err := c.Base.Write(b)
	if c.AfterWrite != nil {
		defer c.AfterWrite(c, b, n, err)
	}
	return n, err
}

// WriteTo writes to the underlying net.PacketConn and invokes relevant hooks
// ('before' and 'after') that were set up.
func (c *Conn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	pconn, implements := c.Base.(net.PacketConn)
	if !implements {
		err = ErrNotPacketConn
		return
	}
	if c.BeforeWriteTo != nil {
		err = c.BeforeWriteTo(c, b, addr)
	}
	if err != nil {
		return
	}
	n, err = pconn.WriteTo(b, addr)
	if c.AfterWriteTo != nil {
		defer c.AfterWriteTo(c, b, addr, n, err)
	}
	return n, err
}

// Close closes the underlying net.Conn and invokes relevant hooks ('before'
// and 'after') that were set up.
func (c *Conn) Close() error {
	if c.BeforeClose != nil {
		if err := c.BeforeClose(c); err != nil {
			return err
		}
	}
	err := c.Base.Close()
	if c.AfterClose != nil {
		defer c.AfterClose(c, err)
	}
	return err
}

// LocalAddr gets the local address from the underlying net.Conn and invokes
// an 'after' hook if it was set up.
func (c *Conn) LocalAddr() net.Addr {
	addr := c.Base.LocalAddr()
	if c.AfterLocalAddr != nil {
		defer c.AfterLocalAddr(c, addr)
	}
	return addr
}

// RemoteAddr gets the remote address from the underlying net.Conn and invokes
// an 'after' hook if it was set up.
func (c *Conn) RemoteAddr() net.Addr {
	addr := c.Base.RemoteAddr()
	if c.AfterRemoteAddr != nil {
		defer c.AfterRemoteAddr(c, addr)
	}
	return addr
}

// SetDeadline sets a deadline on the underlying net.Conn and invokes relevant
// hooks ('before' and 'after') that were set up.
func (c *Conn) SetDeadline(t time.Time) error {
	if c.BeforeSetDeadline != nil {
		if err := c.BeforeSetDeadline(c, t); err != nil {
			return err
		}
	}
	err := c.Base.SetDeadline(t)
	if c.AfterSetDeadline != nil {
		defer c.AfterSetDeadline(c, t, err)
	}
	return err
}

// SetReadDeadline sets a read deadline on the underlying net.Conn and invokes
// relevant hooks ('before' and 'after') that were set up.
func (c *Conn) SetReadDeadline(t time.Time) error {
	if c.BeforeSetReadDeadline != nil {
		if err := c.BeforeSetReadDeadline(c, t); err != nil {
			return err
		}
	}
	err := c.Base.SetReadDeadline(t)
	if c.AfterSetReadDeadline != nil {
		defer c.AfterSetReadDeadline(c, t, err)
	}
	return err
}

// SetWriteDeadline sets a write deadline on the underlying net.Conn and invokes
// relevant hooks ('before' and 'after') that were set up.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	if c.BeforeSetWriteDeadline != nil {
		if err := c.BeforeSetWriteDeadline(c, t); err != nil {
			return err
		}
	}
	err := c.Base.SetWriteDeadline(t)
	if c.AfterSetWriteDeadline != nil {
		defer c.AfterSetWriteDeadline(c, t, err)
	}
	return err
}
