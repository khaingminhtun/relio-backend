package note

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID     int64 `gorm:"primaryKey"`
	UserID int64 `gorm:"not null;index"`

	Title   string `gorm:"type:varchar(255);not null"`
	Content string `gorm:"type:text;not null"`

	Mood *string `gorm:"type:varchar(50)"`

	IsPinned   bool `gorm:"not null;default:false"`
	IsArchived bool `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
