package relationship

import (
	"context"
	"errors"

	transaction "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

type relationshipRepository struct {
	db *gorm.DB
}

func NewRelationshipRepository(
	db *gorm.DB,
) RelationshipRepository {
	return &relationshipRepository{
		db: db,
	}
}

func (r *relationshipRepository) Create(
	ctx context.Context,
	relationship *Relationship,
) error {
	return transaction.DB(ctx, r.db).
		Create(relationship).
		Error
}

func (r *relationshipRepository) GetByID(
	ctx context.Context,
	id int64,
) (*Relationship, error) {
	var relationship Relationship

	err := transaction.DB(ctx, r.db).
		First(&relationship, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeRelationshipNotFound,
				"relationship not found",
				err,
			)
		}

		return nil, err
	}

	return &relationship, nil
}

func (r *relationshipRepository) ListByUserID(
	ctx context.Context,
	userID int64,
) ([]Relationship, error) {
	var relationships []Relationship

	err := transaction.DB(ctx, r.db).
		Joins(
			"JOIN relationship_members rm ON rm.relationship_id = relationships.id",
		).
		Where("rm.user_id = ?", userID).
		Where(
			"rm.status = ?",
			RelationshipMemberStatusActive,
		).
		Find(&relationships).
		Error

	if err != nil {
		return nil, err
	}

	return relationships, nil
}

func (r *relationshipRepository) Update(
	ctx context.Context,
	relationship *Relationship,
) error {
	return transaction.DB(ctx, r.db).
		Save(relationship).
		Error
}

func (r *relationshipRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	result := transaction.DB(ctx, r.db).
		Delete(&Relationship{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
