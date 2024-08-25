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

func getTestData() models.Order {
	var testDelivery = models.Delivery{
		OrderUID: "test_order_uid",
		Name:     "John Doe",
		Phone:    "+1234567890",
		Zip:      "123456",
		City:     "Test City",
		Address:  "123 Test Street",
		Region:   "Test Region",
		Email:    "john.doe@example.com",
	}

	var testPayment = models.Payment{
		OrderUID:     "test_order_uid",
		Transaction:  "test_transaction",
		RequestID:    "test_request_id",
		Currency:     "USD",
		Provider:     "Test Provider",
		Amount:       10000,
		PaymentDt:    1622512320,
		Bank:         "Test Bank",
		DeliveryCost: 500,
		GoodsTotal:   9500,
		CustomFee:    200,
	}

	var testItem1 = models.Item{
		OrderUID:    "test_order_uid",
		ChrtID:      123456,
		TrackNumber: "TESTTRACK123456",
		Price:       1000,
		Rid:         "test_rid_1",
		Name:        "Test Item 1",
		Sale:        10,
		Size:        "L",
		TotalPrice:  900,
		NmID:        111222,
		Brand:       "Test Brand",
		Status:      1,
	}

	var testItem2 = models.Item{
		OrderUID:    "test_order_uid",
		ChrtID:      789012,
		TrackNumber: "TESTTRACK789012",
		Price:       2000,
		Rid:         "test_rid_2",
		Name:        "Test Item 2",
		Sale:        5,
		Size:        "M",
		TotalPrice:  1900,
		NmID:        333444,
		Brand:       "Another Test Brand",
		Status:      2,
	}

	var testOrder = models.Order{
		OrderUID:          "test_order_uid",
		TrackNumber:       "TESTTRACK123",
		Entry:             "WEB",
		Delivery:          testDelivery,
		Payment:           testPayment,
		Items:             []models.Item{testItem1, testItem2},
		Locale:            "en-US",
		InternalSignature: "signature123",
		CustomerID:        "customer123",
		DeliveryService:   "Test Delivery Service",
		Shardkey:          "test_shardkey",
		SmID:              1,
		DateCreated:       time.Now().UTC(),
		OofShard:          "test_oof_shard",
	}
	return testOrder
}

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
	order := getTestData()
	err = repo.AddOrder(ctx, order)
	assert.NoError(t, err, "failed to create item")
	recievedOrder, err := repo.GetOrderByID(ctx, order.OrderUID)
	assert.NoError(t, err, "failed to get order")
	assert.Equal(t, order, recievedOrder, "orders not matches")

}
