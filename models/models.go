package models

import "net"

// Chat interface to initialize, send and receive messages
type Chat interface {
	Send([]byte) error
	Receive() ([]byte, error)
	Init() error
}

// Client struct used by chat users
type Client struct {
	Key  []byte
	Conn *net.TCPConn
}

type Message struct {
	text string
}
