package auth

import (
	"context"
	"errors"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, session *AuthSession) error
	Update(ctx context.Context, session *AuthSession) error
	GetByID(ctx context.Context, id int64) (*AuthSession, error)
	GetByRefreshTokenHash(ctx context.Context, hash string) (*AuthSession, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(
	ctx context.Context,
	session *AuthSession,
) error {

	return r.db.WithContext(ctx).
		Create(session).
		Error
}

func (r *repository) Update(
	ctx context.Context,
	session *AuthSession,
) error {

	return r.db.WithContext(ctx).
		Save(session).
		Error
}

func (r *repository) GetByID(
	ctx context.Context,
	id int64,
) (*AuthSession, error) {

	var session AuthSession

	err := r.db.WithContext(ctx).
		First(&session, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeAuthSessionNotFound,
				"authentication session not found",
				err,
			)
		}

		return nil, err
	}

	return &session, nil
}

func (r *repository) GetByRefreshTokenHash(
	ctx context.Context,
	hash string,
) (*AuthSession, error) {

	var session AuthSession

	err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ?", hash).
		First(&session).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeAuthSessionNotFound,
				"authentication session not found",
				err,
			)
		}

		return nil, err
	}

	return &session, nil
}
