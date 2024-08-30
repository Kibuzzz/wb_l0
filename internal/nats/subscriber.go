package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"wb_l0/internal/models"
	"wb_l0/internal/repository"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	sub  stan.Subscription
	repo repository.Repo
}

func NewSub(conn stan.Conn, repo repository.Repo) (*Subscriber, error) {

	subsc := &Subscriber{repo: repo}
	sub, err := conn.Subscribe("orders", func(msg *stan.Msg) {
		fmt.Println("recieved message")
		subsc.recieveMessage(msg.Data)

	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	subsc.sub = sub

	return subsc, nil
}

func (s *Subscriber) recieveMessage(msg []byte) error {

	var order models.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshall: %w", err)
	}

	err = s.repo.AddOrder(context.TODO(), order)
	if err != nil {
		return fmt.Errorf("failed to add order to repo: %w", err)
	}

	return nil
}

func (s *Subscriber) Close() error {
	return s.sub.Close()
}
