package relationship

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	redisinfra "github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"

	"github.com/google/uuid"
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
}

type invitationService struct {
	invitationRepo InvitationRepository
	emailQueue     redisinfra.EmailQueue
	appBaseURL     string
}

func NewInvitationService(
	invitationRepo InvitationRepository,
	emailQueue redisinfra.EmailQueue,
	appBaseURL string,
) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
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

func generateInvitationToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
