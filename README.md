# Go + React Auth App

Monorepo with a modular Go (Gin + GORM + JWT) backend, React (Vite + TypeScript) frontend, PostgreSQL, and Docker Compose.

## Features

- User signup and signin with JWT authentication
- Request validation (`go-playground/validator`)
- Modular backend (`auth`, `user` modules)
- In-process event bus — `UserSignedUp` event with audit log and welcome email stub subscribers
- Swagger UI at `/swagger/index.html`
- Per-app environment files

## Quick start (local)

From the repo root:

```bash
npm run setup    # copy .env files, install frontend + Go deps
npm run dev      # start Postgres (Docker), backend, and frontend
```

| Service  | URL |
|----------|-----|
| Frontend | http://localhost:5173 |
| API      | http://localhost:8080 |
| Swagger  | http://localhost:8080/swagger/index.html |

### Root scripts

| Script | Description |
|--------|-------------|
| `npm run setup` | Copy `backend/.env` / `frontend/.env` if missing, install dependencies |
| `npm run dev` | Start backend (Docker) + frontend (Vite on host) |
| `npm run dev:backend` | Run Go API in Docker (starts Postgres automatically) |
| `npm run dev:frontend` | Run Vite dev server on the host |
| `npm run db:up` | Start Postgres container only |
| `npm run db:down` | Stop all dev containers |
| `npm run db:reset` | Wipe Postgres volume and restart (migrations re-run on backend start) |
| `npm run db:logs` | Tail Postgres logs |
| `npm run test` | Run backend tests in Docker |
| `npm run swagger` | Regenerate Swagger docs in Docker |
| `npm run stop` | Stop dev containers |

Requires **Docker** and **Node 20+** only — Go runs inside a container, no local Go install needed.

### Troubleshooting

**"Failed to fetch" in the browser, nothing in backend logs**

This is usually CORS or the backend not running yet:

1. **Vite moved to another port** — if you see `Port 5173 is in use, trying another one...`, open the URL Vite prints (e.g. http://localhost:5174). With `CORS_ALLOW_LOCALHOST=true` (default in dev), any localhost port is allowed.

2. **Which frontend URL are you using?**
   - `npm run dev` → http://localhost:5173 (Vite)
   - Full Docker stack → http://localhost:3000 (nginx)
   - Both are allowed when `CORS_ORIGIN` includes `http://localhost:5173,http://localhost:3000` (see `backend/.env.example`).

2. **Backend still starting** — the dev backend runs `go mod download` before the server listens (~20s on first start). Wait until `docker logs go-test-backend-1` shows `listening`.

3. **Conflicting containers** — if you previously ran `docker compose up`, stop the old nginx frontend when using `npm run dev`:
   ```bash
   docker stop go-test-frontend-1
   ```

4. **Verify the API** — `curl http://localhost:8080/healthz` should return HTTP 200.

## Quick start (full Docker stack)

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
docker compose up --build
```

| Service  | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| API      | http://localhost:8080 |
| Swagger  | http://localhost:8080/swagger/index.html |

Docker Compose overrides `DATABASE_URL` for the backend to use the `postgres` service hostname. The frontend is exposed on port **3000** in Docker (Vite dev server uses **5173** locally).

## API endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/signup` | No | Create user, return JWT |
| POST | `/api/v1/auth/signin` | No | Validate credentials, return JWT |
| GET | `/api/v1/me` | Bearer JWT | Current user profile |
| GET | `/healthz` | No | Health check |

## Local development

Prefer the root scripts above (`npm run dev`). Manual steps:

### Backend

```bash
cp backend/.env.example backend/.env
# Start PostgreSQL locally and ensure DATABASE_URL matches
cd backend
go mod tidy
go run ./cmd/server
```

Swagger docs are generated during the Docker build. For local Swagger generation:

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.3
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

### Tests

```bash
npm test
```

### Frontend

```bash
cp frontend/.env.example frontend/.env
cd frontend
npm install
npm run dev
```

## Environment

### Backend (`backend/.env`)

```
PORT=8080
DATABASE_URL=postgres://app:secret@localhost:5432/app?sslmode=disable
JWT_SECRET=change-me-in-production
JWT_EXPIRY=24h
CORS_ORIGIN=http://localhost:5173
```

### Frontend (`frontend/.env`)

```
VITE_API_URL=http://localhost:8080
```

## Project structure

```
backend/
  cmd/server/          # entrypoint
  internal/
    platform/          # config, database, logger, events, authn, middleware, validator
    modules/auth/      # signup, signin, publishes UserSignedUp
    modules/user/      # GET /me
    subscribers/       # audit log + welcome email stub on UserSignedUp
    router/
  migrations/
frontend/
  src/pages/           # SignUp, SignIn, Dashboard
  src/api/client.ts    # API + JWT storage
```
