package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `json:"id" gorm:"type:char(36);not null;primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_email"`
	Phone     string         `json:"phone" gorm:"type:varchar(15)"`
	CreatedAt time.Time      `json:"-" gorm:"type:timestamp;column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"-" gorm:"type:timestamp;column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"type:timestamp null;column:deleted_at;index-deleted-at"`
}

func (u *User) TableName() string {
	return "users"
}

func (s *User) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
