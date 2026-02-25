# Support-Go Runbook (Foundation)

## 1) Поднять инфраструктуру

```powershell
cd deploy
docker compose up -d
```

## 2) Запустить backend

```powershell
cd backend
go run ./cmd/api
```

Проверить:

- `http://localhost:8080/healthz`
- `http://localhost:8080/readyz`

## 3) Запустить frontend

```powershell
cd frontend
npm install
npm run dev
```

Открыть:

- `http://localhost:5173`

