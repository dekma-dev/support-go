# Support-Go Session Log

Этот файл ведем после каждой рабочей сессии.
Цель: в понятном текстовом виде фиксировать, что готово, что реализовано, что в работе и что дальше.

## Формат записи

```text
## Session YYYY-MM-DD #N
- Goal:
- Done:
  - ...
- Current Stage (Roadmap):
  - Stage 1 Foundation: ...
  - Stage 2 Core Ticket Flow: ...
  - Stage 3 Async + Notifications: ...
  - Stage 4 UX + Demo Polish: ...
- Risks / Limits:
  - ...
- Next:
  - ...
```

## Session 2026-02-23 #1
- Goal:
  - Оценить текущее состояние проекта относительно roadmap и подготовить процесс фиксации прогресса.
- Done:
  - Проведен аудит структуры и кода backend/frontend/deploy.
  - Подтверждено: реализован foundation backend + ticket domain (in-memory) + базовый frontend skeleton.
  - Подтверждено: инфраструктурные контейнеры PostgreSQL/Redis/Kafka/Zookeeper описаны в compose.
  - Создан единый журнал сессий: `docs/SESSION_LOG.md`.
- Current Stage (Roadmap):
  - Stage 1 Foundation: частично готово (каркас API/UI есть, auth/миграций/observability нет).
  - Stage 2 Core Ticket Flow: ранняя стадия (ticket endpoints есть, но без PostgreSQL, comments, audit, filters).
  - Stage 3 Async + Notifications: не начато (Kafka в коде не подключен, worker отсутствует).
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - `go` не доступен в PATH текущего окружения, тесты backend локально не запускались.
  - В каталоге проекта пока нет `.git`, публикация в GitHub еще не инициализирована.
- Next:
  - Инициализировать git-репозиторий и выполнить первый push в GitHub.
  - Перейти к следующему приоритету: PostgreSQL repository + migrations для ticket domain.

## Session 2026-02-25 #2
- Goal:
  - Подготовить репозиторий к стабильной командной работе и перевести ticket-хранилище на PostgreSQL.
- Done:
  - Добавлены `.gitignore` и `.gitattributes` в корень проекта.
  - Добавлены миграции таблицы тикетов:
    - `backend/migrations/000001_create_tickets_table.up.sql`
    - `backend/migrations/000001_create_tickets_table.down.sql`
  - Реализован PostgreSQL-репозиторий тикетов: `backend/internal/ticket/postgres/repository.go`.
  - API bootstrap переключен с in-memory на PostgreSQL pool:
    - `backend/cmd/api/main.go`
    - `backend/internal/platform/config/config.go`
    - `backend/go.mod` (добавлен `pgx/v5`)
  - Обновлены инструкции запуска:
    - `backend/README.md`
    - `docs/runbook.md`
- Current Stage (Roadmap):
  - Stage 1 Foundation: расширен (DB-подключение и миграции добавлены).
  - Stage 2 Core Ticket Flow: продолжается (ticket API работает через PostgreSQL, но без RBAC/comments/audit).
  - Stage 3 Async + Notifications: не начато.
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - `go` не доступен в PATH текущего окружения, поэтому `go test ./...` и компиляция не запускались в этой сессии.
  - Миграции пока применяются вручную через `psql`/Docker, автоматизация migration-tooling еще не добавлена.
- Next:
  - Добавить RBAC на операции assign/status.
  - Добавить comments + audit endpoints.
  - Подготовить Kafka producer/consumer каркас.
