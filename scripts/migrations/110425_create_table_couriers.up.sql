CREATE TABLE couriers (
    id SERIAL PRIMARY KEY,
    brand TEXT NOT NULL,
    employer_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    shipping_cost INT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);