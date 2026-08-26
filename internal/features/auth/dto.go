package auth

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type RegisterResponse struct {
	RegistrationID string `json:"registration_id"`
	Message        string `json:"message"`
}

type PendingRegistration struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	OTPHash      string    `json:"otp_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type VerifyRegisterRequest struct {
	RegistrationID string `json:"registration_id" binding:"required"`
	OTP            string `json:"otp" binding:"required"`
}

type VerifyRegisterResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type LoginUserResponse struct {
	ID          int64   `json:"id"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Role        string  `json:"role"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
}

type LoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	TokenType    string            `json:"token_type"`
	ExpiresAt    time.Time         `json:"expires_at"`
	User         LoginUserResponse `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
