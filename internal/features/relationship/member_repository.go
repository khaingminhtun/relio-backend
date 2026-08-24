package relationship

import (
	"context"
	"errors"

	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type RelationshipMemberRepository interface {
	Create(
		ctx context.Context,
		member *RelationshipMember,
	) error

	GetByID(
		ctx context.Context,
		id int64,
	) (*RelationshipMember, error)

	GetByRelationshipAndUser(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) (*RelationshipMember, error)

	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]RelationshipMember, error)

	Update(
		ctx context.Context,
		member *RelationshipMember,
	) error

	Delete(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) error

	Exists(
		ctx context.Context,
		relationshipID int64,
		userID int64,
	) (bool, error)
}

type relationshipMemberRepository struct {
	db *gorm.DB
}

func NewRelationshipMemberRepository(
	db *gorm.DB,
) RelationshipMemberRepository {
	return &relationshipMemberRepository{
		db: db,
	}
}

func (r *relationshipMemberRepository) Create(
	ctx context.Context,
	member *RelationshipMember,
) error {
	return transaction.DB(ctx, r.db).
		Create(member).
		Error
}

func (r *relationshipMemberRepository) GetByID(
	ctx context.Context,
	id int64,
) (*RelationshipMember, error) {
	var member RelationshipMember

	err := transaction.DB(ctx, r.db).
		First(&member, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeRelationshipMemberNotFound,
				"relationship member not found",
				err,
			)
		}

		return nil, err
	}

	return &member, nil
}

func (r *relationshipMemberRepository) GetByRelationshipAndUser(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) (*RelationshipMember, error) {
	var member RelationshipMember

	err := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Where("user_id = ?", userID).
		First(&member).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeRelationshipMemberNotFound,
				"relationship member not found",
				err,
			)
		}

		return nil, err
	}

	return &member, nil
}

func (r *relationshipMemberRepository) ListByRelationshipID(
	ctx context.Context,
	relationshipID int64,
) ([]RelationshipMember, error) {
	var members []RelationshipMember

	err := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Where(
			"status = ?",
			RelationshipMemberStatusActive,
		).
		Find(&members).
		Error

	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *relationshipMemberRepository) Update(
	ctx context.Context,
	member *RelationshipMember,
) error {
	return transaction.DB(ctx, r.db).
		Save(member).
		Error
}

func (r *relationshipMemberRepository) Delete(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) error {
	result := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Where("user_id = ?", userID).
		Delete(&RelationshipMember{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeRelationshipMemberNotFound,
			"relationship member not found",
			gorm.ErrRecordNotFound,
		)
	}

	return nil
}

func (r *relationshipMemberRepository) Exists(
	ctx context.Context,
	relationshipID int64,
	userID int64,
) (bool, error) {
	var count int64

	err := transaction.DB(ctx, r.db).
		Model(&RelationshipMember{}).
		Where("relationship_id = ?", relationshipID).
		Where("user_id = ?", userID).
		Where(
			"status = ?",
			RelationshipMemberStatusActive,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
