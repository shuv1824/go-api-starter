package migration

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed schema/*.sql
var embededSchema embed.FS

func MigrateUp(db *gorm.DB, dbType string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	goose.SetBaseFS(embededSchema)

	if err := goose.SetDialect(dbType); err != nil {
		return fmt.Errorf("failed to set database dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "schema"); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
