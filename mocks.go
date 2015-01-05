package connxray

import (
	"net"
	"time"
)

// mockConn is a mock implementation of net.Conn and net.PacketConn interfaces.
// This is generated manually since there standard mocking solutions like gomock
// do not handle mocking out standard library.
type mockConn struct {
	readHandler             func([]byte) (int, error)
	readFromHandler         func([]byte) (int, net.Addr, error)
	writeHandler            func([]byte) (int, error)
	writeToHandler          func([]byte, net.Addr) (int, error)
	closeHandler            func() error
	localAddrHandler        func() net.Addr
	remoteAddrHandler       func() net.Addr
	setDeadlineHandler      func(time.Time) error
	setReadDeadlineHandler  func(time.Time) error
	setWriteDeadlineHandler func(time.Time) error
}

func (c *mockConn) Read(b []byte) (int, error) {
	return c.readHandler(b)
}

func (c *mockConn) ReadFrom(b []byte) (int, net.Addr, error) {
	return c.readFromHandler(b)
}

func (c *mockConn) Write(b []byte) (int, error) {
	return c.writeHandler(b)
}

func (c *mockConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return c.writeToHandler(b, addr)
}

func (c *mockConn) Close() error {
	return c.closeHandler()
}

func (c *mockConn) LocalAddr() net.Addr {
	return c.localAddrHandler()
}

func (c *mockConn) RemoteAddr() net.Addr {
	return c.remoteAddrHandler()
}

func (c *mockConn) SetDeadline(t time.Time) error {
	return c.setDeadlineHandler(t)
}

func (c *mockConn) SetReadDeadline(t time.Time) error {
	return c.setReadDeadlineHandler(t)
}

func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return c.setWriteDeadlineHandler(t)
}

// mockListener is a mock implementation of net.Listener interface. This is
// generated manually since there standard mocking solutions like gomock do not
// handle mocking out standard library.
type mockListener struct {
	acceptHandler func() (net.Conn, error)
	closeHandler  func() error
	addrHandler   func() net.Addr
}

func (l *mockListener) Accept() (net.Conn, error) {
	return l.acceptHandler()
}

func (l *mockListener) Close() error {
	return l.closeHandler()
}

func (l *mockListener) Addr() net.Addr {
	return l.addrHandler()
}
