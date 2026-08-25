package user

import "time"

// CreateUserRequest is used to create a new user account.
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserRequest is used to update the authenticated user's account.
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=100"`
	Email    string `json:"email" binding:"omitempty,email,max=255"`
}

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID            int64      `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Role          Role       `json:"role"`
	Status        Status     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UserListResponse is used when returning multiple users.
type UserListResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int64          `json:"total"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
}

// Model -> DTO

func toUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func toUserResponseList(users []User) []UserResponse {
	result := make([]UserResponse, 0, len(users))

	for i := range users {
		result = append(result, *toUserResponse(&users[i]))
	}

	return result
}

// user profile
type UserProfileResponse struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	DisplayName   string     `json:"display_name"`
	Bio           *string    `json:"bio,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	Timezone      string     `json:"timezone"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type UpdateUserProfileRequest struct {
	DisplayName *string    `json:"display_name"`
	Bio         *string    `json:"bio"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Timezone    *string    `json:"timezone"`
}

// model to dto
func toUserProfileResponse(
	profile *UserProfile,
) *UserProfileResponse {

	return &UserProfileResponse{
		ID:            profile.ID,
		UserID:        profile.UserID,
		DisplayName:   profile.DisplayName,
		Bio:           profile.Bio,
		DateOfBirth:   profile.DateOfBirth,
		AvatarURL:     profile.AvatarURL,
		CoverImageURL: profile.CoverImageURL,
		Timezone:      profile.Timezone,
		CreatedAt:     profile.CreatedAt,
		UpdatedAt:     profile.UpdatedAt,
	}
}
