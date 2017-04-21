package transport

import (
	"fmt"
	"net"
	"vos"
)

type Connection struct {
	isStreamed bool
	isServer   bool
	transport  string
	localAddr  string
	remoteAddr string

	taskName string
	taskId   vos.TaskId

	baseConnection net.Conn
}

func NewConnection(baseConnection net.Conn, taskName string) *Connection {
	var err error

	var isStreamed bool
	var transport string

	switch baseConnection.(type) {
	case *net.UDPConn:
		isStreamed = false
		transport = "UDP"
	case *net.TCPConn:
		isStreamed = true
		transport = "TCP"
	case *tls.Conn:
		isStreamed = true
		transport = "TLS"
	default:
		fmt.Printlf("Connection %v is not a known connection type. Assume it's a streamed protocol, but this may cause messages to be rejected", baseConnection)
	}

	connection := Connection{baseConnection: baseConnection, isStreamed: isStreamed, transport: transport}
	connection.taskName = taskName
	connection.taskId, _, err = vos.FindTask(taskName)
	if err != nil {
		return nil
	}

	return &connection
}

func (connection *Connection) Send(msg []byte) (err error) {
	n, err := connection.baseConn.Write([]byte(msg))
	if err != nil {
		fmt.Printlf("Connection：send msg failed\n")
		return err
	}

	if n != len(msgData) {
		return fmt.Errorf("not all data was sent when dispatching '%s' to %s",
			msg.Short(), connection.baseConn.RemoteAddr())
	}

	return nil
}
