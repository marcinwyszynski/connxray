package connxray

import (
	"net"
)

// Listener wraps a net.Listener and presents the same interface while allowing
// hook functions to be injected that will be called before and/or after
// the underlying net.Listener calls are invoked. Please see the package
// top-level documentation for more information about hooks.
type Listener struct {
	// Underlying net.Listener.
	Base net.Listener

	// BeforeAccept is a 'before' hook for the Accept method. If it returns
	// an error neither the base method nor the 'after' callback will be
	// called.
	BeforeAccept func(*Listener) error

	// AfterAccept is an 'after' hook for the Accept method.
	AfterAccept func(*Listener, *Conn, error)

	// BeforeClose is a 'before' hook for the Close method.
	BeforeClose func(*Listener) error

	// AfterClose is an 'after' hook for the Close method.
	AfterClose func(*Listener, error)

	// AfterAddr is an 'after' hook for the Addr method.
	AfterAddr func(*Listener, net.Addr)
}

// Accept runs Accept on the underlying net.Listener plus any relevant hooks
// ('before' and 'after') that were set up.
func (l *Listener) Accept() (net.Conn, error) {
	if l.BeforeAccept != nil {
		if err := l.BeforeAccept(l); err != nil {
			return nil, err
		}
	}
	netconn, err := l.Base.Accept()
	conn := &Conn{Base: netconn}
	if l.AfterAccept != nil {
		defer l.AfterAccept(l, conn, err)
	}
	return conn, err
}

// Close runs Close on the underlying net.Listener plus any relevant hooks
// ('before' and 'after') that were set up.
func (l *Listener) Close() error {
	if l.BeforeClose != nil {
		if err := l.BeforeClose(l); err != nil {
			return err
		}
	}
	err := l.Base.Close()
	if l.AfterClose != nil {
		defer l.AfterClose(l, err)
	}
	return err
}

// Addr runs Addr on the underlying net.Listener plus an 'after' hook if it
// was set up.
func (l *Listener) Addr() net.Addr {
	addr := l.Base.Addr()
	if l.AfterAddr != nil {
		defer l.AfterAddr(l, addr)
	}
	return addr
}
