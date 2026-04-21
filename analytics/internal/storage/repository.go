package storage

import (
	"context"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Save(ctx context.Context, event domain.Event) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *GormRepository) FindAll(ctx context.Context, filter domain.Filter) ([]domain.Event, error) {
	var result []domain.Event

	q := r.db.WithContext(ctx).Model(&domain.Event{})

	if filter.EventName != "" {
		q = q.Where("name = ?", filter.EventName)
	}
	if filter.UserID != "" {
		q = q.Where("user_id = ?", filter.UserID)
	}
	if !filter.From.IsZero() {
		q = q.Where("created_at >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		q = q.Where("created_at <= ?", filter.To)
	}

	return result, q.Find(&result).Error
}

func (r *GormRepository) Aggregate(ctx context.Context, query domain.AggregateQuery) ([]domain.AggregateResult, error) {
	var results []domain.AggregateResult

	q := r.db.WithContext(ctx).Table("events").
		Select("name as event_name, count(id) as count").
		Where("name IN ?", query.EventNames).
		Group("name")

	if !query.From.IsZero() {
		q = q.Where("created_at >= ?", query.From)
	}
	if !query.To.IsZero() {
		q = q.Where("created_at <= ?", query.To)
	}

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}

	for i := range results {
		results[i].From = query.From
		results[i].To = query.To
	}

	return results, nil
}

func (r *GormRepository) EventTypeExists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.EventType{}).
		Where("name = ?", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormRepository) CreateEventType(ctx context.Context, eventType domain.EventType) error {
	return r.db.WithContext(ctx).Create(&eventType).Error
}

func (r *GormRepository) GetAllEventTypes(ctx context.Context) ([]domain.EventTypeWithCount, error) {
	var result []domain.EventTypeWithCount
	err := r.db.WithContext(ctx).
		Table("event_types").
		Select("event_types.*, count(events.id) as count").
		Joins("left join events on events.name = event_types.name").
		Group("event_types.id").
		Scan(&result).Error
	return result, err
}

