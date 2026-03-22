package main

import (
	"log"

	_ "github.com/26thavenue/backend_testproj/analytics/internal/api"
	"github.com/26thavenue/backend_testproj/analytics/internal/config"
	_ "github.com/26thavenue/backend_testproj/analytics/internal/events"
	"github.com/26thavenue/backend_testproj/analytics/internal/storage"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	db, err := storage.Open(cfg)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	if err := storage.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "LW Analytics Engine v0.1",
	})

	log.Fatal(app.Listen(cfg.AppAddr))
}
