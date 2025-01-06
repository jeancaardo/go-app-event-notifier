package internalevents

import (
	"context"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
	"time"
)

type (
	Service interface {
		Get(ctx context.Context, id string) (*domain.Event, error)
		GetAll(ctx context.Context, filters Filters) ([]domain.Event, error)
		Store(ctx context.Context, event domain.Event) (*domain.Event, error)
		Update(ctx context.Context, event domain.Event) (*domain.Event, error)
		Delete(ctx context.Context, id string) error
	}

	Filters struct {
		Name     string
		Category string
		Location string
		DateFrom time.Time
		DateTo   time.Time
		NotID    string
		Page     int
		Limit    int
		Sort     string
	}
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Get(ctx context.Context, id string) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEventNotFound
	}
	return event, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters) ([]domain.Event, error) {
	return s.repo.GetAll(ctx, filters)
}

func (s *service) Store(ctx context.Context, event domain.Event) (*domain.Event, error) {
	count, err := s.repo.GetAll(ctx, Filters{Name: event.Name})
	if err != nil {
		return nil, err
	}
	if len(count) > 0 {
		return nil, ErrEventNameAlreadyExists
	}
	err = s.repo.Store(ctx, &event)
	if err != nil {
		return nil, ErrOnStoreEvent
	}
	return &event, nil
}

func (s *service) Update(ctx context.Context, event domain.Event) (*domain.Event, error) {
	count, err := s.repo.GetAll(ctx, Filters{Name: event.Name, NotID: event.ID})
	if err != nil {
		return nil, err
	}
	if len(count) > 0 {
		return nil, ErrEventNameAlreadyExists
	}
	err = s.repo.Update(ctx, &event)
	if err != nil {
		return nil, ErrOnUpdateEvent
	}
	return &event, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	event, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotFound
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return ErrOnDeleteEvent
	}
	return nil
}
