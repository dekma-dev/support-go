# Progress Log

## 2026-02-19

### Done
- Initialized backend API foundation with health endpoints (`/healthz`, `/readyz`).
- Implemented ticket domain foundation in Go (in-memory):
  - model + validation: `backend/internal/ticket/model.go`
  - service layer: `backend/internal/ticket/service.go`
  - repository (memory): `backend/internal/ticket/memory/repository.go`
  - HTTP handlers/routes: `backend/internal/ticket/handler.go`
- Wired ticket routes in API bootstrap: `backend/cmd/api/main.go`.
- Added service tests: `backend/internal/ticket/service_test.go`.
- Updated backend docs with ticket endpoints: `backend/README.md`.

### Current API (ticket foundation)
- `POST /api/v1/tickets`
- `GET /api/v1/tickets`
- `GET /api/v1/tickets/{id}`
- `PATCH /api/v1/tickets/{id}`
- `PATCH /api/v1/tickets/{id}/assign`
- `PATCH /api/v1/tickets/{id}/status`

### Notes / Limits
- Local environment in this session does not have `go` binary in PATH, so `go test ./...` could not be executed here.

### Next Step (planned)
- Replace in-memory repository with PostgreSQL + migrations.
- Add role-based authorization checks for assign/status operations.
- Connect frontend flows to ticket API.

## 2026-02-25

### Done
- Added repository hygiene files:
  - `.gitignore`
  - `.gitattributes`
- Replaced ticket storage from in-memory to PostgreSQL:
  - migrations: `backend/migrations/000001_create_tickets_table.up.sql`, `backend/migrations/000001_create_tickets_table.down.sql`
  - repository: `backend/internal/ticket/postgres/repository.go`
  - API wiring/config: `backend/cmd/api/main.go`, `backend/internal/platform/config/config.go`
- Added role checks for sensitive ticket operations:
  - only `agent/admin` can call `PATCH /api/v1/tickets/{id}/assign`
  - only `agent/admin` can call `PATCH /api/v1/tickets/{id}/status`
  - implementation: `backend/internal/ticket/authorization.go`, `backend/internal/ticket/handler.go`
- Added handler tests for RBAC behavior:
  - `backend/internal/ticket/handler_test.go`
- Updated docs:
  - `backend/README.md`
  - `docs/runbook.md`
  - `docs/SESSION_LOG.md`

### Verification
- `go mod tidy` executed successfully.
- `go test ./...` executed successfully in `backend/`.
