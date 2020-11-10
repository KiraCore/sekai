# INTERX

## Configurations.

#### config.json

- interx configurations

    "mnemonic": "...",
    "status_sync": 5,
    
- caching configurations

    "cache_dir": "cache",
    "max_cache_size": "2 GB",
    "caching_duration": 10,

- faucet configurations

    "faucet": {
        "mnemonic": "...",
        "faucet_amounts": {
            "stake": 100000,
            "...": ...
        },
        "faucet_minimum_amounts": {
            "stake": 100,
            "...": ...
        },
        "time_limit": 20
    },

- RPC configurations (include caching enable/disable for each endpoint)

    "rpc": {
        "API": {
            "GET": {
                "/api/cosmos/status": {
                    "rate_limit": 0.1,
                    "auth_rate_limit": 1,
                    "caching_disable": true
                },
                "...": {...}
            },
            "POST": {
                "...": {...}
            }
        }
    }

#### functions/*.json

functions metadata

    - we have functions metadata in json files in `functions` folder.
        (there can be multiple json files. we can collect all json files to `functions` folder)
    - each json file can have multiple transaction types.
    - each transaction type has `description` and `parameters` field.
    - each parameter has `type` and `description` field.
    
    {
        "{tx_type1}": {
            "description": "Description field for each tx type",
            "parameters": {     // List parameters here.
                "{parameter1}": {
                    "type": "Parameter type field, e.g. bool, string, int, ...",
                    "description": "Description field for each parameter"
                },
                "{parameter2}": {
                    "type": "...",
                    "description": "..."
                },
                ...
            }
        },
        "{tx_type2}": {
            "description": "Description field for each tx type",
            "parameters": {     // List parameters here.
                "{parameter1}": {
                    "type": "...",
                    "description": "..."
                },
                "{parameter2}": {
                    "type": "...",
                    "description": "..."
                },
                ...
            }
        }
    }

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
