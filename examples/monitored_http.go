// In this example HTTP traffic is inspected by introspecting on an underlying
// TCP acceptor. By injecting callbacks on Accept, Read, Write and Close we can
// track stats for each individual connection as it changes state.
//
// This is only one possible use case of the connxray library.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	xray "github.com/marcinwyszynski/connxray"
)

var (
	port = flag.Int("port", 1983, "HTTP port")
)

type ReadCallback func(*xray.Conn, []byte, int, error)
type WriteCallback func(*xray.Conn, []byte, int, error)
type CloseCallback func(*xray.Conn, error)

type stats struct {
	bytesRead    int
	bytesWritten int
	startTime    time.Time
}

func onAccept(_ *xray.Listener, conn *xray.Conn, err error) {
	s := &stats{startTime: time.Now()}
	conn.AfterRead = onRead(s)
	conn.AfterWrite = onWrite(s)
	conn.AfterClose = onClose(s)
	if err != nil {
		glog.Errorf("Error establishing connection: %v", err)
		return
	}
	glog.Infof("%s <-> %s started", conn.LocalAddr(), conn.RemoteAddr())
}

func onRead(s *stats) ReadCallback {
	return func(_ *xray.Conn, _ []byte, n int, _ error) {
		s.bytesRead += n
	}
}

func onWrite(s *stats) WriteCallback {
	return func(_ *xray.Conn, _ []byte, n int, _ error) {
		s.bytesWritten += n
	}
}

func onClose(s *stats) CloseCallback {
	return func(conn *xray.Conn, _ error) {
		msg := "%s closed: %d bytes read, %d bytes written in %d ms"
		glog.Infof(
			msg,
			conn.RemoteAddr(),
			s.bytesRead,
			s.bytesWritten,
			time.Since(s.startTime)/1e6,
		)
	}
}

func main() {
	flag.Parse()
	addr := net.TCPAddr{Port: *port}
	tcpLisetner, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		glog.Fatalf("Error creating a TCP listener: %v", err)
	}
	introspectedListener := &xray.Listener{
		Base:        tcpLisetner,
		AfterAccept: onAccept,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})
	glog.Infof("About to start serving on %s", addr.String())
	glog.Fatal(http.Serve(introspectedListener, nil))
}
