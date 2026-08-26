package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/khaingminhtun/relio-backend/internal/features/user"
	redisinfra "github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"
)

type Service interface {
	Register(
		ctx context.Context,
		req RegisterRequest,
	) (*RegisterResponse, error)

	VerifyRegister(
		ctx context.Context,
		req VerifyRegisterRequest,
	) (*VerifyRegisterResponse, error)

	Authenticate(
		ctx context.Context,
		req LoginRequest,
		userAgent string,
		ipAddress string) (*LoginResponse, error)
	Refresh(
		ctx context.Context,
		req RefreshRequest,
	) (*RefreshResponse, error)
}

type service struct {
	userRepo    user.Repository
	profileRepo user.UserProfileRepository
	authRepo    Repository
	redisStore  redisinfra.RedisStore
	emailQueue  redisinfra.EmailQueue
	jwtManager  *security.JWTManager
	txManger    transaction.Manager
}

func NewService(
	userRepo user.Repository,
	profileRepo user.UserProfileRepository,
	authRepo Repository,
	redisStore redisinfra.RedisStore,
	emailQueue redisinfra.EmailQueue,
	jwtManager *security.JWTManager,
	txManager transaction.Manager,
) Service {
	return &service{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		authRepo:    authRepo,
		redisStore:  redisStore,
		emailQueue:  emailQueue,
		jwtManager:  jwtManager,
		txManger:    txManager,
	}
}

