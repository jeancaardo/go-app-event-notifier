package users

import (
	"context"
	"github.com/go-kit/log"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		GetByID(ctx context.Context, id string) (*domain.User, error)
		Store(ctx context.Context, user *domain.User) error
		Update(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context, filters Filters) ([]domain.User, error)
		Delete(ctx context.Context, id string) error
	}
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

// NewRepository creates a new users repository
func NewRepository(db *gorm.DB, logger log.Logger) Repository {
	return &repo{db, logger}
}

func (r *repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) Store(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *repo) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Omit("created_at").Save(&user).Error
}

func (r *repo) GetAll(ctx context.Context, filters Filters) ([]domain.User, error) {
	var users []domain.User
	query := r.db.WithContext(ctx)
	if filters.NotID != "" {
		query = query.Where("id != ?", filters.NotID)
	}
	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("email LIKE ?", "%"+filters.Email+"%")
	}
	if filters.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+filters.Phone+"%")
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
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}
