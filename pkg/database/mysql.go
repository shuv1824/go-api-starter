package database

import (
	"fmt"
	"log/slog"

	"github.com/shuv1824/go-api-starter/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDB struct{}

func (m *MySQLDB) Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dialector := m.GetDialector(cfg)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	if err := m.testConnection(db); err != nil {
		return nil, err
	}

	slog.Info("successfully connected to MySQL database")
	return db, nil
}

func (m *MySQLDB) GetDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	return mysql.Open(dsn)
}

func (m *MySQLDB) testConnection(db *gorm.DB) error {
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
