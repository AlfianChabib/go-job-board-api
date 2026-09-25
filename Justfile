set dotenv-load
set dotenv-filename := ".env"
# set shell := if os() == "windows" { ["powershell.exe", "-NoProfile", "-Command"] } else { ["sh", "-c"] }

# Jalankan development mode (Docker Compose + Air hot-reload secara paralel)
[parallel]
dev: up start-air

# Build aplikasi Go menjadi binary executable di bin/app.exe
build:
  go build -o bin/app.exe cmd/api/main.go

# Menjalankan Air untuk hot-reload aplikasi Go
start-air:
  air

# Jalankan container Docker Compose di background (PostgreSQL/Redis)
up:
  docker compose up -d

# Hentikan dan hapus container Docker Compose
down:
  docker compose down

# Generate dependency injection dengan Google Wire
wire:
  wire ./cmd/api

database_user := env("DATABASE_USER")
database_password := env("DATABASE_PASSWORD")
database_db := env("DATABASE_DB")
database_host := env("DATABASE_HOST")
database_port := env("DATABASE_PORT")

DB_URL := "postgres://" + database_user + ":" + database_password + "@" + database_host + ":" + database_port + "/" + database_db + "?sslmode=disable"
MIGRATIONS_DIR := "db/migrations"

migrate-new migrate-name:
  migrate create -ext sql -dir db/migrations {{migrate-name}}

# Run pending migrations up: `just migrate-up` atau `just migrate-up 1`
migrate-up steps="":
    migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} up {{steps}}

# Rollback migrations down: `just migrate-down` (default 1) atau `just migrate-down 2`
migrate-down steps="1":
    migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} down {{steps}}

# Check current migration version
migrate-version:
    migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} version

# Force set migration version
migrate-force version:
    migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} force {{version}}