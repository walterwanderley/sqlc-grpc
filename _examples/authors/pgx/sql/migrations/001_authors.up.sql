CREATE TABLE IF NOT EXISTS authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  created_at TIMESTAMP
);