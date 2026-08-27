package relationship

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	userRepo "github.com/khaingminhtun/relio-backend/internal/features/user"
	redisinfra "github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"
)

type InvitationService interface {
	CreateExternalInvitation(
		ctx context.Context,
		relationshipID int64,
		invitedBy int64,
		req ExternalInvitationRequest,
	) (*ExternalInvitationResponse, error)

	GetInvitationByToken(
		ctx context.Context,
		token string,
	) (*ExternalInvitationResponse, error)

	AcceptInvitation(
		ctx context.Context,
		token string,
		userID int64,
	) error
	RejectInvitation(
		ctx context.Context,
		token string,
		userID int64,
	) error
	CancelInvitation(
		ctx context.Context,
		invitationID int64,
		userID int64,
	) error
}

type invitationService struct {
	userRepo       userRepo.Repository
	invitationRepo InvitationRepository
	memberRepo     RelationshipMemberRepository
	emailQueue     redisinfra.EmailQueue
	appBaseURL     string
}

func NewInvitationService(
	userRepo userRepo.Repository,
	invitationRepo InvitationRepository,
	memberRepo RelationshipMemberRepository,
	emailQueue redisinfra.EmailQueue,
	appBaseURL string,
) InvitationService {
	return &invitationService{
		userRepo:       userRepo,
		invitationRepo: invitationRepo,
		memberRepo:     memberRepo,
		emailQueue:     emailQueue,
		appBaseURL:     strings.TrimRight(appBaseURL, "/"),
	}
}

// ============================================================
// Create Invitation
// ============================================================

func (s *invitationService) CreateExternalInvitation(
	ctx context.Context,
	relationshipID int64,
	invitedBy int64,
	req ExternalInvitationRequest,
) (*ExternalInvitationResponse, error) {

	// ============================================================
	// Normalize email
	// ============================================================

	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	// ============================================================
	// Generate secure invitation token
	// ============================================================

	token, err := generateInvitationToken()
	if err != nil {
		return nil, fmt.Errorf(
			"generate invitation token: %w",
			err,
		)
	}

	// ============================================================
	// Hash token
	// ============================================================

	tokenHash, err := security.HashToken(token)
	if err != nil {
		return nil, fmt.Errorf(
			"hash invitation token: %w",
			err,
		)
	}

	// ============================================================
	// Expiration
	// ============================================================

	expiresAt := time.Now().Add(
		7 * 24 * time.Hour,
	)

	// ============================================================
	// Create invitation
	// ============================================================

	invitation := &Invitation{
		RelationshipID: relationshipID,
		InvitedBy:      invitedBy,
		Email:          email,
		TokenHash:      tokenHash,
		Status:         InvitationStatusPending,
		ExpiresAt:      expiresAt,
	}

	if err := s.invitationRepo.Create(
		ctx,
		invitation,
	); err != nil {
		return nil, fmt.Errorf(
			"create invitation: %w",
			err,
		)
	}

	// ============================================================
	// Build invitation URL
	// ============================================================

	inviteURL := fmt.Sprintf(
		"%s/invite/%s",
		s.appBaseURL,
		token,
	)

	// ============================================================
	// Create email job
	// ============================================================

	job := redisinfra.EmailJob{
		ID:        uuid.NewString(),
		To:        email,
		Subject:   "You've been invited to join a relationship on Relio",
		Template:  "relationship_invitation",
		CreatedAt: time.Now(),

		Data: map[string]any{
			"InviteURL": inviteURL,
			"ExpiresIn": "7 days",
		},
	}

	// ============================================================
	// Publish email job
	// ============================================================

	if err := s.emailQueue.Publish(ctx, job); err != nil {

		// The invitation was already created.
		// If email cannot be queued, remove the invitation.
		_ = s.invitationRepo.Delete(
			ctx,
			invitation.ID,
		)

		return nil, fmt.Errorf(
			"queue invitation email: %w",
			err,
		)
	}

	// ============================================================
	// Response
	// ============================================================

	return &ExternalInvitationResponse{
		ID:             invitation.ID,
		RelationshipID: invitation.RelationshipID,
		Status:         string(invitation.Status),
		InviteURL:      inviteURL,
		ExpiresAt:      invitation.ExpiresAt,
	}, nil
}

