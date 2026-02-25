# Support-Go Runbook (Foundation)

## 1) РџРѕРґРЅСЏС‚СЊ РёРЅС„СЂР°СЃС‚СЂСѓРєС‚СѓСЂСѓ

```powershell
docker compose -f deploy/docker-compose.yml up -d
```

## 2) Р—Р°РїСѓСЃС‚РёС‚СЊ backend

```powershell
docker exec -i support_go_postgres psql -U support -d support_go < backend/migrations/000001_create_tickets_table.up.sql

cd backend
go run ./cmd/api
```

РџСЂРѕРІРµСЂРёС‚СЊ:

- `http://localhost:8080/healthz`
- `http://localhost:8080/readyz`

## 3) Р—Р°РїСѓСЃС‚РёС‚СЊ frontend

```powershell
cd frontend
npm install
npm run dev
```

РћС‚РєСЂС‹С‚СЊ:

- `http://localhost:5173`


