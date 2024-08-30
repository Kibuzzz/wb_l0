package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

func Conn() (stan.Conn, error) {
	conn, err := stan.Connect("orders", "test_1")
	if err != nil {
		return nil, fmt.Errorf("failed connet to nats: %w", err)
	}
	return conn, nil
}
