# gofr_layout

A production-oriented Go microservice project layout built on the [gofr](https://github.com/neo532/gofr) framework. It provides a structured, extensible template for building HTTP/gRPC services with clean separation of concerns.

## Features

- **Clean Architecture** — service / domain / repo / data layers with dependency inversion
- **Proto-First API Definition** — services defined in protobuf, code-generated via `protoc-gen-go-svc`
- **Pluggable Middleware** — built-in validation middleware, extensible via the gofr middleware chain
- **Content-Type Codec Negotiation** — JSON by default, plug in XML or any format via codec registry
- **Hot-Reload Configuration** — YAML configs watched at runtime, changes merged via `atomic.Value`
- **Dependency Injection** — [wire](https://github.com/google/wire) for compile-time DI
- **Request Validation** — `protoc-gen-validate` rules enforced at the middleware layer
- **Path Parameter Injection** — typed params (int64, string, bool, float, etc.) binding via reflection
- **Multi-Entrypoint** — `api` (HTTP server), `consumer` (message queue), `script` (one-off tasks)
- **Structured Logging** — slog-based logger with lumberjack rotation
- **Database & Cache** — GORM (MySQL) + Redis support via [gokit](https://github.com/neo532/gokit)

## Project Structure

```
├── cmd/
│   ├── api/          # HTTP server entry point (wire injector)
│   ├── consumer/     # Consumer entry point
│   ├── script/       # Script/task entry point
│   └── cmd.go        # Shared bootstrap (config, logger, context)
│
├── internal/
│   ├── config/       # Config structs (code-generated from YAML)
│   ├── server/       # Server wiring (HTTP server setup, middleware registration)
│   ├── service/      # Application service layer
│   │   ├── api/      #   API service implementations
│   │   ├── consumer/ #   Consumer handlers
│   │   └── script/   #   Script handlers
│   ├── domain/       # Business domain logic
│   ├── repo/         # Repository interfaces
│   ├── data/         # Repository implementations (DB, cache, etc.)
│   │   ├── base/     #   DB/Redis client initialization
│   │   └── model/    #   GORM models
│
├── proto/            # Protobuf definitions & generated code
│   ├── api/          #   Service, struct, enum protos
│   └── third_party/  #   googleapis, validate dependencies
│
├── configs/          # Runtime config YAML files
└── internal/config/dev/  # Dev config templates
```

## Architecture

### Layer Flow

```
HTTP Request
  → gofr HTTP Server
    → Middleware Chain (Validator, Logging, Recovery, Tracing, ...)
      → Service Layer (orchestration, no business logic)
        → Domain Layer (business rules)
          → Repository Interface
            → Data Layer (GORM / Redis)
```

### Dependency Injection

Wire compiles constructors into a single `wire_gen.go` per entrypoint. Each layer provides its own `ProviderSet`:

```go
// internal/service/api/wireProviderSet.go
var ProviderSet = wire.NewSet(
    NewUserApiService,
    // ...
)
```

## Getting Started

### Prerequisites

- Go 1.26+
- Protobuf toolchain (protoc, protoc-gen-go, protoc-gen-validate)
- [protoc-gen-go-svc](https://github.com/neo532/gofr/cmd/protoc-gen-go-svc)

### Quick Start

```bash
# Initialize dev config
make initConfig

# Generate proto code + wire injection
make generate

# Build the api server
make buildApi

# Run with hot-reload
make runApi
```

The server starts at the address configured in `configs/server.yaml` (default `:8000`).

### API Endpoints

The example proto defines:

| Method | Path           | Description    |
| ------ | -------------- | -------------- |
| POST   | `/user`      | Create a user  |
| GET    | `/user/{id}` | Get user by ID |

```bash
curl http://127.0.0.1:8000/user \
  -X POST \
  -H "content-type: application/json" \
  -d '{"id": 1, "name": "alice", "status": "ACTIVE"}'

curl http://127.0.0.1:8000/user/1
```

### Proto Definitions

Services, messages, and enums are defined in `proto/api/`:

```protobuf
service UserApi {
    rpc Post (User) returns (google.protobuf.Empty) {
        option (google.api.http) = { post: "/user"; body: "*"; };
    };
    rpc GetById (GetByIdRequest) returns (User) {
        option (google.api.http) = { get: "/user/{id}"; };
    };
}
```

Validation rules use `protoc-gen-validate` annotations:

```protobuf
message User {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
    string name = 2 [(validate.rules).string = {min_len: 1}];
    user.UserStatus status = 3 [(validate.rules).enum = {defined_only: true}];
}
```

### Adding a New Service

1. Define a proto service in `proto/api/<domain>/v1/`
2. Run `make generate` from the proto directory
3. Implement the generated interface in `internal/service/api/`
4. Register the implementation in `internal/server/api.go`

## Configuration

Configuration is YAML-based with hot-reload. Config structs are auto-generated from YAML files using [config-gen-go-struct](https://github.com/neo532/gokit/cmd/config-gen-go-struct):

```bash
make config
```

All fields use `atomic.Value` or typed atomics for safe concurrent reload. Changes to YAML files are picked up at runtime without restart.

## Middleware

Middleware is registered on the HTTP server in `internal/server/api.go`:

```go
httpSrv := http.NewServer(
    http.Address(cfg.Server.Http.Addr.Load().(string)),
    http.Middleware(
        middleware.Validator(),     // <-- validates incoming requests
        // middleware.Logging(...)
        // middleware.Recovery(...)
    ),
)
```

### Validator Middleware

The `Validator()` middleware automatically calls `Validate()` on incoming requests that implement the `proto.Validator` interface (generated by `protoc-gen-validate`). If validation fails, the request is rejected before reaching the service layer.

## Content-Type Codec

The framework negotiates request decoding and response encoding based on `Content-Type` / `Accept` headers. JSON is built in; custom codecs (XML, YAML, etc.) can be registered:

```go
import "github.com/neo532/gofr/transport/http"

http.RegisterCodec("xml", &http.Codec{
    ContentType: "application/xml",
    Decode:      xml.Unmarshal,
    Encode:      xml.Marshal,
})
```

## Compile-Time Safety

- Service interfaces are enforced at compile time via `var _ pb.UserApiService = (*UserApiService)(nil)`
- The generated `RegisterHTTPServer` panics at startup if any proto-defined service lacks an implementation

## Makefile Targets

| Target            | Description                                   |
| ----------------- | --------------------------------------------- |
| `init`          | Install dev tools (wire, fswatch, config-gen) |
| `config`        | Generate config structs from YAML             |
| `generate`      | Run wire and go:generate                      |
| `build`         | Build all entrypoints                         |
| `buildApi`      | Build api server binary                       |
| `buildConsumer` | Build consumer binary                         |
| `buildScript`   | Build script binary                           |
| `runApi`        | Start api server with hot-reload              |
| `runConsumer`   | Start consumer with hot-reload                |
| `runScript`     | Start script with hot-reload                  |

## Go Workspace

The project uses a `go.work` file for local development alongside sibling modules:

```go
go 1.26.1

use (
    .
    ../gokit
    ../gokit/filepath
    ../gokit/logger/writer/lumberjack
    ./proto
)
```

The gofr framework is fetched from the remote GitHub repository. Set `GOWORK=off` for CI builds to use the go.mod replace directives for local-only modules.

## Dependencies

- **[gofr](https://github.com/neo532/gofr)** — Web framework (HTTP/gRPC transport, middleware)
- **[gokit](https://github.com/neo532/gokit)** — Utility library (ORM, Redis, logging, config)
- **[wire](https://github.com/google/wire)** — Compile-time dependency injection
- **[protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate)** — Request validation from proto annotations
- **[GORM](https://gorm.io/)** — ORM with MySQL driver
