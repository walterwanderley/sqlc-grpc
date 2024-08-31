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