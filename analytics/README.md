# Analytics API

A custom, highly-performant analytics engine written in Go. This service is designed to be a lightweight alternative to heavyweight tracking platforms like Amplitude, PostHog, and Google Analytics. 

It handles event ingestion, multi-event funnel querying, and dynamic data aggregation with an incredibly low memory footprint and high throughput.

## Tech Stack
- **Language**: Go (v1.22+)
- **API Framework**: [Fiber](https://gofiber.io/) - Ultra-fast web framework sitting on top of `fasthttp`.
- **Database Engine**: [SQLite](https://sqlite.org/) via [`glebarez/sqlite`](https://github.com/glebarez/sqlite) (Pure-Go SQLite driver, no CGO dependencies).
- **ORM**: [GORM](https://gorm.io/)

## Architecture 
The codebase is structured around a blend of **Domain-Driven Design (DDD)** and **Clean Architecture (Ports and Adapters)**. 

Feature files are packaged by their domain boundary to keep the dependency graph clean and scalable:
- **`domain/`**: Houses the raw entity models and structs (`Event`, `EventType`, `TrackRequest`, etc.), serving as the absolute core of the application without external dependencies.
- **`events/`**: The Core Domain for analytical tracking. Houses its own HTTP `Handler` implementation and Business Logic `Service`.
- **`storage/`**: The Adapter layer. Implements the data persistence interfaces defined by the domain services using GORM and SQLite.
- **`api/`**: The Presentation layer that wires up HTTP routes.

## Performance Highlights
Built for raw speed and minimal resource usage:
- Core tracking services benchmarked at **~23 nanoseconds per operation**.
- Achieves **0 B/op memory allocations** in core logic paths.
- Compiles down into a single statically linked binary (no messy containers or C compiler requirements thanks to the CGO-free setup).

## Advanced Features

### Multi-Event Funnel Analytics
The API natively handles complex funnel generation out-of-the-box. The frontend can query aggregated drop-off data for multiple event types within a single request instead of iterating multiple calls.

```http
GET /api/v1/analytics?event=page_view&event=signup&event=purchase
```
**Response:**
```json
[
  {
    "event": "page_view",
    "count": 1250,
    "from": "2026-03-01T00:00:00Z",
    "to": "2026-03-31T00:00:00Z"
  },
  {
    "event": "signup",
    "count": 480,
    "from": "2026-03-01T00:00:00Z",
    "to": "2026-03-31T00:00:00Z"
  },
  {
    "event": "purchase",
    "count": 150,
    "from": "2026-03-01T00:00:00Z",
    "to": "2026-03-31T00:00:00Z"
  }
]
```

## Running the Application

### Setup Development Server (Hot Reloading)
We use [Air](https://github.com/cosmtrek/air) to handle live hot-reloading for both Go code and Swagger UI configurations.

1. Install Air globally if you do not have it:
```bash
go install github.com/cosmtrek/air@latest
```
2. Start the watcher:
```bash
air
```

### Static Build
If you're deploying to production:
```bash
go build -ldflags="-s -w" -o ./bin/main server.go
./bin/main
```

## API Documentation
The API contains integrated Swagger UI definitions. When the server is running, visit:
- **Swagger UI**: `http://localhost:8080/api/v1/swagger`
- **Health Check**: `http://localhost:8080/api/v1/health`
