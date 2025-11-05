CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    shop_id INT NOT NULL,
    courier_id INT NOT NULL,
    payment_status payment_status DEFAULT 'pending' NOT NULL,
    status order_status DEFAULT 'pending' NOT NULL,
    total INT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (shop_id) REFERENCES shops(id),
    FOREIGN KEY (courier_id) REFERENCES couriers(id)
);