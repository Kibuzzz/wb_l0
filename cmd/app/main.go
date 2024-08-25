package main

import (
	"context"
	pg "wb_l0/internal/repository/postgres"
	"wb_l0/internal/server"
)

func main() {
	repo := pg.New(context.TODO(), "postgres://test:test@localhost:1234/orders?sslmode=disable")

	server := server.New(repo)
	server.Run()
}
