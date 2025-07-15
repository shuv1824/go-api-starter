package database

import (
	"github.com/shuv1824/go-api-starter/internal/config"
	"gorm.io/gorm"
)

// Database interface abstracts database operations
type Database interface {
	Connect(cfg *config.DatabaseConfig) (*gorm.DB, error)
	GetDialector(cfg *config.DatabaseConfig) gorm.Dialector
}
