package events

import (
	"context"
	"testing"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
)

type noopRepo struct{}

func (n *noopRepo) Save(ctx context.Context, event domain.Event) error {
	return nil
}

func (n *noopRepo) FindAll(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
	return nil, nil
}

func (n *noopRepo) Aggregate(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error) {
	return nil, nil
}

func (n *noopRepo) EventTypeExists(ctx context.Context, name string) (bool, error) {
	return true, nil
}

func (n *noopRepo) CreateEventType(ctx context.Context, eventType domain.EventType) error {
	return nil
}

func (n *noopRepo) GetAllEventTypes(ctx context.Context) ([]domain.EventTypeWithCount, error) {
	return nil, nil
}

func BenchmarkServiceTrack(b *testing.B) {
	svc := NewService(&noopRepo{})
	req := domain.TrackRequest{
		EventName:  "page_view",
		UserID:     "user-1",
		Properties: map[string]any{"page": "/home"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := svc.Track(context.Background(), req); err != nil {
			b.Fatalf("track failed: %v", err)
		}
	}
}