func (s *service) Register(
	ctx context.Context,
	req RegisterRequest,
) (*RegisterResponse, error) {

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// ============================================================
	// Check existing email
	// ============================================================

	_, err := s.userRepo.GetByEmail(ctx, email)

	switch {
	case err == nil:
		return nil, apperror.New(
			apperror.CodeUserAlreadyExists,
			"user with this email already exists",
			nil,
		)

	case apperror.Is(err, apperror.CodeUserNotFound):
		// Expected.
		// Email does not exist, so registration can continue.

	default:
		return nil, err
	}

	// ============================================================
	// Check existing username
	// ============================================================

	_, err = s.userRepo.GetByUsername(ctx, username)

	switch {
	case err == nil:
		return nil, apperror.New(
			apperror.CodeUserAlreadyExists,
			"user with this username already exists",
			nil,
		)

	case apperror.Is(err, apperror.CodeUserNotFound):
		// Expected.
		// Username does not exist, so registration can continue.

	default:
		return nil, err
	}

	// ============================================================
	// Hash password
	// ============================================================

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// ============================================================
	// Generate OTP
	// ============================================================

	otp, err := security.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	otpHash := security.HashOTP(otp)

	// ============================================================
	// Build pending registration
	// ============================================================

	pending := PendingRegistration{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		OTPHash:      otpHash,
	}

	data, err := json.Marshal(pending)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Save pending registration in Redis
	// ============================================================
	registrationID := uuid.NewString()

	key := "auth:register:" + registrationID

	const registrationTTL = 10 * time.Minute

	if err := s.redisStore.Set(
		ctx,
		key,
		string(data),
		registrationTTL,
	); err != nil {
		return nil, fmt.Errorf(
			"save pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Create email job
	// ============================================================

	job := redisinfra.EmailJob{
		ID:        uuid.NewString(),
		To:        email,
		Subject:   "Verify your email",
		Template:  "otp_verification",
		CreatedAt: time.Now(),

		Data: map[string]any{
			"OTP":       otp,
			"ExpiresIn": "10 minutes",
		},
	}

	// ============================================================
	// Publish email job
	// ============================================================

	if err := s.emailQueue.Publish(ctx, job); err != nil {

		// Remove pending registration because
		// the email could not be queued.
		_ = s.redisStore.Delete(ctx, key)

		return nil, fmt.Errorf(
			"queue verification email: %w",
			err,
		)
	}

	// ============================================================
	// Response
	// ============================================================

	return &RegisterResponse{
		RegistrationID: registrationID,
		Message:        "Verification code sent to your email",
	}, nil
}

func (s *service) VerifyRegister(
	ctx context.Context,
	req VerifyRegisterRequest,
) (*VerifyRegisterResponse, error) {

	registrationID := strings.TrimSpace(req.RegistrationID)

	key := "auth:register:" + registrationID

	data, err := s.redisStore.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf(
			"get pending registration: %w",
			err,
		)
	}

	if data == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRequest,
			"registration expired or not found",
			nil,
		)
	}

	// ============================================================
	// Decode pending registration
	// ============================================================

	var pending PendingRegistration

	if err := json.Unmarshal(
		[]byte(data),
		&pending,
	); err != nil {
		return nil, fmt.Errorf(
			"unmarshal pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Verify OTP
	// ============================================================

	if !security.VerifyOTP(
		req.OTP,
		pending.OTPHash,
	) {
		return nil, apperror.New(
			apperror.CodeInvalidVerifyCode,
			"invalid verification code",
			nil,
		)
	}

	// ============================================================
	// Create User and profile in one tranaction
	// ============================================================

	err = s.txManger.WithinTransaction(
		ctx,
		func(txCtx context.Context) error {

			newUser := &user.User{
				Username:      pending.Username,
				Email:         pending.Email,
				PasswordHash:  pending.PasswordHash,
				Role:          user.RoleUser,
				Status:        user.StatusActive,
				EmailVerified: true,
			}

			// --------------------------------------------------------
			// Create User
			// --------------------------------------------------------

			if err := s.userRepo.Create(
				txCtx,
				newUser,
			); err != nil {
				return fmt.Errorf(
					"create user: %w",
					err,
				)
			}

			// --------------------------------------------------------
			// Create User Profile
			// --------------------------------------------------------

			profile := &user.UserProfile{
				UserID:      newUser.ID,
				DisplayName: newUser.Username,
				Timezone:    "Asia/Yangon",
			}

			if err := s.profileRepo.Create(
				txCtx,
				profile,
			); err != nil {
				return fmt.Errorf(
					"create user profile: %w",
					err,
				)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	// ============================================================
	// Delete pending registration
	// ============================================================

	if err := s.redisStore.Delete(
		ctx,
		key,
	); err != nil {

		// User was already created successfully.
		// Don't report registration failure because Redis cleanup
		// failed after the database transaction succeeded.

		return &VerifyRegisterResponse{
			Message: "registration successful",
		}, nil
	}

	// ============================================================
	// Success
	// ============================================================

	return &VerifyRegisterResponse{
		Message: "registration successful",
	}, nil

}

func (s *service) Authenticate(
	ctx context.Context,
	req LoginRequest,
	userAgent string,
	ipAddress string,
) (*LoginResponse, error) {

	// ============================================================
	// Normalize email
	// ============================================================

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// ============================================================
	// Get user
	// ============================================================

	currentUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {

		if apperror.Is(err, apperror.CodeUserNotFound) {
			return nil, apperror.New(
				apperror.CodeInvalidCredentials,
				"invalid email or password",
				nil,
			)
		}

		return nil, err
	}

	// ============================================================
	// Verify password
	// ============================================================

	if err := security.ComparePassword(
		req.Password,
		currentUser.PasswordHash,
	); err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidCredentials,
			"invalid email or password",
			nil,
		)
	}

	// ============================================================
	// Check account status
	// ============================================================

	if currentUser.Status != user.StatusActive {
		return nil, apperror.New(
			apperror.CodeAccountInactive,
			"account is not active",
			nil,
		)
	}

	// ============================================================
	// Check email verification
	// ============================================================

	if !currentUser.EmailVerified {
		return nil, apperror.New(
			apperror.CodeEmailNotVerified,
			"email is not verified",
			nil,
		)
	}

	profile, err := s.profileRepo.GetByUserID(ctx, currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	// ============================================================
	// Create authentication session
	// ============================================================

	session := &AuthSession{
		UserID:    currentUser.ID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(
			s.jwtManager.RefreshExpiration(),
		),
	}

	if err := s.authRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf(
			"create auth session: %w",
			err,
		)
	}

	// ============================================================
	// Generate access token
	// ============================================================

	accessToken, err := s.jwtManager.GenerateAccessToken(
		currentUser.ID,
		string(currentUser.Role),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	// ============================================================
	// Generate refresh token
	// ============================================================

	refreshToken, err := s.jwtManager.GenerateRefreshToken(
		currentUser.ID,
		session.ID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	// ============================================================
	// Hash refresh token
	// ============================================================

	refreshTokenHash, err := security.HashToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf(
			"hash refresh token: %w",
			err,
		)
	}

	session.RefreshTokenHash = refreshTokenHash

	// ============================================================
	// Save refresh token hash
	// ============================================================

	if err := s.authRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf(
			"update auth session: %w",
			err,
		)
	}

	// ============================================================
	// Success
	// ============================================================

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt: time.Now().Add(
			s.jwtManager.AccessExpiration(),
		),
		User: LoginUserResponse{
			ID:          currentUser.ID,
			Email:       currentUser.Email,
			Username:    currentUser.Username,
			DisplayName: profile.DisplayName,
			AvatarURL:   profile.AvatarURL,
			Role:        string(currentUser.Role),
			Status:      string(currentUser.Status),
			Timezone:    profile.Timezone,
		},
	}, nil
}

func (s *service) Refresh(
	ctx context.Context,
	req RefreshRequest,
) (*RefreshResponse, error) {

	// ============================================================
	// Validate refresh token
	// ============================================================

	refreshToken := strings.TrimSpace(req.RefreshToken)

	if refreshToken == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRefreshToken,
			"refresh token is required",
			nil,
		)
	}

	// ============================================================
	// Hash refresh token
	// ============================================================

	refreshTokenHash, err := security.HashToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("hash refresh token: %w", err)
	}

	// ============================================================
	// Find authentication session
	// ============================================================

	session, err := s.authRepo.GetByRefreshTokenHash(
		ctx,
		refreshTokenHash,
	)

	if err != nil {
		return nil, apperror.New(
			apperror.CodeInvalidRefreshToken,
			"invalid refresh token",
			err,
		)
	}

	// ============================================================
	// Check session revoked
	// ============================================================

	if session.RevokedAt != nil {
		return nil, apperror.New(
			apperror.CodeAuthSessionRevoked,
			"authentication session has been revoked",
			nil,
		)
	}

	// ============================================================
	// Check session expiration
	// ============================================================

	if time.Now().After(session.ExpiresAt) {
		return nil, apperror.New(
			apperror.CodeAuthSessionExpired,
			"authentication session has expired",
			nil,
		)
	}

	// ============================================================
	// Get user
	// ============================================================

	currentUser, err := s.userRepo.GetByID(
		ctx,
		session.UserID,
	)

	if err != nil {
		return nil, err
	}

	// ============================================================
	// Check user status
	// ============================================================

	if currentUser.Status != user.StatusActive {
		return nil, apperror.New(
			apperror.CodeUserInactive,
			"user account is inactive",
			nil,
		)
	}

	// ============================================================
	// Generate new access token
	// ============================================================

	accessToken, err := s.jwtManager.GenerateAccessToken(currentUser.ID, string(currentUser.Role))

	if err != nil {
		return nil, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	// ============================================================
	// Generate new refresh token
	// ============================================================

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(currentUser.ID, session.ID)

	if err != nil {
		return nil, fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	// ============================================================
	// Hash new refresh token
	// ============================================================

	newRefreshTokenHash, err := security.HashToken(
		newRefreshToken,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"hash new refresh token: %w", err,
		)
	}

	// ============================================================
	// Rotate refresh token
	// ============================================================

	session.RefreshTokenHash = newRefreshTokenHash
	session.ExpiresAt = time.Now().Add(
		s.jwtManager.RefreshExpiration(),
	)

	if err := s.authRepo.Update(
		ctx,
		session,
	); err != nil {
		return nil, fmt.Errorf(
			"update authentication session: %w",
			err,
		)
	}

	// ============================================================
	// Response
	// ============================================================

	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60,
	}, nil
}
