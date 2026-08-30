package memory

import (
	"time"

	"gorm.io/gorm"
)

type Memory struct {
	ID             int64 `gorm:"primaryKey"`
	RelationshipID int64 `gorm:"not null;index"`
	CreatedBy      int64 `gorm:"not null;index"`

	Title      string     `gorm:"type:varchar(255);not null"`
	Content    *string    `gorm:"type:text"`
	MemoryDate *time.Time `gorm:"type:timestamptz"`

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
