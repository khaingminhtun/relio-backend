package scheduler

import "time"

type ScheduledJob struct {
	ID             int64 `gorm:"primaryKey;autoIncrement"`
	RelationshipID int64 `gorm:"not null;index"`
	EventID        int64 `gorm:"not null;index"`
	ReminderID     int64 `gorm:"not null;index"`

	Type   JobType   `gorm:"size:50;not null"`
	Status JobStatus `gorm:"size:30;not null;default:'pending';index"`

	ScheduledAt time.Time `gorm:"not null;index"`
	StartedAt   *time.Time
	CompletedAt *time.Time

	Attempts    int     `gorm:"not null;default:0"`
	MaxAttempts int     `gorm:"not null;default:3"`
	LastError   *string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
