CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    province TEXT NOT NULL,
    district TEXT NOT NULL,
    subdistrict TEXT NOT NULL,
    zip_code TEXT NOT NULL,
    detail_address TEXT NOT NULL,
    phone TEXT NOT NULL,
    password TEXT NOT NULL,
    type_user user_type DEFAULT 'general' NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TYPE IF NOT EXISTS user_type AS ENUM ('general','shop','admin');