// ============================================================
// Token
// ============================================================

func (s *invitationService) GetInvitationByToken(
	ctx context.Context,
	token string,
) (*ExternalInvitationResponse, error) {

	token = strings.TrimSpace(token)

	if token == "" {
		return nil, apperror.New(
			apperror.CodeInvitationNotFound,
			"invitation not found",
			nil,
		)
	}

	// ============================================================
	// Hash token
	// ============================================================

	tokenHash, err := security.HashToken(token)
	if err != nil {
		return nil, fmt.Errorf(
			"hash invitation token: %w",
			err,
		)
	}

	// ============================================================
	// Find invitation
	// ============================================================

	invitation, err := s.invitationRepo.GetByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return nil, err
	}

	// ============================================================
	// Check expiration
	// ============================================================

	if time.Now().After(invitation.ExpiresAt) {
		return &ExternalInvitationResponse{
			ID:             invitation.ID,
			RelationshipID: invitation.RelationshipID,
			Status:         string(InvitationStatusExpired),
			ExpiresAt:      invitation.ExpiresAt,
			AcceptedAt:     invitation.AcceptedAt,
		}, nil
	}

	// ============================================================
	// Response
	// ============================================================

	return &ExternalInvitationResponse{
		ID:             invitation.ID,
		RelationshipID: invitation.RelationshipID,
		Status:         string(invitation.Status),
		ExpiresAt:      invitation.ExpiresAt,
		AcceptedAt:     invitation.AcceptedAt,
	}, nil
}

func (s *invitationService) AcceptInvitation(
	ctx context.Context,
	token string,
	userID int64,
) error {

	// ============================================================
	// Validate token
	// ============================================================

	token = strings.TrimSpace(token)

	if token == "" {
		return apperror.New(
			apperror.CodeInvitationNotFound,
			"invitation not found",
			nil,
		)
	}

	// ============================================================
	// Hash token
	// ============================================================

	tokenHash, err := security.HashToken(token)
	if err != nil {
		return fmt.Errorf(
			"hash invitation token: %w",
			err,
		)
	}

	// ============================================================
	// Get invitation
	// ============================================================

	invitation, err := s.invitationRepo.GetByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return err
	}

	// ============================================================
	// Check invitation status
	// ============================================================

	if invitation.Status != InvitationStatusPending {
		return apperror.New(
			apperror.CodeInvitationInvalid,
			"invitation is no longer pending",
			nil,
		)
	}

	// ============================================================
	// Check expiration
	// ============================================================

	if time.Now().After(invitation.ExpiresAt) {

		invitation.Status = InvitationStatusExpired

		if err := s.invitationRepo.Update(
			ctx,
			invitation,
		); err != nil {
			return fmt.Errorf(
				"update expired invitation: %w",
				err,
			)
		}

		return apperror.New(
			apperror.CodeInvitationExpired,
			"invitation has expired",
			nil,
		)
	}

	// ============================================================
	// Get authenticated user
	// ============================================================

	currentUser, err := s.userRepo.GetByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	// ============================================================
	// Verify invitation email
	// ============================================================

	if !strings.EqualFold(
		strings.TrimSpace(currentUser.Email),
		strings.TrimSpace(invitation.Email),
	) {
		return apperror.New(
			apperror.CodeInvitationEmailMismatch,
			"invitation email does not match authenticated user",
			nil,
		)
	}

	// ============================================================
	// Create relationship member
	// ============================================================

	now := time.Now()

	member := &RelationshipMember{
		RelationshipID: invitation.RelationshipID,
		UserID:         currentUser.ID,
		Role:           RelationshipMemberRoleMember,
		Status:         RelationshipMemberStatusActive,
		JoinedAt:       &now,
	}

	if err := s.memberRepo.Create(
		ctx,
		member,
	); err != nil {
		return fmt.Errorf(
			"create relationship member: %w",
			err,
		)
	}
	// ============================================================
	// Mark invitation accepted
	// ============================================================

	now = time.Now()

	invitation.Status = InvitationStatusAccepted
	invitation.AcceptedAt = &now

	if err := s.invitationRepo.Update(
		ctx,
		invitation,
	); err != nil {
		return fmt.Errorf(
			"update invitation: %w",
			err,
		)
	}

	return nil
}

