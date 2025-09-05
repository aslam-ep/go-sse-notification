# go-sse-notification

A scalable server-sent events (SSE) notification system with per-user streams, Redis Pub/Sub for horizontal scaling, and Prometheus metrics. Designed for real-time notifications in distributed environments, with Docker and Nginx support for production deployment.

## Features

- **Per-user SSE streams:** Real-time notifications delivered to specific users.
- **Send notification endpoint:** Simple HTTP API to trigger notifications.
- **Redis Pub/Sub:** Enables horizontal scaling and message fan-out across instances.
- **Notification history:** Missed messages are replayed on reconnect.
- **Prometheus metrics:** Track active clients, sent messages, and dropped messages.
- **Dockerized:** Easy to run locally or in production.
- **Nginx reverse proxy:** Handles SSE-specific proxying and buffering.

## Architecture
Client <--SSE-- Nginx <--HTTP-- App (Go) <--Redis Pub/Sub--> Other App Instances

- Clients connect via `/notifications?userID=...` for SSE.
- Notifications are sent via `POST /send`.
- Redis is used for Pub/Sub and notification history.
- Nginx proxies and optimizes SSE connections.

## Endpoints

| Method | Path              | Description                       |
|--------|-------------------|-----------------------------------|
| GET    | `/ping`           | Health check                      |
| GET    | `/notifications`  | SSE stream for a user             |
| POST   | `/send`           | Send notification to a user       |
| GET    | `/metrics`        | Prometheus metrics                |

### Example: Send Notification

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{"userID": "user123", "message": "Hello, User!"}'
```

Example: Receive Notifications (SSE)
```bash
curl -N "http://localhost:8080/notifications?userID=user123"
```

## Configuration
Configuration is loaded from environment variables or a config file. See .sample.env for example values.

### Running Locally
Prerequisites
- Docker
- Docker Compose

Start All Services
```bash
make up
```

- App: http://localhost:8080
- Nginx Proxy: http://localhost:8081
- Redis: localhost:6379

Stop Services
```bash
make down
```

## Metrics
Exposed at /metrics (Prometheus format):

- sse_active_clients — Current number of SSE clients
- sse_message_sent_total{user="..."} — Total messages sent per user
- sse_dropped_message_total — Dropped messages due to slow clients

## Deployment
- Use the provided Dockerfile and docker-compose.yml for containerized deployment.
- Nginx is configured to optimize SSE connections and proxy to the app.
