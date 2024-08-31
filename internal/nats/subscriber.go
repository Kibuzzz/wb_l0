package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb_l0/internal/models"
	"wb_l0/internal/repository"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	sub  stan.Subscription
	repo repository.Repo
}

func NewSub(conn stan.Conn, subj string, repo repository.Repo) error {

	subsc := &Subscriber{repo: repo}
	sub, err := conn.Subscribe(subj, func(msg *stan.Msg) {
		log.Println("recieved message")
		subsc.recieveMessage(msg.Data)

	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	subsc.sub = sub

	return nil
}

func (s *Subscriber) recieveMessage(msg []byte) error {

	var order models.Order
	err := json.Unmarshal(msg, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshall: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = s.repo.AddOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to add order to repo: %w", err)
	}

	return nil
}

func (s *Subscriber) Close() error {
	return s.sub.Close()
}
