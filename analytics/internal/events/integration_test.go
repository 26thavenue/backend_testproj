package events_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
	"github.com/26thavenue/backend_testproj/analytics/internal/events"
	"github.com/26thavenue/backend_testproj/analytics/internal/storage"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEventsIntegration(t *testing.T) {
	db := openTestDB(t)
	repo := storage.NewRepository(db)
	svc := events.NewService(repo)
	handler := events.NewHandler(svc)

	app := fiber.New()
	app.Post("/event-types", handler.CreateEventType)
	app.Post("/events", handler.Track)
	app.Get("/events", handler.List)
	app.Get("/analytics", handler.Aggregate)

	typeBody, _ := json.Marshal(domain.EventType{
		Name:        "page_view",
		Description: "Page view events",
	})
	typeReq := httptest.NewRequest("POST", "/event-types", bytes.NewReader(typeBody))
	typeReq.Header.Set("Content-Type", "application/json")
	typeResp, err := app.Test(typeReq)
	if err != nil {
		t.Fatalf("event type request failed: %v", err)
	}
	if typeResp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected %d, got %d", fiber.StatusCreated, typeResp.StatusCode)
	}

	body, _ := json.Marshal(domain.TrackRequest{
		EventName:  "page_view",
		UserID:     "integration-user",
		Properties: map[string]any{"page": "/pricing"},
	})
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("track request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	aggReq := httptest.NewRequest("GET", "/analytics?event=page_view", nil)
	aggResp, err := app.Test(aggReq)
	if err != nil {
		t.Fatalf("aggregate request failed: %v", err)
	}
	if aggResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected %d, got %d", fiber.StatusOK, aggResp.StatusCode)
	}

	var agg domain.AggregateResult
	if err := json.NewDecoder(aggResp.Body).Decode(&agg); err != nil {
		t.Fatalf("decode aggregate response: %v", err)
	}
	if agg.Count != 1 {
		t.Fatalf("expected count 1, got %d", agg.Count)
	}

	listReq := httptest.NewRequest("GET", "/events?event=page_view", nil)
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if listResp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected %d, got %d", fiber.StatusOK, listResp.StatusCode)
	}

	var items []domain.Event
	if err := json.NewDecoder(listResp.Body).Decode(&items); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 event, got %d", len(items))
	}
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	return db
}
