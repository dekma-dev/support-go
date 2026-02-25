# Support-Go Backend

## Run locally

```powershell
cd backend
go run ./cmd/api
```

API health checks:

- `GET http://localhost:8080/healthz`
- `GET http://localhost:8080/readyz`

Ticket API foundation (in-memory store):

- `POST http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets`
- `GET http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/assign`
- `PATCH http://localhost:8080/api/v1/tickets/{id}/status`
