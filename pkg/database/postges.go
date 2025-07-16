package database

import (
	"fmt"
	"log/slog"

	"github.com/shuv1824/go-api-starter/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct{}

func (p *PostgresDB) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dialector := p.GetDialector(cfg)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := p.testConnection(db); err != nil {
		return nil, err
	}

	slog.Info("successfully connected to PostgreSQL database")
	return db, nil
}

func (p *PostgresDB) GetDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Username, cfg.Password, cfg.DbName, cfg.Port, cfg.SSLMode,
	)
	return postgres.Open(dsn)
}

func (p *PostgresDB) testConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	// sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
