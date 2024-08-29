package main

import (
	"context"
	"wb_l0/internal/repository/inmemory"
	pg "wb_l0/internal/repository/postgres"
	"wb_l0/internal/server"
)

func main() {

	ctx := context.TODO()

	// вынести в конфиг dsn
	repo := pg.New(ctx, "postgres://test:test@localhost:1234/orders?sslmode=disable")

	cache := inmemory.New(ctx, repo)

	server := server.New(repo, cache)
	server.Run()
}
