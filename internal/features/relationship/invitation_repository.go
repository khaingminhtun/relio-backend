package relationship

import (
	"context"
	"errors"

	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type InvitationRepository interface {
	Create(
		ctx context.Context,
		invitation *Invitation,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*Invitation, error)

	GetByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (*Invitation, error)

	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]Invitation, error)

	FindPendingByEmail(
		ctx context.Context,
		relationshipID int64,
		email string,
	) (*Invitation, error)

	Update(
		ctx context.Context,
		invitation *Invitation,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(
	db *gorm.DB,
) InvitationRepository {
	return &invitationRepository{
		db: db,
	}
}

func (r *invitationRepository) Create(
	ctx context.Context,
	invitation *Invitation,
) error {
	return transaction.DB(ctx, r.db).
		Create(invitation).
		Error
}

func (r *invitationRepository) GetByID(
	ctx context.Context,
	id int64,
) (*Invitation, error) {
	var invitation Invitation

	err := transaction.DB(ctx, r.db).
		First(&invitation, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeInvitationNotFound,
				"invitation not found",
				err,
			)
		}

		return nil, err
	}

	return &invitation, nil
}

func (r *invitationRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*Invitation, error) {
	var invitation Invitation

	err := transaction.DB(ctx, r.db).
		Where("token_hash = ?", tokenHash).
		First(&invitation).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeInvitationNotFound,
				"invitation not found",
				err,
			)
		}

		return nil, err
	}

	return &invitation, nil
}

func (r *invitationRepository) ListByRelationshipID(
	ctx context.Context,
	relationshipID int64,
) ([]Invitation, error) {
	var invitations []Invitation

	err := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Order("created_at DESC").
		Find(&invitations).
		Error

	if err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *invitationRepository) FindPendingByEmail(
	ctx context.Context,
	relationshipID int64,
	email string,
) (*Invitation, error) {
	var invitation Invitation

	err := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Where("email = ?", email).
		Where(
			"status = ?",
			InvitationStatusPending,
		).
		First(&invitation).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeInvitationNotFound,
				"pending invitation not found",
				err,
			)
		}

		return nil, err
	}

	return &invitation, nil
}

func (r *invitationRepository) Update(
	ctx context.Context,
	invitation *Invitation,
) error {
	return transaction.DB(ctx, r.db).
		Save(invitation).
		Error
}

func (r *invitationRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	return transaction.DB(ctx, r.db).
		Delete(&Invitation{}, id).
		Error
}
