# Support-Go Backend

## Run locally

```powershell
docker compose -f deploy/docker-compose.yml up -d postgres

docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000001_create_tickets_table.up.sql

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
