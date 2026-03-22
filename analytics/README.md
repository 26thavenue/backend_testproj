## Event Analytics Engine

A lightweight Go service for ingesting events, registering event types, and querying basic analytics. The default storage is SQLite via GORM with migrations.

---

## Quickstart

```bash
go run ./cmd
```

By default the API listens on `:8080` and uses `./data/analytics.db`. You can override:

```bash
set APP_ADDR=:8080
set DB_PATH=.\data\analytics.db
```

---

## Event Types (Required)

Events are validated against the `event_types` table. Create event types before tracking events:

```bash
curl -X POST http://localhost:8080/api/v1/event-types \
  -H "Content-Type: application/json" \
  -d "{\"event_name\":\"page_view\",\"event_description\":\"Page view events\"}"
```

---

## Track an Event

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d "{\"event_name\":\"page_view\",\"user_id\":\"123\",\"properties\":{\"page\":\"/home\"}}"
```

---

## Get Analytics

```bash
curl "http://localhost:8080/api/v1/analytics?event=page_view"
```

---

## API Docs

- OpenAPI spec: `docs/openapi.yaml`
- Postman collection: `docs/postman_collection.json`

You can import the OpenAPI file into Swagger UI or Postman.

---

## Configuration

Environment variables:

- `APP_ADDR` (default `:8080`)
- `DB_PATH` (default `./data/analytics.db`)

---

## Project Structure

```
/cmd
/internal
  /api
  /config
  /domain
  /events
  /storage
```

---

## Testing

```bash
go test ./...
```

Benchmark:

```bash
go test -bench=BenchmarkServiceTrack ./internal/events
```

---

## Coming Soon (v2)

The current implementation is synchronous and does not include goroutines, batching, or a queue/worker model yet. These are planned for v2:

- Ingestion queue + worker pool
- Batch inserts for higher throughput
- Optional async processing toggle
