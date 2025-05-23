# go-http-context-mysql

The purpose of this repo is to test what happens with a go http server when the client cancels a request early, to see if the server will also return early, and if the MySQL driver will appropriately cancel any DB operations.

## Requirements

- MySQL
- Go 1.24.2

## Development Usage

### Setting up the MySQL DB

```sh
mysql -t < setupdb.sql
```

### Running the server

```sh
cd server
MYSQL_USER=CHANGEME MYSQL_PASSWORD=CHANGEME go run server.go
```

### Manual curl commands

```sh
curl localhost:8090/longResponse
curl localhost:8090/longResponseChecksContext
curl localhost:8090/longResponseDB
curl localhost:8090/longResponseDBNoTx
```

You can test Ctrl-C while running the curl commands to see what happens.

`/longResponse` should be the only endpoint where the server keeps processing even after the client has canceled the request.

All the other endpoints are context-aware, and will return early when they notice that the client has canceled the request.

### Using the client

```sh
cd client
go run client.go longResponse
go run client.go longResponseChecksContext
go run client.go longResponseDB
go run client.go longResponseDBNoTx
```

For each of these, the client will call the specified endpoint with a timeout of 1 second, so it cancels the request before the server has completed processing it.

The results are the same as the above: `/longResponse` is the only endpoint where the server keeps processing even after the client has canceled the request.

All the other endpoints are context-aware, and will return early when they notice that the client has canceled the request.
