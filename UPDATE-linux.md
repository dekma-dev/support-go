# Обновление production (Linux)

Инструкция для применения обновлений после `git pull`. Предполагается, что проект уже развёрнут по [INSTALL-linux.md](INSTALL-linux.md).

---

## Быстрая последовательность

```bash
cd /var/www/support-go

DC="docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml"

# 1. Получить обновления
git fetch origin
git status
git pull origin main

# 2. Пересобрать и перезапустить
$DC up -d --build

# 3. Накатить новые миграции (если есть)
for m in backend/migrations/*.up.sql; do
  $DC exec -T postgres sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done

# 4. Проверить
$DC ps
$DC logs --tail=50 api worker
curl -fsS http://localhost/healthz
```

---

## Подробно

### 1. Получить обновления

```bash
cd /var/www/support-go
git fetch origin
git log HEAD..origin/main --oneline   # посмотреть что нового
git pull origin main
```

**Проверка:** убедись что нет конфликтов и `.env.prod` / `Caddyfile` не затронуты (они в `.gitignore` для prod либо не менялись).

```bash
git status
```

Должно быть clean working tree.

---

### 2. Проверить что нового

Посмотреть список файлов, которые изменились:

```bash
git diff HEAD@{1} HEAD --stat
```

Обрати внимание на:
- `backend/migrations/` — новые миграции, накатить руками
- `deploy/docker-compose.prod.yml` — изменения конфига контейнеров
- `backend/.env.example` — могли появиться новые переменные окружения. Сравни со своим `.env.prod`:
  ```bash
  diff backend/.env.example deploy/.env.prod
  ```
- `deploy/Caddyfile` — если менялся, убедись что твои локальные правки не перезаписаны

---

### 3. Пересобрать и перезапустить

**Вариант А — обновить всё:**

```bash
cd /var/www/support-go
DC="docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml"
$DC up -d --build
```

Docker Compose автоматически пересоберёт образы (api, worker, frontend) и перезапустит изменённые контейнеры. Postgres/Kafka/Redis не тронутся.

**Вариант Б — обновить только конкретный сервис:**

```bash
# Только бэкенд
$DC up -d --build api worker

# Только фронтенд
$DC up -d --build frontend

# Только Caddy (если менял Caddyfile)
$DC restart caddy
```

---

### 4. Накатить новые миграции

Посмотреть какие миграции появились:

```bash
ls -la backend/migrations/*.up.sql
```

Накатить всё (идемпотентно — существующие таблицы пропускаются):

```bash
for m in backend/migrations/*.up.sql; do
  echo "Applying $m"
  $DC exec -T postgres sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done
```

Или конкретную миграцию:

```bash
$DC exec -T postgres sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < backend/migrations/000005_some_new_migration.up.sql
```

---

### 5. Верификация

```bash
# Статус контейнеров
$DC ps
```

Все должны быть `Up` или `Up (healthy)`.

```bash
# Логи — ищи ошибки
$DC logs --tail=100 api
$DC logs --tail=100 worker
```

```bash
# Health check
curl -fsS http://localhost/healthz
curl -fsS http://localhost/readyz
```

Открыть сайт в браузере и убедиться что фронтенд работает.

---

## Откат, если обновление сломало

```bash
cd /var/www/support-go
git log --oneline -10                   # найти предыдущий рабочий коммит
git checkout <previous_commit_sha>
$DC up -d --build
```

Если сломала миграция — применить `.down.sql` вручную:

```bash
$DC exec -T postgres sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' \
  < backend/migrations/000005_some_new_migration.down.sql
```

Затем откатить код и пересобрать.

---

## Очистка старых образов (опционально)

После нескольких обновлений накапливаются неиспользуемые образы. Удалить:

```bash
docker image prune -f
```

Осторожно — удалит все образы без тегов. Безопасная проверка:

```bash
docker images
```

---

## Частые проблемы

| Симптом | Решение |
|---------|---------|
| `git pull` выдаёт конфликт | Проверь локальные изменения `git status`, при необходимости `git stash` |
| Новая переменная в `.env.example`, её нет в `.env.prod` | Добавь вручную в `.env.prod` и сделай `$DC up -d` (не `--build`) |
| После `--build` контейнер не стартует | `$DC logs <service>` для диагностики, возможно нужна новая миграция |
| `permission denied` на файлы после pull | `sudo chown -R $USER:$USER /var/www/support-go` |
