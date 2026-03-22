package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type Repository interface {
	Save(ctx context.Context, event domain.Event) error
	FindAll(ctx context.Context, filter domain.Filter) ([]domain.Event, error)
	Aggregate(ctx context.Context, query domain.AggregateQuery) (domain.AggregateResult, error)
	EventTypeExists(ctx context.Context, name string) (bool, error)
	CreateEventType(ctx context.Context, eventType domain.EventType) error
}

func (s *Service) Track(ctx context.Context, req domain.TrackRequest) error {
	eventName := strings.TrimSpace(req.EventName)
	if eventName == "" {
		return fmt.Errorf("event_name is required")
	}
	ok, err := s.repo.EventTypeExists(ctx, eventName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unknown event type: %s", eventName)
	}

	event := domain.Event{
		Name:       eventName,
		UserID:     req.UserID,
		Properties: req.Properties,
		CreatedAt:  time.Now().UTC(),
	}

	return s.repo.Save(ctx, event)
}

func (s *Service) GetAnalytics(ctx context.Context, query domain.AggregateQuery) (domain.AggregateResult, error) {
	eventName := strings.TrimSpace(query.EventName)
	if eventName == "" {
		return domain.AggregateResult{}, fmt.Errorf("event is required")
	}
	ok, err := s.repo.EventTypeExists(ctx, eventName)
	if err != nil {
		return domain.AggregateResult{}, err
	}
	if !ok {
		return domain.AggregateResult{}, fmt.Errorf("unknown event type: %s", eventName)
	}

	query.EventName = eventName
	return s.repo.Aggregate(ctx, query)
}

func (s *Service) GetRawEvents(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *Service) CreateEventType(ctx context.Context, eventType domain.EventType) error {
	name := strings.TrimSpace(eventType.Name)
	if name == "" {
		return fmt.Errorf("event_name is required")
	}

	exists, err := s.repo.EventTypeExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("event type already exists: %s", name)
	}

	eventType.Name = name
	return s.repo.CreateEventType(ctx, eventType)
}
