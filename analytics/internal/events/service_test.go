package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
)

type stubRepo struct {
	saveCalls       int
	findAllCalls    int
	aggregateCalls  int
	existsCalls     int
	createCalls     int
	saveErr         error
	findAllResult   []domain.Event
	aggregateResult []domain.AggregateResult
	existsResult    bool
	existsErr       error
	createErr       error
	eventTypes      []domain.EventTypeWithCount
}

func (s *stubRepo) Save(ctx context.Context, event domain.Event) error {
	s.saveCalls++
	return s.saveErr
}

func (s *stubRepo) FindAll(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
	s.findAllCalls++
	return s.findAllResult, nil
}

func (s *stubRepo) Aggregate(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error) {
	s.aggregateCalls++
	return s.aggregateResult, nil
}

func (s *stubRepo) EventTypeExists(ctx context.Context, name string) (bool, error) {
	s.existsCalls++
	return s.existsResult, s.existsErr
}

func (s *stubRepo) CreateEventType(ctx context.Context, eventType domain.EventType) error {
	s.createCalls++
	return s.createErr
}

func (s *stubRepo) GetAllEventTypes(ctx context.Context) ([]domain.EventTypeWithCount, error) {
	s.findAllCalls++
	return s.eventTypes, nil
}

func TestServiceTrack_ValidEvent(t *testing.T) {
	repo := &stubRepo{existsResult: true}
	svc := NewService(repo)

	err := svc.Track(context.Background(), domain.TrackRequest{
		EventName:  "page_view",
		UserID:     "u1",
		Properties: map[string]any{"page": "/"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.saveCalls != 1 {
		t.Fatalf("expected Save to be called once, got %d", repo.saveCalls)
	}
}

func TestServiceTrack_InvalidEvent(t *testing.T) {
	repo := &stubRepo{existsResult: false}
	svc := NewService(repo)

	err := svc.Track(context.Background(), domain.TrackRequest{
		EventName: "invalid",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.saveCalls != 0 {
		t.Fatalf("did not expect Save to be called")
	}
}

func TestServiceTrack_RepoError(t *testing.T) {
	repo := &stubRepo{saveErr: errors.New("db error"), existsResult: true}
	svc := NewService(repo)

	err := svc.Track(context.Background(), domain.TrackRequest{
		EventName: "page_view",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestServiceGetAnalytics_Valid(t *testing.T) {
	expected := []domain.AggregateResult{{
		EventName: "page_view",
		Count:     5,
		From:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
	}}
	repo := &stubRepo{aggregateResult: expected, existsResult: true}
	svc := NewService(repo)

	result, err := svc.GetAnalytics(context.Background(), domain.AggregateQuery{
		EventNames: []string{"page_view"},
		From:       expected[0].From,
		To:         expected[0].To,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 || result[0].Count != expected[0].Count || result[0].EventName != expected[0].EventName {
		t.Fatalf("unexpected result: %+v", result)
	}
	if repo.aggregateCalls != 1 {
		t.Fatalf("expected Aggregate to be called once, got %d", repo.aggregateCalls)
	}
}

func TestServiceGetAnalytics_InvalidEvent(t *testing.T) {
	repo := &stubRepo{existsResult: false}
	svc := NewService(repo)

	_, err := svc.GetAnalytics(context.Background(), domain.AggregateQuery{
		EventNames: []string{"unknown"},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.aggregateCalls != 0 {
		t.Fatalf("did not expect Aggregate to be called")
	}
}

func TestServiceGetRawEvents(t *testing.T) {
	repo := &stubRepo{
		findAllResult: []domain.Event{{Name: "page_view"}, {Name: "user_login"}},
	}
	svc := NewService(repo)

	result, err := svc.GetRawEvents(context.Background(), domain.Filter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	if repo.findAllCalls != 1 {
		t.Fatalf("expected FindAll to be called once, got %d", repo.findAllCalls)
	}
}

func TestServiceCreateEventType(t *testing.T) {
	repo := &stubRepo{existsResult: false}
	svc := NewService(repo)

	err := svc.CreateEventType(context.Background(), domain.EventType{
		Name:        "page_view",
		Description: "Page view events",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected CreateEventType to be called once, got %d", repo.createCalls)
	}
}

func TestServiceCreateEventType_Duplicate(t *testing.T) {
	repo := &stubRepo{existsResult: true}
	svc := NewService(repo)

	err := svc.CreateEventType(context.Background(), domain.EventType{
		Name: "page_view",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.createCalls != 0 {
		t.Fatalf("did not expect CreateEventType to be called")
	}
}
