# SQLite with LiteFS example

## Steps to generate this code

0. Install the required tools.

```sh
go install github.com/walterwanderley/sqlc-grpc@latest
```
```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

1. Create a directory to store SQL scripts.

```sh
mkdir -p sql/migrations
```

2. Create migrations scripts using [goose](https://github.com/pressly/goose?tab=readme-ov-file#migrations) rules.

```sh
echo "-- +goose Up
CREATE TABLE IF NOT EXISTS authors (
    id   integer    PRIMARY KEY AUTOINCREMENT,
    name text   NOT NULL,
    bio  text
);

-- +goose Down
DROP TABLE IF EXISTS authors;
" > sql/migrations/001_authors.sql
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

### Running a cluster with 3 instances in STATIC leasing mode:

>**Note:** In this mode you configure a single node to be the primary/leader.
The downside of this approach is that you will lose write availability if that node goes down.

- Start instance 1:
```sh
go run . -db authors.db -port 8080 -litefs-redirect http://localhost:8080 -litefs-hostname node1 \
-litefs-config-dir /tmp/cfg1 -litefs-mount-dir /tmp/data1 -litefs-port 20202 \
-litefs-advertise-url http://localhost:20202 -dev
```

- Start instance 2:
```sh
go run . -db authors.db -port 8081 -litefs-redirect http://localhost:8080 -litefs-hostname node1 \
-litefs-config-dir /tmp/cfg2 -litefs-mount-dir /tmp/data2 -litefs-port 20203 \
-litefs-advertise-url http://localhost:20202 \
-litefs-candidate=false -dev
```

- Start instance 3:
```sh
go run . -db authors.db -port 8082 -litefs-redirect http://localhost:8080 -litefs-hostname node1 \
-litefs-config-dir /tmp/cfg3 -litefs-mount-dir /tmp/data3 -litefs-port 20204 \
-litefs-advertise-url http://localhost:20202 \
-litefs-candidate=false -dev
```

### Running a cluster with 3 instances in RAFT-based leasing mode:

>**Note:** In this mode the primary/leader is elected using the RAFT protocol. 
If the leader node goes down another node will be elected.
The downside of this approach is that "[adding distributed consensus to your application nodes can be problematic when they're under high load as they can loose leadership easily](https://github.com/superfly/litefs/issues/259#issuecomment-1398766012)".

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
go run . -db authors.db -port 8082 -litefs-redirect http://localhost:8082 -litefs-hostname node1 \
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

## Explore the API

- [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- [http://localhost:8081/swagger/](http://localhost:8081/swagger/)
- [http://localhost:8082/swagger/](http://localhost:8082/swagger/)

**POST/PUT/DELETE** requests are forwarded to the Leader.

## Litestream

Check out a [Litestream](https://litestream.io) example in [docker-compose.yml](https://github.com/walterwanderley/sqlc-grpc/blob/main/_examples/authors/sqlite/docker-compose.yml).
