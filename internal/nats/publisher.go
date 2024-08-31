package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"wb_l0/internal/models"

	"github.com/nats-io/stan.go"
)

func NewPub(sc stan.Conn, ctx context.Context) {
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			order := models.GenerateRandomOrder()
			bytes, err := json.Marshal(order)
			if err != nil {
				log.Println("error marshalling order: %w", err)
			}
			err = sc.Publish("orders", bytes)
			if err != nil {
				log.Fatal("failed to publish", err.Error())
			} else {
				log.Printf("successfully published: %s\n", order.OrderUID)
			}
		case <-ctx.Done():
			log.Println("publisher shutdowned")
			return
		}
	}
}
