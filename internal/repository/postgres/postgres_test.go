package postgres

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"
	"wb_l0/internal/models"

	"github.com/dvln/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCustomerRepository(t *testing.T) {
	ctx := context.Background()

	dbName := "orders"
	dbUser := "test"
	dbPassword := "test"

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("../../../testdata", "init.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connString, err := postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	assert.NoError(t, err, "failed to connect", connString)
	repo := New(ctx, connString)

	// test one order creation
	order := models.GenerateRandomOrder()
	err = repo.AddOrder(ctx, order)
	assert.NoError(t, err, "failed to create item")

	// test getting one order
	recievedOrder, err := repo.GetOrderByID(ctx, order.OrderUID)
	assert.NoError(t, err, "failed to get order")
	assert.Equal(t, order, recievedOrder, "orders not matches")

	// test All orders
	second_order := models.GenerateRandomOrder()
	err = repo.AddOrder(ctx, second_order)
	assert.NoError(t, err, "failed to create item")
	allOrders, err := repo.AllOrders(ctx)
	assert.NoError(t, err, "failed to get order")
	assert.Len(t, allOrders, 2, "order number is wrong")
}
