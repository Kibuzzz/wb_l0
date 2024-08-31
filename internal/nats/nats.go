package nats

import (
	"fmt"
	"wb_l0/internal/config"

	"github.com/nats-io/stan.go"
)

func Conn(conf config.Config) (stan.Conn, error) {
	conn, err := stan.Connect(conf.NATS.ClusterID, conf.NATS.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed connet to nats: %w", err)
	}
	return conn, nil
}
