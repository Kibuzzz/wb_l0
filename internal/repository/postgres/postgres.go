package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"wb_l0/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(ctx context.Context, connString string) *Postgres {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &Postgres{pool, logger}
}

func (p *Postgres) Close() {
	p.db.Close()
}

func (p *Postgres) AddOrder(ctx context.Context, order models.Order) error {
	sql := `INSERT INTO orders (
		order_uid, 
		track_number, 
		entry, 
		locale, 
		internal_signature, 
		customer_id, 
		delivery_service, 
		shardkey, 
		sm_id, 
		date_created, 
		oof_shard
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	);`

	_, err := p.db.Exec(ctx, sql,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated.UTC(),
		order.OofShard,
	)

	if err != nil {
		fmt.Printf("Error adding order: %v\n", err)
		return err
	}

	for _, item := range order.Items {
		err := p.createItem(ctx, order.OrderUID, item)
		if err != nil {
			fmt.Printf("Error creating item: %v\n", item)
			return err
		}
	}

	err = p.createPayment(ctx, order.Payment)
	if err != nil {
		fmt.Printf("Error creating payment: %v\n", order.Payment)
		return err
	}

	err = p.createDelivery(ctx, order.Delivery)
	if err != nil {
		fmt.Printf("Error creating delivery: %v\n", order.Delivery)
		return err
	}

	return nil
}

func (p *Postgres) GetOrderByID(ctx context.Context, orderUID string) (models.Order, error) {
	sql := `
	select * from orders where order_uid = $1;
	`
	row := p.db.QueryRow(ctx, sql, orderUID)
	order := models.Order{}
	err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		fmt.Printf("Error fetching order by ID: %v\n", err)
		return models.Order{}, err
	}

	// Fetch and assign Delivery
	delivery, err := p.getDelivery(ctx, order.OrderUID)
	if err != nil {
		fmt.Printf("Error fetching delivery: %v\n", err)
		return models.Order{}, err
	}
	order.Delivery = delivery

	// Fetch and assign Payment
	payment, err := p.getPayment(ctx, order.OrderUID)
	if err != nil {
		fmt.Printf("Error fetching payment: %v\n", err)
		return models.Order{}, err
	}
	order.Payment = payment

	// Fetch and assign Items
	items, err := p.getItems(ctx, order.OrderUID) // Use orderUID instead of trackNumber
	if err != nil {
		fmt.Printf("Error fetching items: %v\n", err)
		return models.Order{}, err
	}
	order.Items = items
	// TODO: refactor here and in the tests
	order.DateCreated = order.DateCreated.UTC()
	return order, nil
}

func (p *Postgres) AllOrders(ctx context.Context) ([]models.Order, error) {
	sql := `SELECT * FROM orders;`

	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error getting rows: %v", err)
	}
	defer rows.Close()

	orders := []models.Order{}

	for rows.Next() {
		order := models.Order{}

		// Scan the order data
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			fmt.Printf("Error scanning order: %v\n", err)
			return nil, err
		}

		// Fetch and assign Delivery
		delivery, err := p.getDelivery(ctx, order.OrderUID)
		if err != nil {
			fmt.Printf("Error fetching delivery for order %v: %v\n", order.OrderUID, err)
			return nil, err
		}
		order.Delivery = delivery

		// Fetch and assign Payment
		payment, err := p.getPayment(ctx, order.OrderUID)
		if err != nil {
			fmt.Printf("Error fetching payment for order %v: %v\n", order.OrderUID, err)
			return nil, err
		}
		order.Payment = payment

		// Fetch and assign Items
		items, err := p.getItems(ctx, order.OrderUID)
		if err != nil {
			fmt.Printf("Error fetching items for order %v: %v\n", order.OrderUID, err)
			return nil, err
		}
		order.Items = items

		// Convert date to UTC
		order.DateCreated = order.DateCreated.UTC()

		// Append the order to the orders slice
		orders = append(orders, order)
	}

	// Check if there were any errors during iteration
	if err = rows.Err(); err != nil {
		fmt.Printf("Error during rows iteration: %v\n", err)
		return nil, err
	}

	return orders, nil
}

func (p *Postgres) createDelivery(ctx context.Context, delivery models.Delivery) error {
	sql := `INSERT INTO deliveries (
		order_uid, 
		name, 
		phone, 
		zip, 
		city, 
		address, 
		region, 
		email
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	);`

	_, err := p.db.Exec(ctx, sql,
		delivery.OrderUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)

	return err
}

func (p *Postgres) getDelivery(ctx context.Context, orderUID string) (models.Delivery, error) {
	sql := `SELECT 
		order_uid, 
		name, 
		phone, 
		zip, 
		city, 
		address, 
		region, 
		email 
	FROM deliveries 
	WHERE order_uid = $1;`

	row := p.db.QueryRow(ctx, sql, orderUID)

	var delivery models.Delivery
	err := row.Scan(
		&delivery.OrderUID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}

func (p *Postgres) createPayment(ctx context.Context, payment models.Payment) error {
	sql := `INSERT INTO payments (
        order_uid, 
        transaction, 
        request_id, 
        currency, 
        provider, 
        amount, 
        payment_dt, 
        bank, 
        delivery_cost, 
        goods_total, 
        custom_fee
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
    );`

	_, err := p.db.Exec(ctx, sql,
		payment.OrderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)

	return err
}

func (p *Postgres) getPayment(ctx context.Context, orderUID string) (models.Payment, error) {
	sql := `SELECT 
        order_uid, 
        transaction, 
        request_id, 
        currency, 
        provider, 
        amount, 
        payment_dt, 
        bank, 
        delivery_cost, 
        goods_total, 
        custom_fee 
    FROM payments 
    WHERE order_uid = $1;`

	row := p.db.QueryRow(ctx, sql, orderUID)

	var payment models.Payment
	err := row.Scan(
		&payment.OrderUID,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}

func (p *Postgres) createItem(ctx context.Context, orderUID string, item models.Item) error {
	sql := `INSERT INTO items (
        order_uid, 
        chrt_id, 
        track_number, 
        price, 
        rid, 
        name, 
        sale, 
        size, 
        total_price, 
        nm_id, 
        brand, 
        status
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
    );`

	_, err := p.db.Exec(ctx, sql,
		orderUID,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
	)

	return err
}

func (p *Postgres) getItems(ctx context.Context, orderUID string) ([]models.Item, error) {
	sql := `SELECT 
		order_uid,
        chrt_id, 
        track_number, 
        price, 
        rid, 
        name, 
        sale, 
        size, 
        total_price, 
        nm_id, 
        brand, 
        status 
    FROM items 
    WHERE order_uid = $1;`

	rows, err := p.db.Query(ctx, sql, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.OrderUID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
