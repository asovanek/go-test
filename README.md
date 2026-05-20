# Go + React Auth App

Monorepo with a modular Go (Gin + GORM + JWT) backend, React (Vite + TypeScript) frontend, PostgreSQL, and Docker Compose.

## Features

- User signup and signin with JWT authentication
- Request validation (`go-playground/validator`)
- Modular backend (`auth`, `user` modules)
- In-process event bus — `UserSignedUp` event with audit log and welcome email stub subscribers
- Swagger UI at `/swagger/index.html`
- Per-app environment files

## Quick start (Docker)

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

## Local development (without Docker)

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
