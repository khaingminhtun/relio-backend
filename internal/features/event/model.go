package event

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID             int64      `gorm:"primaryKey;autoIncrement"`
	RelationshipID int64      `gorm:"not null;index"`
	Title          string     `gorm:"size:255;not null"`
	Type           EventType  `gorm:"size:50;not null"`
	Description    *string    `gorm:"type:text"`
	EventDate      time.Time  `gorm:"type:date;not null"`
	EventTime      *time.Time `gorm:"type:time"`
	AllDay         bool       `gorm:"not null;default:false"`
	CreatedBy      int64      `gorm:"not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
