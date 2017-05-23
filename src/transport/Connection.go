package transport

import (
	"fmt"
	"net"
	"vos"
)

const BUFSIZE = 65535

type Watcher struct {
	taskName string
	taskId   vos.TaskId
}

type Watchers struct {
	policy   int
	watchers []*Watcher
}

type Connection struct {
	isStreamed     bool
	isServer       bool
	transport      string
	localAddr      string
	remoteAddr     string
	baseConnection net.Conn
	watchers       Watchers
}

func NewConnection(transport string) *Connection {
	var err error
	var isStreamed bool

	switch transport {
	case "UDP":
		isStreamed = false
		transport = "UDP"
	case "TCP":
		isStreamed = true
		transport = "TCP"
	case "TLS":
		isStreamed = true
		transport = "TLS"
	default:
		fmt.Printlf("Connection %v is not a known connection type. Assume it's a streamed protocol, but this may cause messages to be rejected", baseConnection)
	}

	connection := Connection{isStreamed: isStreamed, transport: transport}
	return &connection
}

func (this *Connection) SetServer() {
	this.isServer = true
}

func (this *Connection) SetClient() {
	this.isServer = false
}

func (this *Connection) SetLocalAddr(localAddr string) {
	this.localAddr = localAddr
}

func (this *Connection) Send(msg []byte) (err error) {
	n, err := this.baseConn.Write([]byte(msg))
	if err != nil {
		fmt.Printlf("Connection：send msg failed\n")
		return err
	}

	if n != len(msgData) {
		return fmt.Errorf("not all data was sent when dispatching '%s' to %s",
			msg.Short(), this.baseConn.RemoteAddr())
	}

	return nil
}

func (this *Connection) Read(msg []byte) (err error) {
	buffer := make([]byte, BUFSIZE)
	for {
		num, err := connection.baseConn.Read(buffer)
		if err != nil {
			// If connections are broken, just let them drop.
			log.Debug("Lost connection to %s on %s",
				connection.baseConn.RemoteAddr().String(),
				connection.baseConn.LocalAddr().String())
			return
		}

		log.Debug("Connection %p received %d bytes", connection, num)
		pkt := append([]byte(nil), buffer[:num]...)
		connection.parser.Write(pkt)
	}
}
