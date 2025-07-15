package database

import (
	"fmt"
	"log/slog"

	"github.com/shuv1824/go-api-starter/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteDB struct{}

func (s *SQLiteDB) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dialector := s.GetDialector(cfg)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	if err := s.testConnection(db); err != nil {
		return nil, err
	}

	slog.Info("successfully connected to SQLite database")
	return db, nil
}

func (s *SQLiteDB) GetDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	// For SQLite, we use DbName as the file path
	// If DbName is empty or ":memory:", use in-memory database
	dbPath := cfg.DbName
	if dbPath == "" || dbPath == ":memory:" {
		dbPath = ":memory:"
	}
	return sqlite.Open(dbPath)
}

func (s *SQLiteDB) testConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
