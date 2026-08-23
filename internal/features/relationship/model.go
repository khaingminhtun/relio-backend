package relationship

import (
	"time"

	"gorm.io/gorm"
)

type Relationship struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	Name        string           `gorm:"size:150;not null"`
	Type        RelationshipType `gorm:"size:50;not null"`
	CustomType  *string          `gorm:"size:100"`
	Description *string
	StartDate   *time.Time
	Timezone    string `gorm:"size:100;not null"`
	CreatedBy   int64  `gorm:"not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type RelationshipMember struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	RelationshipID int64                    `gorm:"not null;index"`
	UserID         int64                    `gorm:"not null;index"`
	Role           RelationshipMemberRole   `gorm:"size:50;not null"`
	Status         RelationshipMemberStatus `gorm:"size:50;not null"`
	JoinedAt       *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Invitation struct {
	ID             int64            `gorm:"primaryKey;autoIncrement"`
	RelationshipID int64            `gorm:"not null;index"`
	InvitedBy      int64            `gorm:"not null;index"`
	Email          string           `gorm:"size:255;not null;index"`
	TokenHash      string           `gorm:"size:255;not null;uniqueIndex"`
	Status         InvitationStatus `gorm:"size:50;not null"`
	ExpiresAt      time.Time        `gorm:"not null"`
	AcceptedAt     *time.Time

	CreatedAt time.Time
}
