-- Create the orders table with a UNIQUE constraint on track_number
CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(50),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(50),
    sm_id INT,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(50)
);

CREATE TABLE deliveries (
    delivery_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

-- Create the payments table
CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(50) NOT NULL,
    request_id VARCHAR(50),
    currency VARCHAR(50),
    provider VARCHAR(50),
    amount INT,
    payment_dt INT,
    bank VARCHAR(50),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE items (
    item_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INT,
    track_number VARCHAR(50),
    price INT,
    rid VARCHAR(50),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(50),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT
);