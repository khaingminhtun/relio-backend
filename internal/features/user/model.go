package user

import (
	"time"

	"gorm.io/gorm"
)

// User has soft delete (DeletedAt) → embed SoftDeleteModel.
type User struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	Email        string `gorm:"size:255;uniqueIndex;not null"`
	Username     string `gorm:"size:100;not null"`
	PasswordHash string `gorm:"not null"`

	Role          Role   `gorm:"size:20;not null;default:'user'"`
	Status        Status `gorm:"size:20;not null;default:'active'"`
	EmailVerified bool   `gorm:"not null;default:false"`

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserProfile struct {
	ID     int64 `gorm:"primaryKey;autoIncrement"`
	UserID int64 `gorm:"uniqueIndex;not null"`

	DisplayName   string     `gorm:"size:150;not null"`
	Bio           *string    `gorm:"type:text"`
	AvatarURL     *string    `gorm:"size:500"`
	CoverImageURL *string    `gorm:"size:500"`
	DateOfBirth   *time.Time `gorm:"type:date"`
	Timezone      string     `gorm:"size:100;not null;default:'Asia/Yangon'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
