package storage

import (
	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20260321_create_events",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&domain.Event{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("events")
			},
		},
		{
			ID: "20260322_create_event_types",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&domain.EventType{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("event_types")
			},
		},
	})

	return m.Migrate()
}
