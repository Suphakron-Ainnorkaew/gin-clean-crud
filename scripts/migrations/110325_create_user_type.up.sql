-- Create Postgres enum type for user.type_user
CREATE TYPE IF NOT EXISTS user_type AS ENUM ('general','shop','admin');
