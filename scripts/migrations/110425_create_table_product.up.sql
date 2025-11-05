CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name TEXT NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL,
    shop_id INT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    FOREIGN KEY (shop_id) REFERENCES shops(id)
);