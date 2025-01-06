package internalevents

import (
	"context"
	"github.com/go-kit/log"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		GetByID(ctx context.Context, id string) (*domain.Event, error)
		Store(ctx context.Context, event *domain.Event) error
		Update(ctx context.Context, event *domain.Event) error
		GetAll(ctx context.Context, filters Filters) ([]domain.Event, error)
		Delete(ctx context.Context, id string) error
	}
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

// NewRepository creates a new internal_events repository
func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return &repo{db, logger}
}

func (r *repo) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repo) Store(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *repo) Update(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Omit("created_at").Save(&event).Error
}

func (r *repo) GetAll(ctx context.Context, filters Filters) ([]domain.Event, error) {
	var events []domain.Event
	query := r.db.WithContext(ctx)
	if filters.NotID != "" {
		query = query.Where("id != ?", filters.NotID)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Location != "" {
		query = query.Where("location = ?", filters.Location)
	}
	if !filters.DateFrom.IsZero() {
		query = query.Where("date >= ?", filters.DateFrom)
	}
	if !filters.DateTo.IsZero() {
		query = query.Where("date <= ?", filters.DateTo)
	}
	if filters.Sort != "" {
		query = query.Order(filters.Sort)
	}
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		query = query.Offset((filters.Page - 1) * filters.Limit)
	}
	err := query.Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Event{}).Error
}
