package user

import (
	"context"
	"fmt"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

type UserProfileService interface {
	Get(
		ctx context.Context,
		userID int64,
	) (*UserProfileResponse, error)

	Update(
		ctx context.Context,
		userID int64,
		req UpdateUserProfileRequest,
	) (*UserProfileResponse, error)
}

type userProfileService struct {
	profileRepo UserProfileRepository
}

func NewUserProfileService(
	profileRepo UserProfileRepository,
) UserProfileService {
	return &userProfileService{
		profileRepo: profileRepo,
	}
}

// ============================================================
// Get Profile
// ============================================================

func (s *userProfileService) Get(
	ctx context.Context,
	userID int64,
) (*UserProfileResponse, error) {

	profile, err := s.profileRepo.GetByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeUserProfileNotFound,
			"user profile not found",
			err,
		)
	}

	return toUserProfileResponse(profile), nil
}

// ============================================================
// Update Profile
// ============================================================

func (s *userProfileService) Update(
	ctx context.Context,
	userID int64,
	req UpdateUserProfileRequest,
) (*UserProfileResponse, error) {

	profile, err := s.profileRepo.GetByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return nil, apperror.New(
			apperror.CodeUserProfileNotFound,
			"user profile not found",
			err,
		)
	}

	// ------------------------------------------------------------
	// Update only fields supplied by the client
	// ------------------------------------------------------------

	if req.DisplayName != nil {
		profile.DisplayName = *req.DisplayName
	}

	if req.Bio != nil {
		profile.Bio = req.Bio
	}

	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}

	if req.Timezone != nil {
		profile.Timezone = *req.Timezone
	}

	// ------------------------------------------------------------
	// Save
	// ------------------------------------------------------------

	if err := s.profileRepo.Update(
		ctx,
		profile,
	); err != nil {
		return nil, fmt.Errorf("faile to update user profile, %w", err)
	}

	return toUserProfileResponse(profile), nil
}

// ============================================================
// Response Mapper
// ============================================================
