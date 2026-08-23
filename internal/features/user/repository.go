package user

import (
	"context"
	"errors"
	"fmt"

	postgres "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// ============================================================
// Create
// ============================================================

func (r *repository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		// Catch database-level constraint failures immediately
		if postgres.IsUniqueViolation(err) && postgres.ConstraintName(err) == "users_email_key" {
			return apperror.New(
				apperror.CodeUserAlreadyExists,
				"user with this email already exists",
				err,
			)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// ============================================================
// Get Methods (Transform RecordNotFound to Domain Errors)
// ============================================================

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeUserNotFound,
				fmt.Sprintf("user with id %d not found", id),
				err,
			)
		}
		return nil, fmt.Errorf("get user by id %d: %w", id, err)
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeUserNotFound,
				"user with this email not found",
				err,
			)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeUserNotFound,
				"user with this username not found",
				err,
			)
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &user, nil
}

// ============================================================
// List
// ============================================================

func (r *repository) List(ctx context.Context, offset, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).
		Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}

// ============================================================
// Update & Delete
// ============================================================

func (r *repository) Update(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if postgres.IsUniqueViolation(err) && postgres.ConstraintName(err) == "users_email_key" {
			return apperror.New(
				apperror.CodeUserAlreadyExists,
				"user with this email already exists",
				err,
			)
		}
		return fmt.Errorf("update user %d: %w", user.ID, err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&User{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete user %d: %w", id, res.Error)
	}

	// GORM Delete doesn't return ErrRecordNotFound if 0 rows are affected.
	// We check RowsAffected to ensure the record actually existed.
	if res.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeUserNotFound,
			fmt.Sprintf("cannot delete user %d: record not found", id),
			gorm.ErrRecordNotFound,
		)
	}

	return nil
}
