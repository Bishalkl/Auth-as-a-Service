package database

import (
	"fmt"
	"log"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBService interface defines the Connect method
type DBService interface {
	Connect() (*gorm.DB, error)
}

// PostgresDB implements the DBService interface
type PostgresDB struct{}

// NewDBService returns a new instance of DBService
func NewDBService() DBService {
	return &PostgresDB{}
}

// Connect opens a connection to the PostgreSQL database
func (p *PostgresDB) Connect() (*gorm.DB, error) {
	cfg := configs.Config

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	// Log safe message
	log.Println("📡 Attempting to connect to PostgreSQL...")

	// Try connecting to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("❌ Failed to connect to PostgreSQL:", err)
		return nil, err
	}

	// Optional: Ping the database
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("❌ Failed to get DB instance:", err)
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("❌ PostgreSQL ping failed:", err)
		return nil, err
	}

	log.Println("✅ Successfully connected to PostgreSQL database.")
	return db, nil
}
