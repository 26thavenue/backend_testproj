package domain

import "time"

type Event struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"index;not null;size:120" json:"event_name"`
	UserID     string         `gorm:"index;size:120" json:"user_id"`
	Properties map[string]any `gorm:"serializer:json" json:"properties"`
	CreatedAt  time.Time      `json:"timestamp"`
}

type EventType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:120" json:"event_name"`
	Description string    `gorm:"size:255" json:"event_description"`
	CreatedAt   time.Time `json:"created_at"`
}

type TrackRequest struct {
	EventName  string         `json:"event_name"`
	UserID     string         `json:"user_id"`
	Properties map[string]any `json:"properties"`
}

type Filter struct {
	EventName string
	UserID    string
	From      time.Time
	To        time.Time
}

type AggregateQuery struct {
	EventNames []string
	From       time.Time
	To         time.Time
}

type AggregateResult struct {
	EventName string    `json:"event"`
	Count     int64     `json:"count"`
	From      time.Time `json:"from,omitempty"`
	To        time.Time `json:"to,omitempty"`
}

type EventTypeWithCount struct {
	EventType
	Count int64 `json:"count"`
}

type ListEventTypesResponse struct {
	EventTypes []EventTypeWithCount `json:"event_types"`
}

type CreateEventTypeRequest struct {
	Name        string `json:"event_name"`
	Description string `json:"event_description"`
}