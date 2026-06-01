CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at timestamp(0) with time zone DEFAULT now(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1

);