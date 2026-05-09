CREATE TABLE IF NOT EXISTS authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  age NUMERIC,
  created_at TIMESTAMP
);