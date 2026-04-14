# Локальная установка (Windows)

Одноразовая настройка окружения. После этого для повседневного запуска используй [RUN.md](RUN.md).

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
docker exec -i support_go_postgres psql -U support -d support_go < C:\wamp\www\support.go\backend\migrations\000004_seed_users.up.sql
```

## 3. Создать .env файл для бэкенда

```
cd C:\wamp\www\support.go\backend
copy .env.example .env
```

## 4. Установить зависимости фронтенда

```
cd C:\wamp\www\support.go\frontend
npm install
```

## Готово

Настройка завершена. Для запуска проекта см. [RUN.md](RUN.md).
