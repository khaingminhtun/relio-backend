package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/khaingminhtun/relio-backend/internal/shared/security"
)

type Service interface {
	CreateUser(
		ctx context.Context,
		req CreateUserRequest,
	) (*UserResponse, error)

	GetUser(
		ctx context.Context,
		id int64,
	) (*UserResponse, error)

	ListUsers(
		ctx context.Context,
		offset, limit int,
	) (*UserListResponse, error)

	UpdateUser(
		ctx context.Context,
		id int64,
		req UpdateUserRequest,
	) (*UserResponse, error)

	DeleteUser(
		ctx context.Context,
		id int64,
	) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// ============================================================
// Create
// ============================================================

func (s *service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (*UserResponse, error) {

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &User{
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		Username:      strings.TrimSpace(req.Username),
		PasswordHash:  passwordHash,
		Role:          RoleUser,
		Status:        StatusActive,
		EmailVerified: false,
	}

	// The repository checks unique index keys and flags errors
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// ============================================================
// Get
// ============================================================

func (s *service) GetUser(
	ctx context.Context,
	id int64,
) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err // Bubbles up configured apperror.CodeUserNotFound
	}

	return toUserResponse(user), nil
}

// ============================================================
// List
// ============================================================

func (s *service) ListUsers(
	ctx context.Context,
	offset, limit int,
) (*UserListResponse, error) {

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	users, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return &UserListResponse{
		Users:  toUserResponseList(users),
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

// ============================================================
// Update
// ============================================================

func (s *service) UpdateUser(
	ctx context.Context,
	id int64,
	req UpdateUserRequest,
) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err // Returns NotFound instantly if record is missing
	}

	if req.Email != "" {
		user.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}

	if req.Username != "" {
		user.Username = strings.TrimSpace(req.Username)
	}

	// Handles collisions via domain wrapper validations
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// ============================================================
// Delete
// ============================================================

func (s *service) DeleteUser(
	ctx context.Context,
	id int64,
) error {

	// Confirms target existence first
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
