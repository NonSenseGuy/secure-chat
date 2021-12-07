package chat

import (
	"net"
)

func (c *Chat) Init() error {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
}
