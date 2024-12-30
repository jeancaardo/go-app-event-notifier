package users

import (
	"context"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
)

type (
	Service interface {
		Get(ctx context.Context, id string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters) ([]domain.User, error)
		Store(ctx context.Context, user domain.User) (*domain.User, error)
		Update(ctx context.Context, user domain.User) (*domain.User, error)
		Delete(ctx context.Context, id string) error
	}

	Filters struct {
		Name  string
		Email string
		Phone string
		NotID string
		Page  int
		Limit int
		Sort  string
	}
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters) ([]domain.User, error) {
	return s.repo.GetAll(ctx, filters)
}

func (s *service) Store(ctx context.Context, user domain.User) (*domain.User, error) {
	count, err := s.repo.GetAll(ctx, Filters{Email: user.Email})
	if err != nil {
		return nil, err
	}
	if len(count) > 0 {
		return nil, ErrUserEmailAlreadyExists
	}
	err = s.repo.Store(ctx, &user)
	if err != nil {
		return nil, ErrOnStoreUser
	}
	return &user, nil
}

func (s *service) Update(ctx context.Context, user domain.User) (*domain.User, error) {
	count, err := s.repo.GetAll(ctx, Filters{Email: user.Email, NotID: user.ID})
	if err != nil {
		return nil, err
	}
	if len(count) > 0 {
		return nil, ErrUserEmailAlreadyExists
	}
	err = s.repo.Update(ctx, &user)
	if err != nil {
		return nil, ErrOnUpdateUser
	}
	return &user, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	user, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return ErrOnDeleteUser
	}
	return nil
}
