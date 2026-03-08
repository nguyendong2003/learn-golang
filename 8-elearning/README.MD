# E-Learning API (Go + Gin + PostgreSQL)

Simple REST API for authentication and user management.

## Tech Stack

- Go `1.25.6` (from `app/go.mod`)
- Gin
- GORM + PostgreSQL
- Docker + Docker Compose

## Project Structure

- `app/`: Go source code
- `script/elearning.sql`: database schema init script
- `docker-compose.yml`: app + postgres services

## Prerequisites

Choose one way to run:

### Option A: Docker (recommended)

- Docker
- Docker Compose

### Option B: Run locally

- Go `>= 1.25.6`
- PostgreSQL `>= 16`

## Configuration

App loads config from `app/config.yaml`.

> Important: based on struct tags in code, JWT keys must be `access_token` and `refresh_token`.

Use this format:

```yaml
database:
  dsn: "host=127.0.0.1 user=postgres password=123456 dbname=elearning port=5432 sslmode=disable"
jwt:
  access_token: "your_access_token_secret"
  refresh_token: "your_refresh_token_secret"
  access_token_expiration: "15m"
  refresh_token_expiration: "24h"
  issuer: "elearning-api"
```

## Run with Docker

From project root:

```bash
docker compose up --build -d
```

Check containers:

```bash
docker compose ps
```

View app logs:

```bash
docker compose logs -f app
```

Stop:

```bash
docker compose down
```

Reset DB volume:

```bash
docker compose down -v
```

## Run Locally (without Docker)

### 1) Start PostgreSQL and create database

Create DB and run schema:

```bash
createdb -U postgres elearning
psql -U postgres -d elearning -f script/elearning.sql
```

### 2) Update `app/config.yaml`

Set correct PostgreSQL DSN and JWT secrets (see config sample above).

### 3) Start API

```bash
cd app
go mod tidy
go run .
```

Server starts on: `http://localhost:8080`

## API Base URL

- `http://localhost:8080/api/v1`

## Health Check

```bash
curl http://localhost:8080/api/v1/health-check
```

## Main Endpoints

### Auth

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh-token`
- `POST /api/v1/auth/forgot-password`
- `POST /api/v1/auth/change-password` (requires Bearer token)

### Users (all require Bearer token)

- `GET /api/v1/users?page=1&limit=10`
- `GET /api/v1/users/filter?page=1&limit=10&name=john&sort=username:desc,created_at:asc`
- `GET /api/v1/users/:id`
- `POST /api/v1/users`
- `PUT /api/v1/users/:id`
- `DELETE /api/v1/users/:id`

## Quick Test Flow

### 1) Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "john",
    "password": "123456"
  }'
```

### 2) Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "123456"
  }'
```

Copy `access_token` from response.

### 3) Call protected endpoint

```bash
curl http://localhost:8080/api/v1/users?page=1&limit=5 \
  -H "Authorization: Bearer <access_token>"
```

## Notes / Troubleshooting

- Docker app healthcheck currently points to `/health-check`, while API route is `/api/v1/health-check`.
  - This can make container health show `unhealthy` even when API is running.
- If token generation/validation fails, verify `app/config.yaml` uses:
  - `jwt.access_token`
  - `jwt.refresh_token`
- Default DB credentials in `docker-compose.yml`:
  - user: `postgres`
  - password: `123456`
  - db: `elearning`
