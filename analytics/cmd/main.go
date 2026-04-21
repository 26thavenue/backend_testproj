package main

import (
	"log"

	"github.com/26thavenue/backend_testproj/analytics/internal/api"
	"github.com/26thavenue/backend_testproj/analytics/internal/config"
	"github.com/26thavenue/backend_testproj/analytics/internal/events"
	"github.com/26thavenue/backend_testproj/analytics/internal/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
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

	app.Use(logger.New())

	repo := storage.NewRepository(db)
	svc := events.NewService(repo)
	handler := events.NewHandler(svc)

	api.Register(app, handler)

	log.Fatal(app.Listen(cfg.AppAddr))
}
