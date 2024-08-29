package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GenerateRandomOrder generates a random Order struct
func GenerateRandomOrder() Order {
	orderUID := randomString(10)
	trackNumber := randomString(14)
	entry := randomString(5)
	locale := "en" // Assuming 'en' for simplicity; it could also be randomly selected.
	customerID := randomString(6)
	deliveryService := randomString(8)
	shardkey := strconv.Itoa(rand.Intn(100))
	smID := rand.Intn(1000)
	oofShard := randomString(5)

	// Generate Delivery
	delivery := Delivery{
		//DeliveryID: rand.Intn(1000),
		OrderUID: orderUID,
		Name:     randomString(8),
		Phone:    fmt.Sprintf("%d", rand.Intn(1000000000)),
		Zip:      fmt.Sprintf("%05d", rand.Intn(100000)),
		City:     randomString(10),
		Address:  randomString(15),
		Region:   randomString(10),
		Email:    randomString(5) + "@example.com",
	}

	// Generate Payment
	payment := Payment{
		//PaymentID:    rand.Intn(1000),
		OrderUID:     orderUID,
		Transaction:  randomString(12),
		RequestID:    randomString(10),
		Currency:     "USD", // Assuming USD for simplicity
		Provider:     randomString(8),
		Amount:       rand.Intn(1000),
		PaymentDt:    int(time.Now().Unix()),
		Bank:         randomString(10),
		DeliveryCost: rand.Intn(50),
		GoodsTotal:   rand.Intn(1000),
		CustomFee:    rand.Intn(100),
	}

	// Generate Items
	items := make([]Item, rand.Intn(5)+1) // Random number of items between 1 and 5
	for i := range items {
		items[i] = Item{
			//ItemID:      rand.Intn(1000),
			OrderUID:    orderUID,
			ChrtID:      rand.Intn(1000),
			TrackNumber: trackNumber,
			Price:       rand.Intn(100),
			Rid:         randomString(10),
			Name:        randomString(12),
			Sale:        rand.Intn(20),
			Size:        randomString(2),
			TotalPrice:  rand.Intn(200),
			NmID:        rand.Intn(1000),
			Brand:       randomString(10),
			Status:      rand.Intn(5),
		}
	}

	// Create and return the Order
	return Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             entry,
		Locale:            locale,
		InternalSignature: randomString(5),
		CustomerID:        customerID,
		DeliveryService:   deliveryService,
		Shardkey:          shardkey,
		SmID:              smID,
		DateCreated:       time.Now(),
		OofShard:          oofShard,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
	}
}

// randomString generates a random string of given length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
