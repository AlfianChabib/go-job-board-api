set dotenv-load
set dotenv-filename := ".env"

database_user := env("DATABASE_USER")
database_password := env("DATABASE_PASSWORD")
database_db := env("DATABASE_DB")
database_host := env("DATABASE_HOST")
database_port := env("DATABASE_PORT")

DB_URL := "postgres://" + database_user + ":" + database_password + "@" + database_host + ":" + database_port + "/" + database_db + "?sslmode=disable"
MIGRATIONS_DIR := "db/migrations"

# Jalankan development mode (Docker Compose + Air hot-reload secara paralel)
# Run development mode (Docker Compose + Air hot-reload in parallel)
[parallel]
dev: up start-air

# Menjalankan Air untuk hot-reload aplikasi Go (Development mode)
# Start application with Air hot-reload (Development mode)
start-air:
  air

# Build binary aplikasi untuk production di bin/app.exe
# Build production application binary to bin/app.exe
build:
  go build -o bin/app.exe ./cmd/api

# Menjalankan binary aplikasi hasil build (Production mode)
# Run production compiled binary executable
start:
  ./bin/app.exe

# Jalankan seeder master data skills ke database
# Run master skills dataset database seeder
seed:
  go run db/seed/seeder.go

# Jalankan container Docker Compose di background (PostgreSQL & MinIO)
# Start Docker Compose containers in background (PostgreSQL & MinIO)
up:
  docker compose up -d

# Hentikan dan hapus container Docker Compose
# Stop and remove Docker Compose containers
down:
  docker compose down

# Melihat real-time logs container Docker Compose (opsional: service spesifik)
# Stream real-time Docker Compose container logs (optional: specific service)
logs service="":
  docker compose logs -f {{service}}

# Menjalankan seluruh unit test
# Run all unit tests
test:
  go test -v ./...

# Rapikan dan verifikasi dependensi Go modules
# Tidy and verify Go module dependencies
tidy:
  go mod tidy

# Generate dependency injection dengan Google Wire
# Generate dependency injection code using Google Wire
wire:
  wire ./cmd/api

# Buat file migrasi baru: `just migrate-new create_users_table`
# Create a new migration file pair: `just migrate-new create_users_table`
migrate-new migrate-name:
  migrate create -ext sql -dir db/migrations {{migrate-name}}

# Run pending migrations up: `just migrate-up` atau `just migrate-up 1`
# Run pending migrations up: `just migrate-up` or `just migrate-up 1`
migrate-up steps="":
  migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} up {{steps}}

# Rollback migrations down: `just migrate-down` (default 1) atau `just migrate-down 2`
# Rollback migrations down: `just migrate-down` (default 1) or `just migrate-down 2`
migrate-down steps="1":
  migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} down {{steps}}

# Check current migration version
# Check current database migration version
migrate-version:
  migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} version

# Force set migration version
# Force set database migration version
migrate-force version:
  migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} force {{version}}

# Reset database: drop all tables, migrate up kembali, dan jalankan seeder
# Reset database: drop all tables, re-run all migrations, and run seeder
db-fresh:
  migrate -database "{{DB_URL}}" -path {{MIGRATIONS_DIR}} drop -f
  just migrate-up
  just seed
