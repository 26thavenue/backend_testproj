## 📊 Event Analytics Engine (Plug & Play)

A scalable, extensible event tracking system built in Go that lets you **ingest, process, and query events** with minimal setup — while allowing you to **swap out the database layer** to fit your stack.

---

## 🚀 Overview

This project is a lightweight analytics engine designed for:

- High-volume event ingestion
- Real-time + batch processing
- Flexible storage backends
- Easy integration into existing systems

You can clone this repo, plug in your preferred database, and start tracking events immediately using pre-defined schemas.

---

## ✨ Features

- ⚡ High-performance event ingestion API
- 🔌 Pluggable database layer (Postgres, ClickHouse, etc.)
- 🧵 Built with Go concurrency (goroutines + channels)
- 📦 Predefined event types (ready to use)
- 📊 Basic analytics endpoints (counts, trends)
- 🧱 Modular architecture (easy to extend)
- 🔄 Async processing support (queue/worker model ready)

---

## 🧩 Predefined Events

Out of the box, the system supports:

- `user_signup`
- `user_login`
- `page_view`
- `button_click`
- `purchase`

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

-> Workers , Event Queues
-> Authentication
-> Streaming
-> Add cahcing for event names so they can check the isValid method from there without making a db call
-> Funnels