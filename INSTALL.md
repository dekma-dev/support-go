# Локальная установка (Windows)

## Требования

- Docker Desktop
- Go 1.26+
- Node.js 24+

## 1. Поднять инфраструктуру

```
cd C:\wamp\www\support.go\deploy
docker compose up -d
```

## 2. Накатить миграции

```
docker exec -i support_go_postgres psql -U support -d support_go < C:\wamp\www\support.go\backend\migrations\000001_create_tickets_table.up.sql
docker exec -i support_go_postgres psql -U support -d support_go < C:\wamp\www\support.go\backend\migrations\000002_create_ticket_comments_events.up.sql
docker exec -i support_go_postgres psql -U support -d support_go < C:\wamp\www\support.go\backend\migrations\000003_create_users_table.up.sql
```

## 3. Запустить бэкенд

```
cd C:\wamp\www\support.go\backend
copy .env.example .env
go run ./cmd/api
```

## 4. Запустить фронтенд (в отдельном терминале)

```
cd C:\wamp\www\support.go\frontend
npm install
npm run dev
```

## 5. Открыть

- Фронтенд: http://localhost:5173
- API health: http://localhost:8080/healthz
