package postgresql

import (
	"expense-tracker/model"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // global DB instance

func ConnectPostgres() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=expense_tracker host=localhost port=5432 sslmode=disable"
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Optional: check if DB is really connected
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database not reachable: %v", err)
	}

	// Auto-migrate your models (creates tables if not exist, does not drop data)
	if err := DB.AutoMigrate(&model.User{}, &model.Expense{}); err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Info("Connected to the database and migration complete.")
}
