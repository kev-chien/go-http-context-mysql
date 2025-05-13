# go-http-context-mysql

## Requirements

- MySQL
- Go 1.24.2

## Development Usage

To set up the MySQL DB:

```sh
mysql -t < setupdb.sql
```

To run the server:

```sh
cd server
MYSQL_USER=CHANGEME MYSQL_PASSWORD=CHANGEME go run server.go
```

Some curl commands:

```sh
curl localhost:8090/longResponse
curl localhost:8090/longResponseChecksContext
curl localhost:8090/longResponseDBNoTx
curl localhost:8090/longResponseDB
```

You can test Ctrl-C while running the curl commands to see what happens.

`/longResponse` should be the only endpoint where the server keeps processing even after the client has canceled the request.

All the other endpoints are context-aware, and will return early when they notice that the client has canceled the request.
