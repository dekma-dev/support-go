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

## Session 2026-02-25 #3
- Goal:
  - Реализовать RBAC для операций assign/status и восстановить рабочее Go-окружение для проверок.
- Done:
  - Добавлен модуль role-check:
    - `backend/internal/ticket/authorization.go`
  - В `backend/internal/ticket/handler.go` добавлены проверки ролей для:
    - `PATCH /api/v1/tickets/{id}/assign`
    - `PATCH /api/v1/tickets/{id}/status`
  - Добавлены тесты RBAC хендлера:
    - `backend/internal/ticket/handler_test.go`
  - Выполнен `go mod tidy`, сгенерирован `backend/go.sum`.
  - Пройдено `go test ./...` в `backend/` (успешно).
  - Go обнаружен в `C:\Program Files\Go\bin\go.exe`; путь добавлен в пользовательский PATH.
- Current Stage (Roadmap):
  - Stage 1 Foundation: частично закрыт (DB + тестовый цикл работают, но auth/JWT и observability еще впереди).
  - Stage 2 Core Ticket Flow: в работе (ticket flow + RBAC для assign/status готовы, comments/audit/filters отсутствуют).
  - Stage 3 Async + Notifications: не начато.
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - RBAC сейчас header-based (временный механизм), без полноценного auth/JWT middleware.
- Next:
  - Реализовать comments endpoints.
  - Реализовать audit trail endpoints.
  - Подготовить каркас Kafka producer/consumer.

## Session 2026-02-26 #4
- Goal:
  - Продвинуть Core Ticket Flow: добавить comments и audit endpoints.
- Done:
  - Добавлены доменные сущности:
    - `backend/internal/ticket/activity.go`
  - Расширен сервис:
    - запись audit-событий при create/update/assign/status/comment
    - методы `AddComment`, `ListComments`, `ListEvents`
    - файл: `backend/internal/ticket/service.go`
  - Добавлены PostgreSQL-репозитории:
    - `backend/internal/ticket/postgres/comment_repository.go`
    - `backend/internal/ticket/postgres/audit_repository.go`
  - Добавлена миграция comments/events:
    - `backend/migrations/000002_create_ticket_comments_events.up.sql`
    - `backend/migrations/000002_create_ticket_comments_events.down.sql`
  - Добавлены/обновлены API endpoints:
    - `POST /api/v1/tickets/{id}/comments`
    - `GET /api/v1/tickets/{id}/comments`
    - `GET /api/v1/tickets/{id}/events`
    - файл: `backend/internal/ticket/handler.go`
  - Обновлен bootstrap зависимостей:
    - `backend/cmd/api/main.go`
  - Добавлены тесты на comments/events + visibility:
    - `backend/internal/ticket/handler_test.go`
  - Обновлена документация:
    - `backend/README.md`
    - `docs/runbook.md`
    - `PROGRESS.md`
- Current Stage (Roadmap):
  - Stage 1 Foundation: частично готово.
  - Stage 2 Core Ticket Flow: существенно продвинут (ticket CRUD/status/assign/comments/events есть, но filters и полноценный auth/JWT отсутствуют).
  - Stage 3 Async + Notifications: не начато.
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - Аутентификация/авторизация пока header-based и не заменяет полноценный JWT+roles middleware.
- Next:
  - Добавить auth/JWT слой с middleware.
  - Добавить ticket filters/search в `GET /api/v1/tickets`.
  - Начать Stage 3: Kafka producer + notification worker skeleton.

## Session 2026-02-27 #5
- Goal:
  - Начать Stage 3 roadmap: добавить асинхронный каркас Kafka (producer + worker consumer).
- Done:
  - Добавлена доменная модель события и интерфейс publisher:
    - `backend/internal/ticket/event.go`
  - Сервис тикетов теперь публикует события в Kafka (через интерфейс) при create/update/assign/status/comment:
    - `backend/internal/ticket/service.go`
  - Добавлен платформенный Kafka модуль:
    - `backend/internal/platform/kafka/publisher.go`
    - `backend/internal/platform/kafka/publisher_test.go`
  - API bootstrap подключен к Kafka publisher (с fallback на noop при пустом `KAFKA_BROKERS`):
    - `backend/cmd/api/main.go`
  - Добавлен Notification Worker skeleton:
    - `backend/cmd/worker/main.go`
    - читает `support.ticket.events` и `support.comment.events`, логирует и коммитит offsets.
  - Обновлен конфиг:
    - `backend/internal/platform/config/config.go`
    - `backend/.env.example` (`KAFKA_NOTIFICATION_GROUP`)
  - Обновлены runbook/backend docs:
    - `backend/README.md`
    - `docs/runbook.md`
  - Прогнаны проверки:
    - `go mod tidy`
    - `go test ./...` (успешно)
- Current Stage (Roadmap):
  - Stage 1 Foundation: частично закрыт.
  - Stage 2 Core Ticket Flow: базовый функционал реализован.
  - Stage 3 Async + Notifications: стартован, есть producer/consumer skeleton.
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - Worker пока только логирует события; реальная логика отправки уведомлений еще не реализована.
  - Публикация событий пока fire-and-forget (ошибки publish не эскалируются в API-ответ).
- Next:
  - Добавить реальную notification-логику (например email/webhook stub + retry policy).
  - Добавить JWT auth middleware и убрать временный header-based RBAC.

## Session 2026-02-27 #6
- Goal:
  - Довести worker до практичного baseline: notification handling + retry + DLQ.
- Done:
  - Добавлен notification service:
    - `backend/internal/notification/service.go`
    - поддерживает маппинг доменных событий в сообщения и отправку через `Sender`.
  - Добавлены тесты notification service:
    - `backend/internal/notification/service_test.go`
  - Worker обновлен:
    - retry/backoff при ошибках обработки
    - отправка в DLQ при ошибке decode или исчерпании retry
    - файлы:
      - `backend/cmd/worker/main.go`
      - `backend/cmd/worker/main_test.go`
  - Добавлены новые env-параметры worker:
    - `NOTIFICATION_RETRY_MAX`
    - `NOTIFICATION_RETRY_BACKOFF_MS`
    - `KAFKA_NOTIFICATION_DLQ_TOPIC`
    - файлы:
      - `backend/internal/platform/config/config.go`
      - `backend/.env.example`
  - Сервис тикетов публикует более полезный payload для downstream notification:
    - requester/assignee/author context
    - файл: `backend/internal/ticket/service.go`
  - Обновлены docs:
    - `backend/README.md`
    - `docs/runbook.md`
  - Прогнаны тесты:
    - `go test ./...` (успешно)
- Current Stage (Roadmap):
  - Stage 1 Foundation: частично закрыт.
  - Stage 2 Core Ticket Flow: реализован базовый объем.
  - Stage 3 Async + Notifications: существенно продвинут (producer + worker + retry/DLQ baseline).
  - Stage 4 UX + Demo Polish: не начато.
- Risks / Limits:
  - Notification sender пока log-based stub, без реальной доставки (email/telegram/webhook).
- Next:
  - Добавить реальный transport для notification (например webhook/email adapter).
  - Добавить JWT auth middleware и заменить временный header-based RBAC.
