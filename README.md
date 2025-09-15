# Polyschedule Backend (stub, module)

Minimal Kafka-compatible service in Go. Can be used standalone or as a module (git submodule) in another repository.

## Environment variables

- `KAFKA_BROKERS` — comma-separated brokers (default `localhost:9092`)
- `KAFKA_GROUP_ID` — consumer group id (default `polyschedule-backend`)
- `KAFKA_INPUT_TOPIC` — input topic for requests (default `schedule.requests`)
- `KAFKA_OUTPUT_TOPIC` — output topic for responses (default `schedule.results`)
- `HTTP_ADDR` — http server address (default `:8080`)
- `LOG_LEVEL` — log level (`debug|info|warn|error`)

`.env` autoload is supported (via godotenv) — if present near the binary or project root, variables will be loaded.

## Run locally

```bash
# Start local broker (Redpanda — Kafka-compatible)
docker compose up -d

# Start the service
export KAFKA_BROKERS=localhost:9092
export KAFKA_INPUT_TOPIC=schedule.requests
export KAFKA_OUTPUT_TOPIC=schedule.results
export KAFKA_GROUP_ID=polyschedule-backend
export HTTP_ADDR=:8080
export LOG_LEVEL=debug

go run ./cmd/polyschedule-backend

# Or via docker-compose (with .env):
# create .env from the example below and start the stack
docker-compose up -d
```

## Quick check

Open Redpanda Console: http://localhost:8081

Or use rpk/Kafka client to send a request message:

```bash
# Sample request (stub):
cat <<'JSON' | rpk topic produce schedule.requests -X brokers=localhost:9092 -k chat:123
{"chat_id":123, "query":"на завтра"}
JSON

# Read response:
rpk topic consume schedule.results -X brokers=localhost:9092 -n 1
```

### .env example

```dotenv
KAFKA_BROKERS=redpanda:9092
KAFKA_GROUP_ID=polyschedule-backend
KAFKA_INPUT_TOPIC=schedule.requests
KAFKA_OUTPUT_TOPIC=schedule.results
HTTP_ADDR=:8080
LOG_LEVEL=debug
APP_HTTP_PORT=8080
REDPANDA_KAFKA_PORT=9092
REDPANDA_CONSOLE_PORT=8081
```

## Use as a module (git submodule)

```go
import (
	"context"
	"github.com/Koshsky/polyschedule-backend/internal/config"
	"github.com/Koshsky/polyschedule-backend/pkg/polyschedule"
)

cfg, _ := config.Load()
svc := polyschedule.New(cfg)
errCh, _ := svc.Start(context.Background())
// ... handle signals and stop: svc.Stop(ctx)
_ = errCh
```

## Structure

- `cmd/polyschedule-backend` — entrypoint
- `internal/config` — env-based configuration
- `internal/kafka` — kafka-go consumer/producer
- `internal/processor` — stub processor
- `internal/httpserver` — HTTP healthcheck `/healthz`
- `pkg/polyschedule` — public package for modular usage


