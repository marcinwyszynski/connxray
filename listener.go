package connxray

import (
	"net"
	"time"
)

// Listener wraps a net.Listener and presents the same interface while allowing
// callback functions to be injected that will be called as the underlying
// net.Listener calls are invoked - with both their arguments and return values
// as callback arguments. Callbacks can also be specified that are passed to
// connections created by 'Accept' calls. All callbacks are void functions.
type Listener struct {
	// Underlying net.Listener.
	NetListener net.Listener

	// AcceptCallback is called every time underlying net.Listener's
	// 'Accept' method is called, with a Conn (wrapping the returned
	// net.Conn) and an error as arguments.
	AcceptCallback func(net.Listener, *Conn, error)

	// CloseCallback is called every time underlying net.Listener's 'Close'
	// method is called, with its return value (error) as argument.
	CloseCallback func(net.Listener, error)

	// AddrCallback is called every time underlying net.Listener's 'Addr'
	// method is called, with its return value (error) as argument.
	AddrCallback func(net.Listener, net.Addr)

	// ConnReadCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time the
	// underlying net.Conn's 'Read' method is called, with its argument
	// ([]byte) and return values (int and error) as arguments.
	ConnReadCallback func(net.Conn, []byte, int, error)

	// ConnWriteCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time the
	// underlying net.Conn's 'Write' method is called, with its argument
	// ([]byte) and return values (int and error) as arguments.
	ConnWriteCallback func(net.Conn, []byte, int, error)

	// ConnCloseCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time the
	// underlying net.Conn's 'Close' method is called, with its its return
	// value (error) as an argument.
	ConnCloseCallback func(net.Conn, error)

	// ConnLocalAddrCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time the
	// underlying net.Conn's 'LocalAddr' method is called, with its its
	// return value (net.Addr) as an argument.
	ConnLocalAddrCallback func(net.Conn, net.Addr)

	// ConnRemoteAddrCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time the
	// underlying net.Conn's 'RemoteAddr' method is called, with its its
	// return value (net.Addr) as an argument.
	ConnRemoteAddrCallback func(net.Conn, net.Addr)

	// ConnRemoteAddrCallback applies to a net.Conn object created by the
	// underlying net Listener's 'Accept' call. It is called every time
	// underlying net.Conn's 'SetDeadline' method is called, with its
	// argument (time.Time) and return value (error) as arguments.
	ConnSetDeadlineCallback func(net.Conn, time.Time, error)

	// ConnSetReadDeadlineCallback applies to a net.Conn object created by
	// the underlying net Listener's 'Accept' call. It is called every time
	// underlying net.Conn's 'SetReadDeadline' method is called, with its
	// argument (time.Time) and return value (error) as arguments.
	ConnSetReadDeadlineCallback func(net.Conn, time.Time, error)

	// ConnSetWriteDeadlineCallback applies to a net.Conn object created by
	// the underlying net Listener's 'Accept' call. It is called every time
	// underlying net.Conn's 'SetWriteDeadline' method is called, with its
	// argument (time.Time) and return value (error) as arguments.
	ConnSetWriteDeadlineCallback func(net.Conn, time.Time, error)
}

func (l *Listener) Accept() (net.Conn, error) {
	netconn, err := l.NetListener.Accept()
	conn := &Conn{
		NetConn:                  netconn,
		ReadCallback:             l.ConnReadCallback,
		WriteCallback:            l.ConnWriteCallback,
		CloseCallback:            l.ConnCloseCallback,
		LocalAddrCallback:        l.ConnLocalAddrCallback,
		RemoteAddrCallback:       l.ConnRemoteAddrCallback,
		SetDeadlineCallback:      l.ConnSetDeadlineCallback,
		SetReadDeadlineCallback:  l.ConnSetReadDeadlineCallback,
		SetWriteDeadlineCallback: l.ConnSetWriteDeadlineCallback,
	}
	if l.AcceptCallback != nil {
		defer l.AcceptCallback(l.NetListener, conn, err)
	}
	return conn, err
}

func (l *Listener) Close() error {
	err := l.NetListener.Close()
	if l.CloseCallback != nil {
		defer l.CloseCallback(l.NetListener, err)
	}
	return err
}

func (l *Listener) Addr() net.Addr {
	addr := l.NetListener.Addr()
	if l.AddrCallback != nil {
		defer l.AddrCallback(l.NetListener, addr)
	}
	return addr
}
