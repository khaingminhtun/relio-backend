package memory

import (
	"context"
	"errors"

	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"gorm.io/gorm"
)

// ============================================================
// Interface
// ============================================================

type MemoryRepository interface {
	Create(
		ctx context.Context,
		memory *Memory,
	) error

	GetByID(
		ctx context.Context,
		relationshipID int64,
		memoryID int64,
	) (*Memory, error)

	ListByRelationshipID(
		ctx context.Context,
		relationshipID int64,
	) ([]Memory, error)

	Update(
		ctx context.Context,
		memory *Memory,
	) error

	Delete(
		ctx context.Context,
		relationshipID int64,
		memoryID int64,
	) error
}

// ============================================================
// Implementation
// ============================================================

type memoryRepository struct {
	db *gorm.DB
}

func NewMemoryRepository(db *gorm.DB) MemoryRepository {
	return &memoryRepository{db: db}
}

func (r *memoryRepository) Create(
	ctx context.Context,
	memory *Memory,
) error {
	return transaction.DB(ctx, r.db).
		Create(memory).
		Error
}

func (r *memoryRepository) GetByID(
	ctx context.Context,
	relationshipID int64,
	memoryID int64,
) (*Memory, error) {
	var memory Memory

	err := transaction.DB(ctx, r.db).
		Where("id = ?", memoryID).
		Where("relationship_id = ?", relationshipID).
		First(&memory).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(
				apperror.CodeMemoryNotFound,
				"memory not found",
				err,
			)
		}

		return nil, err
	}

	return &memory, nil
}

func (r *memoryRepository) ListByRelationshipID(
	ctx context.Context,
	relationshipID int64,
) ([]Memory, error) {
	var memories []Memory

	err := transaction.DB(ctx, r.db).
		Where("relationship_id = ?", relationshipID).
		Order("memory_date DESC, created_at DESC").
		Find(&memories).
		Error

	if err != nil {
		return nil, err
	}

	return memories, nil
}

func (r *memoryRepository) Update(
	ctx context.Context,
	memory *Memory,
) error {
	return transaction.DB(ctx, r.db).
		Save(memory).
		Error
}

func (r *memoryRepository) Delete(
	ctx context.Context,
	relationshipID int64,
	memoryID int64,
) error {
	result := transaction.DB(ctx, r.db).
		Where("id = ?", memoryID).
		Where("relationship_id = ?", relationshipID).
		Delete(&Memory{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.New(
			apperror.CodeMemoryNotFound,
			"memory not found",
			gorm.ErrRecordNotFound,
		)
	}

	return nil
}