func (s *invitationService) RejectInvitation(
	ctx context.Context,
	token string,
	userID int64,
) error {

	// ============================================================
	// Validate token
	// ============================================================

	token = strings.TrimSpace(token)

	if token == "" {
		return apperror.New(
			apperror.CodeInvitationNotFound,
			"invitation not found",
			nil,
		)
	}

	// ============================================================
	// Hash token
	// ============================================================

	tokenHash, err := security.HashToken(token)
	if err != nil {
		return fmt.Errorf(
			"hash invitation token: %w",
			err,
		)
	}

	// ============================================================
	// Get invitation
	// ============================================================

	invitation, err := s.invitationRepo.GetByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		return err
	}

	// ============================================================
	// Check invitation status
	// ============================================================

	if invitation.Status != InvitationStatusPending {
		return apperror.New(
			apperror.CodeInvitationInvalid,
			"invitation is no longer pending",
			nil,
		)
	}

	// ============================================================
	// Check expiration
	// ============================================================

	if time.Now().After(invitation.ExpiresAt) {

		invitation.Status = InvitationStatusExpired

		if err := s.invitationRepo.Update(
			ctx,
			invitation,
		); err != nil {
			return fmt.Errorf(
				"update expired invitation: %w",
				err,
			)
		}

		return apperror.New(
			apperror.CodeInvitationExpired,
			"invitation has expired",
			nil,
		)
	}

	// ============================================================
	// Get authenticated user
	// ============================================================

	currentUser, err := s.userRepo.GetByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	// ============================================================
	// Verify invitation email
	// ============================================================

	if !strings.EqualFold(
		strings.TrimSpace(currentUser.Email),
		strings.TrimSpace(invitation.Email),
	) {
		return apperror.New(
			apperror.CodeInvitationEmailMismatch,
			"invitation email does not match authenticated user",
			nil,
		)
	}

	// ============================================================
	// Reject invitation
	// ============================================================

	invitation.Status = InvitationStatusRejected

	if err := s.invitationRepo.Update(
		ctx,
		invitation,
	); err != nil {
		return fmt.Errorf(
			"update rejected invitation: %w",
			err,
		)
	}

	return nil
}

func (s *invitationService) CancelInvitation(
	ctx context.Context,
	invitationID int64,
	userID int64,
) error {

	// ============================================================
	// Get invitation
	// ============================================================

	invitation, err := s.invitationRepo.GetByID(
		ctx,
		invitationID,
	)
	if err != nil {
		return err
	}

	// ============================================================
	// Check ownership
	// ============================================================

	if invitation.InvitedBy != userID {
		return fmt.Errorf("you are not allowed this access %w", err)
	}

	// ============================================================
	// Check invitation status
	// ============================================================

	if invitation.Status != InvitationStatusPending {
		return apperror.New(
			apperror.CodeInvitationInvalid,
			"only pending invitations can be cancelled",
			nil,
		)
	}

	// ============================================================
	// Cancel invitation
	// ============================================================

	if err := s.invitationRepo.Delete(
		ctx,
		invitation.ID,
	); err != nil {
		return fmt.Errorf(
			"delete invitation: %w",
			err,
		)
	}

	return nil
}

func generateInvitationToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
