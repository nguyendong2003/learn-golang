# E-Learning Backend

Backend service for an **E-Learning platform** written in **Go**, using **PostgreSQL** as the database and running with **Docker**.

---

# Tech Stack

* Go
* Docker
* Docker Compose
* GORM + PostgreSQL
* Gin
* Air (Go hot reload)

---

# Requirements

Make sure the following tools are installed:

* Docker
* Docker Compose
* Make

Verify installation:

```bash
docker -v
docker compose version
make -v
```

---

# Running the Project

## Start development environment

Run the project with hot reload:

```bash
make dev
```

After starting:

* Backend: http://localhost:8080
* PostgreSQL: localhost:5433

Whenever you modify a `.go` file, the application will **automatically rebuild and restart**.

---

## Rebuild containers

If Dockerfile or dependencies change:

```bash
make dev-build
```

---

## View logs

All logs:

```bash
make logs
```

Application logs:

```bash
make logs-app
```

Database logs:

```bash
make logs-db
```

---

## Check running containers

```bash
make ps
```

---

## Restart containers

```bash
make restart
```

---

## Stop the project

```bash
make down
```

---

# Database

The project uses **PostgreSQL**.

Connection details:

```
Host: localhost
Port: 5433
User: postgres
Password: 123456
Database: elearning
```

Database initialization script:

```
script/elearning.sql
```

---

# Reset Database

To remove all database data and recreate it:

```bash
make reset-db
```

---

# Production Build

Build the production image:

```bash
make build
```

Run production containers:

```bash
make prod
```

---

# Health Check

Health check endpoint:

```
GET /health-check
```

Example:

```bash
curl http://localhost:8080/health-check
```

---

# Development Workflow

Typical workflow:

```
make dev
↓
edit code
↓
save file
↓
application auto reload
```

No need to rebuild containers for every code change.

---

# Useful Commands

| Command        | Description                   |
| -------------- | ----------------------------- |
| make dev       | Start development environment |
| make dev-build | Rebuild dev containers        |
| make logs      | Show all logs                 |
| make logs-app  | Show application logs         |
| make logs-db   | Show database logs            |
| make ps        | Show running containers       |
| make restart   | Restart containers            |
| make down      | Stop containers               |
| make reset-db  | Reset database                |

---

# License

MIT
