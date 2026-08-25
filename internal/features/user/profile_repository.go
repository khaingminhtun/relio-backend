package user

import (
	"context"
	"errors"

	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	Create(
		ctx context.Context,
		profile *UserProfile,
	) error

	GetByUserID(
		ctx context.Context,
		userID int64,
	) (*UserProfile, error,
	)

	Update(
		ctx context.Context,
		profile *UserProfile,
	) error
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(
	db *gorm.DB,
) UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (r *userProfileRepository) Create(
	ctx context.Context,
	profile *UserProfile,
) error {
	return transaction.DB(ctx, r.db).
		Create(profile).
		Error
}

func (r *userProfileRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) (*UserProfile, error) {
	var profile UserProfile

	err := transaction.DB(ctx, r.db).
		Where("user_id = ?", userID).
		First(&profile).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeUserProfileNotFound,
				"user profile not found",
				err,
			)
		}

		return nil, err
	}

	return &profile, nil
}

func (r *userProfileRepository) Update(
	ctx context.Context,
	profile *UserProfile,
) error {
	return transaction.DB(ctx, r.db).
		Save(profile).
		Error
}
