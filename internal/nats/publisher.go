package nats

import (
	"encoding/json"
	"log"
	"time"
	"wb_l0/internal/models"

	"github.com/nats-io/stan.go"
)

func NewPub(sc stan.Conn, doneCh chan bool) {
	// move ticker time to config
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			order := models.GenerateRandomOrder()
			bytes, err := json.Marshal(order)
			if err != nil {
				log.Println("error marshalling order: %w", err)
			}
			// move cluster id to config
			err = sc.Publish("orders", bytes)
			if err != nil {
				log.Fatal("failed to publish", err.Error())
			} else {
				log.Println("successfully published")
			}
		case <-doneCh:
			// TODO: complete gracefull shutdown
			return
		}
	}
}
