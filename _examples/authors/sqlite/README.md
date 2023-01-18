# SQLite with LiteFS example

### Start a cluster with 3 instances:

```sh
go run . -db authors.db -port 8080 -litefs-redirect http://localhost:8080 -litefs-hostname node1 -litefs-config-dir /tmp/raft1 -litefs-mount-dir /tmp/data1 -litefs-port 20202 -litefs-advertise-url http://localhost:20202 -litefs-members="node2=localhost:9081, node3=localhost:9082" -litefs-raft-port 9080 -litefs-raft-addr localhost:9080 -dev

go run . -db authors.db -port 8081 -litefs-redirect http://localhost:8080 -litefs-hostname node2 -litefs-config-dir /tmp/raft2 -litefs-mount-dir /tmp/data2 -litefs-port 20203 -litefs-advertise-url http://localhost:20203 -litefs-raft-port 9081 -litefs-raft-addr localhost:9081 -litefs-bootstrap-cluster=false -dev

go run . -db authors.db -port 8082 -litefs-redirect http://localhost:8080 -litefs-hostname node1 -litefs-config-dir /tmp/raft3 -litefs-mount-dir /tmp/data3 -litefs-port 20204 -litefs-advertise-url http://localhost:20204 -litefs-raft-port 9082 -litefs-raft-addr localhost:9082 -litefs-bootstrap-cluster=false -dev

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


## Litestream

Check a [Litestream](https://litestream.io) example in docker-compose.yml