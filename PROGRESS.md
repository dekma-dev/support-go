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

## 2026-02-26

### Done
- Implemented comments and audit functionality for tickets:
  - new domain models: `backend/internal/ticket/activity.go`
  - service methods: add/list comments, list ticket events + audit recording
  - endpoints:
    - `POST /api/v1/tickets/{id}/comments`
    - `GET /api/v1/tickets/{id}/comments`
    - `GET /api/v1/tickets/{id}/events`
- Added PostgreSQL repositories:
  - `backend/internal/ticket/postgres/comment_repository.go`
  - `backend/internal/ticket/postgres/audit_repository.go`
- Added DB migration:
  - `backend/migrations/000002_create_ticket_comments_events.up.sql`
  - `backend/migrations/000002_create_ticket_comments_events.down.sql`
- Wired dependencies in API bootstrap:
  - `backend/cmd/api/main.go`
- Extended handler tests:
  - `backend/internal/ticket/handler_test.go`
  - coverage includes role-based internal comment visibility and events endpoint.

### Notes
- RBAC is still temporary header-based (`X-User-Role`), no JWT middleware yet.

## 2026-02-27

### Done
- Started Stage 3 (Async foundation):
  - Added domain event model + publisher interface:
    - `backend/internal/ticket/event.go`
  - Extended ticket service to publish Kafka domain events for:
    - `ticket.created`
    - `ticket.updated`
    - `ticket.assigned`
    - `ticket.status.changed`
    - `comment.added`
    - file: `backend/internal/ticket/service.go`
- Added Kafka platform module:
  - publisher + noop publisher + broker parser:
    - `backend/internal/platform/kafka/publisher.go`
  - tests:
    - `backend/internal/platform/kafka/publisher_test.go`
- Wired Kafka publisher in API bootstrap:
  - `backend/cmd/api/main.go`
  - falls back to noop publisher when `KAFKA_BROKERS` is empty.
- Added notification worker skeleton (consumer process):
  - `backend/cmd/worker/main.go`
  - consumes `support.ticket.events` and `support.comment.events`
  - logs received events and commits offsets.
- Extended config/env for worker:
  - `backend/internal/platform/config/config.go`
  - `backend/.env.example` (`KAFKA_NOTIFICATION_GROUP`)
- Updated docs:
  - `backend/README.md`
  - `docs/runbook.md`

### Verification
- `go mod tidy` executed successfully.
- `go test ./...` executed successfully in `backend/`.

## 2026-02-27 (part 2)

### Done
- Added notification handling layer:
  - `backend/internal/notification/service.go`
  - maps domain events to notification messages and sends via `Sender`.
- Added tests for notification mapping/flow:
  - `backend/internal/notification/service_test.go`
- Upgraded worker processing:
  - retry with backoff for notification handling
  - DLQ publish on decode/processing failure
  - files:
    - `backend/cmd/worker/main.go`
    - `backend/cmd/worker/main_test.go`
- Extended worker config/env:
  - `NOTIFICATION_RETRY_MAX`
  - `NOTIFICATION_RETRY_BACKOFF_MS`
  - `KAFKA_NOTIFICATION_DLQ_TOPIC`
  - files:
    - `backend/internal/platform/config/config.go`
    - `backend/.env.example`
- Improved event payloads from ticket service for downstream notifications:
  - includes requester/assignee/author context
  - file: `backend/internal/ticket/service.go`
- Updated docs:
  - `backend/README.md`
  - `docs/runbook.md`

### Verification
- `go test ./...` executed successfully in `backend/`.

### Next Session Priority
- Start with deployment to user servers + subdomain setup.
- Produce a clear deployment checklist and runbook for repeatable rollout.

## 2026-03-06

### Done
- Started deployment stage deliverables for repeatable server rollout:
  - production compose stack: `deploy/docker-compose.prod.yml`
  - reverse proxy + TLS routing: `deploy/Caddyfile`
  - frontend nginx SPA config: `deploy/nginx/frontend.conf`
  - production env template: `deploy/.env.prod.example`
- Added container build definitions:
  - `backend/Dockerfile` (builds both `api` and `worker` binaries)
  - `frontend/Dockerfile` (Vite build + nginx runtime)
- Reworked runbook with deployment checklist and verification flow:
  - `docs/runbook.md`
  - includes: server prereqs, DNS/subdomain setup, deploy commands, migration step, post-deploy checks, rollback notes.

### Notes / Limits
- Migration tooling is still manual SQL execution via `psql`; no automated migration runner wired yet.
- Auth remains header-based RBAC; JWT middleware is still pending.

### Next Step (planned)
- Add frontend integration with ticket API and environment-based API base URL.
- Implement JWT auth middleware and replace temporary header-based role handling.

## 2026-03-08

### Done
- Implemented frontend ticket API integration:
  - added typed API client with env-based base URL:
    - `frontend/src/api.ts`
  - integrated tickets page with:
    - list tickets (`GET /api/v1/tickets`)
    - create ticket (`POST /api/v1/tickets`)
    - loading/error handling and manual refresh
    - file: `frontend/src/App.tsx`
- Updated frontend styles for operational form/list UI:
  - `frontend/src/styles.css`
- Added frontend local env template and docs:
  - `frontend/.env.example`
  - `frontend/README.md`
- Improved repo hygiene for frontend TS build artifacts:
  - `.gitignore` now excludes `frontend/*.tsbuildinfo`

### Verification
- Installed frontend dependencies with `npm install`.
- Successfully built frontend with `npm run build`.

### Next Step (planned)
- Implement JWT auth middleware in backend and replace temporary header-based RBAC handling.
