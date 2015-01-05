package connxray

import (
	"errors"
	"net"
	"testing"
)

func TestAcceptWithSucceedingBeforeCallback(t *testing.T) {
	baseCalled, beforeCalled, afterCalled := false, false, false
	expErr := errors.New("chunky bacon")
	ml := &mockListener{
		acceptHandler: func() (net.Conn, error) {
			if !beforeCalled {
				t.Error("Before callback not invoked")
			}
			baseCalled = true
			return nil, expErr
		},
	}
	cl := &Listener{
		Base: ml,
		BeforeAccept: func(_ *Listener) error {
			beforeCalled = true
			return nil
		},
		AfterAccept: func(_ *Listener, _ *Conn, err error) {
			if !baseCalled {
				t.Error("Base method not invoked")
			}
			if err != expErr {
				t.Errorf(
					"Unexpected error %v, expected %v",
					err,
					expErr,
				)
			}
			afterCalled = true
		},
	}
	if _, err := cl.Accept(); err != expErr {
		t.Errorf("Unexpected error %v, expected %v", err, expErr)
	}
	if !afterCalled {
		t.Error("After callback not invoked")
	}
}

func TestAcceptWithFailingBeforeCallback(t *testing.T) {
	baseCalled, beforeCalled, afterCalled := false, false, false
	expErr := errors.New("chunky bacon")
	ml := &mockListener{
		acceptHandler: func() (net.Conn, error) {
			if !beforeCalled {
				t.Error("Before callback not invoked")
			}
			baseCalled = true
			return nil, nil
		},
	}
	cl := &Listener{
		Base: ml,
		BeforeAccept: func(_ *Listener) error {
			beforeCalled = true
			return expErr
		},
		AfterAccept: func(_ *Listener, _ *Conn, _ error) {
			afterCalled = true
		},
	}
	if _, err := cl.Accept(); err != expErr {
		t.Errorf("Unexpected error %v, expected %v", err, expErr)
	}
	if baseCalled {
		t.Error("Base method invoked")
	}
	if afterCalled {
		t.Error("After callback invoked")
	}
}

func TestCloseWithSucceedingBeforeCallback(t *testing.T) {
	baseCalled, beforeCalled, afterCalled := false, false, false
	expErr := errors.New("chunky bacon")
	ml := &mockListener{
		closeHandler: func() error {
			if !beforeCalled {
				t.Error("Before callback not invoked")
			}
			baseCalled = true
			return expErr
		},
	}
	cl := &Listener{
		Base: ml,
		BeforeClose: func(_ *Listener) error {
			beforeCalled = true
			return nil
		},
		AfterClose: func(_ *Listener, err error) {
			if !baseCalled {
				t.Error("Base method not invoked")
			}
			if err != expErr {
				t.Errorf(
					"Unexpected error %v, expected %v",
					err,
					expErr,
				)
			}
			afterCalled = true
		},
	}
	if err := cl.Close(); err != expErr {
		t.Errorf("Unexpected error %v, expected %v", err, expErr)
	}
	if !afterCalled {
		t.Error("After callback not invoked")
	}
}

func TestCloseWithFailingBeforeCallback(t *testing.T) {
	baseCalled, beforeCalled, afterCalled := false, false, false
	expErr := errors.New("chunky bacon")
	ml := &mockListener{
		closeHandler: func() error {
			if !beforeCalled {
				t.Error("Before callback not invoked")
			}
			baseCalled = true
			return nil
		},
	}
	cl := &Listener{
		Base: ml,
		BeforeClose: func(_ *Listener) error {
			beforeCalled = true
			return expErr
		},
		AfterClose: func(_ *Listener, _ error) {
			afterCalled = true
		},
	}
	if err := cl.Close(); err != expErr {
		t.Errorf("Unexpected error %v, expected %v", err, expErr)
	}
	if baseCalled {
		t.Error("Base method invoked")
	}
	if afterCalled {
		t.Error("After callback invoked")
	}
}

func TestAfterAddr(t *testing.T) {
	baseCalled, afterCalled := false, false
	expAddr, _ := net.ResolveTCPAddr("tcp", "localhost:80")
	ml := &mockListener{
		addrHandler: func() net.Addr {
			baseCalled = true
			return expAddr
		},
	}
	cl := &Listener{
		Base: ml,
		AfterAddr: func(_ *Listener, addr net.Addr) {
			if addr != expAddr {
				t.Errorf(
					"Unxpected address: %v, expected %v",
					addr,
					expAddr,
				)
			}
			afterCalled = true
		},
	}
	if addr := cl.Addr(); addr != expAddr {
		t.Errorf("Unxpected address: %v, expected %v", addr, expAddr)
	}
	if !baseCalled {
		t.Error("Base method not invoked")
	}
	if !afterCalled {
		t.Error("After callback not invoked")
	}
}
