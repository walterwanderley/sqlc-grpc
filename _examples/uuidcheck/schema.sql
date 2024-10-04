CREATE TABLE IF NOT EXISTS "locations" ("id" UUID PRIMARY KEY);

CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY,
    "location" UUID REFERENCES "locations" ("id"),
    "name" VARCHAR
);

CREATE TABLE IF NOT EXISTS "category" ("id" SERIAL PRIMARY KEY);

CREATE TABLE IF NOT EXISTS "products" (
    "id" SERIAL PRIMARY KEY,
    "category" INT REFERENCES "category" ("id"),
    "name" VARCHAR
);

CREATE TABLE location_transactions (
    location_id UUID NOT NULL,
    transaction_id UUID NOT NULL
);