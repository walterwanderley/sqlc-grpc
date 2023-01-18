## sqlc-grpc

Create a **gRPC** (and **HTTP/JSON**) **Server** from the generated code by the awesome [sqlc](https://sqlc.dev/) project.

### Requirements

- Go 1.16 or superior
- [sqlc](https://sqlc.dev/)
- [buf](https://buf.build/)

```sh
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
go install github.com/bufbuild/buf/cmd/buf@latest
```

### Installation

```sh
go install github.com/walterwanderley/sqlc-grpc@latest
```

### Example

1. Create a queries.sql file:

```sql
--queries.sql

CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  created_at TIMESTAMP
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, created_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

```

2. Create a sqlc.yaml file

```yaml
version: "1"
packages:
  - path: "internal/author"
    queries: "./queries.sql"
    schema: "./queries.sql"
    engine: "postgresql"

```

3. Execute sqlc

```sh
sqlc generate
```

4. Execute sqlc-grpc

```sh
sqlc-grpc -m "my/module/path"
```

5. Run the generated server

```sh
go run . -db [Database Connection URL] -dev
```

6. Enjoy!

- Swagger UI [http://localhost:5000/swagger](http://localhost:5000/swagger)

### SQLite with LiteFS

Use the **-litefs** command line parameter to replicate SQLite with [LiteFS](https://github.com/superfly/litefs) as a library.
Example: [https://github.com/walterwanderley/sqlc-grpc/tree/main/_examples/authors/sqlite](https://github.com/walterwanderley/sqlc-grpc/tree/main/_examples/authors/sqlite)

### Editing the generated code

- It's safe to edit any generated code that doesn't have the `DO NOT EDIT` indication at the very first line.

- After modify a SQL file, execute these commands below:

```sh
sqlc generate
go generate
```

- After modify a *.proto file, execute `buf generate`.

### Similar Projects

- [xo-grpc](https://github.com/walterwanderley/xo-grpc)
