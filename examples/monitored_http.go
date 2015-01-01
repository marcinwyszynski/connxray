// In this example HTTP traffic is inspected by introspecting on an underlying
// TCP acceptor. By injecting callbacks on Accept, Read, Write and Close we can
// track stats for each individual connection as it changes state.
//
// This is only one possible use case of the connxray library.
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
	xray "github.com/marcinwyszynski/connxray"
)

var (
	port = flag.Int("port", 1983, "HTTP port")
)

type stats struct {
	bytesRead    int
	bytesWritten int
	startTime    time.Time
	guard        sync.Mutex
}

func newStats() *stats {
	return &stats{startTime: time.Now()}
}

type statsTracker struct {
	data      map[net.Conn]*stats
	dataGuard sync.Mutex
}

func newTracker() *statsTracker {
	return &statsTracker{data: make(map[net.Conn]*stats)}
}

func (s *statsTracker) onAccept(_ net.Listener, conn *xray.Conn, err error) {
	if err != nil {
		glog.Errorf("Error establishing connection: %v", err)
		return
	}
	glog.Infof("%s <-> %s started", conn.LocalAddr(), conn.RemoteAddr())
	s.dataGuard.Lock()
	defer s.dataGuard.Unlock()
	s.data[conn.NetConn] = newStats()
}

func (s *statsTracker) onRead(conn net.Conn, _ []byte, n int, err error) {
	if err != nil {
		glog.Errorf("Error reading from connection %#v %v", conn, err)
	}
	data, exists := s.data[conn]
	if !exists {
		glog.Errorf("Connection not tracked: %#v", conn)
		return
	}
	data.guard.Lock()
	defer data.guard.Unlock()
	data.bytesRead += n
}

func (s *statsTracker) onWrite(conn net.Conn, _ []byte, n int, err error) {
	if err != nil {
		glog.Errorf("Error writing to connection %#v %v", conn, err)
	}
	data, exists := s.data[conn]
	if !exists {
		glog.Errorf("Connection not tracked: %#v", conn)
		return
	}
	data.guard.Lock()
	defer data.guard.Unlock()
	data.bytesWritten += n
}

func (s *statsTracker) onClose(conn net.Conn, err error) {
	if err != nil {
		glog.Errorf("Error closing connection %#v %v", conn, err)
	}
	data, exists := s.data[conn]
	if !exists {
		glog.Errorf("Connection not tracked: %#v", conn)
		return
	}
	glog.Infof(
		"%s <-> %s closed: %d bytes read, %d bytes written in %d ms",
		conn.LocalAddr(),
		conn.RemoteAddr(),
		data.bytesRead,
		data.bytesWritten,
		time.Since(data.startTime)/1e6,
	)
	data.guard.Lock()
	defer data.guard.Unlock()
	delete(s.data, conn)
}

func main() {
	flag.Parse()
	addr := net.TCPAddr{Port: *port}
	tcpLisetner, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		glog.Fatalf("Error creating a TCP listener: %v", err)
	}
	tracker := newTracker()
	introspectedListener := &xray.Listener{
		NetListener:       tcpLisetner,
		AcceptCallback:    tracker.onAccept,
		ConnReadCallback:  tracker.onRead,
		ConnWriteCallback: tracker.onWrite,
		ConnCloseCallback: tracker.onClose,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})
	glog.Infof("About to start serving on %s", addr.String())
	glog.Fatal(http.Serve(introspectedListener, nil))
}
