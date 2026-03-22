package storage

import (
	"os"
	"path/filepath"

	"github.com/26thavenue/backend_testproj/analytics/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	if err := ensureDir(cfg.DBPath); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	if err := enableForeignKeys(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func enableForeignKeys(db *gorm.DB) error {
	return db.Exec("PRAGMA foreign_keys = ON;").Error
}
