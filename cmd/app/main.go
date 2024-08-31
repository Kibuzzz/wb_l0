package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"wb_l0/internal/config"
	"wb_l0/internal/nats"
	"wb_l0/internal/repository/inmemory"
	pg "wb_l0/internal/repository/postgres"
	"wb_l0/internal/server"
)

func main() {

	conf := config.Get()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := nats.Conn(conf)
	if err != nil {
		log.Fatalf("failed to connect nats: %s", err.Error())
	}

	repo := pg.New(ctx, conf.Database.DSN)

	go nats.NewPub(conn, ctx)

	err = nats.NewSub(conn, conf.NATS.ClusterID, repo)
	if err != nil {
		log.Fatalf("failed to initialize NATS subscriber: %v", err)
	}

	cache := inmemory.New(ctx, repo)

	addr := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	server := server.New(repo, cache, addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Received shutdown signal, starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown server gracefully: %v", err)
	} else {
		log.Println("server shut down gracefully")
	}

	if err := conn.Close(); err != nil {
		log.Printf("failed to close NATS connection: %v", err)
	} else {
		log.Println("NATS connection closed gracefully")
	}

	repo.Close()
	log.Println("repository closed gracefully")
	log.Println("shut down completed")
}
