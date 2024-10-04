-- name: CreateUser :one
INSERT INTO users (id, location) VALUES ($1, $2) RETURNING id;

-- name: CreateUserReturnPartial :one
INSERT INTO users (id, location) VALUES ($1, $2) RETURNING id, name;

-- name: CreateUserReturnAll :one
INSERT INTO users (id, location) VALUES ($1, $2) RETURNING *;

-----------

-- name: CreateProduct :one
INSERT INTO products (id, category) VALUES ($1, $2) RETURNING id;

-- name: CreateProductReturnPartial :one
INSERT INTO
    products (id, category)
VALUES ($1, $2) RETURNING id,
    name;

-- name: CreateProductReturnAll :one
INSERT INTO products (id, category) VALUES ($1, $2) RETURNING *;

-- name: GetProductsByIds :many
SELECT * FROM products WHERE id = ANY($1::uuid[]);

-----------

-- name: CreateLocationTransactions :exec
INSERT INTO location_transactions
SELECT
    UNNEST($1::UUID[]) as location_id,
    UNNEST($2::UUID[]) as transaction_id;