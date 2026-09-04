# run dev
[parallel]
dev: up start-air

build:
  go build -o bin/app.exe cmd/api/main.go

start-air:
  air

up:
  docker compose up -d

down:
  docker compose down