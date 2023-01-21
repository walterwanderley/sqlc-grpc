# SQLite with LiteFS example

### Start a cluster with 3 instances:

- Start instance 1:
```sh
go run . -db authors.db -port 8080 -litefs-redirect http://localhost:8080 -litefs-hostname node1 \
-litefs-config-dir /tmp/raft1 -litefs-mount-dir /tmp/data1 -litefs-port 20202 \
-litefs-advertise-url http://localhost:20202 -litefs-members="node2=localhost:9081, node3=localhost:9082" \
-litefs-raft-port 9080 -litefs-raft-addr localhost:9080 -dev
```

- Start instance 2:
```sh
go run . -db authors.db -port 8081 -litefs-redirect http://localhost:8081 -litefs-hostname node2 \
-litefs-config-dir /tmp/raft2 -litefs-mount-dir /tmp/data2 -litefs-port 20203 \
-litefs-advertise-url http://localhost:20203 -litefs-raft-port 9081 -litefs-raft-addr localhost:9081 \
-litefs-bootstrap-cluster=false -dev
```

- Start instance 3:
```sh
go run . -db authors.db -port 8082 -litefs-redirect http://localhost:8082 -litefs-hostname node3 \
-litefs-config-dir /tmp/raft3 -litefs-mount-dir /tmp/data3 -litefs-port 20204 \
-litefs-advertise-url http://localhost:20204 -litefs-raft-port 9082 -litefs-raft-addr localhost:9082 \
-litefs-bootstrap-cluster=false -dev
```

- List nodes:
```sh
curl http://localhost:8080/nodes/
```

- Get the leader:
```sh
curl http://localhost:8080/nodes/leader
```

- Add new node to cluster:
```sh
curl -X POST -d '{"id": "{-litefs-hostname}", "addr": "{-litefs-raft-addr}", readOnly: false}' http://localhost:8080/nodes/
```

- Remove a node from cluster:
```sh
curl -X DELETE http://localhost:8080/nodes/{id}
```

### Explore the API

- [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- [http://localhost:8081/swagger/](http://localhost:8080/swagger/)
- [http://localhost:8082/swagger/](http://localhost:8080/swagger/)

**POST/PUT/DELETE** requests are forwarded to the Leader.

### Litestream

Check a [Litestream](https://litestream.io) example in docker-compose.yml

### Steps to generate this code

1. Create a directory to store SQL scripts.

```sh
mkdir -p sql/migrations
```

2. Create migrations scripts using [go-migrate](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) rules.

```sh
echo "CREATE TABLE IF NOT EXISTS authors (
    id   integer    PRIMARY KEY AUTOINCREMENT,
    name text   NOT NULL,
    bio  text
);
" > sql/migrations/001_authors.up.sql
```

```sh
echo "DROP TABLE IF EXISTS authors;" > sql/migrations/001_authors.down.sql
```

3. Create SQL queries

```sh
echo "/* name: GetAuthor :one */
SELECT * FROM authors
WHERE id = ? LIMIT 1;

/* name: ListAuthors :many */
SELECT * FROM authors
ORDER BY name;

/* name: CreateAuthor :execresult */
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ? 
);

/* name: DeleteAuthor :exec */
DELETE FROM authors
WHERE id = ?;
" > sql/queries.sql
```
4. Create sqlc.yaml configuration file

```sh
echo 'version: "1"
packages:
  - path: "internal/authors"
    queries: "./sql/queries.sql"
    schema: "./sql/migrations"
    engine: "sqlite"
' > sqlc.yaml
```

5. Execute sqlc

```sh
sqlc generate
```

6. Execute sqlc-grpc

```sh
sqlc-grpc -m example -migration-path sql/migrations -litefs
```
