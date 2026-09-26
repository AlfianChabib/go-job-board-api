# 🚀 Go Job Board API

A robust, enterprise-grade RESTful API for a Job Board platform built with **Go (Golang)**, **GoFiber v3**, and **GORM**, designed following **Clean Architecture** principles and **Google Wire** Dependency Injection.

---

## 📋 Table of Contents

- [Tech Stack](#-tech-stack)
- [Architecture & Folder Structure](#-architecture--folder-structure)
- [Prerequisites](#-prerequisites)
- [Quick Start Guide](#-quick-start-guide)
- [Database Migrations & Seeding](#-database-migrations--seeding)
- [Project Commands (Justfile)](#-project-commands-justfile)
- [API Endpoints Specification](#-api-endpoints-specification)
- [Environment Configuration](#-environment-configuration)

---

## 🛠 Tech Stack

| Technology | Purpose |
| :--- | :--- |
| **Go 1.27+** | Core Programming Language |
| **GoFiber v3** | High-performance, Express-inspired HTTP Web Framework |
| **PostgreSQL 16** | Relational Database Management System |
| **GORM** | Object-Relational Mapping (ORM) library for Go |
| **Google Wire** | Compile-time Dependency Injection |
| **MinIO** | S3-compatible Object Storage for Resumes/CVs, Avatars, and Company Assets |
| **Golang-Migrate** | Versioned Database Schema Migrations |
| **Go Playground Validator v10** | Request Payload Validation |
| **Viper** | Environment Configuration Management |
| **JWT (golang-jwt/jwt/v5)** | Stateless Authentication (Access & Refresh Tokens) |
| **Air** | Live Hot Reloading during Development |
| **Just** | Command Runner / Task Automation |
| **Docker & Docker Compose** | Containerized Local Infrastructure |

---

## 🏛 Architecture & Folder Structure

This project follows **Clean Architecture** patterns to separate business logic, persistence, and delivery layers:

```text
go-job-board-api/
├── cmd/
│   └── api/
│       ├── main.go               # Application entry point
│       ├── wire.go               # Dependency injection bindings (wireinject)
│       └── wire_gen.go           # Generated wire providers
├── internal/
│   ├── config/                   # Viper environment loader
│   ├── controller/               # Fiber HTTP Handlers
│   ├── database/                 # PostgreSQL & MinIO connection initializers
│   ├── helper/                   # Generic response formatters
│   ├── middleware/               # JWT authentication & role-based authorization
│   ├── model/
│   │   ├── domain/               # Database entities / GORM models
│   │   └── web/                  # Request & Response DTOs
│   ├── repository/               # Data access layer (PostgreSQL queries)
│   ├── router/                   # Fiber route definitions
│   └── service/                  # Core business logic layer
├── db/
│   ├── migrations/               # Up and Down SQL migration files
│   └── seed/                     # Database seeders and dataset CSVs
│       ├── curency-codes.csv     # ISO-4217 Currency dataset
│       ├── seeder.go             # Skills seeder script
│       └── skills-dataset.csv    # 4,300+ skills dataset
├── pkg/
│   ├── errs/                     # Custom AppError & domain error codes
│   ├── utils/                    # Password hashing (bcrypt) & JWT manager
│   └── validator/                # Validation engine
├── embed.go                      # Go embed bridge for static seed assets
├── docker-compose.yml            # Docker services (PostgreSQL & MinIO)
├── Justfile                      # Task automation commands
└── .env.example                  # Environment configuration template
```

---

## ⚙ Prerequisites

Before running the application, ensure the following tools are installed:

- **Go** (v1.27 or higher)
- **Docker & Docker Compose**
- **Just** (recommended command runner): `cargo install just` or `winget install Casey.Just` / `brew install just`
- **golang-migrate CLI**: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- **Google Wire**: `go install github.com/google/wire/cmd/wire@latest`
- **Air** (optional, for hot reload): `go install github.com/air-verse/air@latest`

---

## 🚀 Quick Start Guide

### 1. Clone & Setup Environment

```bash
git clone https://github.com/AlfianChabib/go-job-board-api.git
cd go-job-board-api

# Copy environment template
cp .env.example .env
```

### 2. Start Infrastructure Containers

Start PostgreSQL and MinIO services in the background:

```bash
just up
# or: docker compose up -d
```

### 3. Run Database Migrations

Apply all pending SQL migrations to PostgreSQL:

```bash
just migrate-up
```

### 4. Seed Master Data

Populate the database with over 4,300+ tech skills:

```bash
just seed
# or: go run db/seed/seeder.go
```

### 5. Run Server (Development & Production)

**Development (Hot-Reload):**
```bash
# Run with Air (live hot reload)
just start-air

# Or run both Docker containers & Air in parallel:
just dev
```

**Production (Compiled Binary):**
```bash
# 1. Build binary executable to bin/app.exe
just build

# 2. Start the built production binary
just start
```

The API will be available at: **`http://localhost:8000`**

---

## 🗄 Database Migrations & Seeding

Migrations are managed with `golang-migrate` using raw SQL:

```bash
# Check current migration version
just migrate-version

# Create a new migration file pair (up & down)
just migrate-new create_table_name

# Apply all migrations up
just migrate-up

# Rollback 1 migration step
just migrate-down 1

# Rollback 2 migration steps
just migrate-down 2
```

---

## ⚡ Project Commands (Justfile)

| Command | Description |
| :--- | :--- |
| `just dev` | Runs Docker Compose and Air hot-reload in parallel |
| `just start-air` | Starts the app with Air hot-reload (Development) |
| `just build` | Compiles application into production binary at `bin/app.exe` |
| `just start` | Runs the compiled production binary (`./bin/app.exe`) |
| `just up` | Starts Docker containers (`job-board-postgres`, `minio`) |
| `just down` | Stops and removes Docker containers |
| `just logs` | Streams real-time logs from Docker containers (`just logs postgres`) |
| `just test` | Runs all Go unit tests (`go test -v ./...`) |
| `just tidy` | Cleans up and verifies Go module dependencies (`go mod tidy`) |
| `just seed` | Populates PostgreSQL with master skills dataset (`go run db/seed/seeder.go`) |
| `just wire` | Re-generates Google Wire dependency injection (`cmd/api/wire_gen.go`) |
| `just migrate-new <name>` | Creates a new SQL migration file pair (up & down) |
| `just migrate-up` | Runs all pending database migrations |
| `just migrate-down` | Rolls back the last migration step (default 1) |
| `just migrate-version` | Displays current active database migration version |
| `just db-fresh` | Drops all tables, re-runs all migrations, and re-seeds master data |

---

## 📡 API Endpoints Specification

Base URL: `/api`

### 1. Public & Master Data

| Method | Endpoint | Description | Query Parameters |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/skills` | Get list of skills | `?search=go&limit=10` |
| `GET` | `/api/currencies` | Get ISO currency codes (In-Memory embedded cache) | `?search=idr&limit=10` |
| `GET` | `/api/companies` | Browse registered companies | `?search=tech&industry=IT&page=1&limit=10` |
| `GET` | `/api/companies/:id` | Get company detail by UUID | - |

### 2. Authentication (`/api/auth`)

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Register new user (`CANDIDATE` or `RECRUITER`) | No |
| `POST` | `/api/auth/login` | Login with email & password | No |
| `POST` | `/api/auth/refresh` | Issue new access token using refresh token | No |
| `POST` | `/api/auth/logout` | Revoke session & refresh token | Yes |

### 3. Candidate Profile (`/api/candidate`)

> Requires Header: `Authorization: Bearer <access_token>` with role `CANDIDATE`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/candidate/profile` | Get candidate profile, skills, and work experiences |
| `PUT` | `/api/candidate/profile` | Update candidate headline and phone number |
| `PUT` | `/api/candidate/avatar` | Upload/Replace profile photo (Multipart Form to MinIO) |
| `DELETE` | `/api/candidate/avatar` | Remove profile photo |
| `PUT` | `/api/candidate/skills` | Attach skills to candidate profile |
| `GET` | `/api/candidate/experiences` | Get list of candidate work experiences |
| `POST` | `/api/candidate/experiences` | Add new work experience |
| `PUT` | `/api/candidate/experiences/:experienceId` | Update work experience detail |
| `DELETE` | `/api/candidate/experiences/:experienceId` | Delete work experience entry |

### 4. Recruiter Management (`/api/recruiter`)

> Requires Header: `Authorization: Bearer <access_token>` with role `RECRUITER`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/recruiter/company` | Create company profile |
| `GET` | `/api/recruiter/company` | Get managed company profile |
| `PUT` | `/api/recruiter/company` | Update company information |
| `PUT` | `/api/recruiter/company/logo` | Upload company logo (MinIO) |
| `PUT` | `/api/recruiter/company/banner` | Upload company cover banner (MinIO) |

---

## 🔒 Environment Configuration

Key configuration parameters defined in `.env`:

```ini
APP_ENV=development
PORT=:8000

# PostgreSQL
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_DB=job-board-db
DATABASE_HOST=localhost
DATABASE_PORT=5432

# MinIO Object Storage
MINIO_ENDPOINT=localhost:9000
MINIO_PUBLIC_URL=http://localhost:9000
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_AVATAR_BUCKET=avatar
MINIO_CV_BUCKET=cv-bucket
MINIO_USE_SSL=false

# JWT Secret Keys & Expiry
ACCESS_SECRET_KEY=your_access_secret_key
REFRESH_SECRET_KEY=your_refresh_secret_key
ACCESS_DURATION=900     # 15 minutes
REFRESH_DURATION=604800 # 7 days
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
