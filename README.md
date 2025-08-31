# Auth Service

Example authentication service featuring:
- JWT-based authorization
- Asymmetric-key JWT signing; JWKS intentionally omitted for simplicity
- Redis-backed storage for issued refresh tokens
- Refresh token rotation on each successful refresh
- Structured logging and graceful shutdown
- Integration and unit tests

## Quick start

- Run the service:
  ```bash
  go run main.go serve --config=config.yaml
  ```

- Bring up Redis and Postgres with Docker Compose:
  ```bash
  docker compose -f docker-compose.yaml up -d
  ```

## Features

- Access and refresh tokens with independent TTLs
- Routes exempt from authentication (configurable)
- Redis-backed refresh token store for revocation and rotation
- Postgres connectivity for application data
- Configurable logger

## Configuration

Configuration file path is provided via the `--config` flag. Example filename: `config.yaml`.

The following keys are supported (example values shown).

```yaml
auth: 
  jwt: 
    secret: "very secret key" # Private key content or shared secret (example uses a string for simplicity) 
    issuer: "auth-service" # Issuer claim (iss) value 
    accessTokenTTL: 5s # Access token time-to-live 
    refreshTokenTTL: 1h # Refresh token time-to-live 
    noAuthRoutes: # Routes that bypass authentication middleware 
      - POST /v1/auth/register 
      - POST /v1/auth/login 
      - POST /v1/auth/refresh
http: 
  port: 8080 # HTTP listen port 
  shutdownTimeout: 3s # Graceful shutdown timeout
db: 
  requestTimeout: 3s # Database operation timeout postgres: 
  conn: 
    postgres://{POSTGRES_USER}:{POSTGRES_PASSWORD}@{POSTGRES_HOST}:{POSTGRES_PORT}/{POSTGRES_DB_NAME}?pool_max_conns=10&sslmode={POSTGRES_SSL_MODE} # Postgres connection string with environment interpolation
redis: 
  host: {REDIS_HOST} # Redis host 
  port: {REDIS_PORT} # Redis port
logger: 
  level: info # zap log level: debug, info, warn, error 
  disableCaller: false # omit caller information if true 
  disableStacktrace: true # omit stack traces if true 
  enableDateTime: true # adds human-readable date time to log entries
```

Notes:
- Asymmetric signing without JWKS is assumed for simplicity. The `auth.jwt.secret` field can hold a PEM-encoded private key or a shared secret. Replace with private key material when using an asymmetric algorithm; JWKS distribution is intentionally omitted to keep configuration minimal.
- Ensure environment variables referenced in the config are exported prior to starting the service.

## Token lifecycle

- Access tokens: short-lived, signed JWTs intended for API authorization.
- Refresh tokens: longer-lived, stored in Redis, rotated on refresh. Old refresh tokens are invalidated upon successful rotation.

## Docker Compose

Run existing compose setup:

```bash
docker compose -f test/docker-compose.yaml up -d --build
```

## Environment variables for local run

Example environment setup aligned with the configuration:

```bash
export POSTGRES_USER=auth 
export POSTGRES_PASSWORD=auth 
export POSTGRES_HOST=localhost 
export POSTGRES_PORT=5432 
export POSTGRES_DB_NAME=auth 
export POSTGRES_SSL_MODE=disable
export REDIS_HOST=localhost 
export REDIS_PORT=6379
```

## Project structure
```text
.
├── cmd/                         # CLI entrypoints (cobra commands)
├── internal/                    # Private application modules
│   ├── domain/                  # Core domain DTOs and errors
│   ├── http/                    # HTTP service, routing, middleware, handlers
│   │   ├── handlers/            # Request handlers (+ tests and mocks)
│   │   └── middleware/          # HTTP middlewares
│   ├── jwt/                     # JWT issuer, verifier, authenticator, storage
│   ├── redis/                   # Redis integration
│   └── repository/              # Persistence layer
├── migrations/                  # Database migrations
├── test/                        # Integration tests
├── config.yaml
└── main.go                      
```

## Most important dependencies
- [uber/zap](https://github.com/uber-go/zap) — fast, structured logging
- [sarulabs/di](https://github.com/sarulabs/di) — DI container
- [spf13/cobra](https://github.com/spf13/cobra) — CLI commands and flags
- [spf13/viper](https://github.com/spf13/viper) — configuration loading and environment binding
- [stretchr/testify](https://github.com/stretchr/testify) and [vektra/mockery](https://github.com/vektra/mockery) — powerful tools for testing 
- [jackc/pgx](https://github.com/jackc/pgx) — PostgreSQL Driver and Toolkit
