# Support-Go Runbook

## 1) Local development quick start

```powershell
docker compose -f deploy/docker-compose.yml up -d

docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000001_create_tickets_table.up.sql
docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000002_create_ticket_comments_events.up.sql

cd backend
go run ./cmd/api
go run ./cmd/worker
```

Checks:

- `http://localhost:8080/healthz`
- `http://localhost:8080/readyz`
- `http://localhost:5173`

## 2) Production deployment checklist

Prerequisites:

- Linux server with Docker Engine + Docker Compose Plugin.
- Public DNS zone for your domain.
- Ports `80/tcp` and `443/tcp` open to the server.
- Project files copied to server (for example to `/opt/support-go`).

### 2.1 Prepare environment file

```bash
cd /opt/support-go
cp deploy/.env.prod.example deploy/.env.prod
```

Edit `deploy/.env.prod` and set real values:

- `DOMAIN` (example: `support.example.com`)
- `TLS_EMAIL`
- `POSTGRES_PASSWORD` (strong secret)
- optional notification retry settings

### 2.2 Point subdomain to server

Create DNS records:

- `A` record: `support.example.com -> <server_ipv4>`
- optional `AAAA` for IPv6

Wait for DNS propagation, then verify:

```bash
dig +short support.example.com
```

### 2.3 Build and start containers

```bash
cd /opt/support-go
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml up -d --build
```

### 2.4 Run DB migrations

```bash
cd /opt/support-go
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml exec -T postgres \
  sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < backend/migrations/000001_create_tickets_table.up.sql

docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml exec -T postgres \
  sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < backend/migrations/000002_create_ticket_comments_events.up.sql
```

## 3) Post-deploy verification

Health checks:

```bash
curl -fsS https://support.example.com/healthz
curl -fsS https://support.example.com/readyz
```

API smoke test:

```bash
curl -i -X POST "https://support.example.com/api/v1/tickets" \
  -H "Content-Type: application/json" \
  -d '{"title":"prod smoke","description":"deployment check"}'
```

Container status:

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml ps
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml logs --tail=100 api
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml logs --tail=100 worker
```

## 4) Rollback and recovery

Quick rollback to previous image tag (if you use pinned tags):

1. Update image tags in compose file to last known good versions.
2. Redeploy:
   `docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml up -d`

If only app code broke and DB schema is still compatible:

1. Checkout previous git revision.
2. Redeploy with `--build`.

If migration caused regression:

1. Stop write traffic (maintenance mode at reverse proxy level).
2. Apply matching down migration manually.
3. Restore DB backup if needed.

## 5) Operations notes

- TLS is auto-managed by Caddy (`deploy/Caddyfile`).
- Public entrypoint is `caddy`; `api`, `worker`, `postgres`, `kafka`, `redis` stay internal.
- Current RBAC is still header-based (`X-User-Role`), JWT auth middleware is not implemented yet.
