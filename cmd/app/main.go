package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wb_l0/internal/nats"
	"wb_l0/internal/repository/inmemory"
	pg "wb_l0/internal/repository/postgres"
	"wb_l0/internal/server"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	doneCh := make(chan bool)

	conn, err := nats.Conn()
	if err != nil {
		log.Fatal("failed to connect nats", err)
	}

	// move dsn to config
	repo := pg.New(ctx, "postgres://test:test@localhost:1234/orders?sslmode=disable")

	// todo
	go nats.NewPub(conn, doneCh)

	subscriber, err := nats.NewSub(conn, repo)
	if err != nil {
		// refactor
		panic(err)
	}
	_ = subscriber

	cache := inmemory.New(ctx, repo)

	server := server.New(repo, cache)
	server.Run()

	go func() {
		<-sig
		cancel()
		// stop server
		err := server.Shutdown(context.TODO())
		if err != nil {
			log.Fatalf("failed to shutdown server: %s", err.Error())
		}
		// stop chache
		// stop repo
		// stop nats
		// stop subscriber
		// stop publisher
	}()
}
