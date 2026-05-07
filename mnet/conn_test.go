package mnet

import (
	"net"
	"testing"
	"time"

	"github.com/Hucaru/Valhalla/mpacket"
)

type stubAddr string

func (a stubAddr) Network() string { return string(a) }
func (a stubAddr) String() string  { return string(a) }

type stubConn struct {
	closed bool
}

func (c *stubConn) Read([]byte) (int, error)         { return 0, nil }
func (c *stubConn) Write(b []byte) (int, error)      { return len(b), nil }
func (c *stubConn) Close() error                     { c.closed = true; return nil }
func (c *stubConn) LocalAddr() net.Addr              { return stubAddr("local") }
func (c *stubConn) RemoteAddr() net.Addr             { return stubAddr("remote") }
func (c *stubConn) SetDeadline(time.Time) error      { return nil }
func (c *stubConn) SetReadDeadline(time.Time) error  { return nil }
func (c *stubConn) SetWriteDeadline(time.Time) error { return nil }

func TestServerReaderHandlesPartialTCPReads(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	eRecv := make(chan *Event, 4)
	go serverReader(serverConn, eRecv, 2)

	parts := [][]byte{{3}, {0}, {'a', 'b'}, {'c'}}
	for _, part := range parts {
		if _, err := clientConn.Write(part); err != nil {
			t.Fatalf("write failed: %v", err)
		}
	}

	select {
	case ev := <-eRecv:
		if ev.Type != MEServerConnected {
			t.Fatalf("expected connect event, got %d", ev.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for connect event")
	}

	select {
	case ev := <-eRecv:
		if ev.Type != MEServerPacket {
			t.Fatalf("expected packet event, got %d", ev.Type)
		}
		if got := string(ev.Packet); got != "abc" {
			t.Fatalf("expected payload %q, got %q", "abc", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for packet event")
	}
}

func TestBaseConnSendDisconnectsWhenQueueIsFull(t *testing.T) {
	conn := &stubConn{}
	bc := &baseConn{
		Conn:  conn,
		eSend: make(chan mpacket.Packet, 1),
	}

	bc.eSend <- mpacket.NewPacket()
	bc.Send(mpacket.NewPacket())

	if !bc.closed {
		t.Fatal("expected send queue backpressure to close the connection")
	}
	if !conn.closed {
		t.Fatal("expected underlying connection to be closed")
	}
}
