package api

import (
	"github.com/26thavenue/backend_testproj/analytics/internal/events"
	"github.com/26thavenue/backend_testproj/analytics/internal/health"
	"github.com/26thavenue/backend_testproj/analytics/internal/swagger"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, eventsHandler *events.Handler) {
	api := app.Group("/api/v1")

	// Health
	api.Get("/health", health.Handler)

	// Events
	api.Post("/event-types", eventsHandler.CreateEventType)
	api.Post("/events", eventsHandler.Track)
	api.Get("/events", eventsHandler.List)
	api.Get("/analytics", eventsHandler.Aggregate)

	// Swagger
	api.Get("/swagger", swagger.Handler)
	api.Get("/swagger/openapi.yaml", swagger.OpenAPI)
}
