package events

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type fakeRepo struct {
	saveFn            func(ctx context.Context, event domain.Event) error
	findAllFn         func(ctx context.Context, filter domain.Filter) ([]domain.Event, error)
	aggregateFn       func(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error)
	eventTypeExistsFn func(ctx context.Context, name string) (bool, error)
	createEventTypeFn func(ctx context.Context, eventType domain.EventType) error
	getAllEventTypesFn func(ctx context.Context) ([]domain.EventTypeWithCount, error)
}

func (f *fakeRepo) Save(ctx context.Context, event domain.Event) error {
	if f.saveFn != nil {
		return f.saveFn(ctx, event)
	}
	return nil
}

func (f *fakeRepo) FindAll(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
	if f.findAllFn != nil {
		return f.findAllFn(ctx, filter)
	}
	return nil, nil
}

func (f *fakeRepo) Aggregate(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error) {
	if f.aggregateFn != nil {
		return f.aggregateFn(ctx, query)
	}
	return nil, nil
}

func (f *fakeRepo) EventTypeExists(ctx context.Context, name string) (bool, error) {
	if f.eventTypeExistsFn != nil {
		return f.eventTypeExistsFn(ctx, name)
	}
	return true, nil
}

func (f *fakeRepo) CreateEventType(ctx context.Context, eventType domain.EventType) error {
	if f.createEventTypeFn != nil {
		return f.createEventTypeFn(ctx, eventType)
	}
	return nil
}

func (f *fakeRepo) GetAllEventTypes(ctx context.Context) ([]domain.EventTypeWithCount, error) {
	if f.getAllEventTypesFn != nil {
		return f.getAllEventTypesFn(ctx)
	}
	return nil, nil
}

func TestHandlerTrack_Valid(t *testing.T) {
	var saved domain.Event
	repo := &fakeRepo{
		saveFn: func(ctx context.Context, event domain.Event) error {
			saved = event
			return nil
		},
		eventTypeExistsFn: func(ctx context.Context, name string) (bool, error) {
			return true, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Post("/events", handler.Track)

	body, _ := json.Marshal(domain.TrackRequest{
		EventName:  "page_view",
		UserID:     "user-1",
		Properties: map[string]any{"page": "/"},
	})
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	if saved.Name != "page_view" || saved.UserID != "user-1" {
		t.Fatalf("unexpected saved event: %+v", saved)
	}
	if saved.CreatedAt.IsZero() {
		t.Fatalf("expected created_at to be set")
	}
}

func TestHandlerTrack_InvalidEvent(t *testing.T) {
	repo := &fakeRepo{
		eventTypeExistsFn: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Post("/events", handler.Track)

	body, _ := json.Marshal(domain.TrackRequest{
		EventName: "unknown",
	})
	req := httptest.NewRequest("POST", "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestHandlerAggregate(t *testing.T) {
	repo := &fakeRepo{
		aggregateFn: func(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error) {
			return []domain.AggregateResult{{
				EventName: query.EventNames[0],
				Count:     12,
				From:      query.From,
				To:        query.To,
			}}, nil
		},
		eventTypeExistsFn: func(ctx context.Context, name string) (bool, error) {
			return true, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/analytics", handler.Aggregate)

	from := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)

	req := httptest.NewRequest("GET", "/analytics?event=page_view&from="+from.Format(time.RFC3339)+"&to="+to.Format(time.RFC3339), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result []domain.AggregateResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result) == 0 || result[0].Count != 12 || result[0].EventName != "page_view" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestHandlerList(t *testing.T) {
	repo := &fakeRepo{
		findAllFn: func(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
			return []domain.Event{
				{Name: "page_view", UserID: "u1"},
				{Name: "page_view", UserID: "u2"},
			}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/events", handler.List)

	req := httptest.NewRequest("GET", "/events?event=page_view", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var result []domain.Event
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
}

func TestHandlerCreateEventType(t *testing.T) {
	var created domain.EventType
	repo := &fakeRepo{
		createEventTypeFn: func(ctx context.Context, eventType domain.EventType) error {
			created = eventType
			return nil
		},
		eventTypeExistsFn: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Post("/event-types", handler.CreateEventType)

	body, _ := json.Marshal(domain.EventType{
		Name:        "page_view",
		Description: "Page view events",
	})
	req := httptest.NewRequest("POST", "/event-types", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
	if created.Name != "page_view" {
		t.Fatalf("unexpected created event type: %+v", created)
	}
}
