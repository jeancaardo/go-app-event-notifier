package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID          string    `json:"id" gorm:"type:char(36);not null;primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Location    string    `json:"location" gorm:"type:varchar(255)"`
	Category    string    `json:"category" gorm:"type:varchar(255)"`
	Date        time.Time `json:"date" gorm:"type:timestamp;not null"`
	CreatedAt   time.Time `json:"-" gorm:"type:timestamp;column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"-" gorm:"type:timestamp;column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (e *Event) TableName() string {
	return "internal_events"
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return
}
