# Production Installation (Linux)

Безопасная установка production-окружения на Linux-сервер. Все шаги с проверками — не переходи к следующему, пока не убедился в результате текущего.

Для локальной разработки на Windows см. [INSTALL-windows.md](INSTALL-windows.md).

## Требования

- Linux-сервер (Ubuntu 22.04+ / Debian 12+)
- `sudo` доступ
- Минимум 2 GB RAM, 20 GB диска
- Открытые порты: 80, 443 (или другие по договорённости)
- Домен с A-записью на сервер (опционально — можно по IP)

---

## 1. Обновить систему и поставить зависимости

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y ca-certificates curl git ufw
```

**Проверка:**
```bash
git --version
```

---

## 2. Установить Docker Engine + Compose

```bash
sudo apt install -y docker.io docker-compose-v2
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
```

**ВАЖНО:** выйди и зайди обратно по ssh, чтобы группа `docker` подхватилась.

**Проверка:**
```bash
docker --version
docker compose version
docker ps   # должно работать без sudo
```

---

## 3. Настроить firewall

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable
```

**Проверка:**
```bash
sudo ufw status verbose
```

Должны быть разрешены только SSH, 80, 443.

---

## 4. Склонировать проект

```bash
sudo mkdir -p /var/www/support-go
sudo chown $USER:$USER /var/www/support-go
git clone https://github.com/dekma-dev/support-go.git /var/www/support-go
cd /var/www/support-go
```

**Проверка:**
```bash
ls -la /var/www/support-go
```

---

## 5. Создать .env.prod с сильными секретами

```bash
cd /var/www/support-go/deploy

POSTGRES_PASSWORD=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 32)

cat > .env.prod << EOF
DOMAIN=support.example.com
TLS_EMAIL=admin@example.com

POSTGRES_DB=support_go
POSTGRES_USER=support
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

JWT_SECRET=${JWT_SECRET}

KAFKA_BROKERS=kafka:9092
KAFKA_NOTIFICATION_GROUP=support-go-notification-worker
NOTIFICATION_RETRY_MAX=3
NOTIFICATION_RETRY_BACKOFF_MS=500
KAFKA_NOTIFICATION_DLQ_TOPIC=support.notification.dlq
CORS_ORIGINS=https://support.example.com
EOF

chmod 600 .env.prod
```

**Важно:** отредактируй `DOMAIN`, `TLS_EMAIL`, `CORS_ORIGINS` под свои реальные значения перед следующим шагом.

**Проверка:**
```bash
ls -la .env.prod   # должно быть -rw------- (600)
cat .env.prod | grep -E "POSTGRES_PASSWORD|JWT_SECRET"  # убедись что секреты сгенерированы
```

---

## 6. Настроить DNS

Создай A-запись `support.example.com -> <server_ipv4>` в DNS-провайдере.

**Проверка:**
```bash
dig +short support.example.com
```

Должен вернуть IP сервера. Если нет — подожди пропагации DNS (до 10 минут) перед следующим шагом. Без рабочего DNS Caddy не сможет получить TLS-сертификат.

---

## 7. Собрать и запустить контейнеры

```bash
cd /var/www/support-go

DC="docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml"

$DC up -d --build
```

Подождать 30-60 секунд пока Postgres и Kafka стартанут.

**Проверка:**
```bash
$DC ps
```

Все сервисы должны быть `Up` или `Up (healthy)`. Если `worker` в `Restarting` — см. раздел "Частые проблемы".

---

## 8. Накатить миграции БД

```bash
cd /var/www/support-go

for m in backend/migrations/*.up.sql; do
  echo "Applying $m"
  $DC exec -T postgres sh -lc 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done
```

**Проверка:**
```bash
$DC exec postgres psql -U support -d support_go -c "\dt"
```

Должны быть таблицы: `tickets`, `ticket_comments`, `ticket_events`, `users`.

---

## 9. Верификация развёртывания

```bash
# Если настроен домен с TLS:
curl -fsS https://support.example.com/healthz
curl -fsS https://support.example.com/readyz

# Если по IP без домена:
curl -fsS http://<server_ipv4>/healthz
```

Должно вернуться `{"status":"ok",...}`.

---

## 10. Мониторинг логов

```bash
$DC logs --tail=100 api
$DC logs --tail=100 worker
$DC logs --tail=100 caddy
```

---

## Безопасность: постскриптум

После успешного деплоя проверь:

- [ ] `.env.prod` не в git (добавлен в `.gitignore`)
- [ ] Права на `.env.prod` — `600` (только владелец читает)
- [ ] `ufw status` показывает только нужные порты
- [ ] Postgres/Kafka/Redis **не** проброшены на host (см. `deploy/docker-compose.prod.yml` — у них нет секции `ports`)
- [ ] `docker ps` показывает только Caddy слушающим 80/443 наружу
- [ ] Seed-юзеры с дефолтными паролями (`000004_seed_users`) **не накатаны** в prod — либо удали этот файл перед миграцией, либо сразу смени пароли после первого логина

---

## Частые проблемы

| Симптом | Причина | Решение |
|---------|---------|---------|
| Порт 80/443 занят | На хосте работает nginx/apache | `sudo systemctl stop nginx && sudo systemctl disable nginx` либо пробросить Caddy на другой порт |
| worker в Restarting | Нет `JWT_SECRET` в окружении worker | Проверь что в `docker-compose.prod.yml` в секции `worker.environment` есть `JWT_SECRET: ${JWT_SECRET}` |
| Caddy не получает TLS | DNS ещё не распространился или порт 80 недоступен извне | `dig +short support.example.com`, убедись что 80 открыт в ufw и у облачного провайдера |
| `connection refused` на postgres | Postgres ещё стартует | Подожди 15 сек, `$DC restart api worker` |

---

## Откат и восстановление

Быстрый откат на предыдущий коммит:

```bash
cd /var/www/support-go
git fetch
git checkout <previous_commit_sha>
$DC up -d --build
```

Если миграция сломала БД — применить down-миграцию вручную или восстановить из бэкапа.

Для автоматических бэкапов БД рекомендуется:
```bash
# Cron job (ежедневно в 3:00)
0 3 * * * docker exec support_go_postgres pg_dump -U support support_go | gzip > /var/backups/support_go_$(date +\%F).sql.gz
```
