# About

Booktest example taken from [sqlc][sqlc] Git repository [examples][sqlc-git].

[sqlc]: https://sqlc.dev
[sqlc-git]: https://github.com/sqlc-dev/sqlc/tree/main/examples/booktest

## Running

```sh
./gen.sh
docker compose up
```

### Exploring

- Swagger UI [http://localhost:8080/swagger](http://localhost:8080/swagger)
- Grafana [http://localhost:3000](http://localhost:3000/d/7_VGtoLma/go-grpc1?orgId=1&refresh=10s&from=now-5m&to=now) *user/pass*: **admin/admin**
- Jaeger [http://localhost:16686](http://localhost:16686)

