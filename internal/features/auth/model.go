package auth

import (
	"time"
)

// AuthSession has no soft delete → embed BaseModel.
type AuthSession struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	UserID int64 `gorm:"not null;index"`

	RefreshTokenHash string `gorm:"size:255;uniqueIndex;not null"`

	UserAgent string `gorm:"size:500"`
	IPAddress string `gorm:"size:45"`

	ExpiresAt time.Time `gorm:"not null;index"`
	RevokedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// OAuthAccount has no soft delete → embed BaseModel.
type OAuthAccount struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID int64 `gorm:"not null;index"`

	Provider       string `gorm:"size:50;not null;uniqueIndex:idx_oauth_provider_user"`
	ProviderUserID string `gorm:"size:255;not null;uniqueIndex:idx_oauth_provider_user"`

	Email string `gorm:"size:255"`
}
