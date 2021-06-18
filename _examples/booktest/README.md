# About

Booktest example taken from [sqlc][sqlc] Git repository [examples][sqlc-git].

[sqlc]: https://sqlc.dev
[sqlc-git]: https://github.com/kyleconroy/sqlc/tree/main/examples/booktest

## Running

```sh
./gen.sh
docker-compose up
```

### Exploring

- gRPC UI [http://localhost:8080/grpcui](http://localhost:8080/grpcui)
- Swagger UI [http://localhost:8080/swagger](http://localhost:8080/swagger)
- Grafana [http://localhost:3000](http://localhost:3000/d/7_VGtoLma/go-grpc1?orgId=1&refresh=10s&from=now-5m&to=now) *user/pass*: **admin/admin**
