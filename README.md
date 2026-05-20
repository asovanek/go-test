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
| `npm run dev` | Start Postgres, then backend + frontend with hot reload |
| `npm run dev:backend` | Run Go API only (Postgres must already be up) |
| `npm run dev:frontend` | Run Vite dev server only |
| `npm run db:up` | Start Postgres container |
| `npm run db:down` | Stop Postgres container |
| `npm run db:reset` | Wipe Postgres volume and restart (migrations re-run on backend start) |
| `npm run db:logs` | Tail Postgres logs |
| `npm run test` | Run backend tests (`go test ./...`) |
| `npm run test:docker` | Run backend tests in Docker (no local Go required) |
| `npm run swagger` | Regenerate Swagger docs (requires `swag` CLI) |
| `npm run stop` | Stop Postgres |

Requires **Docker** (for Postgres), **Go 1.23+**, and **Node 20+**.

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
# or without local Go:
npm run test:docker
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
