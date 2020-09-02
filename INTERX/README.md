# INTERX

## Generate go, gRPC-Gateway, OpenAPI output.

- generate go, gRPC-gateway in ./proto-gen.
- generate OpenAPI output in ./third_party/OpenAPI

```bash
make generate
```

## Start gRPC-gateway server

### Start with default settings
```bash
make start
```

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
    "value": "false",
    "required": false
}
```

### GRPC

```json
"GRPC": {
    "description": "GRPC endpoint",
    "value": "dns:///0.0.0.0:9090",
    "required": false
}
```

### RPC

```json
"RPC": {
    "description": "RPC endpoint",
    "value": "http://0.0.0.0:26657",
    "required": false
}
```

```bash
WITH_TRANSPORT_CREDENTIALS=true go run main.go
```

```bash
SERVE_HTTP=false go run main.go
```

```bash
GRPC=dns:///0.0.0.0:9090 RPC=http://0.0.0.0:26657 make start
```
