package database

import (
    "fmt"
    "github.com/shuv1824/go-api-starter/internal/config"
    "gorm.io/gorm"
)

type DatabaseType string

const (
    PostgreSQL DatabaseType = "postgres"
    MySQL      DatabaseType = "mysql"
    SQLite     DatabaseType = "sqlite"
)

func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    var db Database

    switch DatabaseType(cfg.Type) {
    case PostgreSQL:
        db = &PostgresDB{}
    case MySQL:
        db = &MySQLDB{}
    case SQLite:
        db = &SQLiteDB{}
    default:
        return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
    }

    return db.Connect(cfg)
}
