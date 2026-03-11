# Support-Go Backend

## Run locally

```powershell
docker compose -f deploy/docker-compose.yml up -d postgres

docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000001_create_tickets_table.up.sql
docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000002_create_ticket_comments_events.up.sql

cd backend
go run ./cmd/api
go run ./cmd/worker
```

API health checks:

- `GET http://localhost:8080/healthz`
- `GET http://localhost:8080/readyz`

Ticket API foundation (PostgreSQL store):

- `POST http://localhost:8080/api/v1/auth/login`
- `POST http://localhost:8080/api/v1/auth/refresh`
- `POST http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/assign`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/status`
- `POST http://localhost:8080/api/v1/tickets/{id}/comments`
- `GET http://localhost:8080/api/v1/tickets/{id}/comments`
- `GET http://localhost:8080/api/v1/tickets/{id}/events`

RBAC (JWT role claim):

- `POST /api/v1/auth/login` accepts `{ "subject": "agent-1", "role": "agent" }`.
- `POST /api/v1/auth/refresh` accepts `{ "refresh_token": "<jwt>" }` and rotates both tokens.
- Set `JWT_SECRET` in backend environment.
- Pass `Authorization: Bearer <access_token>`; API accepts only access tokens and uses `role` claim (`client`, `agent`, `admin`).
- `PATCH /api/v1/tickets/{id}/assign` and `PATCH /api/v1/tickets/{id}/status` require `agent/admin` role.
- Internal comments (`is_internal=true`) can be created only by `agent/admin`.
- Internal comments are hidden from `client` in `GET /api/v1/tickets/{id}/comments`.
- Invalid bearer token returns `401 Unauthorized`.

Async foundation:

- API publishes domain events to Kafka topics:
  - `support.ticket.events`
  - `support.comment.events`
- Notification worker skeleton consumes both topics:
  - run with `go run ./cmd/worker`
  - retry and DLQ are configurable:
    - `NOTIFICATION_RETRY_MAX` (default `3`)
    - `NOTIFICATION_RETRY_BACKOFF_MS` (default `500`)
    - `KAFKA_NOTIFICATION_DLQ_TOPIC` (default `support.notification.dlq`)
