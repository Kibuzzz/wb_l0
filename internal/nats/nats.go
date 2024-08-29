package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type Connection struct {
	*nats.Conn
}

func Connect() (*Connection, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect nats")
	}
	return &Connection{conn}, nil
}

func (c *Connection) Close() {
	c.Conn.Close()
}
