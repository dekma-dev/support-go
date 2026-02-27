# Support-Go Backend

## Run locally

```powershell
docker compose -f deploy/docker-compose.yml up -d postgres

docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000001_create_tickets_table.up.sql
docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000002_create_ticket_comments_events.up.sql

cd backend
go run ./cmd/api
```

API health checks:

- `GET http://localhost:8080/healthz`
- `GET http://localhost:8080/readyz`

Ticket API foundation (PostgreSQL store):

- `POST http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/assign`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/status`
- `POST http://localhost:8080/api/v1/tickets/{id}/comments`
- `GET http://localhost:8080/api/v1/tickets/{id}/comments`
- `GET http://localhost:8080/api/v1/tickets/{id}/events`

RBAC (temporary header-based):

- For `PATCH /api/v1/tickets/{id}/assign` and `PATCH /api/v1/tickets/{id}/status`, set header `X-User-Role: agent` or `X-User-Role: admin`.
- `client` role (or missing role) receives `403 Forbidden` for these two operations.
- Internal comments (`is_internal=true`) can be created only by `agent/admin`.
- Internal comments are hidden from `client` in `GET /api/v1/tickets/{id}/comments`.
