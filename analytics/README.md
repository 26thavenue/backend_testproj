## 📊 Event Analytics Engine (Plug & Play)

A scalable, extensible event tracking system built in Go that lets you **ingest, process, and query events** with minimal setup — while allowing you to **swap out the database layer** to fit your stack.

---

## 🚀 Overview

This project is a lightweight analytics engine designed for:

* High-volume event ingestion
* Real-time + batch processing
* Flexible storage backends
* Easy integration into existing systems

You can clone this repo, plug in your preferred database, and start tracking events immediately using pre-defined schemas.

---

## ✨ Features

* ⚡ High-performance event ingestion API
* 🔌 Pluggable database layer (Postgres, ClickHouse, etc.)
* 🧵 Built with Go concurrency (goroutines + channels)
* 📦 Predefined event types (ready to use)
* 📊 Basic analytics endpoints (counts, trends)
* 🧱 Modular architecture (easy to extend)
* 🔄 Async processing support (queue/worker model ready)

---

## 🧩 Predefined Events

Out of the box, the system supports:

* `user_signup`
* `user_login`
* `page_view`
* `button_click`
* `purchase`

Each event follows a standard structure:

```json
{
  "event_name": "page_view",
  "user_id": "12345",
  "timestamp": "2026-03-19T12:00:00Z",
  "properties": {
    "page": "/home",
    "device": "mobile"
  }
}
```

You can extend or override these easily.

---

## 🏗️ Architecture

```
Client → API → Event Queue → Workers → Storage
                         ↘ Cache (optional)
```

### Core Components

* **API Layer**: Handles event ingestion and queries
* **Queue Layer**: Buffers events (optional but recommended)
* **Worker Layer**: Processes and persists events
* **Storage Layer**: Pluggable database interface
* **Cache Layer (optional)**: Real-time metrics

---

## ⚙️ Installation

```bash
git clone https://github.com/your-org/event-analytics-engine.git
cd event-analytics-engine
go mod tidy
```

---

## 🔌 Configuration

All configuration is environment-based:

```env
PORT=8080
DB_DRIVER=postgres   # or clickhouse, mysql, etc.
DB_DSN=your_database_connection_string
ENABLE_ASYNC=true
```

---

## 🧠 Pluggable Database Design

The system uses a storage interface:

```go
type EventStore interface {
    SaveEvent(event Event) error
    GetEvents(filter Filter) ([]Event, error)
    Aggregate(query Query) (Result, error)
}
```

### To plug in your own DB:

1. Implement the `EventStore` interface
2. Register your implementation
3. Update config (`DB_DRIVER`)

---

## 📡 API Endpoints

### Ingest Event

```http
POST /event
```

Body:

```json
{
  "event_name": "user_signup",
  "user_id": "123",
  "properties": {}
}
```

---

### Get Analytics

```http
GET /analytics?event=page_view&from=2026-03-01&to=2026-03-19
```

Response:

```json
{
  "event": "page_view",
  "count": 10234
}
```

---

## 🧵 Concurrency Model

* Events are processed using goroutines
* Channels handle communication between ingestion and workers
* Optional batching improves DB performance

---

## 🔄 Async Processing (Optional)

Enable async mode:

```env
ENABLE_ASYNC=true
```

This will:

* Queue incoming events
* Process them in background workers
* Improve ingestion throughput

---

## 📊 Extending Analytics

You can easily add:

* Funnel analysis
* Cohort tracking
* Retention metrics
* Real-time dashboards

---

## 🧪 Running Locally

```bash
go run main.go
```

Test ingestion:

```bash
curl -X POST http://localhost:8080/event \
-H "Content-Type: application/json" \
-d '{"event_name":"page_view","user_id":"1"}'
```

---

## 📁 Project Structure

```
/cmd            → App entrypoint  
/internal
  /api          → HTTP handlers  
  /events       → Event models & validation  
  /storage      → DB interface + implementations  
  /workers      → Async processors  
  /config       → Environment config  
/pkg            → Shared utilities  
```

---

## 🔐 Production Considerations

* Add rate limiting
* Use a message queue (Kafka, RabbitMQ)
* Enable retries and dead-letter queues
* Add observability (logs, metrics, tracing)

---

## 🤝 Contributing

Contributions are welcome:

1. Fork the repo
2. Create a feature branch
3. Submit a PR

---

## 📌 Use Cases

* SaaS product analytics
* Internal event tracking
* Growth & marketing insights
* Backend telemetry systems

---

## 🧭 Roadmap

* [ ] Dashboard UI
* [ ] Multi-tenant support
* [ ] Webhook integrations
* [ ] Event schema validation system
* [ ] Plugin marketplace for storage backends

---

## 🪪 License

MIT License — free to use and modify.

---

## 💡 Philosophy

> Build once. Plug anywhere. Scale when needed.
