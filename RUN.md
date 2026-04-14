# Запуск проекта (Windows)

Повседневный запуск. Для одноразовой настройки см. [INSTALL.md](INSTALL.md).

Нужно три терминала: Docker-инфраструктура, бэкенд, фронтенд.

## 1. Docker-инфраструктура

Открой Docker Desktop (через меню Пуск), дождись пока стартанёт, потом:

```
cd C:\wamp\www\support.go\deploy
docker compose up -d
```

Проверить что всё работает:

```
docker ps
```

Должны быть запущены: `support_go_postgres`, `support_go_kafka`, `support_go_zookeeper`, `support_go_redis`.

## 2. Бэкенд API

В новом терминале cmd:

```
cd C:\wamp\www\support.go\backend
set DATABASE_URL=postgres://support:support@localhost:5432/support_go?sslmode=disable
set JWT_SECRET=change_me
set KAFKA_BROKERS=localhost:9092
go run ./cmd/api
```

Должно вывести:
```
{"level":"INFO","msg":"starting support-go api","port":"8080"}
{"level":"INFO","msg":"postgres connected"}
```

## 3. Фронтенд

В ещё одном терминале:

```
cd C:\wamp\www\support.go\frontend
npm run dev
```

## 4. Открыть

- Фронтенд: http://localhost:5173
- API health: http://localhost:8080/healthz

## Тестовые учётки

См. [CREDENTIALS.md](CREDENTIALS.md) (локальный файл, не в git).

## Остановить всё

- Ctrl+C в терминалах бэкенда и фронтенда
- `docker compose down` в `deploy/` (опционально — контейнеры могут продолжать работать в фоне)
