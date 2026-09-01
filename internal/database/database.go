package database

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// migrate create -ext sql -dir db/migrations
// migrate -database "postgres://postgres:postgres@localhost:5432/job-board-db?sslmode=disable" -path db/migrations up
// migrate -database "postgres://postgres:postgres@localhost:5432/job-board-db?sslmode=disable" -path db/migrations down
// migrate -database "postgres://postgres:postgres@localhost:5432/job-board-db?sslmode=disable" -path db/migrations version
// migrate -database "postgres://postgres:postgres@localhost:5432/job-board-db?sslmode=disable" -path db/migrations force [version]

func OpenConnection() *gorm.DB {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	dsn := env.DatabaseUrl
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}
