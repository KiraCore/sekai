# INTERX


## Config gRPC-gateway server

Default is http server with transport credentials.

### SERVE_HTTP

```json
"SERVE_HTTP": {
    "description": "Whether to use HTTP to serve, instead of HTTPS",
    "value": "true",
    "required": false
},
```

### WITH_TRANSPORT_CREDENTIALS

```json
"WITH_TRANSPORT_CREDENTIALS": {
    "description": "Whether to use transport credentials to connect gRPC server",
    "value": "true",
    "required": false
}
```

## Generate go, gRPC-Gateway, OpenAPI output.

- generate go, gRPC-gateway in ./proto-gen.
- generate OpenAPI output in ./third_party/OpenAPI

```bash
make generate
```

## Start gRPC-gateway server

Default PORT = 11000

API documentation is available http://0.0.0.0:PORT or https://0.0.0.0:PORT

### Run server with default configuration. (HTTP server, use transport credentials)

```bash
go run main.go
```

### Run server with environment configuration.

```bash
PORT=11100 go run main.go
```

```bash
WITH_TRANSPORT_CREDENTIALS=false go run main.go
```

```bash
SERVE_HTTP=false go run main.go
```

## Start example gRPC server

port = 10000

```json
"WITH_TRANSPORT_CREDENTIALS": {
    "description": "Whether to use transport credentials for gRPC server",
    "value": "true",
    "required": false
}
```

### Run server with default configuration. (with transport credentials)

```bash
go run ./example-grpc-server/main.go
```

### Run server with environment configuration.

```bash
WITH_TRANSPORT_CREDENTIALS=false go run ./example-grpc-server/main.go
```
